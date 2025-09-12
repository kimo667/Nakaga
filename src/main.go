package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

type Character struct {
	Name      string
	Class     Classe
	Level     int
	HPMax     int
	HP        int
	Inventory map[string]int
}

// Initialisation du personnage
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	return Character{
		Name:      name,
		Class:     class,
		Level:     level,
		HPMax:     hpMax,
		HP:        hp,
		Inventory: inv,
	}
}

// Affichage info + ASCII art
func displayInfo(c Character) {
	asciiArt := `
⠀⠀⠀⢰⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠘⡇⠀⠀⠀⢠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢷⠀⢠⢣⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢘⣷⢸⣾⣇⣶⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣿⣿⣿⣹⣿⣿⣷⣿⣆⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢼⡇⣿⣿⣽⣶⣶⣯⣭⣷⣶⣿⣿⣶⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠸⠣⢿⣿⣿⣿⣿⡿⣛⣭⣭⣭⡙⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⠿⠿⠿⢯⡛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣾⣿⡿⡷⢿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⡔⣺⣿⣿⣽⡿⣿⣿⣿⣟⡳⠦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⣭⣾⣿⠃⣿⡇⣿⣿⡷⢾⣭⡓⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣾⣿⡿⠷⣿⣿⡇⣿⣿⣟⣻⠶⣭⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣋⣵⣞⣭⣮⢿⣧⣝⣛⡛⠿⢿⣦⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⣀⣀⣠⣶⣿⣿⣿⣿⡿⠟⣼⣿⡿⣟⣿⡇⠀⠙⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⡼⣿⣿⣿⢟⣿⣿⣿⣷⡿⠿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠉⠁⠀⢉⣭⣭⣽⣯⣿⣿⢿⣫⣮⣅⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`
	fmt.Println(asciiArt)
	fmt.Println("=== Informations du personnage ===")
	fmt.Printf("Nom   : %s\n", c.Name)
	fmt.Printf("Classe: %s\n", c.Class)
	fmt.Printf("Niveau: %d\n", c.Level)
	fmt.Printf("PV    : %d / %d\n", c.HP, c.HPMax)
	fmt.Println("Inventaire :")
	if len(c.Inventory) == 0 {
		fmt.Println("  (vide)")
	} else {
		for item, qty := range c.Inventory {
			fmt.Printf("  - %s x%d\n", item, qty)
		}
	}
}

// Utiliser une potion
func TakePot(c *Character) {
	if qty, ok := c.Inventory["RedBull"]; ok && qty > 0 {
		c.HP += 20
		if c.HP > c.HPMax {
			c.HP = c.HPMax
		}
		c.Inventory["RedBull"]--
		fmt.Println("RedBull ! PV =", c.HP)
		return
	}
	fmt.Println("Pas de RedBull dans l'inventaire !")
}

// Ouvrir l'inventaire
func OpenInventory(c Character) {
	if len(c.Inventory) == 0 {
		fmt.Println("L'inventaire est vide.")
		return
	}

	fmt.Println("Inventaire :")
	for item, qty := range c.Inventory {
		fmt.Printf("  - %s x%d\n", item, qty)
	}
}

// Vérifier si mort
func IsDead(c Character) bool {
	return c.HP <= 0
}

func main() {
	c := initCharacter("Yazuo", ClasseSamurai, 1, 100, 40, map[string]int{"RedBull": 3})

	displayInfo(c)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nCommande > ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "boire une redbull":
			TakePot(&c)
		case "ouvrir l'inventaire":
			OpenInventory(c)
		case "information":
			displayInfo(c)
		case "quitter":
			fmt.Println("Au revoir !")
			return
		default:
			fmt.Println("Commande inconnue. Commandes valides : 'boire une redbull', 'ouvrir l'inventaire', 'information', 'quitter'")
		}

		if IsDead(c) {
			fmt.Println("Wasted ! Le personnage est mort.")
			return
		}
	}
}
