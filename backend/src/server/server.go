package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"me.dev/go-board-game/common"
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

    r.HandleFunc("/api/tictactoe/{state}/moves", handleTictactoeMoves)
    r.HandleFunc("/api/tictactoe/{state}/ai", handleTictactoeAi)
    r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowCredentials: true,
    })
    handler := c.Handler(r)

    log.Fatal(http.ListenAndServe(":10000", handler))
}

func handleTictactoeMoves(w http.ResponseWriter, r *http.Request) {
    state := mux.Vars(r)["state"]
    game := tictactoe.NewGame().PlayMovesByString(state)
    status := game.GetGameStatus()
    moves := game.NextAvailableMoves()
    
    var response GameDto
    var movesDto []MoveDto

    if status.IsDone() {
        moves = []common.Move {}
    }

    movesDto = make([]MoveDto, len(moves))
    for i, move := range moves {
        movesDto[i] = MoveDto{Value:move.GetJSON()}
    }
    response =  GameDto{
        State: game.GetJSON(),
        Moves: movesDto}
    
    json.NewEncoder(w).Encode(response)
}

func handleTictactoeAi(w http.ResponseWriter, r *http.Request) {
    state := mux.Vars(r)["state"]
    game := tictactoe.NewGame().PlayMovesByString(state)
    
    _, root := mcts.FindBestMove(game, mcts.BasicConfig())
    
    json.NewEncoder(w).Encode(root.GetJSON())
}