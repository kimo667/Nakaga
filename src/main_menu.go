package main

import (
	"bufio"
	"fmt"
)

// mainMenu affiche le menu principal et route vers les sous-menus.
// Retourne false si l'utilisateur choisit de quitter (pour casser la boucle dans main.go).
func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("5) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)

	case "2", "inventaire":
		inventoryMenu(c, r)

	case "3", "marchand", "shop":
		merchantMenu(c, r)

	case "4", "forgeron", "forge":
		blacksmithMenu(c, r)

	case "5", "q", "quit", "quitter", "exit":
		fmt.Println("Au revoir !")
		return false

	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}
	return true
}
