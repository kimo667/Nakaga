package main

// Capacité & Effets
const (
	BaseInventoryCap     = 10 // capacité de base
	InventoryUpgradeStep = 10 // +10 par upgrade
	MaxInventoryUpgrades = 3  // max 3 upgrades

	HealRedBull        = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
)

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samouraï"
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
