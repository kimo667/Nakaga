package main

import (
	"fmt"
)

/* ====== Inventaire : capacité & upgrades ====== */

func totalItems(c Character) int {
	sum := 0
	for _, q := range c.Inventory {
		if q > 0 {
			sum += q
		}
	}
	return sum
}

func canCarry(c Character, qty int) bool {
	return totalItems(c)+qty <= c.CapMax
}

func addInventory(c *Character, item string, qty int) bool {
	if qty <= 0 {
		return false
	}
	if c.Inventory == nil {
		c.Inventory = make(map[string]int)
	}
	if !canCarry(*c, qty) {
		fmt.Printf(CRed+"Inventaire plein (%d/%d). Impossible d'ajouter %d x %s."+CReset+"\n",
			totalItems(*c), c.CapMax, qty, item)
		return false
	}
	c.Inventory[item] += qty
	return true
}

func removeInventory(c *Character, item string, qty int) bool {
	cur, ok := c.Inventory[item]
	if !ok || qty <= 0 || cur < qty {
		return false
	}
	if cur-qty == 0 {
		delete(c.Inventory, item)
	} else {
		c.Inventory[item] = cur - qty
	}
	return true
}

func upgradeInventorySlot(c *Character) bool {
	if c.InvUpgrades >= MaxInventoryUpgrades {
		fmt.Printf(CYellow+"Capacité déjà au maximum (%d/%d upgrades)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
		return false
	}
	c.InvUpgrades++
	c.CapMax += InventoryUpgradeStep
	fmt.Printf(CGreen+"Capacité augmentée ! Nouvelle capacité : %d (améliorations : %d/%d)"+CReset+"\n",
		c.CapMax, c.InvUpgrades, MaxInventoryUpgrades)
	return true
}

/* ====== Affichage de l’inventaire ====== */

//func displayInventory(c Character) {
//keys := make([]string, 0, len(c.Inventory))
//for k, v := range c.Inventory {
//if v > 0 {
//	keys = append(keys, k)
//}
//}
//sort.Strings(keys)
//fmt.Printf(CYellow+"Inventaire (%d/%d) :"+CReset+"\n", totalItems(c), c.CapMax)
//if len(keys) == 0 {
//fmt.Println("  (vide)")
//	return
//}
//	for _, k := range keys {
//		fmt.Printf("  - %s x%d\n", k, c.Inventory[k])
//	}
//}
