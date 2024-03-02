package util

import (
	"strings"

	"github.com/agnivade/levenshtein"
)

// bestRole := roles[0]
// found := false
// for _, role := range roles {
// 	// check if its a substring match first, if so, use that role, then check for levenshtein distance
// 	if strings.Contains(strings.ToLower(role.Name), strings.ToLower(name)) {
// 		bestRole = role
// 		found = true
// 		break
// 	}
// }
//
// if !found {
// 	for _, role := range roles {
// 		distance := levenshtein.ComputeDistance(strings.ToLower(name), strings.ToLower(role.Name))
// 		if distance == 0 {
// 			bestRole = role
// 			break
// 		}
// 		if distance < levenshtein.ComputeDistance(strings.ToLower(name), strings.ToLower(bestRole.Name)) {
// 			bestRole = role
// 		}
// 	}
// }
//
// convert the above code to a function

func FuzzyFind(name string, options []string) string {
	best := ""
	found := false
	for _, option := range options {
		// check if its a substring match first, if so, use that role, then check for levenshtein distance
		if strings.Contains(strings.ToLower(option), strings.ToLower(name)) {
			best = option
			found = true
			break
		}
	}
	if !found {
		for _, option := range options {
			distance := levenshtein.ComputeDistance(strings.ToLower(name), strings.ToLower(option))
			if distance == 0 {
				best = option
				break
			}
			if distance < levenshtein.ComputeDistance(strings.ToLower(name), strings.ToLower(best)) {
				best = option
			}
		}
	}
	return best
}
