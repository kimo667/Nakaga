package main

// Bonus de PV max selon l’équipement porté
var equipHPBonus = map[string]int{
	"Chapeau de l'aventurier": 10,
	"Tunique de l'aventurier": 25,
	"Bottes de l'aventurier":  15,
}

func equipmentHPBonus(c Character) int {
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
	return bonus
}

// HPMax = BaseHPMax (donnée par la classe) + bonus d’équipement
func recalcHPMax(c *Character) {
	// sécurité : si jamais BaseHPMax n'a pas été posé, on retombe sur l'actuel
	if c.BaseHPMax == 0 {
		c.BaseHPMax = c.HPMax
	}
	c.HPMax = c.BaseHPMax + equipmentHPBonus(*c)
	if c.HP > c.HPMax {
		c.HP = c.HPMax
	}
}
