package main

import (
	"fmt"
	"strings"
)

// --- Mapping item -> slot ---
var equipSlotByItem = map[string]string{
	"Chapeau de l'aventurier": "Head",
	"Tunique de l'aventurier": "Torso",
	"Bottes de l'aventurier":  "Feet",
}

// Normalise les noms (’ -> ')
func normalizeItemName(s string) string {
	return strings.ReplaceAll(s, "’", "'")
}

// Retourne le slot pour un item (si équipable)
func slotForItem(item string) (string, bool) {
	s, ok := equipSlotByItem[normalizeItemName(item)]
	return s, ok
}

// Équiper un objet : retire du sac, place au bon slot, remet l'ancien au sac, puis recalc PV max
func equipItem(c *Character, item string) {
	slot, ok := slotForItem(item)
	if !ok {
		fmt.Println(CRed + "Cet objet n’est pas équipable." + CReset)
		return
	}

	itemNorm := normalizeItemName(item)

	// Retirer 1 du sac : on tente le libellé brut puis normalisé
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

	// Placer au slot et renvoyer l'ancien au sac s'il existe
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

	// Met à jour les PV max selon les équipements portés
	recalcHPMax(c)
	fmt.Printf(CGreen+"%s équipé. PV max = %d"+CReset+"\n", item, c.HPMax)
}

// Déséquiper un slot : renvoie l'objet au sac et recalc PV max
func unequipSlot(c *Character, slot string) {
	switch slot {
	case "Head":
		if c.Equipment.Head != "" {
			addInventory(c, c.Equipment.Head, 1)
			c.Equipment.Head = ""
		}
	case "Torso":
		if c.Equipment.Torso != "" {
			addInventory(c, c.Equipment.Torso, 1)
			c.Equipment.Torso = ""
		}
	case "Feet":
		if c.Equipment.Feet != "" {
			addInventory(c, c.Equipment.Feet, 1)
			c.Equipment.Feet = ""
		}
	default:
		fmt.Println(CRed + "Slot inconnu." + CReset)
		return
	}

	recalcHPMax(c)
	fmt.Printf(CYellow+"Objet retiré. PV max = %d"+CReset+"\n", c.HPMax)
}
