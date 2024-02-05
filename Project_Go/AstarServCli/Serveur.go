package main

import (
	"log"
	"fmt"
	"net"
	"sync"
	"io"
	"time"
	"os"
	"math"
	"regexp"
	"strconv"
)

// -- Constante --
const (

	// Constantes connexion TCP
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

var (

	// Carte du monde
	carte [][]float64 = [][]float64{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
)

// Constructeur de serveur
type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

// Constructeur de coordonnées
type Coordonnees struct {
    X, Y float64
}

// ---------------------------------------
// -- Fonctions liées à l'algorithme A* --
// ---------------------------------------

/*
Manhattan(p0 Coordonnees, p1 Coordonnees) -> Int
	Paramètres :
		p0 : coordonnees d'un premier point
		p1 : coordonnees d'un deuxième point
	Renvoie :
		La distance de manhattan entre p0 et p1
*/
func manhattanDist(p0, p1 Coordonnees) int {
	return int(math.Abs(p0.X-p1.X)) + int(math.Abs(p0.Y-p1.Y))
}

/*
Voisin(p Coordonnees, carte [][]float64) -> []Coordonnees
	Paramètres :
		p : coordonnees du point dont on veut chercher ses voisins
		carte : matrice representant la carte du monde
	Renvoie :
		Liste des coordonnees des voisins du point p
*/
func voisin(p Coordonnees, carte [][]float64) []Coordonnees {
	var res []Coordonnees
	if p.X-1 >= 0 && carte[int(p.X-1)][int(p.Y)] == 0 {
	res = append(res, Coordonnees{p.X-1,p.Y})
	}
	if p.X+1 < float64(len(carte[0])) && carte[int(p.X+1)][int(p.Y)] == 0 {
		res = append(res, Coordonnees{p.X+1,p.Y})
	}
	if p.Y-1 >= 0 && carte[int(p.X)][int(p.Y-1)] == 0 {
		res = append(res, Coordonnees{p.X,p.Y-1})
	}
	if p.Y+1 < float64(len(carte[1])) && carte[int(p.X)][int(p.Y+1)] == 0 {
			res = append(res, Coordonnees{p.X,p.Y+1})
	}
	return res
}
/*
Min(aExplorer map[Coordonnees]bool, arrivee Coordonnees) -> Coordonnees
	Paramètres :
		aExplorer : Map des noeuds qui doivent encore être explorer
		arrivee : coordonnees du point d'arrivee
	Renvoie :
		Coordonnees du noeud qui a la plus petite distance et encore non explore 
*/

func min(aExplorer map[Coordonnees]bool, dist map[Coordonnees]int) Coordonnees  {
	coord_min := Coordonnees{-1,-1}
	dist_min := int(math.Inf(1))
	for coord, explore := range aExplorer {
		if explore && dist[coord] < dist_min {
			coord_min = coord
			dist_min = dist[coord]
		}
	}
	return coord_min
}

/*
AstarNormal(carte [][]float64, depart Coordonnees, arrivee Coordonnees, pathChannel chan []Coordonnees) -> Void
	Paramètres :
		carte :
		depart :
		arrivee :
		pathChannel :
	Algo A* normal qui calcule le chemin entre le depart et l'arrivee 
*/
func astarNormal(carte [][]float64, depart Coordonnees, arrivee Coordonnees, pathChannel chan []Coordonnees, wg *sync.WaitGroup) {

	// Initialisation
	defer wg.Done()
	trouve := false // Booléen pour detecter si on est arrivé
	aExplorer := make(map[Coordonnees]bool) // Map des noeuds à explorer
	dist := make(map[Coordonnees]int) // Map des distances associées à chaques noeuds
	explore := make(map[Coordonnees]bool) // Map des des noeuds déjà explorés
	cameFrom := make(map[Coordonnees]Coordonnees) // Map des noeuds parents

	aExplorer[depart] = true
	dist[depart] = 0

	for len(aExplorer) > 0 && trouve == false { // Tant qu'il y a des noeuds à explorer et que l'on est pas arrivé

		// Selectioner le noeud avec le cout estimé le plus faible qui n'a pas déjà été exploré
		var courant Coordonnees
		courant = min(aExplorer, dist)
		
		// Tester si l'on est arrivé
		if courant == arrivee {
			trouve = true
			fmt.Printf("Chemin trouvé !\n")
		} else {

			// Pour tous les voisins du noeud courant
			for voisinCourant := range voisin(courant, carte) {

				// On met à jour leur coût le plus bas dans dist
				cout := dist[courant] + manhattanDist(courant, voisin(courant, carte)[voisinCourant])
				if !explore[voisin(courant, carte)[voisinCourant]] || cout < dist[voisin(courant, carte)[voisinCourant]] {
					dist[voisin(courant, carte)[voisinCourant]] = cout
					cameFrom[voisin(courant, carte)[voisinCourant]] = courant
				}

					// On le met dans file d'exploration aExplorer si il n'a pas déjà été exploré 
					if !explore[voisin(courant, carte)[voisinCourant]] {
						aExplorer[voisin(courant, carte)[voisinCourant]] = true
					}
			}
		}
		
		// Indiquer que l'on a explore le noeud courant
		delete(aExplorer, courant)
		explore[courant] = true

	}
	
	// Si on a trouve l'arrivée on reconstruit le chemin
	if trouve {

		// Reconstruction du chemin
		chemin := []Coordonnees{arrivee}
		cour := arrivee

		for cour != depart {
		cour = cameFrom[cour]
		chemin = append([]Coordonnees{cour}, chemin...)
		}
		pathChannel <- chemin
	} else {
		pathChannel <- nil
		fmt.Printf("Pas de chemin\n")
	}

}

/*
AstarDouble(carte [][]float64, depart Coordonnees, arrivee Coordonnees, pathChannel chan []Coordonnees, m *sync.Mutex) -> Void
	Paramètres :
		carte :
		depart :
		arrivee :
		pathChannel :
		m :
	Algo A* double qui calcule le chemin entre le depart et l'arrivee 
*/
func astarDouble(carte [][]float64, depart Coordonnees, arrivee Coordonnees, explore1 map[Coordonnees]bool, explore2 map[Coordonnees]bool, pathChannel chan []Coordonnees, RDVChannel chan Coordonnees, m *sync.Mutex, wg *sync.WaitGroup) {

	// Initialisation
	defer wg.Done()
	trouve := false // Booléen pour detecter si on est arrivé
	dist := make(map[Coordonnees]int) // Map des distances associées à chaques 
	aExplorer := make(map[Coordonnees]bool) // Map des noeuds à explorer
	cameFrom := make(map[Coordonnees]Coordonnees) // Map des noeuds parents
	var pointRDV Coordonnees // Noeud de rdv
	aExplorer[depart] = true // On commence par explorer le noeud de départ
	dist[depart] = manhattanDist(depart, arrivee) // On initialise le départ à la distance 0

	for len(aExplorer) > 0 && trouve == false { // Tant qu'il y a des noeuds à explorer et que l'on est pas arrivé

		var courant Coordonnees
		var test bool

		// !!!Passage critique!!!
		test, pointRDV = cleCommun(explore1, explore2, m) // On test si les deux process on exploré deux noeud en commun

		if test { // Si c'est le cas alors les deux process se sont rejoint
			
			trouve = true
			arrivee = pointRDV
			fmt.Println("On s'est rejoint !\n")

		} else {
		
			// Sinon on continue la recherche : on selectione le noeud avec le cout estimé le plus faible qui n'a pas déjà été exploré
			courant = min(aExplorer, dist)

			if courant == arrivee { // Tester si on est arrivé
				
				trouve = true
				fmt.Printf("Chemin trouvé !\n")

			} else { // Sinon on continue la recherche

				// Pour tous les voisins du noeud courant
				for voisinCourant := range voisin(courant, carte) {

 					// !!!Passage critique!!!
					// Si le noeud n'a pas encore été exploré ou que la distance totale est plus petite que celle déjà stockée
					m.Lock()
					if !explore1[voisin(courant, carte)[voisinCourant]] || dist[courant] + manhattanDist(courant, arrivee) < dist[voisin(courant, carte)[voisinCourant]] {
						
						m.Unlock()
						// On met à jour la distance totale la plus petite dans dist
						dist[voisin(courant, carte)[voisinCourant]] = dist[courant] + manhattanDist(courant, arrivee)
						cameFrom[voisin(courant, carte)[voisinCourant]] = courant // On met à jour le prédécesseur
					
					} else {

						m.Unlock()
					}
	
				
					// !!!Passage critique!!!
					// Si le(s) voisin(s) non exploré
					m.Lock()
					if !explore1[voisin(courant, carte)[voisinCourant]] {

						m.Unlock()
						aExplorer[voisin(courant, carte)[voisinCourant]] = true // On le(s) rajoute(s) dans la map des noeuds à explorer
					} else {

						m.Unlock()
					}
				}
			}

			// Indiquer que l'on a explore le noeud courant
			delete(aExplorer, courant)
			// !!!Passage critique!!!
			m.Lock()
			explore1[courant] = true // On indique que l'on a exploré le noeud courant
			m.Unlock()

		}
	}
	// Si on a trouve l'arrivée on reconstruit le chemin
	if trouve {

		// Reconstruction du chemin
		chemin := []Coordonnees{arrivee}
		courrantbis := arrivee

		for courrantbis != depart {
		courrantbis = cameFrom[courrantbis]
		chemin = append([]Coordonnees{courrantbis}, chemin...)
		}
		pathChannel <- chemin
		RDVChannel <- pointRDV
	} else {
		pathChannel <- nil
		fmt.Printf("Pas de chemin\n")
	}
}

/*
cleCommun(map1, map2 map[Coordonnees]bool, m *sync.Mutex) (bool, Coordonnees) --> bool, Coordonnees
	Paramètres :
		map1 :
		map2 :
		m :
	Renvoie :
		Si les map1 et map2 ont une clé en commun et laquelle
*/
func cleCommun(map1, map2 map[Coordonnees]bool, m *sync.Mutex) (bool, Coordonnees) {
	var key Coordonnees
	m.Lock()
	defer m.Unlock()
    for key := range map1 {
        if _, exists := map2[key]; exists {
            return true, key
        }
    }
	return false, key
}

func mergeMaps(map1, map2 []Coordonnees, pointRDV Coordonnees) []Coordonnees {
	mergedMap := make([]Coordonnees, 0)

	// Ajouter les points de la première map
	test := true
	for _, p := range map1 {
		if test {
			mergedMap = append(mergedMap, p)
			if p == pointRDV {
				test = false
			}
		}
	}

	// Ajouter les points de la deuxième map (sauf le premier, qui est déjà inclus dans la première map)
	for i := len(map2)-2; i >= 0; i-- {
		mergedMap = append(mergedMap, map2[i])
	}

	return mergedMap
}

// ----------------------------------------
// -- Fonctions liées à la connexion TCP --
// ----------------------------------------

func NewServer(addr string) *Server {
	s := &Server{
	  quit: make(chan interface{}),
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
	  log.Fatal(err)
	}
	s.listener = l
	s.wg.Add(1)
	go s.serve()
	return s
  }

func (s *Server) serve() {
  defer s.wg.Done()

  for {
    conn, err := s.listener.Accept()
    if err != nil {
      select {
      case <-s.quit:
        return
      default:
        log.Println("accept error", err)
      }
    } else {
      s.wg.Add(1)
      go func() {
        s.handleConection(conn)
        s.wg.Done()
      }()
    }
  }
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
  }

func (s *Server) handleConection(conn net.Conn) {
  defer conn.Close()
  buf := make([]byte, 1024)
  for {
    n, err := conn.Read(buf)
    if err != nil && err != io.EOF {
      log.Println("read error", err)
      return
    }
    if n == 0 {
      return
    }

	// Message recu du client
    fmt.Printf("Received from %v: %s\n", conn.RemoteAddr(), string(buf[:n]))
	depart, arrivee, err := extraireCoordonnees(string(buf[:n]))
	typeAstar, err := extraireType(string(buf[:n]))
	var reponseStr string
	
	debut := time.Now()
	// Elabore la réponse
	if typeAstar == "normal" {

		pathChannel := make(chan []Coordonnees)	// Channel
		var wg sync.WaitGroup // Waitgroup
		wg.Add(1)
		go astarNormal(carte,depart,arrivee, pathChannel, &wg) // A*
		chemin := <- pathChannel // Récupération du chemin (bloquant)

		fin := time.Now()
		temps := fin.Sub(debut)
		reponseStr = fmt.Sprint("Chemin : ", chemin, "Temps : ", temps)

		// Ferme le channel
		go func() {
			wg.Wait()
			close(pathChannel)
		}()

	} else {
		if typeAstar == "double" {

			pathChannel1 := make(chan []Coordonnees) // Channels
			pathChannel2 := make(chan []Coordonnees)
			RDVChannel := make(chan Coordonnees)
			var wg sync.WaitGroup // Waitgroup
			var m sync.Mutex // Mutex
			// Base de donnée des noeuds explorés par chacun des process
			explore1 := make(map[Coordonnees]bool) // Map des noeuds déjà explorés
			explore2 := make(map[Coordonnees]bool) // Map des noeuds déjà explorés

			wg.Add(2)
			go astarDouble(carte, depart, arrivee, explore1, explore2, pathChannel1, RDVChannel, &m, &wg) // A* du départ
			go astarDouble(carte, arrivee, depart, explore2, explore1, pathChannel2, RDVChannel, &m, &wg) // A* de l'arrivée
			chemin1, chemin2, pointRDV := <- pathChannel1, <- pathChannel2, <- RDVChannel // Récupération des chemins
			// Fusionner les chemins des deux process astar
			chemin := mergeMaps(chemin1, chemin2, pointRDV)

			fin := time.Now()
			temps := fin.Sub(debut)
			reponseStr = fmt.Sprint("Chemin : ", chemin, " Temps de la requête : ", temps)

			// Ferme les channels
			go func() {
				wg.Wait()
				close(pathChannel1)
				close(pathChannel2) 
				close(RDVChannel)
			}()

		} else {
			reponseStr = fmt.Sprint("Erreur")
		}
	}

	// Envoie de la réponse au client
	fmt.Printf("Sent to %v: %s\n", conn.RemoteAddr(), reponseStr)
	_, err = conn.Write([]byte(reponseStr))
	if err != nil {
    	fmt.Println("Write data failed:", err.Error())
    	os.Exit(1)
	}

	// Indique au client la fin du message
	_, err = conn.Write([]byte("\n"))
	if err != nil {
    	fmt.Println("Write data failed:", err.Error())
    	os.Exit(1)
	}
  }
}

func extraireCoordonnees(s string) (depart Coordonnees, arrivee Coordonnees, err error) {
	// Utilisation d'une expression régulière pour extraire les coordonnées
	compile := regexp.MustCompile("\\{(\\d+),(\\d+)\\}")

	// Recherche des coordonnées de départ
	matchDepart := compile.FindStringSubmatch(s)
	if len(matchDepart) != 3 {
		err = fmt.Errorf("Les coordonnées de départ ne peuvent pas être extraites")
		return
	}

	// Conversion des coordonnées de départ en float64 puis en coordonnees
	xDepart, errX := strconv.ParseFloat(matchDepart[1], 64)
	yDepart, errY := strconv.ParseFloat(matchDepart[2], 64)
	if errX != nil || errY != nil {
		err = fmt.Errorf("Erreur lors de la conversion des coordonnées de départ en entiers")
		return
	}
	depart = Coordonnees{xDepart, yDepart}

	// Recherche des coordonnées d'arrivée
	matchArrivee := compile.FindAllStringSubmatch(s, -1)
	if len(matchArrivee) != 2 {
		err = fmt.Errorf("Les coordonnées d'arrivée ne peuvent pas être extraites")
		return
	}

	// Conversion des coordonnées d'arrivée en float64 puis en cooordonnees
	xArrivee, errXArr := strconv.ParseFloat(matchArrivee[1][1], 64)
	yArrivee, errYArr := strconv.ParseFloat(matchArrivee[1][2], 64)
	if errXArr != nil || errYArr != nil {
		err = fmt.Errorf("Erreur lors de la conversion des coordonnées d'arrivée en entiers")
		return
	}
	arrivee = Coordonnees{xArrivee, yArrivee}

	return
}

func extraireType(chaine string) (string, error) {
	// Utilisation d'une expression régulière pour extraire le type
	compile := regexp.MustCompile(`Type : ([a-zA-Z]+)`)
	matches := compile.FindStringSubmatch(chaine)

	// Vérifier si une correspondance a été trouvée
	if len(matches) != 2 {
		return "", fmt.Errorf("Aucun type trouvé dans la chaîne")
	}

	return matches[1], nil
}

// -------------------------
// -- Programme principal --
// -------------------------

func main() {

	serveur := NewServer(HOST+":"+PORT)
	time.Sleep(time.Second * 50)
	serveur.Stop()

}