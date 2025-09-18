package main

import (
	"fmt"
	"strings"
)

// Objets équipables → slot
var equipSlotByItem = map[string]string{
	"Chapeau de l'aventurier": "Head",
	"Tunique de l'aventurier": "Torso",
	"Bottes de l'aventurier":  "Feet",
}

func normalizeItemName(s string) string {
	return strings.ReplaceAll(s, "’", "'")
}

func slotForItem(item string) (string, bool) {
	s, ok := equipSlotByItem[normalizeItemName(item)]
	return s, ok
}

// Équiper : enlève 1 du sac, place au slot, rend l'ancien au sac, recalc PV
func equipItem(c *Character, item string) {
	slot, ok := slotForItem(item)
	if !ok {
		fmt.Println(CRed + "Cet objet n’est pas équipable." + CReset)
		return
	}

	itemNorm := normalizeItemName(item)

	// retirer du sac (on essaie la version brute puis normalisée)
	removed := removeInventory(c, item, 1)
	if !removed && itemNorm != item {
		removed = removeInventory(c, itemNorm, 1)
		if removed {
			item = itemNorm
		}
	}
	if !removed {
		fmt.Println(CRed + "Vous ne possédez pas cet objet." + CReset)
		return
	}

	// équiper, rendre l’ancien au sac si besoin
	switch slot {
	case "Head":
		if c.Equipment.Head != "" {
			addInventory(c, c.Equipment.Head, 1)
		}
		c.Equipment.Head = item
	case "Torso":
		if c.Equipment.Torso != "" {
			addInventory(c, c.Equipment.Torso, 1)
		}
		c.Equipment.Torso = item
	case "Feet":
		if c.Equipment.Feet != "" {
			addInventory(c, c.Equipment.Feet, 1)
		}
		c.Equipment.Feet = item
	}

	recalcHPMax(c)
	fmt.Printf(CGreen+"%s équipé. PV max = %d"+CReset+"\n", item, c.HPMax)
}

// Déséquiper : remet l’objet au sac + recalc PV
func unequipSlot(c *Character, slot string) {
	switch slot {
	case "Head":
		if c.Equipment.Head == "" {
			fmt.Println("(rien à déséquiper en tête)")
			return
		}
		addInventory(c, c.Equipment.Head, 1)
		c.Equipment.Head = ""
	case "Torso":
		if c.Equipment.Torso == "" {
			fmt.Println("(rien à déséquiper au torse)")
			return
		}
		addInventory(c, c.Equipment.Torso, 1)
		c.Equipment.Torso = ""
	case "Feet":
		if c.Equipment.Feet == "" {
			fmt.Println("(rien à déséquiper aux pieds)")
			return
		}
		addInventory(c, c.Equipment.Feet, 1)
		c.Equipment.Feet = ""
	default:
		fmt.Println("Slot inconnu.")
		return
	}
	recalcHPMax(c)
}
