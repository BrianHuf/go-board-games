import React from "react";
import Container from "react-bootstrap/Container";
import Table from "react-bootstrap/Table";
import Spinner from "react-bootstrap/Spinner";
import { Circle, CircleFill, Square, SquareFill } from "react-bootstrap-icons";

import rest from "../../api/backend";

export default class TicTacToe extends React.Component {
    state = {
        board: null,
        loading: false,
        isDone: false,
        winner: null,
    };

    componentDidMount() {
        // take care of browser backward/forward
        window.onpopstate = () => {
            const moves = this.getPlayedMoves(this.props.location.pathname);
            this.loadBoard(moves);
        };
    }

    getPlayedMoves(url) {
        return url.substr(url.lastIndexOf("/") + 1);
    }

    async loadBoard(moves) {
        const info = await rest.getMoves("tictactoe", moves);
        const response = info.data.state;
        console.log(response);
        this.setState({
            board: response.state,
            isDone: response.isDone,
            winner: response.winner,
            nextPlayer: response.nextPlayer,
        });
    }

    render() {
        if (this.state.board) {
            const cb =
                this.state.loading || this.state.isDone
                    ? null
                    : this.onPlayMove.bind(this);

            const message = this.getMessage();
            return (
                <TicTacToeBoard
                    board={this.state.board}
                    onPlayMove={cb}
                    message={message}
                />
            );
        } else {
            const moves = this.getPlayedMoves(this.props.location.pathname);
            this.loadBoard(moves);
            return (
                <div>
                    <h1>LOADING</h1>
                    <Spinner animation="border" size="sm" />
                </div>
            );
        }
    }

    getMessage() {
        if (this.state.isDone) {
            if (this.state.winner) {
                return `Winner Player ${this.state.winner.substr(1)}!`;
            }
            return "Tied";
        } else {
            return `Player ${this.state.nextPlayer.substr(1)}`;
        }
    }

    async onPlayMove(index) {
        this.setState({ ...this.state, loading: true });

        const newUrl = this.getUrlForNextMove(index);
        this.props.history.push(newUrl);

        const newMoves = this.getPlayedMoves(newUrl);
        await this.loadBoard(newMoves);

        this.setState({ ...this.state, loading: false });
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
        X: Circle,
        O: CircleFill,
        S: Square,
    };

    state = {
        selected: null,
    };

    render() {
        return (
            <Container>
                <h1 className="pb-2">Tic Tac Toe</h1>
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
