package main

import (
	"bufio"
	"fmt"
)

// Choix classe
func chooseClass(r *bufio.Reader) Classe {
	for {
		fmt.Println(CYellow + "Choisis ta classe :" + CReset)
		fmt.Println("1) Humain  – équilibré")
		fmt.Println("2) Samouraï – PV élevés")
		fmt.Println("3) Ninja    – agile")
		switch readChoice(r) {
		case "1", "humain":
			return ClasseHumain
		case "2", "samourai", "samouraï", "samurai":
			return ClasseSamurai
		case "3", "ninja":
			return ClasseNinja
		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

// Création du perso (init + capacité + skill de base)
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	ch := Character{
		Name:        name,
		Class:       class,
		Level:       level,
		HPMax:       hpMax,
		HP:          clamp(hp, 0, hpMax),
		Inventory:   map[string]int{},
		Skills:      []string{},
		Gold:        100,              // T11
		CapMax:      BaseInventoryCap, // T12+
		InvUpgrades: 0,
	}
	// inventaire de départ (respect capacité)
	for k, v := range inv {
		if v <= 0 {
			continue
		}
		if !addInventory(&ch, k, v) {
			break
		}
	}
	// technique de base
	learnSkill(&ch, "Tempête d'acier")
	return ch
}

// Création interactive (nom + classe)
func createCharacterInteractive(r *bufio.Reader) Character {
	fmt.Println(CYellow + "=== Création de personnage ===" + CReset)
	name := readLine(r, "Entre ton nom: ")
	if name == "" {
		name = "Yazuo"
	}
	class := chooseClass(r)

	// bonus par classe
	hpMax := 100
	switch class {
	case ClasseSamurai:
		hpMax = 110
	case ClasseNinja:
		hpMax = 90
	}
	startHP := hpMax * 40 / 100
	startInv := map[string]int{"RedBull": 3}

	ch := initCharacter(name, class, 1, hpMax, startHP, startInv)
	fmt.Println(CGreen + "Personnage créé !" + CReset)
	return ch
}
