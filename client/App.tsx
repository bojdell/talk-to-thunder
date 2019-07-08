import React from "react";

import { connectGraphQL } from "thunder-react";
import { GraphiQLWithFetcher } from "./graphiql";

import AllSnippets from "./components/all_snippets";
import CreateSnippet from "./components/create_snippet";

import "./App.css";
import "./utils.css";

// TODO
// type State interface {
// }

const App: React.FC = () => {
  const pathname = window.location.pathname;
  if (pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else {
    return (
      <div className="App">
        <h1 className="u-textAlignCenter">⚡️ Talk to Thunder ⚡️</h1>
        <CreateSnippet />
        <AllSnippets />
      </div>
    );
  }
};

export default App;
