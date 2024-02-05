package main

import (
	"log"
	"fmt"
	"net"
	"sync"
	"io"
	"time"
	"os"
	"math/rand"
	"math"
	"regexp"
	"strconv"
)

// Constante 
const (

	// Constantes connexion TCP
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

var (

	// Variables globales
	carte = [][]float64{
		{1,1,1,1,1,1,1},
		{1,0,0,1,0,0,1},
		{1,1,0,1,0,1,1},
		{1,0,0,0,0,1,1},
		{1,1,1,1,1,1,1},
		{1,0,1,0,0,0,1},
		{1,1,1,1,1,1,1}}
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

// Fonctions liées à l'algorithme A*

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
Astar(carte [][]float64, depart Coordonnees, arrivee Coordonnees, pathChannel chan []Coordonnees, wg *sync.WaitGroup) -> Void
	Paramètres :
		carte :
		depart :
		arrivee :
		pathChannel
		wg :
	Algo A* qui calcule le chemin entre le depart et l'arrivee 
*/
func astar(carte [][]float64, depart Coordonnees, arrivee Coordonnees, pathChannel chan []Coordonnees, wg *sync.WaitGroup) {

	// Initialisation
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
				} else {
					if !explore[voisin(courant, carte)[voisinCourant]] || cout == dist[voisin(courant, carte)[voisinCourant]] {
						rand_numb := int(rand.Intn(2))
						if rand_numb == 1 {
							cameFrom[voisin(courant, carte)[voisinCourant]] = courant
						}
					}
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
	defer wg.Done()
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

// Fonctions liées à la connexion TCP

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
	
	// Elabore la réponse (algo A*)
	pathChannel := make(chan []Coordonnees)

	var wg sync.WaitGroup

	wg.Add(1)
	go astar(carte,depart,arrivee, pathChannel, &wg)

	go func() {
		wg.Wait()
		close(pathChannel)
	}()

	chemin := <- pathChannel
	reponseStr := fmt.Sprint("Chemin : ", chemin)

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

func main() {

	serveur := NewServer(HOST+":"+PORT)
	time.Sleep(time.Second * 25)
	serveur.Stop()

}