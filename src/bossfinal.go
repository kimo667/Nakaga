package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/eiannone/keyboard"
)

// // Monster struct (pour rappel)
// type Monster struct {
// 	Name        string
// 	MaxHP       int
// 	CurrentHP   int
// 	AttackPower int
// }

// // ====================
// // Fonctions d’attaque et utilitaires (déjà présentes dans ton code)
// // ====================
// func (m *Monster) Attack(target *Monster) int {
// 	damage := rand.Intn(m.AttackPower) + 1
// 	if rand.Intn(100) < 10 {
// 		damage *= 2
// 		slowPrint("Coup critique !", 50*time.Millisecond)
// 	}
// 	target.CurrentHP -= damage
// 	if target.CurrentHP < 0 {
// 		target.CurrentHP = 0
// 	}
// 	return damage
// }

// func (m *Monster) SpecialAttack(target *Monster) int {
// 	if rand.Intn(100) < 20 {
// 		slowPrint(fmt.Sprintf("%s a raté son attaque spéciale !", m.Name), 30*time.Millisecond)
// 		return 0
// 	}
// 	damage := rand.Intn(m.AttackPower*3) + m.AttackPower*2
// 	if rand.Intn(100) < 20 {
// 		damage *= 2
// 		slowPrint("Coup critique sur l'attaque spéciale !", 50*time.Millisecond)
// 	}
// 	target.CurrentHP -= damage
// 	if target.CurrentHP < 0 {
// 		target.CurrentHP = 0
// 	}
// 	return damage
// }

// // ====================
// // Slow print
// // ====================
// func slowPrint(s string, d time.Duration) {
// 	for _, c := range s {
// 		fmt.Printf("%c", c)
// 		time.Sleep(d)
// 	}
// 	fmt.Println()
// }

// // ====================
// // Affichage HP simple (sans Lipgloss pour simplifier ici)
// // ====================
// func displayHP(player, mob *Monster) {
// 	fmt.Printf("%s : %d/%d PV\n", mob.Name, mob.CurrentHP, mob.MaxHP)
// 	fmt.Printf("%s : %d/%d PV\n\n", player.Name, player.CurrentHP, player.MaxHP)
// }

// // ====================
// // Inventaire
// // ====================
// func showInventory(player *Monster, inventory map[string]int) {
// 	for {
// 		totalItems := 0
// 		for _, qty := range inventory {
// 			totalItems += qty
// 		}
// 		capMax := 10
// 		fmt.Println("===== INVENTAIRE =====")
// 		fmt.Printf("Inventaire (%d/%d) :\n", totalItems, capMax)
// 		if totalItems == 0 {
// 			fmt.Println("  (vide)")
// 		} else {
// 			i := 1
// 			for name, qty := range inventory {
// 				fmt.Printf("  %d) %s x%d\n", i, name, qty)
// 				i++
// 			}
// 		}
// 		fmt.Println("9) Retour")
// 		var choice int
// 		fmt.Print("> ")
// 		fmt.Scan(&choice)
// 		if choice == 9 {
// 			break
// 		}
// 		i := 1
// 		for name := range inventory {
// 			if i == choice {
// 				fmt.Printf("Vous utilisez %s ! (+30 PV)\n", name)
// 				player.CurrentHP += 30
// 				if player.CurrentHP > player.MaxHP {
// 					player.CurrentHP = player.MaxHP
// 				}
// 				inventory[name]--
// 				if inventory[name] <= 0 {
// 					delete(inventory, name)
// 				}
// 				break
// 			}
// 			i++
// 		}
// 	}
// }

// ====================
// Fonction générique pour combattre un boss
// ====================
func StartBossFight(player *Monster, boss *Monster, inventory map[string]int) bool {
	opts := []string{"Attaque normale", "Attaque spéciale", "Inventaire", "Fuir"}
	selected := 0

	err := keyboard.Open()
	keyboardOpened := (err == nil)
	defer func() { _ = keyboard.Close() }()

	for player.CurrentHP > 0 && boss.CurrentHP > 0 {
		fmt.Print("\033[H\033[2J")
		fmt.Println("=== Combat ===")
		displayHP(player, boss)
		fmt.Println("=== À ton tour ! ===")
		for i, o := range opts {
			prefix := "  "
			if i == selected {
				prefix = "→ "
			}
			fmt.Println(prefix + o)
		}

		actionChosen := false
		for !actionChosen {
			if keyboardOpened {
				_, key, _ := keyboard.GetKey()
				switch key {
				case keyboard.KeyArrowUp:
					if selected > 0 {
						selected--
					}
				case keyboard.KeyArrowDown:
					if selected < len(opts)-1 {
						selected++
					}
				case keyboard.KeyEnter:
					actionChosen = true
				}
			} else {
				var ch int
				fmt.Printf("Choisis (1: Attaque, 2: Spéciale, 3: Inventaire, 4: Fuir) > ")
				fmt.Scan(&ch)
				selected = ch - 1
				actionChosen = true
			}
		}

		playerTurn := true
		switch selected {
		case 0:
			damage := player.Attack(boss)
			slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", player.Name, damage, boss.Name), 30*time.Millisecond)
		case 1:
			damage := player.SpecialAttack(boss)
			if damage > 0 {
				slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", player.Name, damage, boss.Name), 30*time.Millisecond)
			}
		case 2:
			showInventory(player, inventory)
			playerTurn = false
		case 3:
			slowPrint("Vous fuyez le combat !", 50*time.Millisecond)
			return false
		}

		if boss.CurrentHP <= 0 {
			slowPrint(fmt.Sprintf("%s a été vaincu !", boss.Name), 50*time.Millisecond)
			return true
		}

		if playerTurn {
			time.Sleep(500 * time.Millisecond)
			if rand.Intn(100) < 40 {
				damage := boss.SpecialAttack(player)
				slowPrint(fmt.Sprintf("%s utilise une attaque spéciale et inflige %d dégâts !", boss.Name, damage), 30*time.Millisecond)
			} else {
				damage := boss.Attack(player)
				slowPrint(fmt.Sprintf("%s attaque et inflige %d dégâts !", boss.Name, damage), 30*time.Millisecond)
			}

			if player.CurrentHP <= 0 {
				slowPrint("Vous avez été vaincu...", 50*time.Millisecond)
				return false
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	return player.CurrentHP > 0
}

// ====================
// Fonction principale des combats de boss
// ====================
func StartAllBossFights() {
	rand.Seed(time.Now().UnixNano())

	player := &Monster{Name: "Joueur", MaxHP: 150, CurrentHP: 150, AttackPower: 25}
	playerInventory := map[string]int{"Potion": 5}

	bosses := []*Monster{
		{Name: "Rogue Ninja", MaxHP: 200, CurrentHP: 200, AttackPower: 20},
		{Name: "Dragon de Feu", MaxHP: 300, CurrentHP: 300, AttackPower: 35},
		{Name: "Yone", MaxHP: 500, CurrentHP: 500, AttackPower: 50},
	}

	for _, boss := range bosses {
		success := StartBossFight(player, boss, playerInventory)
		if !success {
			slowPrint("Vous avez échoué à vaincre le boss. Recommencez !", 50*time.Millisecond)
			return
		}
	}

	slowPrint("Félicitations ! Vous avez vaincu tous les boss !", 50*time.Millisecond)
}
