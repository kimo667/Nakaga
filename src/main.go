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
		fmt.Println("\033[31mVous avez appris : Mur de vent !\033[0m")
	} else {
		fmt.Println("\033[31mVous connaissez déjà : Mur de vent.\033[0m")
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
⠀⠀⠀⠀⢀⣿⣟⣽⣿⣿⣿⣿⣾⣿⣿⣯⡛⠻⢷⣶⣤⣄⡀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⡞⣾⣿⣿⣿⣿⡟⣿⣿⣽⣿⣿⡿⠀⠀⠀⠈⠙⠛⠿⣶⣤⣄⡀⠀⠀
⠀⠀⠀⣾⣸⣿⣿⣷⣿⣿⢧⣿⣿⣿⣿⣿⣷⠁⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⢷⣦
⠀⠀⠀⡿⣛⣛⣛⣛⣿⣿⣸⣿⣿⣿⣻⣿⣿⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⢸⡇⣿⣿⣿⣿⣿⡏⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢰⣶⣶⣶⣾⣿⢃⣿⣿⣿⣿⣯⣿⣭⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`
	fmt.Println(asciiArt)
	fmt.Println("\033[33m=== Informations du personnage ===\u001B[0m")
	fmt.Printf("\033[33mNom   : \u001B[0m%s\n", c.Name)
	fmt.Printf("\033[33mClasse: \u001B[0m%s\n", c.Class)
	fmt.Printf("\033[33mNiveau: \u001B[0m%d\n", c.Level)
	fmt.Printf("\033[32mPV    : \u001B[0m\033[32m%d / %d\n\u001B[0m", c.HP, c.HPMax)

	fmt.Println("\033[33mInventaire : \u001B[0m")
	if len(c.Inventory) == 0 {
		fmt.Println("\033[31m  (vide) \u001B[0m")
	} else {
		for item, qty := range c.Inventory {
			fmt.Printf("  - %s x%d\n", item, qty)
		}
	}

	fmt.Println("\033[33mCompétences : \u001B[0m")
	if len(c.Skills) == 0 {
		fmt.Println("\033[31m(aucune technique) \u001B[0m")
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
		fmt.Println(" \033[36mTu as bu une RedBull ! PV = \u001B[0m", c.HP)
		return
	}
	fmt.Println("\033[31mPas de RedBull dans l'inventaire ! \u001B[0m")
}

// TÂCHE 9 : Potion de poison — 10 dégâts par seconde ×3
// S'arrête immédiatement si mort (après revive T8)
func PoisonPot(c *Character) {
	if !removeInventory(c, "Potion de poison", 1) {
		fmt.Println("\033[31mVous n'avez pas de Potion de poison. \u001B[0m")
		return
	}
	fmt.Println("\033[36mVous utilisez une Potion de poison… \u001B[0m")
	for i := 1; i <= 3; i++ {
		c.HP -= 10
		if c.HP < 0 {
			c.HP = 0
		}
		fmt.Printf("Effet poison %d/3 → PV: %d/%d\n", i, c.HP, c.HPMax)

		if IsDead(c) {
			fmt.Println("\033[31mL'effet du poison est interrompu suite à votre mort. \u001B[0m")
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("\033[36mL'effet du poison est terminé. PV restants : %d/%d\n \u001B[0m", c.HP, c.HPMax)
}

// Utiliser le Livre : Mur de vent
func UseSpellBookWind(c *Character) {
	if c.Inventory["Livre de Sort : Mur de vent"] <= 0 {
		fmt.Println("\033[31mVous n'avez pas de 'Livre de Sort : Mur de vent'.\033[0m")
		return
	}
	if hasSkill(*c, "Mur de vent") {
		fmt.Println("\033[31mVous connaissez déjà 'Mur de vent'. Le livre n'a pas été consommé.\033[0m")
		return
	}
	removeInventory(c, "Livre de Sort : Mur de vent", 1)
	spellBook(c)
}

// ==== Inventaire ====

func OpenInventory(c Character) {
	if len(c.Inventory) == 0 {
		fmt.Println("\033[31mL'inventaire est vide. \u001B[0m")
		return
	}
	fmt.Println("\033[33mInventaire : \u001B[0m")
	for item, qty := range c.Inventory {
		fmt.Printf("  - %s x%d\n", item, qty)
	}
}

// ==== Statut ====
// TÂCHE 8 : si HP <= 0 -> WASTED + revive à 50% PV max (continuer le jeu)
func IsDead(c *Character) bool {
	if c.HP <= 0 {
		fmt.Println("\033[31m\n*** WASTED ***\033[0m")
		c.HP = c.HPMax / 2
		fmt.Printf("\033[32mVous êtes ressuscité avec \033[0m\033[33m%d/%d PV.\n\033[0m", c.HP, c.HPMax)
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
		fmt.Println("\n\033[33m===== INVENTAIRE =====\033[0m")
		OpenInventory(*c)
		fmt.Println("\n\033[36m1) Boire une RedBull (+20 PV)\033[0m")
		fmt.Println("\033[36m2) Utiliser une Potion de poison (10 dmg/s ×3)\033[0m")
		fmt.Println("\033[36m3) Utiliser 'Livre de Sort : Mur de vent'\033[0m")
		fmt.Println("\033[31m9) Retour\033[0m")
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
			fmt.Println("\033[31mChoix invalide.\033[0m")
		}
	}
}

// ==== Marchand ====
// RedBull gratuite une fois
var redbullAvailable = true

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\033[33m\n=== MARCHAND ===\033[0m")
		if redbullAvailable {
			fmt.Println("\033[36m1) RedBull — \033[0m\033[32mGRATUIT\033[0m")
		} else {
			fmt.Println("\033[31m1) RedBull — (ÉPUISÉ)\033[0m")
		}
		fmt.Println("\033[36m2) Potion de poison — \033[0m\033[32mGRATUIT (temporaire)\033[0m")
		fmt.Println("\033[36m3) Livre de Sort : Mur de vent — \033[0m\033[33m0 or \033[O\033[32m(GRATUIT)\033[0m")
		fmt.Println("\033[31m9) Retour\033[0m")

		switch readChoice(r) {
		case "1":
			if redbullAvailable {
				addInventory(c, "RedBull", 1)
				redbullAvailable = false
				fmt.Printf("\033[32mAchat effectué ! Vous avez obtenu : RedBull (total: %d)\n\033[0m", c.Inventory["RedBull"])
			} else {
				fmt.Println("033[31mLa RedBull gratuite n’est plus disponible.\033[0m")
			}
		case "2":
			addInventory(c, "Potion de poison", 1)
			fmt.Printf("\033[32mAchat effectué ! Vous avez obtenu : Potion de poison (total: %d)\n", c.Inventory["Potion de poison \033[0m"])
		case "3":
			addInventory(c, "Livre de Sort : Mur de vent", 1)
			fmt.Println("\033[32mAchat effectué ! Vous avez obtenu : Livre de Sort : Mur de vent\033[0m")
		case "9", "retour", "back":
			return
		default:
			fmt.Println("\033[31mChoix invalide.\033[0m")
		}
	}
}

// ==== Menu principal ====

func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("\033[36m1) Afficher les informations du personnage \u001B[0m")
	fmt.Println("\033[36m2) Accéder au contenu de l’inventaire \u001B[0m")
	fmt.Println("\033[36m3) Marchand \u001B[0m")
	fmt.Println("\033[31m4) Quitter \u001B[0m")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)
	case "2", "inventaire":
		inventoryMenu(c, r)
	case "3", "marchand", "shop":
		merchantMenu(c, r)
	case "4", "q", "quit", "quitter":
		fmt.Println("\037[0mAu revoir !\033[0m")
		return false
	default:
		fmt.Println("\033[31mChoix invalide.\033[0m")
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
