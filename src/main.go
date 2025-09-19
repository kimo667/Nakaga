package main

import (
	"bufio"
	"os"
)

func main() {
	// Lecteur stdin pour tous les menus
	reader := bufio.NewReader(os.Stdin)

	// Création interactive (nom + classe, PV init, inv de départ…)
	player := createCharacterInteractive(reader)

	// Boucle de jeu principale
	for mainMenu(&player, reader) {
		// Sécurité: revive si 0 PV
		isDead(&player)
	}
}
