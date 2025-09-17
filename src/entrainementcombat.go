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

// Affichage d'une barre de PV colorée
func HPBar(current, max int) string {
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
	totalWidth := 20
	filledWidth := current * totalWidth / max
	emptyWidth := totalWidth - filledWidth
	barStyle := lipgloss.NewStyle().Background(lipgloss.Color(color)).Width(filledWidth)
	emptyStyle := lipgloss.NewStyle().Background(lipgloss.Color("#555555")).Width(emptyWidth)
	return barStyle.Render("") + emptyStyle.Render("")
}

// Affichage des PV côte à côte
func DisplayHP(player, monster Monster) {
	sbireLine := fmt.Sprintf("%s : [%s] %d/%d PV", monster.Name, HPBar(monster.CurrentHP, monster.MaxHP), monster.CurrentHP, monster.MaxHP)
	playerLine := fmt.Sprintf("%s : [%s] %d/%d PV", player.Name, HPBar(player.CurrentHP, player.MaxHP), player.CurrentHP, player.MaxHP)

	totalWidth := 60
	space := totalWidth - len(sbireLine) - len(playerLine)
	if space < 2 {
		space = 2
	}
	fmt.Println(sbireLine + strings.Repeat(" ", space) + playerLine + "\n")
}

// Affichage progressif du texte (pour effet ralenti)
func SlowPrint(text string, delay time.Duration) {
	for _, c := range text {
		fmt.Printf("%c", c)
		time.Sleep(delay)
	}
	fmt.Println()
}

func (m *Monster) Attack(target *Monster) {
	damage := rand.Intn(m.AttackPower) + 1
	if rand.Intn(100) < 10 {
		damage *= 2
		SlowPrint("Coup critique !", 50*time.Millisecond)
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	SlowPrint(fmt.Sprintf("%s attaque %s et inflige %d points de dégâts !", m.Name, target.Name, damage), 30*time.Millisecond)
	DisplayHP(player, monster)
}

func (m *Monster) SpecialAttack(target *Monster) {
	if rand.Intn(100) < 25 {
		SlowPrint(fmt.Sprintf("%s a raté son attaque spéciale !", m.Name), 30*time.Millisecond)
		return
	}
	damage := rand.Intn(m.AttackPower*2) + m.AttackPower
	if rand.Intn(100) < 15 {
		damage *= 2
		SlowPrint("Coup critique sur l'attaque spéciale !", 50*time.Millisecond)
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	SlowPrint(fmt.Sprintf("%s utilise une attaque spéciale sur %s et inflige %d points de dégâts !", m.Name, target.Name, damage), 30*time.Millisecond)
	DisplayHP(player, monster)
}

func (m *Monster) Heal(amount int) {
	m.CurrentHP += amount
	if m.CurrentHP > m.MaxHP {
		m.CurrentHP = m.MaxHP
	}
	SlowPrint(fmt.Sprintf("%s se soigne de %d PV !", m.Name, amount), 30*time.Millisecond)
	DisplayHP(player, monster)
}

var player Monster
var monster Monster

func main() {
	rand.Seed(time.Now().UnixNano())

	player = Monster{Name: "Joueur", MaxHP: 100, CurrentHP: 100, AttackPower: 20}
	monster = Monster{Name: "Sbire", MaxHP: 50, CurrentHP: 50, AttackPower: 15}

	fmt.Println("=== Monologue du Sbire ===")
	fmt.Println("Salut, jeune héros ! Je suis là pour t'entraîner.")
	fmt.Println("Es-tu prêt à commencer ton entraînement ? (oui/non)")
	var ready string
	fmt.Scanln(&ready)
	if ready != "oui" && ready != "Oui" {
		fmt.Println("Reviens quand tu seras prêt !")
		return
	}

	menuOptions := []string{"Attaque normale", "Attaque spéciale", "Se soigner"}
	selected := 0

	_ = keyboard.Open()
	defer keyboard.Close()

	for player.CurrentHP > 0 && monster.CurrentHP > 0 {
		fmt.Print("\033[H\033[2J")
		fmt.Println("=== Combat ===")
		DisplayHP(player, monster)

		// Menu joueur
		fmt.Println("=== À ton tour ! ===")
		for i, option := range menuOptions {
			prefix := "  "
			if i == selected {
				prefix = "→ "
			}
			fmt.Println(prefix + option)
		}

		actionChosen := false
		for !actionChosen {
			_, key, _ := keyboard.GetKey()
			switch key {
			case keyboard.KeyArrowUp:
				if selected > 0 {
					selected--
				}
			case keyboard.KeyArrowDown:
				if selected < len(menuOptions)-1 {
					selected++
				}
			case keyboard.KeyEnter:
				switch selected {
				case 0:
					player.Attack(&monster)
				case 1:
					player.SpecialAttack(&monster)
				case 2:
					player.Heal(20)
				}
				actionChosen = true
			}

			// Rafraîchir menu
			fmt.Print("\033[H\033[2J")
			fmt.Println("=== Combat ===")
			DisplayHP(player, monster)
			fmt.Println("=== À ton tour ! ===")
			for i, option := range menuOptions {
				prefix := "  "
				if i == selected {
					prefix = "→ "
				}
				fmt.Println(prefix + option)
			}
		}

		if monster.CurrentHP <= 0 {
			SlowPrint(monster.Name+" a été vaincu !", 50*time.Millisecond)
			break
		}

		// Tour du Sbire au ralenti
		SlowPrint("=== Tour du Sbire ===", 50*time.Millisecond)
		time.Sleep(300 * time.Millisecond) // pause avant l'action
		monsterChoice := rand.Intn(3)
		switch monsterChoice {
		case 0:
			monster.Attack(&player)
		case 1:
			monster.SpecialAttack(&player)
		case 2:
			monster.Heal(10)
		}

		if player.CurrentHP <= 0 {
			SlowPrint(player.Name+" a été vaincu !", 50*time.Millisecond)
			break
		}
		selected = 0
		time.Sleep(300 * time.Millisecond)
	}
}
