package main

import (
	"fmt"
	"time"
)

/* ====== Consommables & Effets ====== */

func takeRedBull(c *Character) {
	if !removeInventory(c, "RedBull", 1) {
		fmt.Println(CRed + "Pas de RedBull dans l'inventaire !" + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+HealRedBull, 0, c.HPMax)
	fmt.Printf(CGreen+"Glou glou ! PV: %d → %d"+CReset+"\n", before, c.HP)
}

func usePotionVie(c *Character) {
	if !removeInventory(c, "Potion de vie", 1) {
		fmt.Println(CRed + "Pas de Potion de vie !" + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+HealPotionVie, 0, c.HPMax)
	fmt.Printf(CGreen+"Gorgée de potion ! PV: %d → %d"+CReset+"\n", before, c.HP)
}

func useManaPotion(c *Character) { // Mission 4
	if !removeInventory(c, "Potion de mana", 1) {
		fmt.Println(CRed + "Pas de Potion de mana !" + CReset)
		return
	}
	before := c.Mana
	c.Mana = clamp(c.Mana+ManaPotRestore, 0, c.ManaMax)
	fmt.Printf(CGreen+"Surcharge d’éther ! Mana: %d → %d"+CReset+"\n", before, c.Mana)
}

// Potion de poison : applique 3 ticks de 10 dégâts sur ~1.5s
func poisonPot(c *Character) {
	if !removeInventory(c, "Potion de poison", 1) {
		fmt.Println(CRed + "Pas de Potion de poison !" + CReset)
		return
	}
	fmt.Println(CYellow + "Vous vous empoisonnez volontairement (pour tester)..." + CReset)
	for i := 0; i < PoisonTicks; i++ {
		time.Sleep(500 * time.Millisecond)
		c.HP -= PoisonDamagePerSec
		fmt.Printf(CRed+"[Poison] -%d PV (reste %d/%d)"+CReset+"\n", PoisonDamagePerSec, c.HP, c.HPMax)
		if isDead(c) {
			break
		}
	}
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
