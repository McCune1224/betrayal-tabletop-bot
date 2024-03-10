package view

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/mccune1224/betrayal-tabletop-bot/util"
	"github.com/zekrotja/ken"
)

const infinity = "∞"

type View struct {
	models *data.Models
}

func (v *View) Initialize(models *data.Models) {
	v.models = models
}

var _ ken.SlashCommand = (*View)(nil)

// Description implements ken.SlashCommand.
func (*View) Description() string {
	return "is this thing on?"
}

// Name implements ken.SlashCommand.
func (*View) Name() string {
	return "view"
}

// Options implements ken.SlashCommand.
func (v *View) Options() []*discordgo.ApplicationCommandOption {
	statusChoices := []*discordgo.ApplicationCommandOptionChoice{}
	statuses, _ := v.models.Statuses.GetAll()
	for _, s := range statuses {
		statusChoices = append(statusChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  s.Name,
			Value: s.Name,
		})
	}

	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "role",
			Description: "View a role",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("name", "Name of the role", true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "ability",
			Description: "View an ability",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("name", "Name of the role", true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "passive",
			Description: "View a passive",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("name", "Name of the passive", true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "item",
			Description: "View an item",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("name", "Name of the role", true),
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "status",
			Description: "View a status",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Name of the status",
					Required:    true,
					Choices:     statusChoices,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "all_statuses",
			Description: "View all status",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "duel",
			Description: "View how minigame 'Duel' works",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "care_package",
			Description: "Learn about 'Care Package' Drops",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "night_order",
			Description: "View the order of how actions are processed during the night",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "actions",
			Description: "View actions that can be taken during a cycle",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "deceased",
			Description: "View all players marked as deceased",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "daily_events",
			Description: "View all daily events.",
		},
	}
}

// Run implements ken.SlashCommand.
func (v *View) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "role", Run: v.viewRole},
		ken.SubCommandHandler{Name: "ability", Run: v.viewAbility},
		ken.SubCommandHandler{Name: "passive", Run: v.viewPassive},
		ken.SubCommandHandler{Name: "item", Run: v.viewItem},
		ken.SubCommandHandler{Name: "status", Run: v.viewStatus},
		ken.SubCommandHandler{Name: "all_statuses", Run: v.viewAllStatuses},
		ken.SubCommandHandler{Name: "duel", Run: v.viewDuel},
		ken.SubCommandHandler{Name: "care_package", Run: v.viewCarePackage},
		ken.SubCommandHandler{Name: "night_order", Run: v.viewNightOrder},
		ken.SubCommandHandler{Name: "actions", Run: v.viewActions},
		ken.SubCommandHandler{Name: "deceased", Run: v.viewDeceased},
		ken.SubCommandHandler{Name: "daily_events", Run: v.viewDailyEvents},
	)
	return err
}

func (v *View) viewRole(ctx ken.SubCommandContext) error {
	name := ctx.Options().GetByName("name").StringValue()

	roles, err := v.models.Roles.GetAll()
	if err != nil {
		return discord.AlexError(ctx, "Idk lol")
	}
	roleNames := []string{}
	for _, v := range roles {
		roleNames = append(roleNames, v.Name)
	}

	roleName := util.FuzzyFind(name, roleNames)
	var role *data.Role
	for _, r := range roles {
		if r.Name == roleName {
			role = r
			break
		}
	}

	embed, err := v.roleEmbed(role)
	if err != nil {
		log.Println(err)
		return discord.AlexError(ctx, "Idk lol")
	}

	return ctx.RespondEmbed(embed)
}

// Version implements ken.SlashCommand.
func (*View) Version() string {
	return "1.0.0"
}

func (v *View) roleEmbed(role *data.Role) (*discordgo.MessageEmbed, error) {
	color := 0x000000
	switch role.Alignment {
	case "Lawful":
		color = discord.ColorThemeGreen
	case "Chaotic":
		color = discord.ColorThemeRed
	case "Outlander":
		color = discord.ColorThemeYellow
	}

	var abilities []*data.Ability
	var passives []*data.Passive

	for _, id := range role.AbilityIDs {
		abilitie, err := v.models.Abilities.Get(int(id))
		if err != nil {
			return nil, err
		}
		abilities = append(abilities, abilitie)
	}

	for _, id := range role.PassiveIDs {
		passive, err := v.models.Passives.Get(int(id))
		if err != nil {
			return nil, err
		}
		passives = append(passives, passive)
	}

	var embededAbilitiesFields []*discordgo.MessageEmbedField
	embededAbilitiesFields = append(embededAbilitiesFields, &discordgo.MessageEmbedField{
		Name:   "\n\n" + discord.Underline("Abilities") + "\n",
		Value:  "",
		Inline: false,
	})
	for _, ability := range abilities {
		title := ability.Name
		fStr := "%s [%d] - %s"

		categories := strings.Join(ability.Categories, ", ")
		if ability.Charges == -1 {
			title = fmt.Sprintf("%s [%s] - %s", ability.Name, infinity, categories)
		} else {
			title = fmt.Sprintf(fStr, ability.Name, ability.Charges, categories)
		}
		embededAbilitiesFields = append(
			embededAbilitiesFields,
			&discordgo.MessageEmbedField{
				Name:   title,
				Value:  ability.Description,
				Inline: false,
			},
		)
	}
	embededAbilitiesFields = append(embededAbilitiesFields, &discordgo.MessageEmbedField{
		Name:  "\n\n",
		Value: "\n",
	})

	var embededPassiveFields []*discordgo.MessageEmbedField
	embededAbilitiesFields = append(embededAbilitiesFields, &discordgo.MessageEmbedField{
		Name:   discord.Underline("Passives"),
		Value:  "",
		Inline: false,
	})

	for _, passives := range passives {
		embededPassiveFields = append(
			embededPassiveFields,
			&discordgo.MessageEmbedField{
				Name:   passives.Name,
				Value:  passives.Description + "\n",
				Inline: false,
			},
		)
	}

	embed := &discordgo.MessageEmbed{
		Title:  role.Name,
		Color:  color,
		Fields: append(embededAbilitiesFields, embededPassiveFields...),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Alignment: " + role.Alignment,
		},
	}
	return embed, nil
}

func (v *View) itemEmbed(item *data.Item) (*discordgo.MessageEmbed, error) {
	categories := strings.Join(item.Categories, ", ")
	fields := []*discordgo.MessageEmbedField{}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Categories",
		Value:  categories,
		Inline: false,
	})

	costStr := ""
	if item.Cost == -1 {
		costStr = infinity
	} else {
		costStr = fmt.Sprint(item.Cost)
	}

	embed := &discordgo.MessageEmbed{
		Title:       item.Name,
		Description: item.Description,
		Fields:      fields,
		Color:       discord.ComponentColorByRarity(item.Rarity),
		Footer: &discordgo.MessageEmbedFooter{
			Text: costStr,
		},
	}
	return embed, nil
}

func (v *View) viewItem(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}
	name := c.Options().GetByName("name").StringValue()
	items, err := v.models.Items.GetAll()
	if err != nil {
		return discord.AlexError(c, "Idk lol")
	}

	itemNames := []string{}
	for _, v := range items {
		itemNames = append(itemNames, v.Name)
	}
	var item *data.Item
	itemName := util.FuzzyFind(name, itemNames)
	for _, v := range items {
		if v.Name == itemName {
			item = v
			break
		}
	}

	embed, err := v.itemEmbed(item)
	if err != nil {
		log.Println(err)
		return discord.AlexError(c, "Idk lol")
	}
	return c.RespondEmbed(embed)
}

func (v *View) viewStatus(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}
	name := c.Options().GetByName("name").StringValue()
	status, err := v.models.Statuses.GetByName(name)
	if err != nil {
		log.Println(err)
		return discord.AlexError(c, "Idk lol")
	}

	embed := &discordgo.MessageEmbed{
		Title:       status.Name,
		Description: status.Description,
		Color:       discord.ColorThemeWhite,
	}
	return c.RespondEmbed(embed)
}

func (v *View) viewDuel(ctx ken.SubCommandContext) (err error) {
	gameText := []string{
		fmt.Sprintf("In %s players will present one out of nine number tiles and the player who presented the higher numbered tile wins.", discord.Bold("Black and White")),
		fmt.Sprintf("The players will each receive 99 tiles are divided into black and white colors. %s", discord.Bold("Even numbers 0, 2, 4, 6 and 8 are black. Odd numbers 1, 3, 5 and 7 are white.\n")),
		fmt.Sprintf("The starting player will first choose a number from 0 to 8 (selecting the number in their confessional), The host will announce publicly %s. The following player will then present their tile. Only hosts will see numbers used, and the player who put a higher number will win and gain one point. %s.", discord.Bold("what color was used"), discord.Bold("Used numbers will not be revealed even after the results are announced")),
		"Example: Sophia begins the game and uses a 3. The host will announce: Sophia has used a white tile. Lindsey will place a black tile, a 0. Host will announce a black tile was used. Host will announce that Sophia has won. Both tiles/numbers are taken away and a new round begins, the winner goes first in presenting the tile for the next round. Lindsey can infer very little from her loss because any white tile can beat a black 0, but Sophia will know that she used either a 0 or a 2 based on her win.",
		"The player with more points after 9th round will win, the loser will be eliminated.",
	}

	return ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Game Duel - Black and White",
		Color:       discord.ColorThemePearl,
		Description: gameText[0],
		Fields: []*discordgo.MessageEmbedField{
			{
				Value: gameText[1],
			},
			{
				Value: gameText[2],
			},
			{
				Value: gameText[3],
			},
			{
				Value: gameText[4],
			},
		},
	})
}

func (v *View) viewPassive(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		return err
	}
	nameArg := ctx.Options().GetByName("name").StringValue()
	passives, err := v.models.Passives.GetAll()
	if err != nil {
		ctx.RespondError("Unable to find Passive",
			fmt.Sprintf("Unable to find Passive: %s", nameArg),
		)
		return err
	}

	passiveNames := []string{}
	for _, v := range passives {
		passiveNames = append(passiveNames, v.Name)
	}
	passiveName := util.FuzzyFind(nameArg, passiveNames)

	var passive *data.Passive
	for _, v := range passives {
		if v.Name == passiveName {
			passive = v
			break
		}
	}
	associatedRoles, err := v.models.Roles.GetAllByPassiveID(passive.ID)
	if err != nil {
		log.Println(err)
		discord.ErrorMessage(ctx,
			"Error Finding Role",
			fmt.Sprintf("Unable to find Associated Role for Ability: %s", nameArg))
		return err
	}

	rolenames := []string{}
	for _, v := range associatedRoles {
		rolenames = append(rolenames, v.Name)
	}

	passiveEmbed := &discordgo.MessageEmbed{
		Title:       passive.Name,
		Description: passive.Description,
		Color:       discord.ColorThemeWhite,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Associated Roles: %s", strings.Join(rolenames, ", ")),
		},
	}
	return ctx.RespondEmbed(passiveEmbed)
}

func (v *View) viewAbility(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}
	nameArg := ctx.Options().GetByName("name").StringValue()
	// ability, err := v.models.Abilities.GetByFuzzy(nameArg)

	abilities, err := v.models.Abilities.GetAll()
	if err != nil {
		discord.ErrorMessage(ctx,
			"Error Finding Ability",
			fmt.Sprintf("Unable to find Ability: %s", nameArg),
		)
		return err
	}

	abilityNames := []string{}
	for _, v := range abilities {
		abilityNames = append(abilityNames, v.Name)
	}

	var ability *data.Ability

	abilityName := util.FuzzyFind(nameArg, abilityNames)
	for _, v := range abilities {
		if v.Name == abilityName {
			ability = v
			break
		}
	}

	associatedRoles, err := v.models.Roles.GetAllByAbilityID(ability.ID)
	if err != nil {
		log.Println(err)
		return discord.ErrorMessage(ctx,
			"Error Finding Role",
			fmt.Sprintf("Unable to find Associated Role for Ability: %s", nameArg))
	}

	roleNames := []string{}
	for _, v := range associatedRoles {
		roleNames = append(roleNames, v.Name)
	}

	abilityEmbed := &discordgo.MessageEmbed{
		Title:       ability.Name,
		Description: ability.Description,
		Color:       discord.ComponentColorByRarity(ability.Rarity),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Categories",
				Value:  strings.Join(ability.Categories, ", "),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Associated Roles: %s", strings.Join(roleNames, ", ")),
		},
	}

	return ctx.RespondEmbed(abilityEmbed)
}

func (v *View) viewCarePackage(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}

	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Care Package",
		Description: "Granted 1 random item and Any Ability, higher luck increases the chance of higher rarity of items and abilities",
	})
}

func (v *View) viewAllStatuses(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}
	statuses, err := v.models.Statuses.GetAll()
	if err != nil {
		return discord.AlexError(c, "idk lol")
	}

	msg := discordgo.MessageEmbed{
		Title:       "Statuses",
		Description: "All Statuses",
	}

	for _, s := range statuses {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  s.Name,
			Value: s.Description,
		})
	}

	return c.RespondEmbed(&msg)
}

func (v *View) viewNightOrder(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}
	order := []string{"Alteration", "Reactive", "Redirection", "Protection", "Visit Blocking", "Vote Manipulation", "Support", "Debuff", "Theft", "Healing", "Destruction", "Killing"}
	msg := discordgo.MessageEmbed{}
	msg.Title = "Night Order"
	msg.Description = "The order in which actions are processed during the night"
	for i, o := range order {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s", i+1, o),
			Value: "",
		})
	}
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Value: "During the Night, players are able to do actions if they have enough action slots remaining to permit it (each player may do max three per cycle). All AAs, Abilities, and Items are codified by one of the twelve types above.",
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Value: "The Judge ability/AA Forced Trial, for example, is codified as such: (Alteration/Neutral/Non-visiting/Night). This means that after the Host has woken everyone and seen who is doing what, this ability will trigger very early in the Night Order (likely, first, although if multiple people do Alteration abilities, it will be rng'd which occurs first, in the case of Forced Trial, whether it is the first or third alteration ability to be processed will likely not matter, but for other types it can be more important, the order in which they are processed).",
	})
	return c.RespondEmbed(&msg)
}

func (v *View) viewActions(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}

	msg := discordgo.MessageEmbed{}
	msg.Title = "Actions"
	msg.Description = "Actions that can be taken during a cycle"

	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Players can take up to three actions per cycle (Day 1 + Night 1 = one cycle).",
		Value: "Types of actions include:\n",
	})

	actions := []string{"Using an ability from role card", "Using an AA (Any Ability)", "Using an item", "Purchasing an item", "Sending an item to another player", "Switching seats with another player", "Creating an alliance/coalition", "Accepting an alliance/coalition", "Leaving an alliance/coalition", "Using a Betrayal token"}
	for i, a := range actions {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%d. %s", i+1, a),
			Value: "",
		})
	}

	return c.RespondEmbed(&msg)
}

func (v *View) viewDeceased(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}
	deceased := discord.GetMembersWithRoleName(c.GetSession(), c.GetEvent(), "Deceased")

	msg := discordgo.MessageEmbed{}
	msg.Title = fmt.Sprintf("%d Deceased Players", len(deceased))
	for _, m := range deceased {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  m.User.Username,
			Value: "",
		})
	}

	return c.RespondEmbed(&msg)
}

func (v *View) viewDailyEvents(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		return err
	}
	msg := discordgo.MessageEmbed{}
	msg.Title = "Daily Events"
	msg.Description = "In addition to the day specific events below, a player will be voted out every day. Following the day phase, the game enters night. Abilities that are codified as Night, can only be used during this period."

	events := []string{
		"Day 1 (game start): Care package (one AA + one item, rng’d)",
		"Day 2: Valentines Day. Each player must choose another, the most voted gets elimination immunity (cannot vote for yourself).",
		"Day 3: Item Rain (1-3 random items, rarity is dependent on luck)",
		"Day 4: Power Drop (All players gain an AA, rarity is dependent on luck)",
		"Day 5: Rock Paper Scissors tournament. Everyone plays rock, paper, scissors. Winner gets a special prize. Host creates the prize package prior to game start.",
		"Day 6: Item Rain (1-3 random items, rarity is dependent on luck)",
		"Day 7: Coin Rain (All Players get x2 coins) & Power Drop (All players gain an AA, rarity is dependent on luck)",
		"Day 8:  Valentines Day. Each player must choose another, the most voted gets elimination immunity (cannot vote for yourself).",
		"Day 9: Item Rain (1-3 random items, rarity is dependent on luck). Revival opportunities for graveyard end Day 9.",
		"Day 10:  Duels. Choose to challenge someone to a duel. Life is at stake.",
	}
	for _, e := range events {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  e,
			Value: "",
		})
	}
	return c.RespondEmbed(&msg)
}
