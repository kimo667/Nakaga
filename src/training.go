package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/eiannone/keyboard"
)

// ====================
// Définition du monstre/joueur
// ====================
type Monster struct {
	Name        string
	MaxHP       int
	CurrentHP   int
	AttackPower int
}

// ====================
// UI PV
// ====================
func hpBar(current, max int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}
	percent := float64(current) / float64(max)
	var color string
	switch {
	case percent > 0.5:
		color = "#00FF00"
	case percent > 0.2:
		color = "#FFFF00"
	default:
		color = "#FF0000"
	}
	const totalWidth = 20
	filled := current * totalWidth / max
	empty := totalWidth - filled
	bar := lipgloss.NewStyle().Background(lipgloss.Color(color)).Width(filled).Render("")
	rest := lipgloss.NewStyle().Background(lipgloss.Color("#555555")).Width(empty).Render("")
	return bar + rest
}

func displayHP(player, mob *Monster) {
	left := fmt.Sprintf("%s : [%s] %d/%d PV", mob.Name, hpBar(mob.CurrentHP, mob.MaxHP), mob.CurrentHP, mob.MaxHP)
	right := fmt.Sprintf("%s : [%s] %d/%d PV", player.Name, hpBar(player.CurrentHP, player.MaxHP), player.CurrentHP, player.MaxHP)
	totalWidth := 60
	space := totalWidth - len(left) - len(right)
	if space < 2 {
		space = 2
	}
	fmt.Println(left + strings.Repeat(" ", space) + right + "\n")
}

// ====================
// Effets visuels
// ====================
func slowPrint(s string, d time.Duration) {
	for _, c := range s {
		fmt.Printf("%c", c)
		time.Sleep(d)
	}
	fmt.Println()
}

// ====================
// Actions
// ====================
func (m *Monster) Attack(target *Monster) int {
	damage := rand.Intn(m.AttackPower) + 1
	if rand.Intn(100) < 10 {
		damage *= 2
		slowPrint("Coup critique !", 50*time.Millisecond)
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	return damage
}

func (m *Monster) SpecialAttack(target *Monster) int {
	if rand.Intn(100) < 25 {
		slowPrint(fmt.Sprintf("%s a raté son attaque spéciale !", m.Name), 30*time.Millisecond)
		return 0
	}
	damage := rand.Intn(m.AttackPower*2) + m.AttackPower
	if rand.Intn(100) < 15 {
		damage *= 2
		slowPrint("Coup critique sur l'attaque spéciale !", 50*time.Millisecond)
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	return damage
}

// ====================
// Système d’inventaire
// ====================
func showInventory(player *Monster, inventory map[string]int) {
	for {
		totalItems := 0
		for _, qty := range inventory {
			totalItems += qty
		}
		capMax := 10
		fmt.Println("===== INVENTAIRE =====")
		fmt.Printf("Inventaire (%d/%d) :\n", totalItems, capMax)
		if totalItems == 0 {
			fmt.Println("  (vide)")
		} else {
			i := 1
			names := []string{}
			for name, qty := range inventory {
				fmt.Printf("  %d) %s x%d\n", i, name, qty)
				names = append(names, name)
				i++
			}
		}
		fmt.Println("9) Retour")
		var choice int
		fmt.Print("> ")
		fmt.Scan(&choice)
		if choice == 9 {
			break // on sort de l'inventaire sans passer le tour
		}
		// utilisation d'un objet
		i := 1
		for name := range inventory {
			if i == choice {
				fmt.Printf("Vous utilisez %s ! (+20 PV)\n", name)
				player.CurrentHP += 20
				if player.CurrentHP > player.MaxHP {
					player.CurrentHP = player.MaxHP
				}
				inventory[name]--
				if inventory[name] <= 0 {
					delete(inventory, name)
				}
				break
			}
			i++
		}
	}
}

// ====================
// Combat d’entraînement
// ====================
func StartTraining() {
	rand.Seed(time.Now().UnixNano())

	player := &Monster{Name: "Joueur", MaxHP: 100, CurrentHP: 100, AttackPower: 20}
	playerInventory := map[string]int{
		"RedBull": 3, // exemple d'inventaire
	}

	mob := &Monster{Name: "Gobelin d’entrainement", MaxHP: 50, CurrentHP: 50, AttackPower: 15}

	fmt.Println("=== Monologue du Gobelin ===")
	slowPrint("Salut ! Je vais t'entraîner au combat.", 50*time.Millisecond)
	slowPrint("Es-tu prêt ? (oui/non)", 50*time.Millisecond)
	var ready string
	fmt.Scanln(&ready)
	if strings.ToLower(strings.TrimSpace(ready)) != "oui" {
		fmt.Println("Reviens quand tu seras prêt !")
		return
	}

	opts := []string{"Attaque normale", "Attaque spéciale", "Inventaire", "Fuir"}
	selected := 0

	err := keyboard.Open()
	keyboardOpened := (err == nil)
	defer func() { _ = keyboard.Close() }()

	for player.CurrentHP > 0 && mob.CurrentHP > 0 {
		fmt.Print("\033[H\033[2J")
		fmt.Println("=== Combat ===")
		displayHP(player, mob)
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
				fmt.Print("\033[H\033[2J")
				fmt.Println("=== Combat ===")
				displayHP(player, mob)
				fmt.Println("=== À ton tour ! ===")
				for i, o := range opts {
					prefix := "  "
					if i == selected {
						prefix = "→ "
					}
					fmt.Println(prefix + o)
				}
			} else {
				var ch int
				fmt.Print("Choisis (1: Attaque, 2: Spéciale, 3: Inventaire, 4: Fuir) > ")
				fmt.Scan(&ch)
				selected = ch - 1
				actionChosen = true
			}
		}

		// Actions du joueur
		playerTurn := true // par défaut, on fait le tour du gobelin après le joueur
		switch selected {
		case 0:
			damage := player.Attack(mob)
			slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", player.Name, damage, mob.Name), 30*time.Millisecond)
		case 1:
			damage := player.SpecialAttack(mob)
			if damage > 0 {
				slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", player.Name, damage, mob.Name), 30*time.Millisecond)
			}
		case 2:
			showInventory(player, playerInventory)
			playerTurn = false // ouvrir l'inventaire ne passe pas le tour
		case 3:
			slowPrint("Vous fuyez le combat !", 50*time.Millisecond)
			return
		}

		if mob.CurrentHP <= 0 {
			slowPrint(mob.Name+" a été vaincu !", 50*time.Millisecond)
			break
		}

		// Tour du gobelin si le joueur a fait une action offensive
		if playerTurn {
			slowPrint(fmt.Sprintf("=== Tour du %s ===", mob.Name), 50*time.Millisecond)
			time.Sleep(300 * time.Millisecond)
			damage := mob.Attack(player)
			slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", mob.Name, damage, player.Name), 30*time.Millisecond)

			if player.CurrentHP <= 0 {
				slowPrint(player.Name+" a été vaincu !", 50*time.Millisecond)
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}

	fmt.Println("\nFin de l’entraînement.")
}
