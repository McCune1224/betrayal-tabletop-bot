package roll

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestRollAtRaritySimulated(t *testing.T) {
	rarityPriorities := []string{"common", "uncommon", "rare", "epic", "legendary", "mythical"}

	file := "outcomes.txt"
	// create a file
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

  iterations := 50000
	for l := 0; l < 100; l++ {
		outcomeTally := make(map[string]int)
		for _, rarity := range rarityPriorities {
			outcomeTally[rarity] = 0
		}
		for i := 0; i < iterations; i++ {
			rarity := rollAtRarity(float64(l), rarityPriorities)
			outcomeTally[rarity]++
		}
		// write to csv
		f.WriteString("Luck Level: " + fmt.Sprintf("%-8d ", (l)))
		for _, rarity := range rarityPriorities {
			f.WriteString(fmt.Sprintf("%s: %%%-6.2f, ", rarity, (float64(outcomeTally[rarity])/float64(iterations))*100))
		}
		f.WriteString("\n")
	}

	f.WriteString("------------------------------------------------------------------------------------------\n")

	for l := 0; l < 100; l++ {
		outcomeTally := make(map[string]int)
		for _, rarity := range rarityPriorities {
			outcomeTally[rarity] = 0
		}
		for i := 0; i < iterations; i++ {
			rarity := rollAtRarity(float64(l), rarityPriorities[3:])
			outcomeTally[rarity]++
		}
		// write to csv
		f.WriteString("Luck Level: " + fmt.Sprintf("%-8d ", (l)))
		for _, rarity := range rarityPriorities {
			f.WriteString(fmt.Sprintf("%s: %%%-6.2f, ", rarity, (float64(outcomeTally[rarity])/float64(iterations))*100))
		}
		f.WriteString("\n")
	}
}
