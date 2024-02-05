package main

import (
	"net"
	"os"
	"strings"
	"fmt"
	"regexp"
	"strconv"
)

const (
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

	depart string
	arrivee string
	typeAlgo string 

)

func afficherCarte(carte [][]float64) {

	println("")
	println("             -- Voici la carte du monde --")
	println("")
	print("     x\n")
	print("     0---------5---------10-------15--------20--------\n")
	print("\n")
	for i := 0; i < len(carte); i++ {
		echelley := ""
		if i == 0 {
			echelley = "y 0 "
		} else {
			if i%5 == 0 {
				if i == 5 {
					echelley = fmt.Sprint(" ", i, "  ")
				} else {
					echelley = fmt.Sprint(" ", i, " ")
					}
			} else {
				echelley = "  | "
			}
		}
		for j := 0; j < len(carte[i]); j++ {
			if j == 0 {
				fmt.Printf("%s %v ", echelley, carte[i][j])
			} else {
				fmt.Printf("%v ", carte[i][j])
			}
		}
		fmt.Println()
	}
	print("\n")
}

func afficherCarteAvecChemin(carte [][]float64, chemin string) {
	// Convertir la chaîne de caractères en une liste de coordonnées
	coordRegex := regexp.MustCompile(`\{(\d+) (\d+)\}`)
	matches := coordRegex.FindAllStringSubmatch(chemin, -1)
	var cheminCoord [][]int

	for _, match := range matches {
		x, _ := strconv.Atoi(match[1])
		y, _ := strconv.Atoi(match[2])
		cheminCoord = append(cheminCoord, []int{x, y})
	}

	// Créer une copie de la carte pour éviter de modifier la carte originale
	carteAvecChemin := make([][]float64, len(carte))
	for i := range carte {
		carteAvecChemin[i] = make([]float64, len(carte[i]))
		copy(carteAvecChemin[i], carte[i])
	}

	// Placer le chemin sur la carte
	for _, coord := range cheminCoord {
		carteAvecChemin[coord[0]][coord[1]] = 2
	}

	// Afficher la carte avec le chemin
	print("             -- Tracé du chemin --\n")
	print("\n")
	for i := 0; i < len(carteAvecChemin); i++ {
		for j := 0; j < len(carteAvecChemin[i]); j++ {
			switch carteAvecChemin[i][j] {
			case 1:
				fmt.Print("1 ")
			case 2:
				fmt.Print("\033[91m2\033[0m ") // Séquence d'échappement ANSI pour le texte en rouge
			default:
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

func main() {
	
	afficherCarte(carte)

	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Quel est votre point de départ (format {x,y}) : ")
	fmt.Scan(&depart)
	fmt.Println("Quel est votre point de départ (format {x,y}) : ")
	fmt.Scan(&arrivee)
	fmt.Println("Quel type de recherche (normal ou double) : ")
	fmt.Scan(&typeAlgo)

	_, err = conn.Write([]byte("Départ : " + depart + ", Arrivée : " + arrivee + ", Type : " + typeAlgo))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	received := ""
	var builder strings.Builder

	for {
		n, err := conn.Read(buf)
		if err != nil {
			println("Read data failed:", err.Error())
			os.Exit(1)
		}

		received += string(buf[:n])

		builder.Write(buf[:n])
		if strings.Contains(builder.String(), "\n") {
			print("\n")
			println(received)
			conn.Close()
			break
		}
	}
	afficherCarteAvecChemin(carte, received)

}
