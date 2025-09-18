package main

// Bonus de PV max par équipement
var equipHPBonus = map[string]int{
	"Chapeau de l'aventurier": 10,
	"Tunique de l'aventurier": 25,
	"Bottes de l'aventurier":  15,
}

// Somme des bonus selon les pièces portées
func equipmentHPBonus(c Character) int {
	b := 0
	if v, ok := equipHPBonus[c.Equipment.Head]; ok {
		b += v
	}
	if v, ok := equipHPBonus[c.Equipment.Torso]; ok {
		b += v
	}
	if v, ok := equipHPBonus[c.Equipment.Feet]; ok {
		b += v
	}
	return b
}

// Recalcule HPMax (base + bonus) et “clamp” HP si besoin
func recalcHPMax(c *Character) {
	if c.BaseHPMax == 0 { // sécurité si vieux perso
		c.BaseHPMax = c.HPMax
	}
	c.HPMax = c.BaseHPMax + equipmentHPBonus(*c)
	if c.HP > c.HPMax {
		c.HP = c.HPMax
	}
}
