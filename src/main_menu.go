package main

import (
	"bufio"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

type character struct {
	Name  string
	Level int
	HP    int
}

// Menu item
type menuItem string

func (i menuItem) Title() string       { return string(i) }
func (i menuItem) Description() string { return "" }
func (i menuItem) FilterValue() string { return string(i) }

// Bubble Tea model
type model struct {
	list      list.Model
	character Character
}

// Menu principal
func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("5) Entraînement (combat d'essai)")
	fmt.Println("6) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)

	case "2", "inventaire":
		inventoryMenu(c, r)

	case "3", "marchand", "shop":
		merchantMenu(c, r)

	case "4", "forgeron", "forge":
		blacksmithMenu(c, r)

	case "5", "entrainement", "entraînement", "training", "combat":
		// Lancement du mode entraînement
		StartTraining()

	case "6", "q", "quit", "quitter", "exit":
		fmt.Println("Au revoir !")
		return false

	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}

	return true

}
