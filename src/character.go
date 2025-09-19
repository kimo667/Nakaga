package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// ====== Initialisation du personnage ======
func initCharacter(name string, class Classe, level, hpMax, hp int, inv map[string]int) Character {
	c := Character{
		Name:      name,
		Class:     class,
		Level:     level,
		HPMax:     hpMax,
		HP:        hp,
		BaseHPMax: hpMax,

		Initiative: 0,
		XP:         0,
		XPMax:      50,
		Mana:       0,
		ManaMax:    0,

		Inventory: map[string]int{},
		Skills:    []string{},
		Gold:      100,
		CapMax:    BaseInventoryCap,
	}
	for it, q := range inv {
		_ = addInventory(&c, it, q)
	}
	_ = learnSkill(&c, "Tempête d'acier")

	// Mana de base selon classe (Mission 4)
	switch class {
	case ClasseHumain:
		c.ManaMax = 30
	case ClasseSamurai:
		c.ManaMax = 20
	case ClasseNinja:
		c.ManaMax = 40
	}
	c.Mana = c.ManaMax

	// Initiative de départ
	rand.Seed(time.Now().UnixNano())
	c.Initiative = rollInitiative(10)

	ensureMissions(&c)
	return c
}

// jet d’initiative
func rollInitiative(base int) int { return base + rand.Intn(10) }

// Gain d’XP + montée de niveau (Mission 2)
func addXP(c *Character, amount int) {
	if amount <= 0 {
		return
	}
	c.XP += amount
	lvlUp := false
	for c.XP >= c.XPMax {
		c.XP -= c.XPMax
		c.Level++
		lvlUp = true
		// boost simple
		c.BaseHPMax += 10
		recalcHPMax(c)
		c.HP = c.HPMax
		c.XPMax = int(float64(c.XPMax)*1.3 + 5)
	}
	if lvlUp {
		fmt.Printf(CGreen+"Niveau augmenté ! Vous êtes niv. %d. PV max: %d"+CReset+"\n", c.Level, c.HPMax)
	}
}

// ====== Création interactive ======
func createCharacterInteractive(r *bufio.Reader) Character {
	fmt.Println(CCyan + "=== Création du personnage ===" + CReset)
	name := readLine(r, "Ton nom, voyageur du vent ? ")
	for name == "" {
		fmt.Println(CRed + "Le nom ne doit pas être vide." + CReset)
		name = readLine(r, "Ton nom, voyageur du vent ? ")
	}
	name = strings.TrimSpace(name)

	fmt.Println("Choisis ta voie : 1) Humain  2) Samurai  3) Ninja")
	var class Classe
	switch readChoice(r) {
	case "1", "humain":
		class = ClasseHumain
	case "2", "samurai", "samouraï":
		class = ClasseSamurai
	case "3", "ninja":
		class = ClasseNinja
	default:
		class = ClasseHumain
	}

	base := 0
	switch class {
	case ClasseHumain:
		base = 100
	case ClasseSamurai:
		base = 120
	case ClasseNinja:
		base = 90
	}
	hpMax := base
	hp := base
	name = sanitizeName(name)

	inv := map[string]int{
		"Potion de vie": 1,
		"RedBull":       1,
	}

	player := initCharacter(name, class, 1, hpMax, hp, inv)
	fmt.Println(CGreen + "Personnage créé !" + CReset)
	return player
}

func sanitizeName(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' || r == '_' {
			b.WriteRune(r)
		}
	}
	out := strings.TrimSpace(b.String())
	if out == "" {
		return "Rônin"
	}
	return out
}
