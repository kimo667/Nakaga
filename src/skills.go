package main

import (
	"fmt"
	"sort"
)

func hasSkill(c Character, s string) bool {
	for _, k := range c.Skills {
		if k == s {
			return true
		}
	}
	return false
}

func learnSkill(c *Character, s string) bool {
	if hasSkill(*c, s) {
		return false
	}
	c.Skills = append(c.Skills, s)
	sort.Strings(c.Skills)
	return true
}

func useSpellBookWind(c *Character) {
	if c.Inventory["Livre de Sort : Mur de vent"] <= 0 {
		fmt.Println(CRed + "Vous n'avez pas de 'Livre de Sort : Mur de vent'." + CReset)
		return
	}
	if hasSkill(*c, "Mur de vent") {
		fmt.Println(CRed + "Vous connaissez déjà 'Mur de vent'. Le livre n'a pas été consommé." + CReset)
		return
	}
	removeInventory(c, "Livre de Sort : Mur de vent", 1)
	if learnSkill(c, "Mur de vent") {
		fmt.Println(CGreen + "Vous avez appris : Mur de vent !" + CReset)
	}
}
