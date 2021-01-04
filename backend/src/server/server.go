package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"me.dev/go-board-game/common"
	"me.dev/go-board-game/games"
	"me.dev/go-board-game/games/tictactoe"
	"me.dev/go-board-game/mcts"
)

// Start the http (html + rest) server
func Start() {
	handleRequests()
}

func handleRequests() {
	path := http.Dir("/work/git/go-board-games/frontend/build")
	fs := http.FileServer(path)

	r := mux.NewRouter().StrictSlash(true)

	addGame(r, "tictactoe", tictactoe.NewTicTacToeGame)
	addGame(r, "siam", games.NewGame)

	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":10000", handler))
}

func addGame(r *mux.Router, name string, gameCtor func() common.Move) {
	handleMoves := func(w http.ResponseWriter, r *http.Request) {
		state := mux.Vars(r)["state"]
		game := gameCtor().PlayMovesByString(state)
		status := game.GetGameStatus()
		moves := game.NextAvailableMoves()

		var response GameDto
		var movesDto []MoveDto

		if status.IsDone() {
			moves = []common.Move{}
		}

		movesDto = make([]MoveDto, len(moves))
		for i, move := range moves {
			movesDto[i] = MoveDto{Value: move.GetJSON()}
		}
		response = GameDto{
			State: game.GetJSON(),
			Moves: movesDto}

		json.NewEncoder(w).Encode(response)
	}

	handleAi := func(w http.ResponseWriter, r *http.Request) {
		state := mux.Vars(r)["state"]
		game := gameCtor().PlayMovesByString(state)

		_, root := mcts.FindBestMove(game, mcts.MultithreadedConfig())

		json.NewEncoder(w).Encode(root.GetJSON())
	}

	r.HandleFunc("/api/"+name+"/{state}/moves", handleMoves)
	r.HandleFunc("/api/"+name+"/{state}/ai", handleAi)
}
