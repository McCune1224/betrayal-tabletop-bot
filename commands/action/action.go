package action

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/mccune1224/betrayal-tabletop-bot/util"
	"github.com/zekrotja/ken"
)

const funnelChannelID = "1211657721376411648"

type Action struct {
	models *data.Models
}

var _ ken.SlashCommand = (*Action)(nil)

func (a *Action) Initialize(models *data.Models) {
	a.models = models
}

// Description implements ken.SlashCommand.
func (*Action) Description() string {
	return "Send action request to admin for approval"
}

// Name implements ken.SlashCommand.
func (*Action) Name() string {
	return "action"
}

// Options implements ken.SlashCommand.
func (*Action) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "request",
			Description: "Request an action to be performed",
			Options: []*discordgo.ApplicationCommandOption{
				discord.StringCommandArg("action", "Action to be preformed", true),
			},
		},
	}
}

// Run implements ken.SlashCommand.
func (af *Action) Run(ctx ken.Context) (err error) {
	return ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "request", Run: af.request},
	)
}

// Version implements ken.SlashCommand.
func (*Action) Version() string {
	return "1.0.0"
}

func (af *Action) request(ctx ken.SubCommandContext) (err error) {
	if err = ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}
	event := ctx.GetEvent()

	reqArg := ctx.Options().GetByName("action").StringValue()
	// East coast time babyyy
	humanReqTime := util.GetEstTimeStamp()

	// get the user's name in guild
	guildMember, err := ctx.GetSession().
		GuildMember(event.GuildID, event.Member.User.ID)
	if err != nil {
		log.Println(err)
		return discord.ErrorMessage(ctx, "Error getting guild member",
			"There was an error getting the guild member. Let Alex know he's a bad programmer.")
	}

	if guildMember.Nick == "" {
		guildMember.Nick = guildMember.User.Username
	}

	actionLog := fmt.Sprintf(
		"%s - %s - %s",
		guildMember.Nick,
		reqArg,
		humanReqTime,
	)

	// maybe will do something else with this but code block gives nice formatting
	// similar to that of what a logger would be...
	actionLog = discord.Code(actionLog)
	_, err = ctx.GetSession().ChannelMessageSend(funnelChannelID, actionLog)
	if err != nil {
		log.Println(err)
		return discord.ErrorMessage(
			ctx,
			"Error sending action request",
			"Alex is a bad programmer.",
		)
	}

	return discord.SuccessfulMessage(ctx, "Action Requested", fmt.Sprintf("Request '%s' sent for processing", reqArg))
}
