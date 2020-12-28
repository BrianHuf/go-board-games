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
    selectedTarget: null,
  };

  RESERVE = -51;

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
                  {this.renderReserve(0)}
                </div>
              </td>
              <td>
                <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
                  {this.renderReserve(1)}
                </div>
              </td>
            </tr>
          </tbody>
        </Table>
      </Container>
    );
  }

  renderReserve(player) {
    const isCurrentPlayer = this.props.board.nextPlayer === `p${player + 1}`;
    const Icon = player ? Circle : CircleFill;

    const selector = isCurrentPlayer
      ? this.getCellSelector(this.RESERVE)
      : null;
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

    const selector = this.getCellSelector(index);
    return (
      <td>
        <div style={{ width: this.ICON_SIZE, height: this.ICON_SIZE }}>
          <SiamSquare style={style} size={this.ICON_SIZE} />
          {selector}
        </div>
      </td>
    );
  }

  getCellSelector(index) {
    if (this.props.onPlayMove == null) {
        return null;
    }
    
    if (index === this.state.selectedSource) {
      return this.renderSelector(index, true);
    } else if (index === this.state.selectedTarget) {
      return this.renderDirection(index);
    }

    let stage = 0;
    if (this.state.selectedSource == null) {
      stage = 0;
    } else {
      stage = 1;
    }

    const matches = this.props.moves.filter(
      (m) => m.value.lastMove.charCodeAt(stage) - 97 === index
    );
    if (matches.length) {
      return this.renderSelector(index, false);
    }

    return null;
  }

  renderSelector(index, alreadySelected) {
    const onClick = (e) => this.onClickSquare(index);
    const selectorColor = alreadySelected ? "red" : "gray";
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

  renderDirection(index) {
    const style = {
      opacity: "70%",
      position: "relative",
      top: -this.ICON_SIZE,
      bottom: -this.ICON_SIZE,
    };

    const selectorColor = "gray";
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
                  backgroundColor: selectorColor,
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
                  backgroundColor: selectorColor,
                  width: this.ICON_SIZE * THICK,
                }}
                onClick={onClickLeft}
              />
              <td className="p-0" style={{ width: this.ICON_SIZE * PAD }} />
              <td
                className="p-0"
                style={{
                  backgroundColor: selectorColor,
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
                  backgroundColor: selectorColor,
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
    console.log(`onClickDirection ${index} ${direction}`);
    const from =
      this.state.selectedSource === this.RESERVE
        ? "."
        : String.fromCharCode(97 + this.state.selectedSource);
    const to =
      this.state.selectedTarget === this.RESERVE
        ? "."
        : String.fromCharCode(97 + this.state.selectedTarget);

    this.setState({
      selectedSource: null,
      selectedTarget: null,
    });

    this.props.onPlayMove(from + to + direction);
  }

  onClickSquare(index) {
    console.log(`onClickSquare ${index}`);
    if (this.state.selectedSource === index) {
      this.setState({
        ...this.state,
        selectedSource: null,
        selectedTarget: null,
      });
    } else if (this.state.selectedTarget === index) {
      this.setState({ ...this.state, selectedTarget: null });
    } else if (this.state.selectedSource == null) {
      this.setState({ ...this.state, selectedSource: index });
    } else {
      this.setState({ ...this.state, selectedTarget: index });
    }
  }
}
