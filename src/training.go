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
}

// ===== UI PV =====

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

// ===== Actions =====

func (m *Monster) Attack(target *Monster) {
	damage := rand.Intn(m.AttackPower) + 1
	if rand.Intn(100) < 10 {
		damage *= 2
		slowPrint("Coup critique !", 50*time.Millisecond)
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
}

func (m *Monster) SpecialAttack(target *Monster) {
	if rand.Intn(100) < 25 {
		slowPrint(fmt.Sprintf("%s a raté son attaque spéciale !", m.Name), 30*time.Millisecond)
		return
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
}

func (m *Monster) Heal(amount int) {
	m.CurrentHP += amount
	if m.CurrentHP > m.MaxHP {
		m.CurrentHP = m.MaxHP
	}
}

// ===== Entrée entraînement, appelée depuis le menu =====

func StartTraining() {
	rand.Seed(time.Now().UnixNano())

	player := &Monster{Name: "Joueur", MaxHP: 100, CurrentHP: 100, AttackPower: 20}
	mob := &Monster{Name: "Sbire", MaxHP: 50, CurrentHP: 50, AttackPower: 15}

	fmt.Println("=== Monologue du Sbire ===")
	fmt.Println("Salut, jeune héros ! Je suis là pour t'entraîner.")
	fmt.Println("Es-tu prêt à commencer ton entraînement ? (oui/non)")
	var ready string
	fmt.Scanln(&ready)
	if strings.ToLower(strings.TrimSpace(ready)) != "oui" {
		fmt.Println("Reviens quand tu seras prêt !")
		return
	}

	opts := []string{"Attaque normale", "Attaque spéciale", "Se soigner"}
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

		if keyboardOpened {
			actionChosen := false
			for !actionChosen {
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
					switch selected {
					case 0:
						player.Attack(mob)
						slowPrint(fmt.Sprintf("%s attaque %s !", player.Name, mob.Name), 30*time.Millisecond)
					case 1:
						player.SpecialAttack(mob)
						slowPrint(fmt.Sprintf("%s utilise une attaque spéciale !", player.Name), 30*time.Millisecond)
					case 2:
						player.Heal(20)
						slowPrint(fmt.Sprintf("%s se soigne de 20 PV !", player.Name), 30*time.Millisecond)
					}
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
			}
		} else {
			var ch int
			fmt.Print("Choisis (1: Attaque, 2: Spéciale, 3: Soin) > ")
			fmt.Scan(&ch)
			switch ch {
			case 1:
				player.Attack(mob)
				slowPrint(fmt.Sprintf("%s attaque %s !", player.Name, mob.Name), 30*time.Millisecond)
			case 2:
				player.SpecialAttack(mob)
				slowPrint(fmt.Sprintf("%s utilise une attaque spéciale !", player.Name), 30*time.Millisecond)
			default:
				player.Heal(20)
				slowPrint(fmt.Sprintf("%s se soigne de 20 PV !", player.Name), 30*time.Millisecond)
			}
		}

		if mob.CurrentHP <= 0 {
			slowPrint(mob.Name+" a été vaincu !", 50*time.Millisecond)
			break
		}

		// Tour du sbire
		slowPrint("=== Tour du Sbire ===", 50*time.Millisecond)
		time.Sleep(300 * time.Millisecond)
		switch rand.Intn(3) {
		case 0:
			mob.Attack(player)
			slowPrint(fmt.Sprintf("%s attaque %s !", mob.Name, player.Name), 30*time.Millisecond)
		case 1:
			mob.SpecialAttack(player)
			slowPrint(fmt.Sprintf("%s lance une attaque spéciale !", mob.Name), 30*time.Millisecond)
		case 2:
			mob.Heal(10)
			slowPrint(fmt.Sprintf("%s se soigne de 10 PV !", mob.Name), 30*time.Millisecond)
		}

		if player.CurrentHP <= 0 {
			slowPrint(player.Name+" a été vaincu !", 50*time.Millisecond)
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println("\nFin de l’entraînement.")
}
