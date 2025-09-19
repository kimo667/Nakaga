package main

import "fmt"

/* ====== Mort / revive (T8) ====== */

func isDead(c *Character) bool {
	if c.HP <= 0 {
		fmt.Println(CRed + "\n*** WASTED ***" + CReset)
		c.HP = clamp(c.HPMax/2, 1, c.HPMax)
		fmt.Printf(CGreen+"Vous êtes ressuscité avec %d/%d PV."+CReset+"\n", c.HP, c.HPMax)
		return true
	}
	return false
}
