package main

// Bonus de PV max selon l’équipement porté
var equipHPBonus = map[string]int{
	"Chapeau de l'aventurier": 10,
	"Tunique de l'aventurier": 25,
	"Bottes de l'aventurier":  15,
}

// Recalcule HPMax à partir de la base + bonus d’équipement, et ajuste HP si nécessaire
func recalcHPMax(c *Character) {
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

	// HPMax = base nue + bonus d’équipement
	c.HPMax = c.BaseHPMax + bonus

	// Si les PV actuels dépassent le nouveau max, on les réduit
	if c.HP > c.HPMax {
		c.HP = c.HPMax
	}
}
