package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Structure Monster
type Monster struct {
	Name        string
	MaxHP       int
	CurrentHP   int
	AttackPower int
}

// Attaque normale
func (m *Monster) Attack(target *Monster) {
	damage := rand.Intn(m.AttackPower) + 1
	if rand.Intn(100) < 10 { // 10% de chance de critique
		damage *= 2
		fmt.Println("Coup critique !")
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	fmt.Printf("%s attaque %s et inflige %d points de dégâts !\n", m.Name, target.Name, damage)
	fmt.Printf("%s a maintenant %d/%d PV\n\n", target.Name, target.CurrentHP, target.MaxHP)
}

// Attaque spéciale
func (m *Monster) SpecialAttack(target *Monster) {
	if rand.Intn(100) < 25 { // 25% de chance de rater
		fmt.Println(m.Name, "a raté son attaque spéciale !")
		return
	}
	damage := rand.Intn(m.AttackPower*2) + m.AttackPower
	if rand.Intn(100) < 15 { // 15% de critique sur attaque spéciale
		damage *= 2
		fmt.Println("Coup critique sur l'attaque spéciale !")
	}
	target.CurrentHP -= damage
	if target.CurrentHP < 0 {
		target.CurrentHP = 0
	}
	fmt.Printf("%s utilise une attaque spéciale sur %s et inflige %d points de dégâts !\n", m.Name, target.Name, damage)
	fmt.Printf("%s a maintenant %d/%d PV\n\n", target.Name, target.CurrentHP, target.MaxHP)
}

// Soigne le monstre
func (m *Monster) Heal(amount int) {
	m.CurrentHP += amount
	if m.CurrentHP > m.MaxHP {
		m.CurrentHP = m.MaxHP
	}
	fmt.Printf("%s se soigne de %d PV !\n", m.Name, amount)
	fmt.Printf("%s a maintenant %d/%d PV\n\n", m.Name, m.CurrentHP, m.MaxHP)
}

func entrainement() {
	rand.Seed(time.Now().UnixNano())

	player := Monster{Name: "Joueur", MaxHP: 100, CurrentHP: 100, AttackPower: 20}
	monster := Monster{Name: "Gobelin", MaxHP: 50, CurrentHP: 50, AttackPower: 15}

	var choice int

	for player.CurrentHP > 0 && monster.CurrentHP > 0 {
		fmt.Println("=== Tour du joueur ===")
		fmt.Println("1. Attaque normale")
		fmt.Println("2. Attaque spéciale")
		fmt.Println("3. Se soigner")
		fmt.Print("Choisissez une action : ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			player.Attack(&monster)
		case 2:
			player.SpecialAttack(&monster)
		case 3:
			player.Heal(20)
		default:
			fmt.Println("Action invalide !")
			continue
		}

		if monster.CurrentHP <= 0 {
			fmt.Println(monster.Name, "a été vaincu !")
			break
		}

		// Tour du monstre (simple AI)
		fmt.Println("=== Tour du monstre ===")
		monsterChoice := rand.Intn(3) + 1
		switch monsterChoice {
		case 1:
			monster.Attack(&player)
		case 2:
			monster.SpecialAttack(&player)
		case 3:
			monster.Heal(10)
		}

		if player.CurrentHP <= 0 {
			fmt.Println(player.Name, "a été vaincu !")
			break
		}
	}
}
