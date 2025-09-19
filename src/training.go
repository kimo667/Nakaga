package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/eiannone/keyboard"
)

type Monster struct {
	Name        string
	MaxHP       int
	CurrentHP   int
	AttackPower int
	Initiative  int
}

// ===== UI =====
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

func slowPrint(s string, d time.Duration) {
	for _, c := range s {
		fmt.Printf("%c", c)
		time.Sleep(d)
	}
	fmt.Println()
}

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

// inventaire d'entraînement
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
			for name, qty := range inventory {
				fmt.Printf("  %d) %s x%d\n", i, name, qty)
				i++
			}
		}
		fmt.Println("9) Retour")
		var choice int
		fmt.Print("> ")
		fmt.Scan(&choice)
		if choice == 9 {
			break
		}
		// simplifié : n'importe quel item soigne 20
		player.CurrentHP += 20
		if player.CurrentHP > player.MaxHP {
			player.CurrentHP = player.MaxHP
		}
		fmt.Println("Vous utilisez un objet ! (+20 PV)")
		return
	}
}

// ===== Combat d’entraînement =====
func StartTrainingFight() bool {
	rand.Seed(time.Now().UnixNano())

	player := &Monster{Name: "Joueur", MaxHP: 100, CurrentHP: 100, AttackPower: 20, Initiative: rollInitiative(10)}
	mob := &Monster{Name: "Gobelin d’entrainement", MaxHP: 50, CurrentHP: 50, AttackPower: 15, Initiative: rollInitiative(8)}

	// Mana d’entraînement local (Mission 3/4)
	playerMana := 20
	playerManaMax := 20

	fmt.Println("=== Monologue du Gobelin ===")
	slowPrint("Salut ! Je vais t'entraîner au combat.", 40*time.Millisecond)
	slowPrint("Es-tu prêt ? (oui/non)", 40*time.Millisecond)
	var ready string
	fmt.Scanln(&ready)
	if strings.ToLower(strings.TrimSpace(ready)) != "oui" {
		fmt.Println("Reviens quand tu seras prêt !")
		return false
	}

	opts := []string{"Attaque normale", "Attaque spéciale", "Sorts", "Inventaire", "Fuir"}
	selected := 0

	err := keyboard.Open()
	keyboardOpened := (err == nil)
	defer func() { _ = keyboard.Close() }()

	round := 1
	for player.CurrentHP > 0 && mob.CurrentHP > 0 {
		// initiative par tour (Mission 1)
		player.Initiative = rollInitiative(10)
		mob.Initiative = rollInitiative(8)
		playerFirst := player.Initiative >= mob.Initiative

		fmt.Print("\033[H\033[2J")
		fmt.Printf("=== Combat (Tour %d) ===\n", round)
		fmt.Printf("Initiative -> Joueur:%d  Gobelin:%d  (%s commence)\n",
			player.Initiative, mob.Initiative, map[bool]string{true: "Joueur", false: "Gobelin"}[playerFirst])
		displayHP(player, mob)
		fmt.Println("=== À ton tour ! ===")
		for i, o := range opts {
			prefix := "  "
			if i == selected {
				prefix = "→ "
			}
			fmt.Println(prefix + o)
		}

		actionIdx := -1
		for actionIdx < 0 {
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
					actionIdx = selected
				}
				fmt.Print("\033[H\033[2J")
				fmt.Printf("=== Combat (Tour %d) ===\n", round)
				fmt.Printf("Initiative -> Joueur:%d  Gobelin:%d  (%s commence)\n",
					player.Initiative, mob.Initiative, map[bool]string{true: "Joueur", false: "Gobelin"}[playerFirst])
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
				fmt.Print("Choisis (1: Attaque, 2: Spéciale, 3: Sorts, 4: Inventaire, 5: Fuir) > ")
				fmt.Scan(&ch)
				actionIdx = ch - 1
			}
		}

		attackPhase := func(att, def *Monster) {
			switch actionIdx {
			case 0:
				d := att.Attack(def)
				slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", att.Name, d, def.Name), 25*time.Millisecond)
			case 1:
				d := att.SpecialAttack(def)
				if d > 0 {
					slowPrint(fmt.Sprintf("%s inflige %d dégâts à %s !", att.Name, d, def.Name), 25*time.Millisecond)
				}
			case 2:
				// Sorts (Mission 3/4)
				fmt.Println("Sorts: 1) Coup de poing (8 dmg, 2 mana)  2) Boule de feu (18 dmg, 5 mana)  9) Retour")
				var s int
				fmt.Print("> ")
				fmt.Scan(&s)
				if s == 1 {
					if playerMana < 2 {
						fmt.Println(CRed + "Mana insuffisant." + CReset)
					} else {
						playerMana -= 2
						def.CurrentHP -= 8
						if def.CurrentHP < 0 {
							def.CurrentHP = 0
						}
						fmt.Printf("Coup de poing ! -8 PV (Mana %d/%d)\n", playerMana, playerManaMax)
					}
				} else if s == 2 {
					if playerMana < 5 {
						fmt.Println(CRed + "Mana insuffisant." + CReset)
					} else {
						playerMana -= 5
						def.CurrentHP -= 18
						if def.CurrentHP < 0 {
							def.CurrentHP = 0
						}
						fmt.Printf("Boule de feu ! -18 PV (Mana %d/%d)\n", playerMana, playerManaMax)
					}
				} else {
					fmt.Println("(annulé)")
				}
			case 3:
				showInventory(player, map[string]int{"RedBull": 1})
			case 4:
				fmt.Println("Vous fuyez l'entraînement.")
				player.CurrentHP = 0
			}
		}

		// ordre selon initiative
		if playerFirst {
			attackPhase(player, mob)
			if mob.CurrentHP > 0 && player.CurrentHP > 0 && actionIdx != 4 {
				_ = mob.Attack(player)
			}
		} else {
			_ = mob.Attack(player)
			if mob.CurrentHP > 0 && player.CurrentHP > 0 && actionIdx != 4 {
				attackPhase(player, mob)
			}
		}

		if mob.CurrentHP <= 0 || player.CurrentHP <= 0 {
			break
		}

		round++
		time.Sleep(250 * time.Millisecond)
	}

	if player.CurrentHP <= 0 {
		fmt.Println(CRed + "Défaite en entraînement." + CReset)
		return false
	}
	fmt.Println(CGreen + "Victoire en entraînement !" + CReset)
	return true
}
