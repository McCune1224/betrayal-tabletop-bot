package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

// Current roles with eleveted permissions.
var AdminRoles = []string{
	"Host",
	"Co-Host",
	"Bot Developer",
}

// Check if user who invoked command has required role
func IsAdminRole(ctx ken.Context, adminRoles ...string) bool {
	event := ctx.GetEvent()
	guildRoles, _ := ctx.GetSession().GuildRoles(event.GuildID)

	// Nothing screams "I'm a good programmer" more than a triple for loop.
	for _, rid := range event.Member.Roles {
		for _, r := range guildRoles {
			for _, ar := range adminRoles {
				if rid == r.ID && r.Name == ar {
					return true
				}
			}
		}
	}
	return false
}

func GetAdminRoleUsers(s *discordgo.Session, e *discordgo.InteractionCreate, adminRoles ...string) []string {
	guildRoles, _ := s.GuildRoles(e.GuildID)
	var users []string
	for _, rid := range e.Member.Roles {
		for _, r := range guildRoles {
			for _, ar := range adminRoles {
				if rid == r.ID && r.Name == ar {
					users = append(users, e.Member.User.ID)
				}
			}
		}
	}
	return users
}

// Get all the players within a guild that have the specified role
func GetMembersWithRole(s *discordgo.Session, e *discordgo.InteractionCreate, roleID string) []*discordgo.Member {
	members, _ := s.GuildMembers(e.GuildID, "", 1000)
	var roleMembers []*discordgo.Member
	for _, m := range members {
		for _, rid := range m.Roles {
			if rid == roleID {
				roleMembers = append(roleMembers, m)
			}
		}
	}
	return roleMembers
}

func GetMembersWithRoleName(s *discordgo.Session, e *discordgo.InteractionCreate, roleName string) []*discordgo.Member {
	guildRoles, _ := s.GuildRoles(e.GuildID)
	var roleID string
	for _, r := range guildRoles {
		if r.Name == roleName {
			roleID = r.ID
		}
	}
	return GetMembersWithRole(s, e, roleID)
}
