import React from "react";

import { GraphiQLWithFetcher } from "./graphiql";

import Home from "./components/home";
import Trash from "./components/trash";
import { withMenu } from "./components/menu";

import "./App.css";
import "./utils.css";

// TODO
// type State interface {
// }

const App: React.FC = () => {
  const route = window.location.pathname;
  switch (route) {
    case "/graphiql":
      return withMenu(route, <GraphiQLWithFetcher />);
    case "/trash":
      return withMenu(route, <Trash />);
    default:
      return withMenu("/", <Home />);
  }
};

export default App;
