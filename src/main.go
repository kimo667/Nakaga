package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

/* ================== Couleurs ================== */

const (
	CReset  = "\033[0m"
	CRed    = "\033[31m"
	CGreen  = "\033[32m"
	CYellow = "\033[33m"
	CCyan   = "\033[36m"
)

/* ================== Constantes ================== */

const (
	InventoryCapacity  = 10
	HealRedBull        = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
)

/* ================== Types ================== */

type Classe string

const (
	ClasseHumain Classe = "Humain"
	ClasseElfe   Classe = "Elfe"
	ClasseNain   Classe = "Nain"
	// (On conserve ton univers samouraï, mais pour T11 on propose ces 3 classes.)
)

type Character struct {
	Name      string
	Class     Classe
	Level     int
	HPMax     int
	HP        int
	Inventory map[string]int
	Skills    []string
	Gold      int // juste initialisé/affiché (T13 sera pour l’économie)
}

/* ================== Utils ================== */

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

/* =========== Inventaire (capacité 10 items) =========== */

func totalItems(c Character) int {
	sum := 0
	for _, q := range c.Inventory {
		if q > 0 {
			sum += q
		}
	}
	return sum
}

func addInventory(c *Character, item string, qty int) int {
	if qty <= 0 {
		return 0
	}
	if c.Inventory == nil {
		c.Inventory = make(map[string]int)
	}
	used := totalItems(*c)
	free := InventoryCapacity - used
	if free <= 0 {
		fmt.Printf(CRed+"Inventaire plein (%d/%d) — impossible d'ajouter %s."+CReset+"\n", used, InventoryCapacity, item)
		return 0
	}
	add := qty
	if add > free {
		add = free
	}
	c.Inventory[item] += add
	if add < qty {
		fmt.Printf(CYellow+"Ajout partiel: +%d %s (cap %d/%d)."+CReset+"\n", add, item, used+add, InventoryCapacity)
	} else {
		fmt.Printf(CGreen+"Ajouté: +%d %s (cap %d/%d)."+CReset+"\n", add, item, used+add, InventoryCapacity)
	}
	return add
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

/* ================== Skills ================== */

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

/* =============== TÂCHE 11 : characterCreation =============== */

func isLetters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && r != '-' && r != ' ' {
			return false
		}
	}
	return len(strings.TrimSpace(s)) > 0
}

func normalizeName(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.ToLower(raw)
	// Met majuscule sur première lettre de chaque mot
	parts := strings.Fields(raw)
	for i, p := range parts {
		runes := []rune(p)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
			parts[i] = string(runes)
		}
	}
	return strings.Join(parts, " ")
}

func askClass(r *bufio.Reader) Classe {
	for {
		fmt.Println("Choisis ta classe :")
		fmt.Println("1) Humain (100 PV max)")
		fmt.Println("2) Elfe   (80  PV max)")
		fmt.Println("3) Nain   (120 PV max)")
		ch := readChoice(r)
		switch ch {
		case "1", "humain":
			return ClasseHumain
		case "2", "elfe":
			return ClasseElfe
		case "3", "nain":
			return ClasseNain
		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

func baseStatsForClass(cl Classe) (hpMax int) {
	switch cl {
	case ClasseHumain:
		return 100
	case ClasseElfe:
		return 80
	case ClasseNain:
		return 120
	default:
		return 100
	}
}

func characterCreation(r *bufio.Reader) Character {
	// Nom (lettres uniquement) + normalisation
	var name string
	for {
		raw := readLine(r, "Entre le nom de ton personnage (lettres, espaces, tirets) : ")
		if !isLetters(raw) {
			fmt.Println(CRed + "Nom invalide. Lettres/espaces/tirets uniquement." + CReset)
			continue
		}
		name = normalizeName(raw)
		if name != "" {
			break
		}
	}

	// Classe + PV selon classe (PV actuels = 50% PV max), lvl=1, skill de base
	cl := askClass(r)
	hpMax := baseStatsForClass(cl)
	hp := hpMax / 2

	ch := Character{
		Name:      name,
		Class:     cl,
		Level:     1,
		HPMax:     hpMax,
		HP:        hp,
		Inventory: map[string]int{},
		Skills:    []string{},
		Gold:      100, // *** juste initialisé (pas d’économie ici) ***
	}
	learnSkill(&ch, "Coup de Poing") // base demandé par T11

	// on te met 3 RedBull pour le côté humour :)
	addInventory(&ch, "RedBull", 3)

	return ch
}

/* ================== Affichages ================== */

func displayInventory(c Character) {
	keys := make([]string, 0, len(c.Inventory))
	for k, v := range c.Inventory {
		if v > 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	fmt.Printf(CYellow+"Inventaire (%d/%d) :"+CReset+"\n", totalItems(c), InventoryCapacity)
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

/* ================== Consommables ================== */

func takeRedBull(c *Character) {
	if !removeInventory(c, "RedBull", 1) {
		fmt.Println(CRed + "Pas de RedBull dans l'inventaire !" + CReset)
		return
	}
	before := c.HP
	c.HP = clamp(c.HP+HealRedBull, 0, c.HPMax)
	fmt.Printf(CCyan+"Tu as bu une RedBull !"+CReset+" PV: %d → "+CGreen+"%d/%d"+CReset+"\n", before, c.HP, c.HPMax)
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

/* =========== Livre de sort (version actuelle « Mur de vent ») =========== */

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

/* ================== Mort / revive (T8) ================== */

func isDead(c *Character) bool {
	if c.HP <= 0 {
		fmt.Println(CRed + "\n*** WASTED ***" + CReset)
		c.HP = clamp(c.HPMax/2, 1, c.HPMax)
		fmt.Printf(CGreen+"Vous êtes ressuscité avec %d/%d PV."+CReset+"\n", c.HP, c.HPMax)
		return true
	}
	return false
}

/* ================== Menus ================== */

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

// Marchand version « gratuit » (T7/T9), pas d’économie ici
var redbullFreeOnce = true
var poisonFreeOnce = true
var windBookFreeOnce = true

func merchantMenu(c *Character, r *bufio.Reader) {
	for {
		fmt.Println("\n" + CYellow + "=== MARCHAND ===" + CReset)
		if redbullFreeOnce {
			fmt.Println("1) RedBull — " + CGreen + "GRATUIT" + CReset)
		} else {
			fmt.Println("1) RedBull — (ÉPUISÉ)")
		}
		if poisonFreeOnce {
			fmt.Println("2) Potion de poison — " + CGreen + "GRATUIT" + CReset)
		} else {
			fmt.Println("2) Potion de poison — (ÉPUISÉ)")
		}
		if windBookFreeOnce {
			fmt.Println("3) Livre de Sort : Mur de vent — " + CGreen + "GRATUIT" + CReset)
		} else {
			fmt.Println("3) Livre de Sort : Mur de vent — (ÉPUISÉ)")
		}
		fmt.Println("9) Retour")

		switch readChoice(r) {
		case "1":
			if redbullFreeOnce {
				if addInventory(c, "RedBull", 1) > 0 {
					redbullFreeOnce = false
				}
			} else {
				fmt.Println(CRed + "Plus de RedBull gratuite." + CReset)
			}
		case "2":
			if poisonFreeOnce {
				if addInventory(c, "Potion de poison", 1) > 0 {
					poisonFreeOnce = false
				}
			} else {
				fmt.Println(CRed + "Plus de Potion de poison gratuite." + CReset)
			}
		case "3":
			if windBookFreeOnce {
				if addInventory(c, "Livre de Sort : Mur de vent", 1) > 0 {
					windBookFreeOnce = false
				}
			} else {
				fmt.Println(CRed + "Livre déjà pris." + CReset)
			}
		case "9", "retour", "back":
			return
		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
		}
	}
}

func mainMenu(c *Character, r *bufio.Reader) bool {
	fmt.Println("\n\033[1;33m===== MENU =====\033[0m")
	fmt.Println("1) Afficher les informations du personnage")
	fmt.Println("2) Accéder au contenu de l’inventaire")
	fmt.Println("3) Marchand (gratuit)")
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
		fmt.Println(CRed + "Choix invalide." + CReset)
	}
	return true
}

/* ================== main ================== */

func main() {
	reader := bufio.NewReader(os.Stdin)
	// TÂCHE 11 remplace l’init en dur dans main:
	c := characterCreation(reader)

	for mainMenu(&c, reader) {
		isDead(&c) // revive auto si besoin
	}
}
