package main

type Classe string

const (
	ClasseHumain  Classe = "Humain"
	ClasseSamurai Classe = "Samurai"
	ClasseNinja   Classe = "Ninja"
)

// Inventaire
const (
	BaseInventoryCap     = 10
	InventoryUpgradeStep = 10
	MaxInventoryUpgrades = 3
)

// Consommables / effets
const (
	HealRedBull        = 20
	HealPotionVie      = 20
	PoisonTicks        = 3
	PoisonDamagePerSec = 10
	ManaPotRestore     = 15 // (Mission 4)
)

// États de mission
type MissionState int

const (
	NotStarted MissionState = iota
	InProgress
	Completed
)

type Equipment struct {
	Head  string // Tête
	Torso string // Torse
	Feet  string // Pieds
}

type Mission struct {
	ID          int
	Title       string
	Description string
	State       MissionState

	RequiredTrainKills int
	TrainKills         int

	RewardGold        int
	RewardItem        string
	RewardSkill       string
	RewardUpgradeSlot bool
}

type Character struct {
	Name  string
	Class Classe
	Level int

	// Pv
	HPMax     int
	HP        int
	BaseHPMax int

	// Initiative (Mission 1)
	Initiative int

	// XP / niveau (Mission 2)
	XP    int
	XPMax int

	// Mana (Mission 4)
	Mana    int
	ManaMax int

	Inventory   map[string]int
	Skills      []string
	Gold        int
	CapMax      int
	InvUpgrades int

	Equipment Equipment
	Missions  []Mission
}
