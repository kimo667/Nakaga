package main

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

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

	ch.BaseHPMax = hpMax
	// Copie sécurisée de l’inventaire (respect de l
	// a capacité)
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

// Création interactive d’un personnage
func createCharacterInteractive(r *bufio.Reader) Character {
	// 1) Nom du perso (prompt unique via readLine)
	rawName := readLine(r, "Entrez le nom de votre personnage : ")
	name := strings.TrimSpace(rawName)
	if name == "" {
		name = "Yazuo" // nom par défaut si l'utilisateur valide sans rien
	} else {
		// Majuscule sur la première lettre, le reste inchangé
		runes := []rune(name)
		runes[0] = unicode.ToUpper(runes[0])
		name = string(runes)
	}

	// 2) Choix de classe — texte “comme avant”
	fmt.Println("Choisissez une classe :")
	fmt.Println("1) Humain ")
	fmt.Println("2) Samurai ")
	fmt.Println("3) Ninja   ")
	classChoice := readChoice(r)

	var class Classe
	switch classChoice {
	case "1", "humain":
		class = ClasseHumain
	case "2", "samurai":
		class = ClasseSamurai
	case "3", "ninja":
		class = ClasseNinja
	default:
		fmt.Println(CRed + "Choix invalide. Classe par défaut : Humain." + CReset)
		class = ClasseHumain
	}

	// 3) PV init selon classe (tu peux remettre tes valeurs d’avant)
	hpMax := 110
	switch class {
	case ClasseSamurai:
		hpMax = 175
	case ClasseNinja:
		hpMax = 90
	}

	// 4) Inventaire de départ (à ta convenance)
	inv := map[string]int{
		"Potion de vie": 1,
	}

	// 5) Création du perso
	player := initCharacter(name, class, 1, hpMax, hpMax, inv)

	// (si tu veux un petit message fin)
	fmt.Println(CGreen + "Personnage créé !" + CReset)

	return player
}
