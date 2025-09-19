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

// Bonus de PV max selon l’équipement porté
var equipHPBonus = map[string]int{
	"Chapeau de l'aventurier": 10,
	"Tunique de l'aventurier": 25,
	"Bottes de l'aventurier":  15,
}

// Normalise les noms (’ -> ')
func normalizeItemName(s string) string {
	s = strings.ReplaceAll(s, "’", "'")
	return s
}

// Renvoie le slot pour un item si équipable
func slotForItem(item string) (string, bool) {
	item = normalizeItemName(item)
	slot, ok := equipSlotByItem[item]
	return slot, ok
}

// Recalcule les PV max
func recalcHPMax(c *Character) {
	if c.BaseHPMax == 0 {
		c.BaseHPMax = c.HPMax
	}
	// Bonus
	bonus := 0
	if c.Equipment.Head != "" {
		bonus += equipHPBonus[normalizeItemName(c.Equipment.Head)]
	}
	if c.Equipment.Torso != "" {
		bonus += equipHPBonus[normalizeItemName(c.Equipment.Torso)]
	}
	if c.Equipment.Feet != "" {
		bonus += equipHPBonus[normalizeItemName(c.Equipment.Feet)]
	}
	c.HPMax = c.BaseHPMax + bonus
	if c.HP > c.HPMax {
		c.HP = c.HPMax
	}
}

// Équipe un objet depuis l'inventaire
func equipItem(c *Character, item string) {
	slot, ok := slotForItem(item)
	if !ok {
		fmt.Println(CRed + "Cet objet ne peut pas être équipé." + CReset)
		return
	}
	if c.Inventory[item] <= 0 {
		fmt.Println(CRed + "Vous ne possédez pas cet objet." + CReset)
		return
	}

	// retirer 1 de l'inventaire
	removeInventory(c, item, 1)

	// si un objet est déjà porté sur ce slot, le remettre dans l'inventaire
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
	fmt.Printf(CGreen+"%s équipé sur %s. PV max = %d"+CReset+"\n", item, slot, c.HPMax)
}

// Déséquipe un slot
func unequipSlot(c *Character, slot string) {
	switch slot {
	case "Head":
		if c.Equipment.Head == "" {
			fmt.Println(CYellow + "Rien sur la tête." + CReset)
			return
		}
		addInventory(c, c.Equipment.Head, 1)
		c.Equipment.Head = ""
	case "Torso":
		if c.Equipment.Torso == "" {
			fmt.Println(CYellow + "Rien sur le torse." + CReset)
			return
		}
		addInventory(c, c.Equipment.Torso, 1)
		c.Equipment.Torso = ""
	case "Feet":
		if c.Equipment.Feet == "" {
			fmt.Println(CYellow + "Rien aux pieds." + CReset)
			return
		}
		addInventory(c, c.Equipment.Feet, 1)
		c.Equipment.Feet = ""
	default:
		fmt.Println(CRed + "Slot inconnu." + CReset)
		return
	}

	recalcHPMax(c)
	fmt.Printf(CYellow+"Objet retiré. PV max = %d"+CReset+"\n", c.HPMax)
}
