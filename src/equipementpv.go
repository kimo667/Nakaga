package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Type d'équipement
type EquipmentType int

const (
	Head EquipmentType = iota
	Body
	Feet
)

// Structure d’un équipement
type Equipment struct {
	Name    string
	Type    EquipmentType
	HPBonus int
}

// Structure du joueur
type Player struct {
	MaxHP     int
	BaseHP    int
	Inventory []Equipment
	Equipment map[EquipmentType]Equipment
}

// Crée un nouveau joueur
func NewPlayer(baseHP int) *Player {
	return &Player{
		BaseHP:    baseHP,
		MaxHP:     baseHP,
		Inventory: []Equipment{},
		Equipment: make(map[EquipmentType]Equipment),
	}
}

// Ajouter un équipement à l’inventaire
func (p *Player) AddToInventory(e Equipment) {
	p.Inventory = append(p.Inventory, e)
}

// Équiper un équipement par son nom
func (p *Player) EquipByName(name string) {
	var item *Equipment
	for i, invItem := range p.Inventory {
		if strings.EqualFold(invItem.Name, name) {
			item = &p.Inventory[i]
			break
		}
	}
	if item == nil {
		fmt.Println("❌ Équipement introuvable dans l’inventaire :", name)
		return
	}
	e := *item

	// Si déjà un équipement dans cette section
	if old, exists := p.Equipment[e.Type]; exists {
		// On le remet dans l’inventaire
		p.Inventory = append(p.Inventory, old)
		fmt.Println("↩️ Vous avez retiré :", old.Name)
	}

	// On équipe le nouvel objet
	p.Equipment[e.Type] = e

	// On l’enlève de l’inventaire
	for i, invItem := range p.Inventory {
		if invItem.Name == e.Name {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			break
		}
	}

	p.UpdateMaxHP()
	fmt.Println("✅ Vous avez équipé :", e.Name)
}

// Recalcul des PV max
func (p *Player) UpdateMaxHP() {
	p.MaxHP = p.BaseHP
	for _, eq := range p.Equipment {
		p.MaxHP += eq.HPBonus
	}
}

// Affiche l’inventaire
func (p *Player) ShowInventory() {
	if len(p.Inventory) == 0 {
		fmt.Println("🎒 Inventaire vide.")
		return
	}
	fmt.Println("🎒 Inventaire :")
	for _, item := range p.Inventory {
		fmt.Printf(" - %s (+%d HP)\n", item.Name, item.HPBonus)
	}
}

// Affiche les équipements portés
func (p *Player) ShowEquipment() {
	if len(p.Equipment) == 0 {
		fmt.Println("⚔️ Aucun équipement équipé.")
		return
	}
	fmt.Println("⚔️ Équipements équipés :")
	for _, eq := range p.Equipment {
		fmt.Printf(" - %s (+%d HP)\n", eq.Name, eq.HPBonus)
	}
}

// Affiche le statut du joueur
func (p *Player) ShowStatus() {
	fmt.Println("❤️  PV Max :", p.MaxHP)
	p.ShowEquipment()
}

func Equipement() {
	// Création joueur
	player := NewPlayer(100)

	// Équipements de départ
	player.AddToInventory(Equipment{Name: "Capuche du Shinobi", Type: Head, HPBonus: 10})
	player.AddToInventory(Equipment{Name: "Veste du Shinobi", Type: Body, HPBonus: 25})
	player.AddToInventory(Equipment{Name: "Tabi du Shinobi", Type: Feet, HPBonus: 15})

	fmt.Println("👋 Bienvenue dans le jeu d’équipement !")
	fmt.Println("Commandes : inventory | equip <nom> | status | quit")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		command := strings.ToLower(parts[0])

		switch command {
		case "inventory":
			player.ShowInventory()
		case "equip":
			if len(parts) < 2 {
				fmt.Println("⚠️ Utilisation : equip <nom de l’objet>")
				continue
			}
			player.EquipByName(parts[1])
		case "status":
			player.ShowStatus()
		case "quit":
			fmt.Println("👋 Au revoir !")
			return
		default:
			fmt.Println("❓ Commande inconnue :", command)
		}
	}
}
