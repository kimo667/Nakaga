# 🌸 Nakaga – RPG en ligne de commande (Projet RED)

**Nakaga** est un petit RPG en Go inspiré de l’univers de Yasuo (LoL).
Développé dans le cadre du projet RED (YNOV) par une équipe de 3 étudiants.

---

## 🚀 Lancement

### Pré-requis

* Go 1.22+
* Modules externes :

  * `github.com/charmbracelet/lipgloss` – affichage stylé
  * `github.com/eiannone/keyboard` – navigation clavier

### Exécution

Depuis la racine du projet :

```bash
go run ./src
```

---

## 🧩 Fonctionnalités

### 🎭 Personnage

* Création : nom + classe (Humain, Samouraï, Ninja).
* Statistiques : PV, Mana, Or, XP, Niveau.
* Mort → revive automatique à 50% PV.

### 🎒 Inventaire

* Capacité limitée, améliorable auprès du marchand.
* Consommables :

  * **RedBull** (+20 PV) – une gratuite au premier achat.
  * **Potion de vie** (+30 PV).
  * **Potion de poison** (inflige un DoT en combat).
  * **Potion de mana** (+15 Mana).
  * **Grimoire** (apprend le sort *Mur de vent*).

### ⚔️ Équipement

* Slots : Tête, Torse, Pieds.
* Chaque pièce peut donner des bonus (ex: +PVMax).
* Menu dédié pour équiper/déséquiper.
* Affichage dans la fiche du personnage.

### 📜 Missions (obligatoires)

6 missions scénarisées, avec états (Non commencée / En cours / Terminée) :

1. **Premier sang** : vaincre 1 sbire à l’entraînement.
2. **Épreuve de force** : accumuler de l’XP et monter de niveau.
3. **Éveil mystique** : lancer un sort en combat.
4. **Énergie d’éther** : utiliser une potion de mana.
5. **Artisan en herbe** : forger un objet.
6. **Entraînement complet** : terminer toutes les étapes précédentes.

> Toutes les missions doivent être validées pour débloquer le boss final.

### 🛒 Marchand

* Achats d’objets (consommables, ressources).
* RedBull gratuite (1x).
* Grimoire unique.
* Potion de mana.
* Amélioration de la capacité de l’inventaire.

### 🔨 Forgeron

* Recettes simples (armes, armures).
* Consomme matériaux + or.

### 🥋 Entraînement

* Combat tour par tour contre un sbire.
* Navigation au clavier (flèches + Entrée) ou fallback en saisie numérique.
* Initiative calculée à chaque tour.
* Récompense : XP et progression de mission.

### 👹 Boss final

* Débloqué après avoir complété toutes les missions.
* Combat scénarisé contre le **Frère d’Acier**.
* Tour par tour avec attaques spéciales, Mur de vent, sorts et mana.
* Récompense : Or, XP, compétence *Souffle de l’Acier* et écran de fin.

---

## 🎮 Menu principal

```
1) Informations du personnage
2) Inventaire
3) Marchand
4) Forgeron
5) Entraînement
6) Missions
7) Équipement
8) Boss final (verrouillé tant que missions incomplètes)
9) Quitter
```

---

## 🖼️ Aperçu (ASCII Art)

L’écran d’informations affiche un visuel ASCII pour renforcer l’ambiance :

```
⠀⠀⠀⠀⠀⠀⢰⡆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠘⡇⠀⠀⠀⢠⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢷⠀⢠⢣⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢘⣷⢸⣾⣇⣶⣦⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⣿⣿⣿⣹⣿⣿⣷⣿⣆⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢼⡇⣿⣿⣽⣶⣶⣯⣭⣷⣶⣿⣿⣶⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠸⠣⢿⣿⣿⣿⣿⡿⣛⣭⣭⣭⡙⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⠿⠿⠿⢯⡛⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣾⣿⡿⡷⢿⡄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
```

---

## 📂 Structure du projet

Tous les fichiers sources sont dans `/src` (package `main`).
19 fichiers principaux :

* `types.go` – structures et constantes
* `character.go` – création personnage, init missions
* `inventory.go` – logique inventaire
* `inventory_menu.go` – menu Inventaire
* `consumables.go` – consommables et effets
* `skills.go` – gestion compétences
* `merchant.go` – marchand
* `forge.go` – forgeron
* `display.go` – affichages (infos, inventaire, équipement)
* `input.go` – helpers I/O
* `colors.go` – codes ANSI
* `main_menu.go` – menu principal
* `training.go` – combat d’entraînement
* `status.go` – mort & revive
* `bossfinal.go` – boss final
* `missions.go` – système de missions
* `equipment_menu.go` – menu équipement
* `equipement.go` – logique équipement
* `main.go` – point d’entrée

---

## 👨‍💻 Auteurs

* **Kimo** – 
* **Gabi** – 
* **Nat** –

---

## 📝 Notes

* Projet scolaire – pas optimisé pour la production.
* Le code est clair et découpé pour faciliter la lecture.
* Inspiré de Yasuo (LoL), mais libre et créatif.
