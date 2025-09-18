package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/eiannone/keyboard"
)

/* ====== Marchand (économie + upgrades) ====== */

func canAfford(c *Character, price int) bool { return c.Gold >= price }

func buyItem(c *Character, item string, price, qty int) bool {
	if price < 0 || qty <= 0 {
		return false
	}
	// vérifier la place avant de payer
	if !canCarry(*c, qty) {
		fmt.Printf(CRed+"Inventaire plein (%d/%d). Libérez de la place avant d’acheter %s."+CReset+"\n",
			totalItems(*c), c.CapMax, item)
		return false
	}
	if !canAfford(c, price) {
		fmt.Printf(CRed+"Or insuffisant pour %s (coût %d). Solde: %d"+CReset+"\n", item, price, c.Gold)
		return false
	}
	c.Gold -= price
	if !addInventory(c, item, qty) {
		// sécurité
		c.Gold += price
		return false
	}
	fmt.Printf(CGreen+"Achat effectué ! %s x%d (−%d or). Or restant: %d"+CReset+"\n", item, qty, price, c.Gold)
	return true
}

var redbullFreeOnce = true
var windBookStock = 1 // stock unique du livre

func merchantMenu(c *Character) {
	bufio.NewReader(os.Stdin)
	options := []string{
		"RedBull — GRATUIT",
		"Potion de vie — 3 or",
		"Potion de poison — 6 or",
		"Livre de Sort : Mur de vent — 25 or",
		"Fourrure de Loup — 4 or",
		"Peau de Troll — 7 or",
		"Cuir de Sanglier — 3 or",
		"Plume de Corbeau — 1 or",
		"Augmentation d’inventaire — 30 or",
		"Retour",
	}

	selected := 0

	if err := keyboard.Open(); err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	for {
		// Clear screen
		fmt.Print("\033[H\033[2J")
		fmt.Println(CYellow + "=== MARCHAND ===" + CReset)
		for i, option := range options {
			prefix := "  "
			if i == selected {
				prefix = "> "
			}
			fmt.Println(prefix + option)
		}

		_, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch key {
		case keyboard.KeyArrowUp:
			if selected > 0 {
				selected--
			}
		case keyboard.KeyArrowDown:
			if selected < len(options)-1 {
				selected++
			}
		case keyboard.KeyEnter:
			switch selected {
			case 0: // RedBull gratuite 1 fois
				if redbullFreeOnce {
					if addInventory(c, "RedBull", 1) {
						redbullFreeOnce = false
						fmt.Printf(CGreen+"RedBull reçue ! (total: %d)"+CReset+"\n", c.Inventory["RedBull"])
					}
				} else {
					fmt.Println(CRed + "La RedBull gratuite n’est plus disponible." + CReset)
				}
			case 1:
				_ = buyItem(c, "Potion de vie", 3, 1)
			case 2:
				_ = buyItem(c, "Potion de poison", 6, 1)
			case 3:
				if windBookStock <= 0 {
					fmt.Println(CRed + "Le Livre de Sort : Mur de vent n’est plus disponible." + CReset)
					break
				}
				if buyItem(c, "Livre de Sort : Mur de vent", 25, 1) {
					windBookStock--
				}
			case 4:
				_ = buyItem(c, "Fourrure de Loup", 4, 1)
			case 5:
				_ = buyItem(c, "Peau de Troll", 7, 1)
			case 6:
				_ = buyItem(c, "Cuir de Sanglier", 3, 1)
			case 7:
				_ = buyItem(c, "Plume de Corbeau", 1, 1)
			case 8: // Augmentation d’inventaire
				if c.InvUpgrades >= MaxInventoryUpgrades {
					fmt.Printf(CYellow+"Limite d’améliorations atteinte (%d/%d)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
					break
				}
				if !canAfford(c, 30) {
					fmt.Printf(CRed+"Or insuffisant pour l’amélioration (coût 30). Solde: %d"+CReset+"\n", c.Gold)
					break
				}
				c.Gold -= 30
				if !upgradeInventorySlot(c) {
					c.Gold += 30
				}
			case 9: // Retour
				return
			}
			fmt.Println("\nAppuyez sur une touche pour continuer...")
			keyboard.GetSingleKey()
		}
	}
}
