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
	for mainMenu(&player, reader) { // ici on appelle la fonction menu, pas le type menuItem
		// Sécurité: si le joueur tombe à 0 PV quelque part, on applique le revive T8
		isDead(&player)
	}
}
