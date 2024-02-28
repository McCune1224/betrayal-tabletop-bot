package view

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/zekrotja/ken"
)

const infinity = "∞"

type View struct {
	models *data.Models
}

func (v *View) Initialize(models *data.Models) {
	v.models = models
}

// Description implements ken.SlashCommand.
func (*View) Description() string {
	return "is this thing on?"
}

// Name implements ken.SlashCommand.
func (*View) Name() string {
	return "view"
}

// Options implements ken.SlashCommand.
func (*View) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "role",
			Description: "View a role",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("name", "Name of the role", true),
			},
		},
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "ability",
		// 	Description: "View an ability",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		discord.StringCommandArg("name", "Name of the role", true),
		// 	},
		// },
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "passive",
		// 	Description: "View a passive",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		discord.StringCommandArg("name", "Name of the role", true),
		// 	},
		// },
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "item",
		// 	Description: "View an item",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		discord.StringCommandArg("name", "Name of the role", true),
		// 	},
		// },
	}
}

// Run implements ken.SlashCommand.
func (v *View) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "role", Run: v.viewRole},
		// ken.SubCommandHandler{Name: "ability", Run: v.viewAbility},
		// ken.SubCommandHandler{Name: "perk", Run: v.viewPerk},
		// ken.SubCommandHandler{Name: "item", Run: v.viewItem},
		// ken.SubCommandHandler{Name: "status", Run: v.viewStatus},
		// ken.SubCommandHandler{Name: "duel", Run: v.viewDuel},
	)
	return err
}

func (v *View) viewRole(ctx ken.SubCommandContext) error {
	name := ctx.Options().GetByName("name").StringValue()
	role, err := v.models.Roles.GetByName(name)
	if err != nil {
		return discord.AlexError(ctx, "Idk lol")
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

var _ ken.SlashCommand = (*View)(nil)

func (v *View) roleEmbed(role *data.Role) (*discordgo.MessageEmbed, error) {
	color := 0x000000
	switch role.Alignment {
	case "LAWFUL":
		color = discord.ColorThemeGreen
	case "CHAOTIC":
		color = discord.ColorThemeRed
	case "OUTLANDER":
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
