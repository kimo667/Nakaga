package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"time"
)

type Boss struct {
	Name        string
	MaxHP       int
	HP          int
	AttackPower int
}

func startBossFinal(c *Character, r *bufio.Reader) {
	if !allMissionsCompleted(c) {
		fmt.Println(CRed + "Le boss est verrouillé. Termine d'abord toutes les missions." + CReset)
		return
	}

	fmt.Println(CCyan + "Le voile se déchire... Ton rival, « Frère d'Acier », s'avance." + CReset)
	fmt.Println("Son regard rappelle celui d'un certain épéiste d'outre-monde...")

	boss := Boss{Name: "Frère d'Acier", MaxHP: 180, HP: 180, AttackPower: 20}
	playerHP := clamp(c.HP, 1, c.HPMax)

	rand.Seed(time.Now().UnixNano())
	round := 1
	for boss.HP > 0 && playerHP > 0 {
		fmt.Printf("\n--- Tour %d ---\n", round)
		fmt.Printf("%s: %d/%d PV | %s: %d/%d PV  | Mana: %d/%d\n", c.Name, playerHP, c.HPMax, boss.Name, boss.HP, boss.MaxHP, c.Mana, c.ManaMax)
		fmt.Println("1) Attaquer  2) Mur de vent (si appris)  3) Utiliser RedBull  4) Sorts  5) Fuir")

		choice := readChoice(r)
		skipBossTurn := false

		switch choice {
		case "1":
			dmg := rand.Intn(15) + 8
			if c.Equipment.Torso == "Tunique de l'aventurier" {
				dmg += 3
			}
			boss.HP -= dmg
			if boss.HP < 0 {
				boss.HP = 0
			}
			fmt.Printf("Vous frappez : -%d PV à %s.\n", dmg, boss.Name)

		case "2":
			if hasSkill(*c, "Mur de vent") {
				block := rand.Intn(10) + 10
				fmt.Printf("Vous érigez un Mur de vent ! Vous bloquerez jusqu'à %d dégâts ce tour.\n", block)
				dmg := rand.Intn(boss.AttackPower) + 8
				if dmg <= block {
					fmt.Println("Le Mur de vent dissipe l'assaut ! 0 dégât.")
					dmg = 0
				} else {
					dmg -= block
					fmt.Printf("Le Mur réduit les dégâts à %d.\n", dmg)
				}
				playerHP -= dmg
				if playerHP < 0 {
					playerHP = 0
				}
				skipBossTurn = true
			} else {
				fmt.Println(CRed + "Vous ne connaissez pas « Mur de vent »." + CReset)
				continue
			}

		case "3":
			if c.Inventory["RedBull"] > 0 {
				takeRedBull(c)
				playerHP = clamp(playerHP+HealRedBull, 1, c.HPMax)
			} else {
				fmt.Println(CRed + "Aucune RedBull disponible." + CReset)
				continue
			}

		case "4":
			fmt.Println("Sorts: 1) Coup de poing (8 dmg, 2 mana)  2) Boule de feu (18 dmg, 5 mana)  9) Retour")
			s := readChoice(r)
			if s == "1" {
				if c.Mana < 2 {
					fmt.Println(CRed + "Mana insuffisant." + CReset)
					continue
				}
				c.Mana -= 2
				boss.HP -= 8
				if boss.HP < 0 {
					boss.HP = 0
				}
				fmt.Printf("Coup de poing ! -8 PV à %s (Mana %d/%d)\n", boss.Name, c.Mana, c.ManaMax)
			} else if s == "2" {
				if c.Mana < 5 {
					fmt.Println(CRed + "Mana insuffisant." + CReset)
					continue
				}
				c.Mana -= 5
				boss.HP -= 18
				if boss.HP < 0 {
					boss.HP = 0
				}
				fmt.Printf("Boule de feu ! -18 PV à %s (Mana %d/%d)\n", boss.Name, c.Mana, c.ManaMax)
			} else {
				continue
			}

		case "5":
			fmt.Println("Vous reculez... le duel reste inachevé.")
			return

		default:
			fmt.Println(CRed + "Choix invalide." + CReset)
			continue
		}

		if boss.HP <= 0 {
			fmt.Println(CGreen + "Victoire ! Le Frère d'Acier s'incline. Tu gagnes 100 or et 60 XP." + CReset)
			c.Gold += 100
			addXP(c, 60) // Mission 2
			learnSkill(c, "Souffle de l'Acier")
			fmt.Println(CGreen + "Compétence signature apprise : Souffle de l'Acier." + CReset)
			fmt.Println(CCyan + "=== FIN ===" + CReset)
			return
		}

		if !skipBossTurn {
			dmg := rand.Intn(boss.AttackPower) + 8
			if rand.Intn(100) < 15 {
				dmg += 10
				fmt.Println("Le boss déclenche une coupe transversale fulgurante !")
			}
			playerHP -= dmg
			if playerHP < 0 {
				playerHP = 0
			}
			fmt.Printf("%s riposte : -%d PV.\n", boss.Name, dmg)
		}

		if playerHP <= 0 {
			fmt.Println(CRed + "Tu t'effondres... mais l'histoire ne s'arrête pas là." + CReset)
			isDead(c)
			return
		}

		round++
	}
}
