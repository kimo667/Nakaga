package main

import "fmt"

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ----- capacité -----

func totalItems(c Character) int {
	n := 0
	for _, q := range c.Inventory {
		if q > 0 {
			n += q
		}
	}
	return n
}

func canCarry(c Character, qty int) bool {
	return totalItems(c)+qty <= c.CapMax
}

func upgradeInventorySlot(c *Character) bool {
	if c.InvUpgrades >= MaxInventoryUpgrades {
		fmt.Printf(CYellow+"Capacité déjà au max (%d/%d)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
		return false
	}
	c.InvUpgrades++
	c.CapMax += InventoryUpgradeStep
	fmt.Printf(CGreen+"Capacité augmentée → %d (améliorations: %d/%d)"+CReset+"\n",
		c.CapMax, c.InvUpgrades, MaxInventoryUpgrades)
	return true
}

// ----- add / remove -----

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
