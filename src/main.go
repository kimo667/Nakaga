package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

/* ====== Couleurs ====== */

const (
	CReset  = "\033[0m"
	CRed    = "\033[31m"
	CGreen  = "\033[32m"
	CYellow = "\033[33m"
	CCyan   = "\033[36m"
)

/* ====== Capacité & Effets ====== */

const (
	BaseInventoryCap     = 10 // capacité de base
	InventoryUpgradeStep = 10 // +10 par upgrade
	MaxInventoryUpgrades = 3  // max 3 upgrades

	HealRedBull        = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
)

/* ====== Types ====== */

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

type Character struct {
	Name        string
	Class       Classe
	Level       int
	HPMax       int
	HP          int
	Inventory   map[string]int
	Skills      []string
	Gold        int
	CapMax      int
	InvUpgrades int
}

/* ====== Utils ====== */

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func readLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func readChoice(r *bufio.Reader) string {
	return strings.ToLower(strings.TrimSpace(readLine(r, "> ")))
}

/* ====== Inventaire (capacité & upgrades) ====== */

func totalItems(c Character) int {
	sum := 0
	for _, q := range c.Inventory {
		if q > 0 {
			sum += q
		}
	}
	return sum
}

func canCarry(c Character, qty int) bool {
	return totalItems(c)+qty <= c.CapMax
}

func addInventory(c *Character, item string, qty int) bool {
	if qty <= 0 {
		return false
	}
	if c.Inventory == nil {
		c.Inventory = make(map[string]int)
	}
	if !canCarry(*c, qty) {
		fmt.Printf(CRed+"Inventaire plein (%d/%d). Impossible d'ajouter %d x %s."+CReset+"\n",
			totalItems(*c), c.CapMax, qty, item)
		return false
	}
	c.Inventory[item] += qty
	return true
}

func removeInventory(c *Character, item string, qty int) bool {
	cur, ok := c.Inventory[item]
	if !ok || qty <= 0 || cur < qty {
		return false
	}
	if cur-qty == 0 {
		delete(c.Inventory, item)
	} else {
		c.Inventory[item] = cur - qty
	}
	return true
}

func upgradeInventorySlot(c *Character) bool {
	if c.InvUpgrades >= MaxInventoryUpgrades {
		fmt.Printf(CYellow+"Capacité déjà au maximum (%d/%d upgrades)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
		return false
	}
	c.InvUpgrades++
	c.CapMax += InventoryUpgradeStep
	fmt.Printf(CGreen+"Capacité augmentée ! Nouvelle capacité : %d (améliorations : %d/%d)"+CReset+"\n",
		c.CapMax, c.InvUpgrades, MaxInventoryUpgrades)
	return true
}

/* ====== Skills ====== */

func hasSkill(c Character, s string) bool {
	for _, k := range c.Skills {
		if k == s {
			return true
		}
	}
	return false
}

func learnSkill(c *Character, s string) bool {
	if hasSkill(*c, s) {
		return false
	}
	c.Skills = append(c.Skills, s)
	sort.Strings(c.Skills)
	return true
}

/* ====== Initialisation ====== */

func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	ch := Character{
		Name:        name,
		Class:       class,
		Level:       level,
		HPMax:       hpMax,
		HP:          clamp(hp, 0, hpMax),
		Inventory:   map[string]int{},
		Skills:      []string{},
		Gold:        100,              // or de départ
		CapMax:      BaseInventoryCap, // capacité de base
		InvUpgrades: 0,
	}
	// copie inventaire de départ en respectant la capacité
	for k, v := range inv {
		if v <= 0 {
			continue
		}
		if !addInventory(&ch, k, v) {
			break
		}
	}
	// technique de base
	learnSkill(&ch, "Tempête d'acier")
	return ch
}

/* ====== Création interactive du personnage ====== */

func chooseClass(r *bufio.Reader) Classe {
	for {
		fmt.Println(CYellow + "Choisis ta classe :" + CReset)
		fmt.Println("1) Humain  – équilibré")
		fmt.Println("2) Samouraï – PV élevés")
		fmt.Println("3) Ninja    – agile")
		choice := readChoice(r)
		switch choice {
		case "1", "humain":
			return ClasseHumain
		case "2", "samourai", "samouraï", "samurai":
			return ClasseSamurai
		case "3", "ninja":
			return ClasseNinja
		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

func createCharacterInteractive(r *bufio.Reader) Character {
	fmt.Println(CYellow + "=== Création de personnage ===" + CReset)
	name := readLine(r, "Entre ton nom: ")
	if name == "" {
		name = "Yazuo"
	}
	class := chooseClass(r)

	// petits bonus par classe (optionnel, simple)
	hpMax := 100
	switch class {
	case ClasseSamurai:
		hpMax = 110
	case ClasseNinja:
		hpMax = 90
	default:
		hpMax = 100
	}

	// HP de départ 40% des PV max (comme avant ~40/100)
	startHP := hpMax * 40 / 100

	// inventaire de départ (humour 3 RedBull)
	startInv := map[string]int{
		"RedBull": 3,
	}

	ch := initCharacter(name, class, 1, hpMax, startHP, startInv)
	fmt.Println(CGreen + "Personnage créé !" + CReset)
	return ch
}

/* ====== Affichages ====== */

func displayInventory(c Character) {
	keys := make([]string, 0, len(c.Inventory))
	for k, v := range c.Inventory {
		if v > 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	fmt.Printf(CYellow+"Inventaire (%d/%d) :"+CReset+"\n", totalItems(c), c.CapMax)
	if len(keys) == 0 {
		fmt.Println("  (vide)")
		return
	}
	for _, k := range keys {
		fmt.Printf("  - %s x%d\n", k, c.Inventory[k])
	}
}

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
	fmt.Println(CYellow + "=== Informations du personnage ===" + CReset)
	fmt.Printf(CYellow+"Nom   : "+CReset+"%s\n", c.Name)
	fmt.Printf(CYellow+"Classe: "+CReset+"%s\n", c.Class)
	fmt.Printf(CYellow+"Niveau: "+CReset+"%d\n", c.Level)
	fmt.Printf(CGreen+"PV    : "+CReset+CGreen+"%d/%d"+CReset+"\n", c.HP, c.HPMax)
	fmt.Printf(CYellow+"Or    : "+CReset+"%d\n", c.Gold)
	fmt.Printf(CYellow+"Capacité: "+CReset+"%d/%d (améliorations: %d/%d)\n",
		totalItems(c), c.CapMax, c.InvUpgrades, MaxInventoryUpgrades)

	fmt.Println(CYellow + "Compétences :" + CReset)
	if len(c.Skills) == 0 {
		fmt.Println("  (aucune technique)")
	} else {
		for _, s := range c.Skills {
			fmt.Println("  - " + s)
		}
	}

	displayInventory(c)
}

/* ====== Consommables & Effets ====== */

func takeRedBull(c *Character) {
	if !removeInventory(c, "RedBull", 1) {
		fmt.Println(CRed + "Pas de RedBull dans l'inventaire !" + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+HealRedBull, 0, c.HPMax)
	fmt.Printf(CCyan+"Tu as bu une RedBull !"+CReset+" PV: %d → "+CGreen+"%d/%d"+CReset+"\n", before, c.HP, c.HPMax)
}

func usePotionVie(c *Character) {
	if !removeInventory(c, "Potion de vie", 1) {
		fmt.Println(CRed + "Pas de Potion de vie." + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+20, 0, c.HPMax)
	fmt.Printf(CCyan+"Vous buvez une Potion de vie."+CReset+" PV: %d → "+CGreen+"%d/%d"+CReset+"\n", before, c.HP, c.HPMax)
}

func poisonPot(c *Character) {
	if !removeInventory(c, "Potion de poison", 1) {
		fmt.Println(CRed + "Vous n'avez pas de Potion de poison." + CReset)
		return
	}
	fmt.Println(CCyan + "Vous utilisez une Potion de poison…" + CReset)
	for i := 1; i <= PoisonTicks; i++ {
		before := c.HP
		c.HP = clamp(c.HP-PoisonDamagePerSec, 0, c.HPMax)
		fmt.Printf("Effet poison %d/%d → PV: %d → %d/%d\n", i, PoisonTicks, before, c.HP, c.HPMax)
		if isDead(c) { // revive + stop poison
			fmt.Println(CRed + "L'effet du poison est interrompu suite à votre mort." + CReset)
			return
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf(CCyan+"L'effet du poison est terminé. PV restants : %d/%d"+CReset+"\n", c.HP, c.HPMax)
}

func useSpellBookWind(c *Character) {
	if c.Inventory["Livre de Sort : Mur de vent"] <= 0 {
		fmt.Println(CRed + "Vous n'avez pas de 'Livre de Sort : Mur de vent'." + CReset)
		return
	}
	if hasSkill(*c, "Mur de vent") {
		fmt.Println(CRed + "Vous connaissez déjà 'Mur de vent'. Le livre n'a pas été consommé." + CReset)
		return
	}
	removeInventory(c, "Livre de Sort : Mur de vent", 1)
	if learnSkill(c, "Mur de vent") {
		fmt.Println(CGreen + "Vous avez appris : Mur de vent !" + CReset)
	}
}

/* ====== Mort / revive (T8) ====== */

func isDead(c *Character) bool {
	if c.HP <= 0 {
		fmt.Println(CRed + "\n*** WASTED ***" + CReset)
		c.HP = clamp(c.HPMax/2, 1, c.HPMax)
		fmt.Printf(CGreen+"Vous êtes ressuscité avec %d/%d PV."+CReset+"\n", c.HP, c.HPMax)
		return true
	}
	return false
}

/* ====== Marchand (économie + upgrades) ====== */

func canAfford(c Character, price int) bool { return c.Gold >= price }

func buyItem(c *Character, item string, price, qty int) bool {
	if price < 0 || qty <= 0 {
		return false
	}
	// vérifier la place avant de payer
	if !canCarry(*c, qty) {
		fmt.Printf(CRed+"Inventaire plein (%d/%d). Libérez de la place avant d’acheter %s."+CReset+"\n",
			totalItems(*c), c.CapMax, item)
		return false
	}
	if !canAfford(*c, price) {
		fmt.Printf(CRed+"Or insuffisant pour %s (coût %d). Solde: %d"+CReset+"\n", item, price, c.Gold)
		return false
	}
	c.Gold -= price
	if !addInventory(c, item, qty) {
		// sécurité
		c.Gold += price
		return false
	}
	fmt.Printf(CGreen+"Achat effectué ! %s x%d (−%d or). Or restant: %d"+CReset+"\n", item, qty, price, c.Gold)
	return true
}

var redbullFreeOnce = true
var windBookStock = 1 // stock unique du livre

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n" + CYellow + "=== MARCHAND ===" + CReset)
		if redbullFreeOnce {
			fmt.Println("1) RedBull — " + CGreen + "GRATUIT" + CReset)
		} else {
			fmt.Println("1) RedBull — (ÉPUISÉ)")
		}
		// Tarifs (version pote)
		fmt.Println("2) Potion de vie — 3 or")
		fmt.Println("3) Potion de poison — 6 or")
		if windBookStock > 0 {
			fmt.Println("4) Livre de Sort : Mur de vent — 25 or")
		} else {
			fmt.Println("4) Livre de Sort : Mur de vent — (ÉPUISÉ)")
		}
		// Matériaux
		fmt.Println("5) Fourrure de Loup — 4 or")
		fmt.Println("6) Peau de Troll — 7 or")
		fmt.Println("7) Cuir de Sanglier — 3 or")
		fmt.Println("8) Plume de Corbeau — 1 or")
		// Upgrade capacité
		fmt.Printf("10) Augmentation d’inventaire — 30 or (utilisée %d/%d)\n", c.InvUpgrades, MaxInventoryUpgrades)

		fmt.Println("9) Retour")

		switch readChoice(r) {
		case "1": // RedBull gratuite 1 fois
			if redbullFreeOnce {
				if addInventory(c, "RedBull", 1) {
					redbullFreeOnce = false
					fmt.Printf(CGreen+"RedBull reçue ! (total: %d)"+CReset+"\n", c.Inventory["RedBull"])
				}
			} else {
				fmt.Println(CRed + "La RedBull gratuite n’est plus disponible." + CReset)
			}

		case "2":
			_ = buyItem(c, "Potion de vie", 3, 1)

		case "3":
			_ = buyItem(c, "Potion de poison", 6, 1)

		case "4":
			if windBookStock <= 0 {
				fmt.Println(CRed + "Le Livre de Sort : Mur de vent n’est plus disponible." + CReset)
				break
			}
			if buyItem(c, "Livre de Sort : Mur de vent", 25, 1) {
				windBookStock--
			}

		case "5":
			_ = buyItem(c, "Fourrure de Loup", 4, 1)
		case "6":
			_ = buyItem(c, "Peau de Troll", 7, 1)
		case "7":
			_ = buyItem(c, "Cuir de Sanglier", 3, 1)
		case "8":
			_ = buyItem(c, "Plume de Corbeau", 1, 1)

		case "10":
			if c.InvUpgrades >= MaxInventoryUpgrades {
				fmt.Printf(CYellow+"Limite d’améliorations atteinte (%d/%d)."+CReset+"\n", c.InvUpgrades, MaxInventoryUpgrades)
				break
			}
			if !canAfford(*c, 30) {
				fmt.Printf(CRed+"Or insuffisant pour l’amélioration (coût 30). Solde: %d"+CReset+"\n", c.Gold)
				break
			}
			c.Gold -= 30
			if !upgradeInventorySlot(c) {
				c.Gold += 30
			}

		case "9", "retour", "back":
			return

		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

/* ====== Menus ====== */

func inventoryMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n" + CYellow + "===== INVENTAIRE =====" + CReset)
		displayInventory(*c)

		type opt struct {
			key, label string
			run        func()
		}
		opts := []opt{}
		add := func(label string, fn func()) {
			opts = append(opts, opt{key: fmt.Sprintf("%d", len(opts)+1), label: label, run: fn})
		}

		if c.Inventory["RedBull"] > 0 {
			add("Boire une RedBull (+20 PV)", func() { takeRedBull(c); isDead(c) })
		}
		if c.Inventory["Potion de vie"] > 0 {
			add("Boire une Potion de vie (+20 PV)", func() { usePotionVie(c) })
		}
		if c.Inventory["Potion de poison"] > 0 {
			add("Utiliser une Potion de poison (10 dmg/s ×3)", func() { poisonPot(c) })
		}
		if c.Inventory["Livre de Sort : Mur de vent"] > 0 {
			add("Utiliser 'Livre de Sort : Mur de vent'", func() { useSpellBookWind(c) })
		}

		if len(opts) == 0 {
			fmt.Println("(Aucune action disponible)")
		} else {
			for _, o := range opts {
				fmt.Printf("%s) %s\n", o.key, o.label)
			}
		}
		fmt.Println("9) Retour")

		ch := readChoice(r)
		if ch == "9" || ch == "retour" || ch == "back" {
			return
		}
		ok := false
		for _, o := range opts {
			if ch == o.key {
				o.run()
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand")
	fmt.Println("4) Forgeron")
	fmt.Println("4) Quitter")

	switch readChoice(r) {
	case "1", "infos", "information":
		displayInfo(*c)
	case "2", "inventaire":
		inventoryMenu(c, r)
	case "3", "marchand", "shop":
		merchantMenu(c, r)
	case "4", "forgeron", "forge":
		blacksmithMenu(c, r) //
	case "5", "q", "quit", "quitter":
		fmt.Println("Au revoir !")
		return false
	default:
		fmt.Println(CRed + "Choix invalide." + CReset)
	}
	return true
}

/* ====== main ====== */

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Création interactive (nom + classe)
	c := createCharacterInteractive(reader)

	for mainMenu(&c, reader) {
		isDead(&c) // revive auto si besoin
	}
}

// ==== Forgeron ====

// Liste des recettes : item → ressources nécessaires
var forgeRecipes = map[string]map[string]int{
	"Chapeau de l’aventurier": {
		"Plume de Corbeau": 1,
		"Cuir de Sanglier": 1,
	},
	"Tunique de l’aventurier": {
		"Fourrure de loup": 2,
		"Peau de Troll":    1,
	},
	"Bottes de l’aventurier": {
		"Fourrure de loup": 1,
		"Cuir de Sanglier": 1,
	},
}

func hasResources(c *Character, recipe map[string]int) bool {
	for item, qty := range recipe {
		if c.Inventory[item] < qty {
			return false
		}
	}
	return true
}

// Consomme les ressources du joueur
func consumeResources(c *Character, recipe map[string]int) {
	for item, qty := range recipe {
		removeInventory(c, item, qty)
	}
}
func blacksmithMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n=== FORGERON ===")
		fmt.Printf("Votre or : %d\n", c.Gold)
		fmt.Println("Que voulez-vous fabriquer ? (5 or par objet)")

		idx := 1
		itemList := []string{}
		for item := range forgeRecipes {
			fmt.Printf("%d) %s\n", idx, item)
			itemList = append(itemList, item)
			idx++
		}
		fmt.Println("9) Retour")

		choice := readChoice(r)
		if choice == "9" || choice == "retour" || choice == "back" {
			return
		}

		var selected int
		fmt.Sscanf(choice, "%d", &selected)

		if selected >= 1 && selected <= len(itemList) {
			item := itemList[selected-1]
			recipe := forgeRecipes[item]

			// Vérification or + ressources
			if c.Gold < 5 {
				fmt.Println("Pas assez d’or pour fabriquer cet objet !")
				continue
			}
			if !hasResources(c, recipe) {
				fmt.Println("Ressources insuffisantes pour fabriquer :", item)
				continue
			}

			// Consommer or + ressources
			c.Gold -= 5
			consumeResources(c, recipe)

			// Ajouter l'objet crafté
			addInventory(c, item, 1)
			fmt.Printf("Vous avez fabriqué : %s (reste %d or).\n", item, c.Gold)
		} else {
			fmt.Println("Choix invalide.")
		}
	}
}
