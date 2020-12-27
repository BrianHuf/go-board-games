import React from "react";
import { Router, Route, Link } from "react-router-dom";

import history from "./History";
import TicTacToe from "./games/TicTacToe";
import Siam from "./games/Siam";
import Welcome from "./Welcome";
import "./App.css";

const App = () => {
  return (
    <Router history={history}>
      <div>
        <Link to="/">
          <p className="p-2">Home</p>
        </Link>
        <div className="App">
          <Route exact path="/" component={Welcome} />
          <Route
            exact
            path="/game/tictactoe/:playedMoves"
            component={TicTacToe}
          />
          <Route
            exact
            path="/game/siam/:playedMoves"
            component={Siam}
          />          
        </div>
      </div>
    </Router>
  );
};

export default App;
