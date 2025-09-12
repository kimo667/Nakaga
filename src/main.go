package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

// ==== Utils inventaire ====

func addInventory(c *Character, item string, qty int) {
	if qty <= 0 {
		return
	}
	if c.Inventory == nil {
		c.Inventory = make(map[string]int)
	}
	c.Inventory[item] += qty
}

func removeInventory(c *Character, item string, qty int) bool {
	cur, ok := c.Inventory[item]
	if !ok || qty <= 0 || cur < qty {
		return false
	}
	newQty := cur - qty
	if newQty == 0 {
		delete(c.Inventory, item)
	} else {
		c.Inventory[item] = newQty
	}
	return true
}

// ==== Initialisation du personnage ====

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

// ==== Affichage info + ASCII art ====

func displayInfo(c Character) {
	asciiArt := `
⠀⠀⠀⠀⠀⠀⢰⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
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
⠀⠀⠀⠀⢀⣿⣟⣽⣿⣿⣿⣿⣾⣿⣿⣯⡛⠻⢷⣶⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⡞⣾⣿⣿⣿⣿⡟⣿⣿⣽⣿⣿⡿⠀⠀⠀⠈⠙⠛⠿⣶⣤⣄⡀⠀⠀
⠀⠀⠀⣾⣸⣿⣿⣷⣿⣿⢧⣿⣿⣿⣿⣿⣷⠁⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⢷⣦
⠀⠀⠀⡿⣛⣛⣛⣛⣿⣿⣸⣿⣿⣿⣻⣿⣿⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢸⡇⣿⣿⣿⣿⣿⡏⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢰⣶⣶⣶⣾⣿⢃⣿⣿⣿⣿⣯⣿⣭⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
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

// ==== Consommables ====

func TakePot(c *Character) {
	if removeInventory(c, "RedBull", 1) {
		c.HP += 20
		if c.HP > c.HPMax {
			c.HP = c.HPMax
		}
		fmt.Println("Tu as bu une RedBull ! PV =", c.HP)
		return
	}
	fmt.Println("Pas de RedBull dans l'inventaire !")
}

// TÂCHE 9 : Potion de poison — inflige 10 dégâts par seconde pendant 3s
func PoisonPot(c *Character) {
	if !removeInventory(c, "Potion de poison", 1) {
		fmt.Println("Vous n'avez pas de Potion de poison.")
		return
	}
	fmt.Println("Vous utilisez une Potion de poison…")
	for i := 1; i <= 3; i++ {
		c.HP -= 10
		if c.HP < 0 {
			c.HP = 0
		}
		fmt.Printf("Effet poison %d/3 → PV: %d/%d\n", i, c.HP, c.HPMax)

		// Vérifie si mort → revive + stop poison
		if IsDead(c) {
			fmt.Println("L'effet du poison est interrompu suite à votre mort.")
			return
		}

		time.Sleep(1 * time.Second)
	}
	// ✅ Affichage final après la fin de l'effet
	fmt.Printf("L'effet du poison est terminé. PV restants : %d/%d\n", c.HP, c.HPMax)
}

// ==== Inventaire ====

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

// ==== Statut ====

// TÂCHE 8 : si HP <= 0 -> WASTED + revive à 50% PV max (et on continue le jeu)
func IsDead(c *Character) bool {
	if c.HP <= 0 {
		fmt.Println("\n*** WASTED ***")
		c.HP = c.HPMax / 2
		fmt.Printf("Vous êtes ressuscité avec %d/%d PV.\n", c.HP, c.HPMax)
		return true
	}
	return false
}

// ==== Lecture entrée ====

func readChoice(r *bufio.Reader) string {
	fmt.Print("> ")
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(line))
}

// ==== Sous-menu Inventaire ====

func inventoryMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n--- INVENTAIRE ---")
		OpenInventory(*c)
		fmt.Println("\n1) Boire une RedBull (+20 PV)")
		fmt.Println("2) Utiliser une Potion de poison (10 dmg/s ×3)")
		fmt.Println("9) Retour")
		switch readChoice(r) {
		case "1":
			TakePot(c)
			IsDead(c) // le perso revive à 50%, on CONTINUE le jeu
		case "2":
			PoisonPot(c)
		case "9", "retour", "back":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// ==== Marchand ====

var redbullAvailable = true // disponibilité RedBull gratuite

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n=== MARCHAND ===")
		if redbullAvailable {
			fmt.Println("1) RedBull — GRATUIT")
		} else {
			fmt.Println("1) RedBull — (ÉPUISÉ)")
		}
		fmt.Println("2) Potion de poison — GRATUIT (temporaire)")
		fmt.Println("9) Retour")
		switch readChoice(r) {
		case "1":
			if redbullAvailable {
				addInventory(c, "RedBull", 1)
				redbullAvailable = false
				fmt.Printf("Achat effectué ! Vous avez obtenu : RedBull (total: %d)\n", c.Inventory["RedBull"])
			} else {
				fmt.Println("La RedBull gratuite n’est plus disponible.")
			}
		case "2":
			addInventory(c, "Potion de poison", 1)
			fmt.Printf("Achat effectué ! Vous avez obtenu : Potion de poison (total: %d)\n", c.Inventory["Potion de poison"])
		case "9", "retour", "back":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// ==== Menu principal ====

func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n===== MENU =====")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)
	case "2", "inventaire":
		inventoryMenu(c, r)
	case "3", "marchand", "shop":
		merchantMenu(c, r)
	case "4", "q", "quit", "quitter":
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println("Choix invalide.")
	}
	return true
}

func main() {
	c := initCharacter("Yazuo", ClasseSamurai, 1, 100, 40, map[string]int{
		"RedBull":          3,
		"Potion de poison": 0, // tu peux laisser 0; on en récupère chez le marchand
	})

	reader := bufio.NewReader(os.Stdin)

	for mainMenu(&c, reader) {
		IsDead(&c) // pas besoin de quitter : on revive et on continue
	}
}
