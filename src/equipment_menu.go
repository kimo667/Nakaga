package main

import (
	"bufio"
	"fmt"
)

func equipment_menu(c *Character, r *bufio.Reader) { // alias
	equipmentMenu(c, r)
}

func equipmentMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n===== ÉQUIPEMENT =====")
		// état actuel
		fmt.Println("Actuellement portés :")
		fmt.Printf("  Tête : %s\n", orNone(c.Equipment.Head))
		fmt.Printf("  Torse: %s\n", orNone(c.Equipment.Torso))
		fmt.Printf("  Pieds: %s\n", orNone(c.Equipment.Feet))

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

		// proposer d'équiper chaque item équipable possédé
		for item, qty := range c.Inventory {
			if qty <= 0 {
				continue
			}
			if _, ok := slotForItem(item); ok {
				it := item
				add("Équiper "+it, func() { equipItem(c, it) })
			}
		}
		// déséquiper
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

		choice := readChoice(r)
		if choice == "9" {
			return
		}
		ok := false
		for _, o := range opts {
			if choice == o.key {
				o.run()
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println("Choix invalide.")
		}
	}
}

func orNone(s string) string {
	if s == "" {
		return "(aucun)"
	}
	return s
}
