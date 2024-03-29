package roll

import (
	"fmt"
	"log"
	"math/rand"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/zekrotja/ken"
)

type Roll struct {
	models *data.Models
}

var _ ken.SlashCommand = (*Roll)(nil)

func (v *Roll) Initialize(models *data.Models) {
	v.models = models
}

// Description implements ken.SlashCommand.
func (*Roll) Description() string {
	return "roll different items and abilities"
}

// Name implements ken.SlashCommand.
func (*Roll) Name() string {
	return "roll"
}

// Options implements ken.SlashCommand.
func (*Roll) Options() []*discordgo.ApplicationCommandOption {
	minRarityOpts := []*discordgo.ApplicationCommandOptionChoice{}
	for _, r := range rarityPriorities {
		minRarityOpts = append(minRarityOpts, &discordgo.ApplicationCommandOptionChoice{
			Name:  r,
			Value: r,
		})
	}

	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "item",
			Description: "roll an item.",
			Options: []*discordgo.ApplicationCommandOption{
				discord.IntCommandArg("luck", "luck level", true),
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "minimum rarity to roll for",
					Required:    false,
					Choices:     minRarityOpts,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "exact",
					Description: "exact rarity to roll for",
					Required:    false,
					// drop unique and mythical (currently not a possible roll)
					Choices: minRarityOpts[:len(minRarityOpts)-1],
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "ability",
			Description: "roll an any ability.",
			Options: []*discordgo.ApplicationCommandOption{
				discord.IntCommandArg("luck", "luck level", true),
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "minimum rarity to roll for",
					Required:    false,
					// drop unique and mythical (currently not a possible roll)
					Choices: minRarityOpts[:len(minRarityOpts)-1],
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "exact",
					Description: "exact rarity to roll for",
					Required:    false,
					// drop unique and mythical (currently not a possible roll)
					Choices: minRarityOpts[:len(minRarityOpts)-1],
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "care_package",
			Description: "roll a random item and any ability",
			Options: []*discordgo.ApplicationCommandOption{
				discord.IntCommandArg("luck", "luck level", true),
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "minimum rarity to roll for",
					Required:    false,
					Choices:     minRarityOpts,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "item_rain",
			Description: "Make it rain 1-3 random items",
			Options: []*discordgo.ApplicationCommandOption{
				discord.IntCommandArg("luck", "luck level", true),
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "minimum rarity to roll for",
					Required:    false,
					Choices:     minRarityOpts,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "power_drop",
			Description: "roll an any ability.",
			Options: []*discordgo.ApplicationCommandOption{
				discord.IntCommandArg("luck", "luck level", true),
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "rarity",
					Description: "minimum rarity to roll for",
					Required:    false,
					// drop unique and mythical (currently not a possible roll)
					Choices: minRarityOpts[:len(minRarityOpts)-1],
				},
			},
		},
	}
}

// Run implements ken.SlashCommand.
func (r *Roll) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "item", Run: r.rollItem},
		ken.SubCommandHandler{Name: "ability", Run: r.rollAbility},
		ken.SubCommandHandler{Name: "care_package", Run: r.rollCarePackage},
		ken.SubCommandHandler{Name: "item_rain", Run: r.rollItemRain},
		ken.SubCommandHandler{Name: "power_drop", Run: r.rollPowerDrop},
	)
	return err
}

// Version implements ken.SlashCommand.
func (*Roll) Version() string {
	return "1.0.0"
}

func (r *Roll) rollItem(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}
	luck := c.Options().GetByName("luck").IntValue()
	rOpt, ok := c.Options().GetByNameOptional("rarity")
	eOpt, eOk := c.Options().GetByNameOptional("exact")
	byRarity := ok
	rarity := ""
	if byRarity {
		rarity = rOpt.StringValue()
	}

	if eOk {
		item, err := r.models.Items.GetRandomByRarity(eOpt.StringValue())
		if err != nil {
			return discord.ErrorMessage(c, "Error getting item", err.Error())
		}
		return c.RespondEmbed(&discordgo.MessageEmbed{
			Title:       "Item Roll",
			Description: fmt.Sprintf("You rolled a %s item: %s", eOpt.StringValue(), item.Name),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  item.Name,
					Value: item.Description,
				},
			},
			Color: discord.ColorThemeWhite,
		})
	}

	rarityRoll := rollAtRarity(float64(luck), rarityPriorities)
	if byRarity {
		minOption := slices.Index(rarityPriorities, rarity)
		if minOption == -1 {
			return discord.ErrorMessage(c, "Invalid rarity type", fmt.Sprintf("%s is not a valid rarity", rarity))
		}
		for minOption > slices.Index(rarityPriorities, rarityRoll) {
			// reroll if the roll is less than the minimum rarity
			rarityRoll = rollAtRarity(float64(luck), rarityPriorities[:minOption+1])
		}
	}
	item, err := r.models.Items.GetRandomByRarity(rarityRoll)
	if err != nil {
		return discord.ErrorMessage(c, "Error getting item", err.Error())
	}

	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Item Roll",
		Description: fmt.Sprintf("You rolled a %s item: %s", rarityRoll, item.Name),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  item.Name,
				Value: item.Description,
			},
		},
		Color: discord.ColorThemeWhite,
	})
}

func (r *Roll) rollAbility(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}
	argLuck := c.Options().GetByName("luck").IntValue()
	argRarity, ok := c.Options().GetByNameOptional("rarity")
	argExact, eok := c.Options().GetByNameOptional("exact")
	minRarity := ""
	if ok {
		minRarity = argRarity.StringValue()
	}

	if eok {
		minRarity = argExact.StringValue()
		ability, err := r.models.Abilities.GetRandomByRarity(minRarity)
		if err != nil {
			log.Println(err)
			return discord.AlexError(c, "Lol idk")
		}
		return c.RespondEmbed(&discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Rolled ability '%s' - %s", ability.Name, ability.Rarity),
			Description: ability.Description,
		})
	}

	var ability *data.Ability
	if minRarity != "" {
		rIdx := slices.Index(rarityPriorities, minRarity)
		if rIdx == -1 {
			return discord.ErrorMessage(c, "Invalid rarity", fmt.Sprintf("%s is not a valid rarity", minRarity))
		}
		// Drop the last 2 rarities, as we ability only has rarities up to legendary
		choices := rarityPriorities[rIdx:(len(rarityPriorities) - 1)]
		rarityRoll := rollAtRarity(float64(argLuck), choices)
		ability, err = r.models.Abilities.GetRandomByRarity(rarityRoll)
		if err != nil {
			log.Println(err)
			return discord.AlexError(c, "Lol idk")
		}
	} else {
		rarityRoll := rollAtRarity(float64(argLuck), rarityPriorities)
		ability, err = r.models.Abilities.GetRandomByRarity(rarityRoll)
		if err != nil {
			log.Println(err)
			return discord.AlexError(c, "Lol idk")
		}
	}

	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Rolled ability '%s' - %s", ability.Name, ability.Rarity),
		Description: ability.Description,
	})
}

// rollCarePackage rolls a random item and ability
func (r *Roll) rollCarePackage(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}
	argLuck := c.Options().GetByName("luck").IntValue()
	item, err := r.models.Items.GetRandomByRarity(rollAtRarity(float64(argLuck), rarityPriorities))
	if err != nil {
		log.Println(err)
		return discord.AlexError(c, "Lol idk")
	}

	var ability *data.Ability
	rarityRoll := rollAtRarity(float64(argLuck), rarityPriorities)
	if rarityRoll == "mythical" || rarityRoll == "unique" {
		rarityRoll = "legendary"
	}
	ability, err = r.models.Abilities.GetRandomByRarity(rarityRoll)
	if err != nil {
		log.Println(err)
		return discord.AlexError(c, "Lol idk")
	}
	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Care Package",
		Description: "You rolled a care package!",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Item",
				Value: fmt.Sprintf("%s - %s\n%s", item.Name, item.Rarity, item.Description),
			},
			{
				Name:  "Ability",
				Value: fmt.Sprintf("%s - %s\n%s", ability.Name, ability.Rarity, ability.Description),
			},
		},
	})
}

func (r *Roll) rollItemRain(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}

	luck := c.Options().GetByName("luck").IntValue()
	rOpt, ok := c.Options().GetByNameOptional("rarity")
	minRarity := ""
	if ok {
		minRarity = rOpt.StringValue()
	}

	items := []*data.Item{}
	// Roll 1-3 items
	rolls := 1 + rand.Intn(3)
	for i := 0; i < rolls; i++ {
		rarityRoll := rollAtRarity(float64(luck), rarityPriorities)
		if minRarity != "" {
			minOption := slices.Index(rarityPriorities, minRarity)
			if minOption == -1 {
				return discord.ErrorMessage(c, "Invalid rarity type", fmt.Sprintf("%s is not a valid rarity", minRarity))
			}
			for minOption > slices.Index(rarityPriorities, rarityRoll) {
				// reroll if the roll is less than the minimum rarity
				rarityRoll = rollAtRarity(float64(luck), rarityPriorities[:minOption+1])
			}
		}
		item, err := r.models.Items.GetRandomByRarity(rarityRoll)
		if err != nil {
			return discord.ErrorMessage(c, "Error getting item", err.Error())
		}
		items = append(items, item)
	}
	fields := []*discordgo.MessageEmbedField{}
	for _, item := range items {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s - %s", item.Name, item.Rarity),
			Value:  item.Description,
			Inline: true,
		})
	}
	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Item Rain",
		Description: fmt.Sprintf("You rolled %d item(s)", rolls),
		Fields:      fields,
	})
}

func (r *Roll) rollPowerDrop(c ken.SubCommandContext) (err error) {
	if err = c.Defer(); err != nil {
		log.Println(err)
		return err
	}
	argLuck := c.Options().GetByName("luck").IntValue()
	argRarity, ok := c.Options().GetByNameOptional("rarity")
	minRarity := ""
	if ok {
		minRarity = argRarity.StringValue()
	}

	var ability *data.Ability
	if minRarity != "" {
		rIdx := slices.Index(rarityPriorities, minRarity)
		if rIdx == -1 {
			return discord.ErrorMessage(c, "Invalid rarity", fmt.Sprintf("%s is not a valid rarity", minRarity))
		}
		// Drop the last 2 rarities, as we ability only has rarities up to legendary
		choices := rarityPriorities[rIdx:(len(rarityPriorities) - 1)]
		rarityRoll := rollAtRarity(float64(argLuck), choices)
		ability, err = r.models.Abilities.GetRandomByRarity(rarityRoll)
		if err != nil {
			log.Println(err)
			return discord.AlexError(c, "Lol idk")
		}
	} else {
		rarityRoll := rollAtRarity(float64(argLuck), rarityPriorities)
		ability, err = r.models.Abilities.GetRandomByRarity(rarityRoll)
		if err != nil {
			log.Println(err)
			return discord.AlexError(c, "Lol idk")
		}
	}

	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title: "Power Drop Rolled!",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("%s - %s", ability.Name, ability.Rarity),
				Value: ability.Description,
			},
		},
	})
}
