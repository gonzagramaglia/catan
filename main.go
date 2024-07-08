package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// ELO constants
const (
	initialRanking = 1200
	kFactor        = 32
)

// Player represents a player with a name and ranking
type Player struct {
	name    string
	ranking float64
}

// calculateExpectedScore calculates the expected score for a player
func calculateExpectedScore(rankingA, rankingB float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (rankingB-rankingA)/400))
}

// calculateNewRanking calculates the new ranking for a player
func calculateNewRanking(currentRanking, expectedScore, actualScore float64) float64 {
	return currentRanking + kFactor*(actualScore-expectedScore)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n")

	// Initialize players
	var players []Player
	for i := 1; i <= 3; i++ {
		j := ""
		switch i {
		case 1:
			j = "A"
		case 2:
			j = "B"
		case 3:
			j = "C"
		}
		fmt.Printf("Ingrese el nombre del jugador %s: ", j)
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Printf("Ingrese el ranking del jugador %s: ", name)
		rankingInput, _ := reader.ReadString('\n')
		rankingInput = strings.TrimSpace(rankingInput)
		ranking, err := strconv.ParseFloat(rankingInput, 64)
		if err != nil || ranking > 1800 || ranking < 800 {
			fmt.Println("Por favor, ingrese un ranking válido (entre 800 y 1800).")
			i--
			continue
		}

		players = append(players, Player{name: name, ranking: ranking})
	}

	fmt.Print("¿Quiere agregar un cuarto jugador? (s/n): ")
	morePlayers, _ := reader.ReadString('\n')
	morePlayers = strings.TrimSpace(morePlayers)
	if morePlayers == "s" || morePlayers == "S" {
		fmt.Print("Ingrese el nombre del jugador D: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Printf("Ingrese el ranking del jugador %s: ", name)
		rankingInput, _ := reader.ReadString('\n')
		rankingInput = strings.TrimSpace(rankingInput)
		ranking, err := strconv.ParseFloat(rankingInput, 64)
		if err != nil || ranking > 1800 || ranking < 800 {
			fmt.Println("Por favor, ingrese un ranking válido (entre 800 y 1800).")
		} else {
			players = append(players, Player{name: name, ranking: ranking})
		}
	}

	numPlayersInMatch := len(players)

	fmt.Println("")

	for {
		// Gather positions for each player in the match
		results := make(map[string]float64)
		playersInMatch := make([]*Player, numPlayersInMatch)

		for i := 1; i <= numPlayersInMatch; i++ {
			fmt.Printf("¿Quién salió en la posición %d?\n", i)
			for index, player := range players {
				fmt.Printf("%d: %s\n", index+1, player.name)
			}
			fmt.Print("Ingrese el número del jugador: ")
			playerIndexInput, _ := reader.ReadString('\n')
			playerIndexInput = strings.TrimSpace(playerIndexInput)
			playerIndex, err := strconv.Atoi(playerIndexInput)
			if err != nil || playerIndex < 1 || playerIndex > len(players) {
				fmt.Println("Por favor, ingrese un número válido.")
				i--
				continue
			}

			player := &players[playerIndex-1]

			// Assign positions to playersInMatch
			playersInMatch[i-1] = player

			// Assign scores based on position
			score := 1.0 - float64(i-1)/float64(numPlayersInMatch-1)
			results[player.name] = score
		}

		// Print table of current positions after gathering inputs
		fmt.Println("\nTabla de posiciones:")
		for index, player := range playersInMatch {
			position := ""
			switch index + 1 {
			case 1:
				position = "1ro"
			case 2:
				position = "2do"
			case 3:
				position = "3ro"
			case 4:
				position = "4to"
			}
			fmt.Printf("%s: %s\n", position, player.name)
		}

		// Ask user if the current table of positions is correct
		fmt.Print("¿Esto es correcto? (s/n): ")
		correctTableInput, _ := reader.ReadString('\n')
		correctTableInput = strings.TrimSpace(correctTableInput)
		if correctTableInput == "s" || correctTableInput == "S" {
			// Store initial rankings to calculate ELO variation later
			initialRankings := make(map[string]float64)
			for _, player := range playersInMatch {
				initialRankings[player.name] = player.ranking
			}

			// Process each player against all others in the match
			for i := 0; i < numPlayersInMatch; i++ {
				for j := 0; j < numPlayersInMatch; j++ {
					if i == j {
						continue
					}

					playerA := playersInMatch[i]
					playerB := playersInMatch[j]

					expectedScoreA := calculateExpectedScore(playerA.ranking, playerB.ranking)
					actualScoreA := results[playerA.name]

					playerA.ranking = calculateNewRanking(playerA.ranking, expectedScoreA, actualScoreA)
				}
			}
			fmt.Println("")

			// Print updated player rankings with ELO variation rounded up
			for _, player := range players {
				newRanking := math.Ceil(player.ranking)
				delta := newRanking - initialRankings[player.name]
				deltaRounded := int(math.Ceil(delta))
				fmt.Printf("Nuevo ranking de %s: %d (%d)\n", player.name, int(newRanking), deltaRounded)
			}

			fmt.Println("\n\n")
			break
		} else {
			fmt.Println("Vuelva a ingresar las posiciones.")
			fmt.Println("")
		}
	}
}
