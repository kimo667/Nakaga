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
		"Fourrure d'Okami": 2,
		"Peau d'Oni":       1,
	},
	"Bottes du Shinobi": {
		"Cuir d'Inoshishi": 1,
		"Plume de Karasu":  1,
	},
	// Set aventurier (équipable)
	"Chapeau de l'aventurier": {
		"Plume de Karasu":  1,
		"Cuir d'Inoshishi": 1,
	},
	"Tunique de l'aventurier": {
		"Fourrure d'Okami": 1,
		"Peau d'Oni":       1,
	},
	"Bottes de l'aventurier": {
		"Fourrure d'Okami": 1,
		"Plume de Karasu":  1,
	},
}

// Vérifie si toutes les ressources sont présentes
func hasResources(c Character, recipe map[string]int) bool {
	for it, q := range recipe {
		if c.Inventory[it] < q {
			return false
		}
	}
	return true
}

// Consomme les ressources
func consumeResources(c *Character, recipe map[string]int) {
	for it, q := range recipe {
		removeInventory(c, it, q)
	}
}

func blacksmithMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println(CYellow + "\n=== FORGERON ===" + CReset)
		fmt.Println("Chaque craft coûte 5 or en plus des matériaux.")
		fmt.Println("Ressources :")
		for _, res := range []string{"Fourrure d'Okami", "Peau d'Oni", "Cuir d'Inoshishi", "Plume de Karasu"} {
			fmt.Printf("  - %-20s x%d\n", res, c.Inventory[res])
		}
		fmt.Println("\nRecettes disponibles :")
		i := 1
		items := []string{}
		for item := range forgeRecipes {
			fmt.Printf("  %d) %s\n", i, item)
			items = append(items, item)
			i++
		}
		fmt.Println("9) Retour")

		ch := readChoice(r)
		if ch == "9" {
			return
		}
		var idx int
		fmt.Sscanf(ch, "%d", &idx)
		if idx >= 1 && idx <= len(items) {
			item := items[idx-1]
			recipe := forgeRecipes[item]
			if !hasResources(*c, recipe) {
				fmt.Println(CRed + "Matériaux insuffisants." + CReset)
				continue
			}
			if c.Gold < 5 {
				fmt.Println(CRed + "Or insuffisant." + CReset)
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
