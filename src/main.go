package main

import (
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Création interactive (nom + classe)
	c := createCharacterInteractive(reader)

	for mainMenu(&c, reader) {
		isDead(&c) // revive auto si besoin
	}
}
