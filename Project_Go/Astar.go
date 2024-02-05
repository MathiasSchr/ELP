package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Coordonnees struct {
    X, Y float64
}

func manhattanDist(p0, p1 Coordonnees) int {
	return int(math.Abs(p0.X-p1.X)) + int(math.Abs(p0.Y-p1.Y))
}

func voisin(p Coordonnees, laby [][]float64) []Coordonnees {
	var res []Coordonnees
	if p.X-1 >= 0 && laby[int(p.X-1)][int(p.Y)] == 0 {
	res = append(res, Coordonnees{p.X-1,p.Y})
	}
	if p.X+1 < float64(len(laby[0])) && laby[int(p.X+1)][int(p.Y)] == 0 {
		res = append(res, Coordonnees{p.X+1,p.Y})
	}
	if p.Y-1 >= 0 && laby[int(p.X)][int(p.Y-1)] == 0 {
		res = append(res, Coordonnees{p.X,p.Y-1})
	}
	if p.Y+1 < float64(len(laby[1])) && laby[int(p.X)][int(p.Y+1)] == 0 {
			res = append(res, Coordonnees{p.X,p.Y+1})
	}
	return res
}

func min(aExplorer map[Coordonnees]bool, coord_arrivee Coordonnees) Coordonnees  {
	
	coord_min := Coordonnees{-1,-1}
	dist_min := int(math.Inf(1))
	for coord, explore := range aExplorer {
		if explore && manhattanDist(coord, coord_arrivee) < dist_min {
			coord_min = coord
			dist_min = manhattanDist(coord, coord_arrivee)
		}
	}
	return coord_min
}

func reconstructionChemin(coord_depart, coord_arrivee Coordonnees, cameFrom map[Coordonnees]Coordonnees) []Coordonnees {
	chemin := []Coordonnees{coord_arrivee}
	courant := coord_arrivee

	for courant != coord_depart {
		courant = cameFrom[courant]
		chemin = append([]Coordonnees{courant}, chemin...)
	}

	return chemin
}

func astar(laby [][]float64, depart []float64, arrivee []float64) []Coordonnees {
	
	// Initialisation
	trouve := false // Booléen pour detecter si on est arrivé
	aExplorer := make(map[Coordonnees]bool) // Map des noeuds à explorer
	dist := make(map[Coordonnees]int) // Map des distances associées à chaques noeuds
	explore := make(map[Coordonnees]bool) // Map des des noeuds déjà explorés
	cameFrom := make(map[Coordonnees]Coordonnees) // Map des noeuds parents

	coord_depart := Coordonnees{depart[0],depart[1]}
	coord_arrivee := Coordonnees{arrivee[0],arrivee[1]}
	aExplorer[coord_depart] = true
	dist[coord_depart] = 0

	for len(aExplorer) > 0 && trouve == false { // Tant qu'il y a des noeuds à explorer et que l'on est pas arrivé
		
		// Selectioner le noeud avec le cout estimé le plus faible qui n'a pas déjà été exploré
		var courant Coordonnees
		courant = min(aExplorer, coord_arrivee)
		
		// Tester si l'on est arrivé
		if courant == coord_arrivee {
			trouve = true
			fmt.Println("Chemin trouvé")
		} else {

			// Pour tous les voisins du noeud courant
			for voisinCourant := range voisin(courant, laby) {

				// On met à jour leur coût le plus bas dans dist
				cout := dist[courant] + manhattanDist(courant, voisin(courant, laby)[voisinCourant])
				if !explore[voisin(courant, laby)[voisinCourant]] || cout < dist[voisin(courant, laby)[voisinCourant]] {
					dist[voisin(courant, laby)[voisinCourant]] = cout
					cameFrom[voisin(courant, laby)[voisinCourant]] = courant
				} else {
					if !explore[voisin(courant, laby)[voisinCourant]] || cout == dist[voisin(courant, laby)[voisinCourant]] {
						rand_numb := int(rand.Intn(2))
						if rand_numb == 1 {
							cameFrom[voisin(courant, laby)[voisinCourant]] = courant
						}
					}
				}

					// On le met dans file d'exploration aExplorer si il n'a pas déjà été exploré 
					if !explore[voisin(courant, laby)[voisinCourant]] {
						aExplorer[voisin(courant, laby)[voisinCourant]] = true
					}
			}
		}
		
		// Indiquer que l'on a explore le noeud courant
		aExplorer[courant] = false
		explore[courant] = true
	}

	// Si on a trouve l'arrivée on reconstruit le chemin
	if trouve {
		return reconstructionChemin(coord_depart, coord_arrivee, cameFrom)
	}

	return nil
}

func genererLabyrinthe(n int) [][]float64 {
	// Création d'un labyrinthe vide de taille n x n.
	laby := make([][]float64, n)
	for i := range laby {
		laby[i] = make([]float64, n)
	}

	return laby
}

func afficherLabyrinthe(laby [][]float64) {
	// Affichage du labyrinthe.
	for _, ligne := range laby {
		fmt.Println(ligne)
	}
}

func genererLabyrintheFaisable(n int) [][]float64 {
	// Création d'un labyrinthe avec des murs représentés par des 1.0.
	laby := make([][]float64, n)
	for i := range laby {
		laby[i] = make([]float64, n)
		for j := range laby[i] {
			laby[i][j] = 1.0
		}
	}

	// Initialisation du générateur de nombres aléatoires.
	rand.Seed(time.Now().UnixNano())

	// Appel de la fonction de génération du labyrinthe.
	genererLabyrintheRecursiveBacktracking(laby, 0, 0, n-1, n-1)

	return laby
}

func genererLabyrintheRecursiveBacktracking(laby [][]float64, startX, startY, endX, endY int) {
	// Algorithme de génération de labyrinthe par le "recursive backtracking".

	if startX > endX || startY > endY {
		return
	}

	// Choix aléatoire d'une direction (haut, bas, gauche, droite).
	directions := rand.Perm(4)

	for _, dir := range directions {
		nextX, nextY := startX, startY

		switch dir {
		case 0: // Haut
			nextY -= 2
		case 1: // Bas
			nextY += 2
		case 2: // Gauche
			nextX -= 2
		case 3: // Droite
			nextX += 2
		}

		if nextX >= 0 && nextX <= endX && nextY >= 0 && nextY <= endY && laby[nextY][nextX] == 1.0 {
			// Détruire le mur entre les cellules actuelles et suivantes.
			laby[startY+(nextY-startY)/2][startX+(nextX-startX)/2] = 0.0
			laby[nextY][nextX] = 0.0

			// Récursion pour la prochaine cellule.
			genererLabyrintheRecursiveBacktracking(laby, nextX, nextY, endX, endY)
		}
	}
}

func tracerChemin(laby [][]float64, chemin []Coordonnees) {
	for coord := range chemin {
		laby[int(chemin[coord].X)][int(chemin[coord].Y)] = 2.0
	}
}

func main() {

	// Définir la taille du labyrinthe (n x n).
	tailleLabyrinthe := 10

	// Générer le labyrinthe faisable.
	labyrinthe := genererLabyrintheFaisable(tailleLabyrinthe)

	// Afficher le labyrinthe.
	afficherLabyrinthe(labyrinthe)

	depart := []float64{0,0}
	arrivee := []float64{float64(tailleLabyrinthe-2), float64(tailleLabyrinthe-2)}

	fmt.Println(astar(labyrinthe,depart,arrivee))
	chemin := astar(labyrinthe,depart,arrivee)

	// Tracer le chemin dans le labyrinthe avec des 2.
	tracerChemin(labyrinthe, chemin)

	// Afficher le labyrinthe.
	afficherLabyrinthe(labyrinthe)

}
