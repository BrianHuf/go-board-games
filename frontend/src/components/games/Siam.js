import React from "react";
import { Link } from "react-router-dom";
import Container from "react-bootstrap/Container";
import Table from "react-bootstrap/Table";
import Spinner from "react-bootstrap/Spinner";
import Form from "react-bootstrap/Form";
import {
  Square,
  Dot,
  SquareFill,
  CircleFill,
  ArrowDownCircleFill,
  ArrowUpCircleFill,
  ArrowLeftCircleFill,
  ArrowRightCircleFill,
  Circle,
  ArrowDownCircle,
  ArrowUpCircle,
  ArrowLeftCircle,
  ArrowRightCircle,
} from "react-bootstrap-icons";

import rest from "../../api/backend";

export default class Siam extends React.Component {
  state = {
    board: null,
    loading: true,
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

    const info = await rest.getMoves("siam", moves);
    const response = info.data;
    console.log(response);
    this.setState({
      ...this.state,
      board: response.state,
      moves: response.moves,
      isDone: response.state.isDone,
      winner: response.state.winner,
      nextPlayer: response.state.nextPlayer,
      loading: false,
    });

    // if (this.isComputerMoveNext()) {
    //     this.onPlayAi();
    // }
  }

  async onPlayAi() {
    console.log("onPlayAI");
    const moves = this.getPlayedMoves(this.props.location.pathname);
    const aiMoves = await rest.getAiMove("Siam", moves);

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
    if (!this.state.board) {
      return (
        <div>
          <h1>LOADING</h1>
          <Spinner animation="border" size="sm" />
        </div>
      );
    }

    const isComputerNext = this.isComputerMoveNext();

    const cb =
      this.state.loading || this.state.isDone || isComputerNext
        ? null
        : this.onPlayMove.bind(this);

    const message = this.getMessage(isComputerNext);

    return (
      <>
        <h1 className="pb-2">Siam</h1>
        {this.renderPlayMode()}
        <SiamBoard
          board={this.state.board}
          moves={this.state.moves}
          onPlayMove={cb}
          message={message}
        />
      </>
    );
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
          <Link onClick={this.onRestart.bind(this)} to="/game/Siam/-">
            Restart
          </Link>
        </div>
      );
    } else {
      const isComputerSuffix = isComputerNext ? " (computer)" : "";
      return `Player ${this.state.nextPlayer.substr(1)}${isComputerSuffix}`;
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

  async onPlayMove(move) {
    console.log("onPlayMove");
    const newUrl = this.getUrlForNextMove(move);
    this.props.history.push(newUrl);

    const newMoves = this.getPlayedMoves(newUrl);
    await this.onLoadBoard(newMoves);
  }

  getUrlForNextMove(move) {
    var url = this.props.location.pathname;
    if (url.endsWith("/-")) {
      url = url.substr(0, url.length - 1);
    }

    return url + move;
  }
}

class SiamBoard extends React.Component {
  ICON_SIZE = 50;
  MAP_VALUE_TO_ICON = {
    ".": Dot,
    M: SquareFill,
    D: ArrowDownCircleFill,
    U: ArrowUpCircleFill,
    L: ArrowLeftCircleFill,
    R: ArrowRightCircleFill,
    d: ArrowDownCircle,
    u: ArrowUpCircle,
    l: ArrowLeftCircle,
    r: ArrowRightCircle,
  };

  state = {
    selectedSource: null,
    selectedReserveDirection: null,
    selectedTarget: null,
  };

  OFFBOARD = -51;

  onClear() {
    console.log("ON CLEAR");
    this.setState({
      selectedSource: null,
      selectedReserveDirection: null,
      selectedTarget: null,
    });
  }

  render() {
    return (
      <Container>
        <h5>{this.props.message}</h5>
        <Table bordered className="mx-auto" style={{ width: "1%" }}>
          <tbody>
            {this.renderRow(0)}
            {this.renderRow(1)}
            {this.renderRow(2)}
            {this.renderRow(3)}
            {this.renderRow(4)}
          </tbody>
        </Table>
        <Table bordered className="mx-auto" style={{ width: "1%" }}>
          <tbody>
            <tr>
              <td>
                <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
                  {this.renderOffboard(0)}
                </div>
              </td>
              <td>
                <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
                  {this.renderOffboard(1)}
                </div>
              </td>
            </tr>
          </tbody>
        </Table>
      </Container>
    );
  }

  renderOffboard(player) {
    const isCurrentPlayer = this.props.board.nextPlayer === `p${player + 1}`;
    const Icon = player ? Circle : CircleFill;

    const selector = isCurrentPlayer ? this.renderOffboardSelector() : null;

    if (!selector) {
      return (
        <div onClick={this.onClear.bind(this)}>
          <Icon size={this.ICON_SIZE} />
        </div>
      );
    }

    return (
      <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
        <Icon size={this.ICON_SIZE} />
        {selector}
      </div>
    );
  }

  renderRow(row) {
    return (
      <tr>
        {this.renderSquare(5 * row)}
        {this.renderSquare(5 * row + 1)}
        {this.renderSquare(5 * row + 2)}
        {this.renderSquare(5 * row + 3)}
        {this.renderSquare(5 * row + 4)}
      </tr>
    );
  }

  renderSquare(index) {
    const squareState = this.props.board.state[index];
    const SiamSquare = this.MAP_VALUE_TO_ICON[squareState];
    const style = squareState === "." ? { visibility: "hidden" } : {};

    const selector = this.renderSquareSelector(index);

    const onClear = selector ? null : this.onClear.bind(this);

    return (
      <td onClick={onClear}>
        <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
          <SiamSquare style={style} size={this.ICON_SIZE} />
          {selector}
        </div>
      </td>
    );
  }

  renderOffboardSelector() {
    if (this.state.selectedSource === this.OFFBOARD) {
      return this.renderSelectDirection(
        this.OFFBOARD,
        this.state.selectedReserveDirection
      );
    }
    
    if (this.state.selectedSource) {
      return null;
    }

    if (this.state.selectedTarget != null) {
      return null;
    }

    const matches = this.props.moves.filter(
      (m) => m.value.lastMove.charCodeAt(0) < 97
    );
    if (matches.length) {
      return this.renderSelectSquare(
        this.OFFBOARD,
        this.state.selectedSource != null
      );
    }
  }

  renderSquareSelector(index) {
    const matches = this.props.moves
      .map((m) => this.getRenderMatch(m, index))
      .filter((m) => m != null);
    if (matches.length) {
      return matches[0];
    }

    return null;
  }

  getRenderMatch(move, index) {
    if (this.props.onPlayMove == null) {
      return null;
    }

    if (this.state.selectedTarget) {
      if (index === this.state.selectedTarget) {
        return this.renderSelectDirection(index);
      }
    } else if (this.state.selectedSource) {
      if (this.state.selectedSource === this.OFFBOARD) {
        const dir = this.state.selectedReserveDirection;
        if (dir == null) {
          return null;
        }

        if (
          move.value.lastMove[0] === dir &&
          move.value.lastMove.charCodeAt(1) - 97 === index
        ) {
          return this.renderSelectSquare(index, false);
        }
      } else if (
        move.value.lastMove.charCodeAt(0) - 97 === this.state.selectedSource &&
        move.value.lastMove.charCodeAt(1) - 97 === index
      ) {
        return this.renderSelectSquare(index, true);
      }
    } else if (move.value.lastMove.charCodeAt(0) - 97 === index) {
      return this.renderSelectSquare(index, false);
    }
    return null;
  }

  renderSelectSquare(index, alreadySelected) {
    const onClick = (e) => this.onClickSquare(index);
    const selectorColor = alreadySelected ? "red" : "yellow";
    const style = {
      opacity: "70%",
      position: "relative",
      top: -this.ICON_SIZE,
      bottom: -this.ICON_SIZE,
    };

    return (
      <Square
        color={selectorColor}
        style={style}
        size={this.ICON_SIZE}
        onClick={onClick}
      />
    );
  }

  renderSelectDirection(index, selectedDir = null) {
    const style = {
      opacity: "70%",
      position: "relative",
      top: -this.ICON_SIZE,
      bottom: -this.ICON_SIZE,
    };

    const selectorColors = ["yellow", "red"];
    const THICK = 0.3;
    const PAD = 1.0 - THICK - THICK;

    const onClickUp = (e) => this.onClickDirection(index, "U");
    const onClickDown = (e) => this.onClickDirection(index, "D");
    const onClickLeft = (e) => this.onClickDirection(index, "L");
    const onClickRight = (e) => this.onClickDirection(index, "R");

    return (
      <div style={style}>
        <Table>
          <tbody>
            <tr style={{ height: this.ICON_SIZE * THICK }}>
              <td className="p-0" style={{ width: this.ICON_SIZE * THICK }} />
              <td
                className="p-0"
                style={{
                  backgroundColor: selectorColors[selectedDir === "U" ? 1 : 0],
                  width: this.ICON_SIZE * PAD,
                }}
                onClick={onClickUp}
              />
              <td className="p-0" style={{ width: this.ICON_SIZE * THICK }} />
            </tr>
            <tr style={{ height: this.ICON_SIZE * PAD }}>
              <td
                className="p-0"
                style={{
                  backgroundColor: selectorColors[selectedDir === "L" ? 1 : 0],
                  width: this.ICON_SIZE * THICK,
                }}
                onClick={onClickLeft}
              />
              <td className="p-0" style={{ width: this.ICON_SIZE * PAD }} />
              <td
                className="p-0"
                style={{
                  backgroundColor: selectorColors[selectedDir === "R" ? 1 : 0],
                  width: this.ICON_SIZE * THICK,
                }}
                onClick={onClickRight}
              />
            </tr>
            <tr style={{ height: this.ICON_SIZE * THICK }}>
              <td className="p-0" style={{ width: this.ICON_SIZE * THICK }} />
              <td
                className="p-0"
                style={{
                  backgroundColor: selectorColors[selectedDir === "D" ? 1 : 0],
                  width: this.ICON_SIZE * PAD,
                }}
                onClick={onClickDown}
              />
              <td className="p-0" style={{ width: this.ICON_SIZE * THICK }} />
            </tr>
          </tbody>
        </Table>
      </div>
    );
  }

  onClickDirection(index, direction) {
    console.log(
      `onClickDirection ${index} ${direction} ${this.props.board.nextPlayer}`
    );

    if (index === this.OFFBOARD) {
      this.setState({ ...this.state, selectedReserveDirection: direction });
      return;
    }

    const playerDirection =
      this.props.board.nextPlayer === "p1"
        ? direction.toUpperCase()
        : direction.toLowerCase();

    const from =
      this.state.selectedSource === this.OFFBOARD
        ? this.state.selectedReserveDirection
        : String.fromCharCode(97 + this.state.selectedSource);

    const to =
      this.state.selectedTarget === this.OFFBOARD
        ? "."
        : String.fromCharCode(97 + this.state.selectedTarget);

    this.setState({
      selectedSource: null,
      selectedReserveDirection: null,
      selectedTarget: null,
    });

    this.props.onPlayMove(from + to + playerDirection);
  }

  onClickSquare(index) {
    console.log(
      `onClickSquare ${index} ${this.state.selectedSource} ${this.state.selectedTarget}`
    );
    if (this.state.selectedSource === index) {
      this.setState({
        ...this.state,
        selectedSource: index,
        selectedTarget: index,
      });
    } else if (this.state.selectedSource == null) {
      this.setState({ ...this.state, selectedSource: index });
    } else if (index === this.OFFBOARD) {
      this.setState({
        selectedSource: null,
        selectedReserveDirection: null,
        selectedTarget: null,
      });

      const from = String.fromCharCode(97 + this.state.selectedSource);
      this.props.onPlayMove(from + "..");
    } else {
      this.setState({ ...this.state, selectedTarget: index });
    }
  }
}
