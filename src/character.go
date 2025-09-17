package main

import (
	"bufio"
	"fmt"
)

// ====== Choix de classe ======

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

// ====== Initialisation du personnage ======
//
// - Copie l’inventaire de départ SANS dépasser la capacité
// - Donne 100 or au départ (T11)
// - Capacité de base = BaseInventoryCap (T12+)
// - Apprend la technique de base "Tempête d'acier"
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	ch := Character{
		Name:        name,
		Class:       class,
		Level:       level,
		HPMax:       hpMax,
		HP:          clamp(hp, 0, hpMax),
		Inventory:   map[string]int{},
		Skills:      []string{},
		Gold:        100,
		CapMax:      BaseInventoryCap,
		InvUpgrades: 0,
	}

	// Copie sécurisée de l’inventaire (respect de la capacité)
	for k, v := range inv {
		if v <= 0 {
			continue
		}
		if !addInventory(&ch, k, v) { // addInventory sera défini dans inventory.go
			break
		}
	}

	// Technique de base
	_ = learnSkill(&ch, "Tempête d'acier") // learnSkill sera défini dans skills.go
	return ch
}

// ====== Création interactive (nom + classe) ======

func createCharacterInteractive(r *bufio.Reader) Character {
	fmt.Println(CYellow + "=== Création de personnage ===" + CReset)

	name := readLine(r, "Entre ton nom: ")
	if name == "" {
		name = "Yazuo"
	}

	class := chooseClass(r)

	// Petits bonus simples par classe
	hpMax := 100
	switch class {
	case ClasseSamurai:
		hpMax = 110
	case ClasseNinja:
		hpMax = 90
	}

	startHP := hpMax * 40 / 100
	startInv := map[string]int{
		"RedBull": 3,
	}

	ch := initCharacter(name, class, 1, hpMax, startHP, startInv)
	fmt.Println(CGreen + "Personnage créé !" + CReset)
	return ch
}
