import React from "react";

import AllSnippets from "./all_snippets";
import CreateSnippet from "./create_snippet";

import "./home.css";

const Home: React.FC = () => (
  <div>
    <div className="Home-createSnippet">
      <CreateSnippet />
    </div>
    <AllSnippets />
  </div>
);

export default Home;
