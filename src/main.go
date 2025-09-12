package main

import (
	"fmt"
)

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

// TÂCHE 1 : structure Character
type Character struct {
	Name      string
	Class     Classe
	Level     int
	HPMax     int
	HP        int
	Inventory map[string]int // inventaire simple: nom d'objet -> quantité
}

// TÂCHE 2 : initCharacter
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	return Character{
		Name:      name,
		Class:     class,
		Level:     level,
		HPMax:     hpMax,
		HP:        hp,
		Inventory: inv,
	}
}

// TÂCHE 3 : displayInfo
func displayInfo(c Character) {
	fmt.Println("=== Informations du personnage ===")
	fmt.Printf("Nom   : %s\n", c.Name)
	fmt.Printf("Classe: %s\n", c.Class)
	fmt.Printf("Niveau: %d\n", c.Level)
	fmt.Printf("PV    : %d / %d\n", c.HP, c.HPMax)
	fmt.Println("Inventaire :")
	if len(c.Inventory) == 0 {
		fmt.Println("  (vide)")
	} else {
		for item, qty := range c.Inventory {
			fmt.Printf("  - %s x%d\n", item, qty)
		}
	}
}

func main() {
	// TÂCHE 2 : initialisation demandée par l’énoncé
	// Remplace "Kimo" par ton prénom si tu veux.
	c1 := initCharacter(
		"Yazuo",
		ClasseSamurai,                      // Classe: Samurai
		1,                                  // Niveau: 1
		100,                                // PV max: 100
		40,                                 // PV actuels: 40
		map[string]int{"Potion de vie": 3}, // Inventaire: 3 potions
	)

	// TÂCHE 3 : affichage
	displayInfo(c1)
}
