package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Player struct {
	Number      int
	Dices       []int // Last dices the player has
	RollResults []int // Result of dices that are rolled by player
	Point       int
}

func (p Player) RollDice() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(6) + 1
}

func (p *Player) IncrementPoint() {
	p.Point += 1
}

func initGame(totalPlayers, totalDices int) []Player {
	var players = make([]Player, totalPlayers)

	for i := 0; i < totalPlayers; i++ {
		var dices = make([]int, totalDices)
		var player = Player{
			Number: i + 1,
			Dices:  dices,
			Point:  0,
		}
		players[i] = player
	}

	return players
}

// TODO: Handle the case when multiple players have the same highest point
func getWinner(players []Player) Player {
	var winner Player = players[0]
	for _, p := range players {
		if p.Point > winner.Point {
			winner = p
		}
	}

	return winner
}

func isPlayerAbleToRoll(player Player) bool {
	return len(player.Dices) > 0
}

func isGameFinished(players []Player) bool {
	var totalPlayerHasDices = len(players)
	for _, player := range players {
		if len(player.Dices) == 0 {
			totalPlayerHasDices -= 1
		}
	}
	return totalPlayerHasDices <= 1
}

func evaluateGame(players []Player) {
	var additionalDicesByPlayer = make(map[int][]int, 0)
	for i := 0; i < len(players); i++ {
		p := &players[i]

		var evaluatedDices []int
		for _, r := range p.RollResults {
			switch r {
			case 6:
				p.IncrementPoint()
			case 1:
				var nextPlayerIndex int
				if i == len(players)-1 {
					nextPlayerIndex = 0 // first player
				} else {
					nextPlayerIndex = i + 1
				}
				additionalDicesByPlayer[nextPlayerIndex] = append(additionalDicesByPlayer[nextPlayerIndex], r)
			default:
				evaluatedDices = append(evaluatedDices, r)
			}
		}

		p.Dices = evaluatedDices
		p.RollResults = nil
	}

	for i := 0; i < len(players); i++ {
		players[i].Dices = append(players[i].Dices, additionalDicesByPlayer[i]...)
	}
}

func printRollDiceResults(player Player) {
	var resultsBuilder strings.Builder
	for i := 0; i < len(player.RollResults); i++ {
		resultsBuilder.WriteString(
			fmt.Sprintf("%d ", player.RollResults[i]),
		)
	}
	fmt.Printf("Player #%d (%d): %s\n", player.Number, player.Point, resultsBuilder.String())
}

func printEvaluationResults(players []Player) {
	for _, player := range players {
		var resultsBuilder strings.Builder
		for _, dice := range player.Dices {
			resultsBuilder.WriteString(fmt.Sprintf("%d ", dice))
		}
		fmt.Printf("Player #%d (%d): %s\n", player.Number, player.Point, resultsBuilder.String())
	}
}

func main() {
	var totalPlayers, totalDices int
	flag.IntVar(&totalPlayers, "pemain", 3, "Jumlah pemain")
	flag.IntVar(&totalDices, "dadu", 4, "Jumlah Dadu")
	flag.Parse()
	fmt.Printf("Pemain = %d, Jumlah = %d\n", totalPlayers, totalDices)

	var players = initGame(totalPlayers, totalDices)
	var isFinished = false
	var count = 1
	for !isFinished {
		fmt.Printf("Percobaan lempar dadu ke-%d:\n", count)
		for ip, p := range players {
			if isPlayerAbleToRoll(p) {
				for i := 0; i < len(p.Dices); i++ {
					p.RollResults = append(p.RollResults, p.RollDice())
				}
				players[ip] = p
				printRollDiceResults(p)
			}
		}

		evaluateGame(players)
		fmt.Println("Setelah Evaluasi: ")
		printEvaluationResults(players)
		isFinished = isGameFinished(players)
		count += 1
	}

	player := getWinner(players)
	fmt.Printf("Winner: Pemain #%d\n", player.Number)
}
