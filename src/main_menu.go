package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

<<<<<<< HEAD
// Character struct
type character struct {
	Name  string
	Level int
	HP    int
}

// Menu item
type menuItem string

func (i menuItem) Title() string       { return string(i) }
func (i menuItem) Description() string { return "" }
func (i menuItem) FilterValue() string { return string(i) }

// Bubble Tea model
type model struct {
	list      list.Model
	character Character
}

func mainMenu() {
	menuItems := []list.Item{
		menuItem("Afficher les informations du personnage"),
		menuItem("Accéder au contenu de l’inventaire"),
		menuItem("Marchand"),
		menuItem("Forgeron"),
		menuItem("Entraînement"),
		menuItem("Quitter"),
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(menuItems, delegate, 50, 10)
	l.Title = "===== MENU ====="
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	m := model{
		list:      l,
		character: Character{Name: "Héros", Level: 1, HP: 100},
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Erreur: %v\n", err)
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				switch i {
				case "Afficher les informations du personnage":
					fmt.Printf("\nNom: %s\nLevel: %d\nHP: %d\n", m.character.Name, m.character.Level, m.character.HP)
				case "Accéder au contenu de l’inventaire":
					fmt.Println("\n[Inventaire placeholder]")
				case "Marchand":
					fmt.Println("\n[Marchand placeholder]")
				case "Forgeron":
					fmt.Println("\n[Forgeron placeholder]")
				case "Entraînement":
					fmt.Println("\n[Entraînement placeholder]")
				case "Quitter":
					return m, tea.Quit
				}
				fmt.Println("\nAppuyez sur Entrée pour revenir au menu...")
				fmt.Scanln()
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
=======
// Menu principal
func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("5) Entraînement (combat d'essai)")
	fmt.Println("6) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)

	case "2", "inventaire":
		inventoryMenu(c, r)

	case "3", "marchand", "shop":
		merchantMenu(c, r)

	case "4", "forgeron", "forge":
		blacksmithMenu(c, r)

	case "5", "entrainement", "entraînement", "training", "combat":
		// Lancement du mode entraînement
		StartTraining()

	case "6", "q", "quit", "quitter", "exit":
		fmt.Println("Au revoir !")
		return false

	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}

	return true
>>>>>>> 469c8fe15dcbdd88bf7b51ebb19de423d7f498da
}
