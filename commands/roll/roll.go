package roll

import (
	"fmt"
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
			},
		},
		// {
		// 	Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 	Name:        "ability",
		// 	Description: "roll an ability.",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		discord.IntCommandArg("luck", "luck level", true),
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionString,
		// 			Name:        "rarity",
		// 			Description: "minimum rarity to roll for",
		// 			Required:    true,
		// 			Choices:     minRarityOpts,
		// 		},
		// 	},
		// },
	}
}

// Run implements ken.SlashCommand.
func (r *Roll) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "item", Run: r.rollItem},
	)
	return err
}

// Version implements ken.SlashCommand.
func (*Roll) Version() string {
	return "1.0.0"
}

func (r *Roll) rollItem(c ken.SubCommandContext) (err error) {
	luck := c.Options().GetByName("luck").IntValue()
	rOpt, ok := c.Options().GetByNameOptional("rarity")
	byRarity := ok
	rarity := ""
	if byRarity {
		rarity = rOpt.StringValue()
	}

	rarityRoll := rollAtRarity(float64(luck), rarityPriorities)
	if byRarity {
		minOption := slices.Index(rarityPriorities, rarity)
		if minOption == -1 {
			return discord.ErrorMessage(c, "Invalid rarity type", fmt.Sprintf("%s is not a valid rarity", rOpt))
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
        Name: item.Name,
        Value: item.Description,
      },
		},
		Color: discord.ColorThemeWhite,
	})
}
