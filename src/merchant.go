package main

import (
	"bufio"
	"fmt"
)

// Économie
func canAfford(c Character, price int) bool { return c.Gold >= price }

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
	if !canAfford(*c, price) {
		fmt.Printf(CRed+"Or insuffisant pour %s (coût %d). Solde: %d"+CReset+"\n",
			item, price, c.Gold)
		return false
	}
	c.Gold -= price
	if !addInventory(c, item, qty) {
		// sécurité
		c.Gold += price
		return false
	}
	fmt.Printf(CGreen+"Achat effectué ! %s x%d (−%d or). Or restant: %d"+CReset+"\n",
		item, qty, price, c.Gold)
	return true
}

var redbullFreeOnce = true
var windBookStock = 1 // stock unique du livre

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n" + CYellow + "=== MARCHAND ===" + CReset)
		if redbullFreeOnce {
			fmt.Println("1) RedBull — " + CGreen + "GRATUIT" + CReset)
		} else {
			fmt.Println("1) RedBull — (ÉPUISÉ)")
		}
		// Tarifs
		fmt.Println("2) Potion de vie — 3 or")
		fmt.Println("3) Potion de poison — 6 or")
		if windBookStock > 0 {
			fmt.Println("4) Livre de Sort : Mur de vent — 25 or")
		} else {
			fmt.Println("4) Livre de Sort : Mur de vent — (ÉPUISÉ)")
		}
		// Matériaux
		fmt.Println("5) Fourrure de Loup — 4 or")
		fmt.Println("6) Peau de Troll — 7 or")
		fmt.Println("7) Cuir de Sanglier — 3 or")
		fmt.Println("8) Plume de Corbeau — 1 or")
		// Upgrade capacité
		fmt.Printf("10) Augmentation d’inventaire — 30 or (utilisée %d/%d)\n", c.InvUpgrades, MaxInventoryUpgrades)

		fmt.Println("9) Retour")

		switch readChoice(r) {
		case "1": // RedBull gratuite 1 fois
			if redbullFreeOnce {
				if addInventory(c, "RedBull", 1) {
					redbullFreeOnce = false
					fmt.Printf(CGreen+"RedBull reçue ! (total: %d)"+CReset+"\n", c.Inventory["RedBull"])
				}
			} else {
				fmt.Println(CRed + "La RedBull gratuite n’est plus disponible." + CReset)
			}

		case "2":
			_ = buyItem(c, "Potion de vie", 3, 1)
		case "3":
			_ = buyItem(c, "Potion de poison", 6, 1)
		case "4":
			if windBookStock <= 0 {
				fmt.Println(CRed + "Le Livre de Sort : Mur de vent n’est plus disponible." + CReset)
				break
			}
			if buyItem(c, "Livre de Sort : Mur de vent", 25, 1) {
				windBookStock--
			}
		case "5":
			_ = buyItem(c, "Fourrure de Loup", 4, 1)
		case "6":
			_ = buyItem(c, "Peau de Troll", 7, 1)
		case "7":
			_ = buyItem(c, "Cuir de Sanglier", 3, 1)
		case "8":
			_ = buyItem(c, "Plume de Corbeau", 1, 1)
		case "10":
			if c.InvUpgrades >= MaxInventoryUpgrades {
				fmt.Printf(CYellow+"Limite d’améliorations atteinte (%d/%d)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
				break
			}
			if !canAfford(*c, 30) {
				fmt.Printf(CRed+"Or insuffisant pour l’amélioration (coût 30). Solde: %d"+CReset+"\n", c.Gold)
				break
			}
			c.Gold -= 30
			if !upgradeInventorySlot(c) {
				c.Gold += 30
			}
		case "9", "retour", "back":
			return
		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}
