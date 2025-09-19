package main

import (
	"bufio"
	"fmt"
)

func defaultMissions() []Mission {
	return []Mission{
		{
			ID:                 1,
			Title:              "Initiative",
			Description:        "Ajoute et utilise l'initiative en combat (gagne 1 combat d'entraînement).",
			State:              NotStarted,
			RequiredTrainKills: 1,
			RewardGold:         15,
			RewardItem:         "RedBull",
		},
		{
			ID:          2,
			Title:       "Expérience",
			Description: "Système d’XP : monte d’un niveau.",
			State:       NotStarted,
			RewardGold:  20,
			RewardItem:  "Potion de vie",
		},
		{
			ID:          3,
			Title:       "Combat Magique",
			Description: "Utilise un sort en combat (coup de poing ou boule de feu).",
			State:       NotStarted,
			RewardGold:  10,
		},
		{
			ID:          4,
			Title:       "Ressource de mana",
			Description: "Utilise une potion de mana ou lance un sort (coûte du mana).",
			State:       NotStarted,
			RewardGold:  10,
			RewardItem:  "Potion de mana",
		},
		{
			ID:                5,
			Title:             "Améliorer le jeu",
			Description:       "Atteins le niveau 2 (récompenses de boss/entraînement donnent de l’XP).",
			State:             NotStarted,
			RewardGold:        25,
			RewardUpgradeSlot: true,
		},
		{
			ID:          6,
			Title:       "Qui sont-ils ?",
			Description: "Nouveau menu dans Missions → « Qui sont-ils ? »",
			State:       NotStarted,
			RewardGold:  5,
		},
	}
}

func ensureMissions(c *Character) {
	if len(c.Missions) == 0 {
		c.Missions = defaultMissions()
	}
}

func allMissionsCompleted(c *Character) bool {
	ensureMissions(c)
	for _, m := range c.Missions {
		if m.State != Completed {
			return false
		}
	}
	return true
}

func applyMissionRewards(c *Character, m *Mission) {
	if m.RewardGold > 0 {
		c.Gold += m.RewardGold
		fmt.Printf(CGreen+"Récompense : +%d or (total %d)."+CReset+"\n", m.RewardGold, c.Gold)
	}
	if m.RewardItem != "" {
		addInventory(c, m.RewardItem, 1)
		fmt.Printf(CGreen+"Récompense : %s x1."+CReset+"\n", m.RewardItem)
	}
	if m.RewardSkill != "" && learnSkill(c, m.RewardSkill) {
		fmt.Printf(CGreen+"Récompense : compétence '%s' apprise."+CReset+"\n", m.RewardSkill)
	}
	if m.RewardUpgradeSlot {
		_ = upgradeInventorySlot(c)
	}
}

func markCompleted(c *Character, id int) {
	for i := range c.Missions {
		if c.Missions[i].ID == id {
			if c.Missions[i].State != Completed {
				c.Missions[i].State = Completed
				applyMissionRewards(c, &c.Missions[i])
			}
			return
		}
	}
}

func missionsMenu(c *Character, r *bufio.Reader) {
	ensureMissions(c)

	// auto-checks simples
	// M2 & M5: niveau >= 2
	if c.Level >= 2 {
		markCompleted(c, 2)
		markCompleted(c, 5)
	}
	// M4: si de la mana a été dépensée
	if c.Mana < c.ManaMax {
		markCompleted(c, 4)
	}

	for {
		fmt.Println(CYellow + "\n=== MISSIONS ===" + CReset)
		for _, m := range c.Missions {
			state := map[MissionState]string{NotStarted: "Non commencée", InProgress: "En cours", Completed: "Terminée"}[m.State]
			fmt.Printf("%d) %s — [%s]\n   %s\n", m.ID, m.Title, state, m.Description)
		}
		fmt.Println("0) Qui sont-ils ?  |  9) Retour")

		choice := readChoice(r)
		if choice == "9" {
			return
		}
		if choice == "0" { // Mission 6
			fmt.Println(CCyan + "Artistes cachés :")
			fmt.Println(" - Akira Toriyama")
			fmt.Println(" - Eiichiro Oda")
			fmt.Println(" - Hayao Miyazaki" + CReset)
			markCompleted(c, 6)
			continue
		}

		var id int
		fmt.Sscanf(choice, "%d", &id)
		var found *Mission
		for i := range c.Missions {
			if c.Missions[i].ID == id {
				found = &c.Missions[i]
				break
			}
		}
		if found == nil {
			fmt.Println(CRed + "Mission inconnue." + CReset)
			continue
		}

		switch found.ID {
		case 1:
			// Victoires d'entraînement comptées via main_menu.go
			if found.TrainKills >= found.RequiredTrainKills {
				found.State = Completed
				fmt.Println(CGreen + "Mission « Initiative » terminée !" + CReset)
				applyMissionRewards(c, found)
			} else {
				found.State = InProgress
				fmt.Println(CYellow + "Gagne un combat d'entraînement (les victoires sont comptées automatiquement)." + CReset)
			}
		case 2:
			if c.Level >= 2 {
				found.State = Completed
				fmt.Println(CGreen + "Mission « Expérience » terminée !" + CReset)
				applyMissionRewards(c, found)
			} else {
				found.State = InProgress
				fmt.Println(CYellow + "Monte au niveau 2 (entraînement/boss donnent de l’XP)." + CReset)
			}
		case 3:
			// Validation manuelle simple
			fmt.Println("As-tu utilisé un sort en combat (entraînement ou boss) ? (oui/non)")
			if readChoice(r) == "oui" {
				found.State = Completed
				fmt.Println(CGreen + "Mission « Combat Magique » terminée !" + CReset)
				applyMissionRewards(c, found)
			} else {
				found.State = InProgress
			}
		case 4:
			if c.Mana < c.ManaMax || c.Inventory["Potion de mana"] > 0 {
				found.State = Completed
				fmt.Println(CGreen + "Mission « Ressource de mana » terminée !" + CReset)
				applyMissionRewards(c, found)
			} else {
				found.State = InProgress
				fmt.Println(CYellow + "Utilise une potion de mana ou lance un sort pour dépenser du mana." + CReset)
			}
		case 5:
			if c.Level >= 2 {
				found.State = Completed
				fmt.Println(CGreen + "Mission « Améliorer le jeu » terminée !" + CReset)
				applyMissionRewards(c, found)
			} else {
				found.State = InProgress
				fmt.Println(CYellow + "Atteins le niveau 2." + CReset)
			}
		case 6:
			fmt.Println("Ouvre l’option « 0) Qui sont-ils ? » ci-dessus.")
		}
	}
}
