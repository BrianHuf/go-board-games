import React from "react";
import { Link } from "react-router-dom";
import Container from "react-bootstrap/Container";
import Table from "react-bootstrap/Table";
import Spinner from "react-bootstrap/Spinner";
import Form from "react-bootstrap/Form";
import { Circle, CircleFill, Square, SquareFill } from "react-bootstrap-icons";

import rest from "../../api/backend";

export default class TicTacToe extends React.Component {
    state = {
        board: null,
        loading: false,
        isDone: false,
        winner: null,
        mode: "pvp",
    };

    componentDidMount() {
        // take care of browser backward/forward
        window.onpopstate = () => {
            const moves = this.getPlayedMoves(this.props.location.pathname);
            this.onLoadBoard(moves);
        };

        if (!this.state.board) {
            const moves = this.getPlayedMoves(this.props.location.pathname);
            this.onLoadBoard(moves);
        }
    }

    getPlayedMoves(url) {
        return url.substr(url.lastIndexOf("/") + 1);
    }

    async onLoadBoard(moves) {
        this.setState({ ...this.state, loading: true });

        const info = await rest.getMoves("tictactoe", moves);
        const response = info.data.state;
        this.setState({
            ...this.state,
            board: response.state,
            isDone: response.isDone,
            winner: response.winner,
            nextPlayer: response.nextPlayer,
            loading: false,
        });

        if (this.isComputerMoveNext()) {
            this.onPlayAi();
        }
    }

    async onPlayAi() {
        console.log("onPlayAI");
        const moves = this.getPlayedMoves(this.props.location.pathname);
        const aiMoves = await rest.getAiMove("tictactoe", moves);

        let maxScore = -1;
        let best = null;
        aiMoves.data.children.forEach((child) => {
            if (best == null || child.score > maxScore) {
                best = child;
                maxScore = child.score;
            }
        });
        this.onPlayMove(best.move.lastMove);
    }

    onRestart() {
        console.log(this.state.mode);
        this.setState({
            ...this.state,
            board: null,
            loading: false,
            isDone: false,
            winner: null,
        });
        this.onLoadBoard("-");
    }

    render() {
        console.log(`RENDER ${Date.now()} ${JSON.stringify(this.state)}`);
        if (this.state.board) {
            const isComputerNext = this.isComputerMoveNext();
            const cb =
                this.state.loading || this.state.isDone || isComputerNext
                    ? null
                    : this.onPlayMove.bind(this);

            const message = this.getMessage(isComputerNext);

            return (
                <>
                    <h1 className="pb-2">Tic Tac Toe</h1>
                    {this.renderPlayMode()}
                    <TicTacToeBoard
                        board={this.state.board}
                        onPlayMove={cb}
                        message={message}
                    />
                </>
            );
        } else {
            return (
                <div>
                    <h1>LOADING</h1>
                    <Spinner animation="border" size="sm" />
                </div>
            );
        }
    }

    renderPlayMode() {
        return (
            <div
                className="pb-4 col text-left"
                style={{ display: "inline-block", width: "auto" }}
            >
                <Form onChange={this.onChangeMode.bind(this)}>
                    <Form.Check
                        name="mode"
                        type="radio"
                        id="pvp"
                        label="Player v Player"
                    />
                    <Form.Check
                        name="mode"
                        type="radio"
                        id="pvc"
                        label="Player v Computer"
                    />
                    <Form.Check
                        name="mode"
                        type="radio"
                        id="cvp"
                        label="Computer v Player"
                    />
                    <Form.Check
                        name="mode"
                        type="radio"
                        id="cvc"
                        label="Computer v Computer"
                    />
                </Form>
            </div>
        );
    }

    onChangeMode(e) {
        this.setState({ ...this.state, mode: e.target.id, selected: null });
        if (this.isComputerMoveNext(e.target.id)) {
            this.onPlayAi();
        }
    }

    getMessage(isComputerNext) {
        if (this.state.isDone) {
            const msg = this.state.winner
                ? `Winner Player ${this.state.winner.substr(1)}!`
                : "Tied";
            return (
                <div>
                    <p>{msg}</p>
                    <Link
                        onClick={this.onRestart.bind(this)}
                        to="/game/tictactoe/-"
                    >
                        Restart
                    </Link>
                </div>
            );
        } else {
            const isComputerSuffix = isComputerNext ? " (computer)" : "";
            return `Player ${this.state.nextPlayer.substr(
                1
            )}${isComputerSuffix}`;
        }
    }

    isComputerMoveNext(mode) {
        if (this.state.loading || this.state.isDone) {
            return false;
        }

        if (!mode) {
            mode = this.state.mode;
        }

        const currentPlayer = parseInt(this.state.nextPlayer.substr(1));
        const index = 2 * (currentPlayer - 1);
        const playerMode = mode.substr(index, index + 1);
        console.log(
            `playerMode ${index} ${this.state.mode} ${currentPlayer} ${playerMode}`
        );
        return playerMode === "c";
    }

    async onPlayMove(index) {
        console.log("onPlayMove");
        const newUrl = this.getUrlForNextMove(index);
        this.props.history.push(newUrl);

        const newMoves = this.getPlayedMoves(newUrl);
        await this.onLoadBoard(newMoves);
    }

    getUrlForNextMove(index) {
        var url = this.props.location.pathname;
        if (url.endsWith("/-")) {
            url = url.substr(0, url.length - 1);
        }

        return url + index;
    }
}

class TicTacToeBoard extends React.Component {
    ICON_SIZE = 50;
    MAP_VALUE_TO_ICON = {
        " ": SquareFill,
        X: CircleFill,
        O: Circle,
        S: Square,
    };

    state = {
        selected: null,
    };

    render() {
        return (
            <Container>
                <h5>{this.props.message}</h5>
                <Table bordered className="mx-auto" style={{ width: "1%" }}>
                    <tbody>
                        {this.renderRow(0)}
                        {this.renderRow(1)}
                        {this.renderRow(2)}
                    </tbody>
                </Table>
            </Container>
        );
    }

    renderRow(row) {
        return (
            <tr>
                {this.renderSquare(3 * row)}
                {this.renderSquare(3 * row + 1)}
                {this.renderSquare(3 * row + 2)}
            </tr>
        );
    }

    renderSquare(index) {
        const key =
            this.state.selected === index ? "S" : this.props.board[index];

        const TicTacToeSquare = this.MAP_VALUE_TO_ICON[key];

        if (key === " ") {
            const onClick = this.props.onPlayMove
                ? this.onClickSquare.bind(this, index)
                : null;

            return (
                <td onClick={onClick}>
                    <TicTacToeSquare
                        style={{ visibility: "hidden" }}
                        size={this.ICON_SIZE}
                        color="yellow"
                    />
                </td>
            );
        }

        const color = key === "S" ? "yellow" : "black";

        return (
            <td onClick={this.onClickSquare.bind(this, index)}>
                <TicTacToeSquare size={this.ICON_SIZE} color={color} />
            </td>
        );
    }

    onClickSquare(index) {
        if (this.state.selected === index) {
            this.props.onPlayMove(index);
            index = undefined;
        }

        if (this.props.board[index] !== " ") {
            index = undefined;
        }

        this.setState({ ...this.state, selected: index });
    }
}
