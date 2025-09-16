package main

// ====== Types & constantes globales ======

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

// Capacité d’inventaire et upgrades (T12+)
const (
	BaseInventoryCap     = 10 // capacité de base
	InventoryUpgradeStep = 10 // +10 par upgrade
	MaxInventoryUpgrades = 3  // max 3 upgrades
)

// Effets & valeurs de gameplay
const (
	HealRedBull        = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
)

// Personnage
type Character struct {
	Name        string
	Class       Classe
	Level       int
	HPMax       int
	HP          int
	Inventory   map[string]int
	Skills      []string
	Gold        int // (T11) initialisé plus tard
	CapMax      int // capacité actuelle
	InvUpgrades int // nb d’upgrades de capacité
}
