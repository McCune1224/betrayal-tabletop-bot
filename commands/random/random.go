package random

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/zekrotja/ken"
)

type Random struct {
	models *data.Models
}

var _ ken.SlashCommand = (*Random)(nil)

func (v *Random) Initialize(models *data.Models) {
	v.models = models
}

// Description implements ken.SlashCommand.
func (r *Random) Description() string {
	return "Get random outcomes for things like players, roles, etc."
}

// Name implements ken.SlashCommand.
func (r *Random) Name() string {
	return "random"
}

// Options implements ken.SlashCommand.
func (r *Random) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "role",
			Description: "Get a random role.",
		},
	}
}

// Run implements ken.SlashCommand.
func (r *Random) Run(ctx ken.Context) (err error) {
	err = ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "role", Run: r.randomRole},
	)
	return err
}

// Version implements ken.SlashCommand.
func (r *Random) Version() string {
	return "1.0.0"
}

func (r *Random) randomRole(c ken.SubCommandContext) error {
	role, err := r.models.Roles.GetRandomRole()
	if err != nil {
		log.Println(err)
		return discord.AlexError(c, "idk lol")
	}

	return c.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Random Role Rolled",
		Description: fmt.Sprintf("Role: %s", role.Name),
	})
}
