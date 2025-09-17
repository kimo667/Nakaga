package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Character struct
type Character struct {
	Name  string
	Level int
	HP    int
}

// Menu options
var menuOptions = []string{
	"Afficher les informations du personnage",
	"Accéder au contenu de l’inventaire",
	"Marchand",
	"Forgeron",
	"Entraînement",
	"Quitter",
}

// Bubble Tea model
type model struct {
	cursor    int
	character Character
}

// Main
func main() {
	m := model{
		cursor:    0,
		character: Character{Name: "Héros", Level: 1, HP: 100},
	}

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Println("Erreur:", err)
		os.Exit(1)
	}
}

// Bubble Tea interface
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(menuOptions)-1 {
				m.cursor++
			}
		case "enter":
			fmt.Println("\nVous avez sélectionné:", menuOptions[m.cursor])
			fmt.Println("Appuyez sur Entrée pour revenir au menu...")
			fmt.Scanln()
		}
	}
	return m, nil
}

func (m model) View() string {
	s := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")).Bold(true).Render("===== MENU =====") + "\n\n"
	for i, choice := range menuOptions {
		prefix := "  "
		if i == m.cursor {
			prefix = "→ "
			choice = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")).Bold(true).Render(choice)
		}
		s += fmt.Sprintf("%s%s\n", prefix, choice)
	}
	s += "\nUtilisez les flèches ↑ ↓ et Entrée pour sélectionner, q pour quitter."
	return s
}
