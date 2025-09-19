package main

import (
	"bufio"
	"fmt"
)

// Menu principal
func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Informations du personnage")
	fmt.Println("2) Inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("5) Entraînement")
	fmt.Println("6) Missions")
	fmt.Println("7) Équipement")
	if allMissionsCompleted(c) {
		fmt.Println("8) Boss final")
	} else {
		fmt.Println("8) Boss final (verrouillé)")
	}
	fmt.Println("9) Quitter")

	switch readChoice(r) {
	case "1":
		displayInfo(*c)
	case "2":
		inventoryMenu(c, r)
	case "3":
		merchantMenu(c)
	case "4":
		blacksmithMenu(c, r)
	case "5":
		if StartTrainingFight() {
			// Victoire → récompense XP (Mission 2) et progression Mission 1
			addXP(c, 20)
			for i := range c.Missions {
				if c.Missions[i].ID == 1 {
					c.Missions[i].TrainKills++
					break
				}
			}
		}
	case "6":
		missionsMenu(c, r)
	case "7":
		equipmentMenu(c, r)
	case "8":
		if allMissionsCompleted(c) {
			startBossFinal(c, r)
		} else {
			fmt.Println(CRed + "Terminez d'abord toutes les missions !" + CReset)
		}
	case "9", "q", "quit", "quitter", "exit":
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}

	return true
}
