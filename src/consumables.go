package main

import (
	"fmt"
	"time"
)

func takeRedBull(c *Character) {
	if !removeInventory(c, "RedBull", 1) {
		fmt.Println(CRed + "Pas de RedBull dans l'inventaire !" + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+HealRedBull, 0, c.HPMax)
	fmt.Printf(CCyan+"Tu as bu une RedBull !"+CReset+" PV: %d → "+CGreen+"%d/%d"+CReset+"\n",
		before, c.HP, c.HPMax)
}

func usePotionVie(c *Character) {
	if !removeInventory(c, "Potion de vie", 1) {
		fmt.Println(CRed + "Pas de Potion de vie." + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+20, 0, c.HPMax)
	fmt.Printf(CCyan+"Vous buvez une Potion de vie."+CReset+" PV: %d → "+CGreen+"%d/%d"+CReset+"\n",
		before, c.HP, c.HPMax)
}

func poisonPot(c *Character) {
	if !removeInventory(c, "Potion de poison", 1) {
		fmt.Println(CRed + "Vous n'avez pas de Potion de poison." + CReset)
		return
	}
	fmt.Println(CCyan + "Vous utilisez une Potion de poison…" + CReset)
	for i := 1; i <= PoisonTicks; i++ {
		before := c.HP
		c.HP = clamp(c.HP-PoisonDamagePerSec, 0, c.HPMax)
		fmt.Printf("Effet poison %d/%d → PV: %d → %d/%d\n", i, PoisonTicks, before, c.HP, c.HPMax)
		if isDead(c) {
			fmt.Println(CRed + "L'effet du poison est interrompu suite à votre mort." + CReset)
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf(CCyan+"L'effet du poison est terminé. PV restants : %d/%d"+CReset+"\n", c.HP, c.HPMax)
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
