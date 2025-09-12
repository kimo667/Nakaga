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
	Skills    []string // techniques/sorts connus
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

// ==== Utils skills ====

func hasSkill(c Character, spell string) bool {
	for _, s := range c.Skills {
		if s == spell {
			return true
		}
	}
	return false
}

func learnSkill(c *Character, spell string) bool {
	if hasSkill(*c, spell) {
		return false
	}
	c.Skills = append(c.Skills, spell)
	return true
}

func spellBook(c *Character) {
	if learnSkill(c, "Mur de vent") {
		fmt.Println("Vous avez appris : Mur de vent !")
	} else {
		fmt.Println("Vous connaissez déjà : Mur de vent.")
	}
}

// ==== Initialisation du personnage ====
// Ajoute automatiquement la technique de base "Tempête d'acier"
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	ch := Character{
		Name:      name,
		Class:     class,
		Level:     level,
		HPMax:     hpMax,
		HP:        hp,
		Inventory: inv,
		Skills:    []string{},
	}
	learnSkill(&ch, "Tempête d'acier")
	return ch
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

	fmt.Println("Compétences :")
	if len(c.Skills) == 0 {
		fmt.Println("  (aucune technique)")
	} else {
		for _, s := range c.Skills {
			fmt.Printf("  - %s\n", s)
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

// TÂCHE 9 : Potion de poison — 10 dégâts par seconde ×3
// S'arrête immédiatement si mort (après revive T8)
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

		if IsDead(c) {
			fmt.Println("L'effet du poison est interrompu suite à votre mort.")
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("L'effet du poison est terminé. PV restants : %d/%d\n", c.HP, c.HPMax)
}

// Utiliser le Livre : Mur de vent
func UseSpellBookWind(c *Character) {
	if c.Inventory["Livre de Sort : Mur de vent"] <= 0 {
		fmt.Println("Vous n'avez pas de 'Livre de Sort : Mur de vent'.")
		return
	}
	if hasSkill(*c, "Mur de vent") {
		fmt.Println("Vous connaissez déjà 'Mur de vent'. Le livre n'a pas été consommé.")
		return
	}
	removeInventory(c, "Livre de Sort : Mur de vent", 1)
	spellBook(c)
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
// TÂCHE 8 : si HP <= 0 -> WASTED + revive à 50% PV max (continuer le jeu)
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
		fmt.Println("3) Utiliser 'Livre de Sort : Mur de vent'")
		fmt.Println("9) Retour")
		switch readChoice(r) {
		case "1":
			TakePot(c)
			IsDead(c)
		case "2":
			PoisonPot(c)
		case "3":
			UseSpellBookWind(c)
		case "9", "retour", "back":
			return
		default:
			fmt.Println("Choix invalide.")
		}
	}
}

// ==== Marchand ====
// RedBull gratuite une fois
var redbullAvailable = true

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n=== MARCHAND ===")
		if redbullAvailable {
			fmt.Println("1) RedBull — GRATUIT")
		} else {
			fmt.Println("1) RedBull — (ÉPUISÉ)")
		}
		fmt.Println("2) Potion de poison — GRATUIT (temporaire)")
		fmt.Println("3) Livre de Sort : Mur de vent — 0 or (GRATUIT)")
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
		case "3":
			addInventory(c, "Livre de Sort : Mur de vent", 1)
			fmt.Println("Achat effectué ! Vous avez obtenu : Livre de Sort : Mur de vent")
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
		"RedBull":                     3,
		"Potion de poison":            0,
		"Livre de Sort : Mur de vent": 0,
	})
	reader := bufio.NewReader(os.Stdin)

	for mainMenu(&c, reader) {
		IsDead(&c) // revive auto si besoin, on continue le jeu
	}
}
