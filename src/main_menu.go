package main

import (
	"bufio"
	"fmt"
)

// Menu principal
func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("5) Entraînement (combat d'essai)")
	fmt.Println("6) Accéder au donjon (Boss Final : Yone)")
	fmt.Println("7) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)

	case "2", "inventaire":
		inventoryMenu(c, r)

	case "3", "marchand", "shop":
		merchantMenu(c)

	case "4", "forgeron", "forge":
		blacksmithMenu(c, r)

	case "5", "entrainement", "entraînement", "training", "combat":
		StartTrainingFight()

	case "6", "donjon", "dungeon", "boss", "yone":
		StartAllBossFights()

	case "7", "q", "quit", "quitter", "exit":
		fmt.Println("Au revoir !")
		return false

	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}

	return true
}
