package main

import (
	"bufio"
	"fmt"
)

func inventoryMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n===== INVENTAIRE =====")
		displayInventory(*c)

		type opt struct {
			key, label string
			run        func()
		}
		opts := []opt{}
		add := func(label string, fn func()) {
			opts = append(opts, opt{
				key:   fmt.Sprintf("%d", len(opts)+1),
				label: label,
				run:   fn,
			})
		}

		// Consommables / livres : n’afficher que si possédés
		if c.Inventory["RedBull"] > 0 {
			add("Boire une RedBull (+20 PV)", func() {
				takeRedBull(c)
				isDead(c)
			})
		}
		if c.Inventory["Potion de vie"] > 0 {
			add("Boire une Potion de vie (+20 PV)", func() {
				usePotionVie(c)
				isDead(c)
			})
		}
		if c.Inventory["Potion de mana"] > 0 {
			add("Boire une Potion de mana (+15 Mana)", func() {
				useManaPotion(c)
			})
		}
		if c.Inventory["Potion de poison"] > 0 {
			add("Utiliser une Potion de poison (10 dmg/s ×3)", func() {
				poisonPot(c)
			})
		}
		if c.Inventory["Livre de Sort : Mur de vent"] > 0 {
			add("Utiliser 'Livre de Sort : Mur de vent'", func() {
				useSpellBookWind(c)
			})
		}

		// ÉQUIPEMENT : “Équiper …” pour chaque objet équipable possédé
		for item, qty := range c.Inventory {
			if qty <= 0 {
				continue
			}
			if _, ok := slotForItem(item); ok {
				it := item
				add("Équiper "+it, func() { equipItem(c, it) })
			}
		}

		// Déséquiper si quelque chose est porté
		if c.Equipment.Head != "" {
			add("Déséquiper (Tête)", func() { unequipSlot(c, "Head") })
		}
		if c.Equipment.Torso != "" {
			add("Déséquiper (Torse)", func() { unequipSlot(c, "Torso") })
		}
		if c.Equipment.Feet != "" {
			add("Déséquiper (Pieds)", func() { unequipSlot(c, "Feet") })
		}

		if len(opts) == 0 {
			fmt.Println("(Aucune action disponible)")
		} else {
			for _, o := range opts {
				fmt.Printf("%s) %s\n", o.key, o.label)
			}
		}
		fmt.Println("9) Retour")

		ch := readChoice(r)
		if ch == "9" || ch == "retour" || ch == "back" {
			return
		}
		ok := false
		for _, o := range opts {
			if ch == o.key {
				o.run()
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}
