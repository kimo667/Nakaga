package main

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

const (
	BaseInventoryCap     = 10
	InventoryUpgradeStep = 10
	MaxInventoryUpgrades = 3

	HealRedBull        = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
)

// --------- ÉQUIPEMENT (slots) ----------
type Equipment struct {
	Head  string // tête
	Torso string // torse
	Feet  string // pieds
}

// --------- PERSONNAGE ----------
type Character struct {
	Name        string
	Class       Classe
	Level       int
	HPMax       int
	HP          int
	BaseHPMax   int // base “nue” (sans bonus d’équipement)
	Inventory   map[string]int
	Skills      []string
	Gold        int
	CapMax      int
	InvUpgrades int

	Equipment Equipment // <<< DOIT exister
}
