import axios from "axios";

const restUrl =
    process.env.NODE_ENV === "production"
        ? `${window.location.origin}/api`
        : "http://localhost:10000/api";

class BackendApi {
    backend = axios.create({
        baseURL: restUrl,
    });

    getMoves(game, playedMoves) {
        return this.backend.get(`${game}/${playedMoves}/moves`);
    }

    async getAiMove(game, playedMoves) {
        console.log(`getAiMove ${game} ${playedMoves}`);
        return await this.backend.get(`${game}/${playedMoves}/ai`);
    }
}

class BackendApiMock {
    getMoves(game, playedMoves) {
        // playedMoves == 46231
        const data = {
            state: { state: " XXOX O  ", lastMove: "1", nextPlayer: "p2" },
            moves: [
                {
                    value: {
                        state: "OXXOX O  ",
                        lastMove: "0",
                        nextPlayer: "p1",
                    },
                },
                {
                    value: {
                        state: " XXOXOO  ",
                        lastMove: "5",
                        nextPlayer: "p1",
                    },
                },
                {
                    value: {
                        state: " XXOX OO ",
                        lastMove: "7",
                        nextPlayer: "p1",
                    },
                },
                {
                    value: {
                        state: " XXOX O O",
                        lastMove: "8",
                        nextPlayer: "p1",
                    },
                },
            ],
        };

        return this.sleepAndReturn(1000, { data: data });
    }

    getAiMove(game, playedMoves) {
        const data = {
            move: { state: " XXOX O  ", lastMove: "1", nextPlayer: "p2" },
            children: [
                {
                    move: {
                        state: "OXXOX O  ",
                        lastMove: "0",
                        nextPlayer: "p1",
                    },
                    children: null,
                    visits: 6228,
                    score: 6228,
                },
                {
                    move: {
                        state: " XXOXOO  ",
                        lastMove: "5",
                        nextPlayer: "p1",
                    },
                    children: [
                        {
                            move: {
                                state: "XXXOXOO  ",
                                lastMove: "0",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 402,
                            score: 402,
                        },
                        {
                            move: {
                                state: " XXOXOOX ",
                                lastMove: "7",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 413,
                            score: 413,
                        },
                        {
                            move: {
                                state: " XXOXOO X",
                                lastMove: "8",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 136,
                            score: 32,
                        },
                    ],
                    visits: 951,
                    score: 104,
                },
                {
                    move: {
                        state: " XXOX OO ",
                        lastMove: "7",
                        nextPlayer: "p1",
                    },
                    children: [
                        {
                            move: {
                                state: "XXXOX OO ",
                                lastMove: "0",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 1360,
                            score: 1360,
                        },
                        {
                            move: {
                                state: " XXOXXOO ",
                                lastMove: "5",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 156,
                            score: 0,
                        },
                        {
                            move: {
                                state: " XXOX OOX",
                                lastMove: "8",
                                nextPlayer: "p2",
                            },
                            children: [
                                {
                                    move: {
                                        state: "OXXOX OOX",
                                        lastMove: "0",
                                        nextPlayer: "p1",
                                    },
                                    children: null,
                                    visits: 229,
                                    score: 229,
                                },
                                {
                                    move: {
                                        state: " XXOXOOOX",
                                        lastMove: "5",
                                        nextPlayer: "p1",
                                    },
                                    children: [
                                        {
                                            move: {
                                                state: "XXXOXOOOX",
                                                lastMove: "0",
                                                nextPlayer: "p2",
                                            },
                                            children: null,
                                            visits: 133,
                                            score: 66.5,
                                        },
                                    ],
                                    visits: 133,
                                    score: 66.5,
                                },
                            ],
                            visits: 362,
                            score: 66.5,
                        },
                    ],
                    visits: 1878,
                    score: 451.5,
                },
                {
                    move: {
                        state: " XXOX O O",
                        lastMove: "8",
                        nextPlayer: "p1",
                    },
                    children: [
                        {
                            move: {
                                state: "XXXOX O O",
                                lastMove: "0",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 434,
                            score: 434,
                        },
                        {
                            move: {
                                state: " XXOXXO O",
                                lastMove: "5",
                                nextPlayer: "p2",
                            },
                            children: [
                                {
                                    move: {
                                        state: "OXXOXXO O",
                                        lastMove: "0",
                                        nextPlayer: "p1",
                                    },
                                    children: null,
                                    visits: 40,
                                    score: 40,
                                },
                                {
                                    move: {
                                        state: " XXOXXOOO",
                                        lastMove: "7",
                                        nextPlayer: "p1",
                                    },
                                    children: null,
                                    visits: 42,
                                    score: 42,
                                },
                            ],
                            visits: 82,
                            score: 0,
                        },
                        {
                            move: {
                                state: " XXOX OXO",
                                lastMove: "7",
                                nextPlayer: "p2",
                            },
                            children: null,
                            visits: 427,
                            score: 427,
                        },
                    ],
                    visits: 943,
                    score: 82,
                },
            ],
            visits: 10000,
            score: 9673.5,
        };

        return this.sleepAndReturn(1000, { data: data });
    }

    sleepAndReturn(ms, value) {
        return new Promise((resolve) =>
            setTimeout(() => {
                console.log(`sleepAndReturn ${JSON.stringify(value)}`);
                resolve(value);
            }, 2000)
        );
    }
}

const rest =
    process.env.REACT_APP_MODE === "test"
        ? new BackendApiMock()
        : new BackendApi();

export default rest;
