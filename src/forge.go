package main

import (
	"bufio"
	"fmt"
)

// Liste des recettes : item → ressources nécessaires
var forgeRecipes = map[string]map[string]int{
	"Capuche du Shinobi": {
		"Plume de Karasu":  1,
		"Cuir d'Inoshishi": 1,
	},
	"Veste du Shinobi": {
		"Fourrure d'Okam": 2,
		"Peau d'Oni":      1,
	},
	"Tabi du Shinobi": {
		"Fourrure d'Okam":  1,
		"Cuir d'Inoshishi": 1,
	},
}

// Vérifie si le joueur a les ressources nécessaires
func hasResources(c *Character, recipe map[string]int) bool {
	for item, qty := range recipe {
		if c.Inventory[item] < qty {
			return false
		}
	}
	return true
}

// Consomme les ressources du joueur
func consumeResources(c *Character, recipe map[string]int) {
	for item, qty := range recipe {
		removeInventory(c, item, qty)
	}
}

// Menu du forgeron
func blacksmithMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n=== FORGERON ===")
		fmt.Printf("Votre or : %d\n", c.Gold)
		fmt.Println("Que voulez-vous fabriquer ? (5 or par objet)")

		idx := 1
		itemList := []string{}
		for item := range forgeRecipes {
			fmt.Printf("%d) %s\n", idx, item)
			itemList = append(itemList, item)
			idx++
		}
		fmt.Println("9) Retour")

		choice := readChoice(r)
		if choice == "9" || choice == "retour" || choice == "back" {
			return
		}

		var selected int
		fmt.Sscanf(choice, "%d", &selected)

		if selected >= 1 && selected <= len(itemList) {
			item := itemList[selected-1]
			recipe := forgeRecipes[item]

			if c.Gold < 5 {
				fmt.Println("Pas assez d’or pour fabriquer cet objet !")
				continue
			}
			if !hasResources(c, recipe) {
				fmt.Println("Ressources insuffisantes pour fabriquer :", item)
				continue
			}

			// Consommer or + ressources
			c.Gold -= 5
			consumeResources(c, recipe)

			// Ajouter l'objet crafté
			addInventory(c, item, 1)
			fmt.Printf("Vous avez fabriqué : %s (reste %d or).\n", item, c.Gold)
		} else {
			fmt.Println("Choix invalide.")
		}
	}
}
