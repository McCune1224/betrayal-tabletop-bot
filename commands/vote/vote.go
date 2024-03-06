package vote

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mccune1224/betrayal-tabletop-bot/data"
	"github.com/mccune1224/betrayal-tabletop-bot/discord"
	"github.com/zekrotja/ken"
)

const voteChannelID = "1211657759238524939"

type Vote struct {
	models *data.Models
}

// Initialize implements main.BetrayalCommand.
func (v *Vote) Initialize(models *data.Models) {
	v.models = models
}

var _ ken.SlashCommand = (*Vote)(nil)

// Description implements ken.SlashCommand.
func (*Vote) Description() string {
	return "Vote a player"
}

// Name implements ken.SlashCommand.
func (*Vote) Name() string {
	return "vote"
}

// Options implements ken.SlashCommand.
func (*Vote) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "batch",
			Description: "Batch vote players up to 5 players",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user2",
					Description: "User to vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user3",
					Description: "User to vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user4",
					Description: "User to vote",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user5",
					Description: "User to vote",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "player",
			Description: "Vote a single player",
			Options: []*discordgo.ApplicationCommandOption{
				discord.UserCommandArg(true),
				discord.StringCommandArg("context", "Additional Context/Details to provide (i.e using Gold Card)", false),
			},
		},
	}
}

// Run implements ken.SlashCommand.
func (v *Vote) Run(ctx ken.Context) (err error) {
	return ctx.HandleSubCommands(
		ken.SubCommandHandler{Name: "batch", Run: v.batch},
		ken.SubCommandHandler{Name: "player", Run: v.player},
	)
}

func (v *Vote) batch(ctx ken.SubCommandContext) (err error) {
	if err := ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}

	users := []*discordgo.User{ctx.Options().GetByName("user").UserValue(ctx)}

	for i := 2; i <= 5; i++ {
		user, ok := ctx.Options().GetByNameOptional(fmt.Sprintf("user%d", i))
		if ok {
			users = append(users, user.UserValue(ctx))
		}
	}

	voteMsg := fmt.Sprintf("%s voted for", ctx.User().Username)
	for _, user := range users {
		voteMsg += fmt.Sprintf(" %s", user.Username)
	}

	sesh := ctx.GetSession()
	_, err = sesh.ChannelMessageSend(voteChannelID, discord.Code(voteMsg))
	if err != nil {
		return discord.AlexError(ctx, "Failed to send vote message")
	}

	votedFor := ""
	for i, user := range users {
		if i == len(users)-1 {
			votedFor += discord.MentionUser(user.ID)
			continue
		}
		votedFor += fmt.Sprintf("%s, ", discord.MentionUser(user.ID))
	}
	return discord.SuccessfulMessage(ctx, "Vote Sent", fmt.Sprintf("Voted for %s", votedFor))
}

func (v *Vote) player(ctx ken.SubCommandContext) (err error) {
	if err := ctx.Defer(); err != nil {
		log.Println(err)
		return err
	}
	voteUser := ctx.Options().GetByName("user").UserValue(ctx)
	voteContext, ok := ctx.Options().GetByNameOptional("context")

	voteMsg := ""

	if ok {
		voteContext := voteContext.StringValue()
		voteMsg = fmt.Sprintf("%s voted for %s with context: %s", ctx.User().Username, voteUser.Username, voteContext)
	} else {
		voteMsg = fmt.Sprintf("%s voted for %s", ctx.User().Username, voteUser.Username)
	}

	sesh := ctx.GetSession()
	_, err = sesh.ChannelMessageSend(voteChannelID, discord.Code(voteMsg))
	if err != nil {
		return discord.AlexError(ctx, "Failed to send vote message")
	}

	return discord.SuccessfulMessage(ctx, "Vote Sent for Processing.", fmt.Sprintf("Voted for %s", discord.MentionUser(voteUser.ID)))
}

// Version implements ken.SlashCommand.
func (*Vote) Version() string {
	return "1.0.0"
}
