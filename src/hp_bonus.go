package main

// Bonus de PV max selon l’équipement porté
var equipHPBonus = map[string]int{
	"Chapeau de l'aventurier": 10,
	"Tunique de l'aventurier": 25,
	"Bottes de l'aventurier":  15,
}

// Calcule la somme des bonus d'équipement
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

// Recalcule HPMax = base (100) + bonus équipements
func recalcHPMax(c *Character) {
	if c.BaseHPMax == 0 {
		c.BaseHPMax = 100 // sécurité : base fixée à 100
	}
	c.HPMax = c.BaseHPMax + equipmentHPBonus(*c)

	// Clamp des PV actuels si nécessaire
	if c.HP > c.HPMax {
		c.HP = c.HPMax
	}
}
