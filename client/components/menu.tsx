import React from "react";

import "./menu.css";

// TODO: use react-router
type Route = "/" | "/trash" | "/graphiql";

interface Page {
  displayName: string;
  route: Route;
}

export const Pages: Page[] = [
  {
    displayName: "⚡️ Talk to Thunder",
    route: "/"
  },
  {
    displayName: "GraphiQL",
    route: "/graphiql"
  },
  {
    displayName: "Trash",
    route: "/trash"
  }
];

interface Props {
  currentRoute: Route;
}

interface State {
  isOpen?: boolean;
}

export function withMenu(currentRoute: Route, component: JSX.Element) {
  const withAppWrapper = (c: JSX.Element) =>
    currentRoute === "/graphiql" ? (
      c
    ) : (
      <div className="App">
        <div className="App-body">{c}</div>
      </div>
    );

  return (
    <div style={{ height: "100%" }}>
      <Menu currentRoute={currentRoute} />
      {withAppWrapper(
        <div style={{ height: "100%" }}>
          <div style={{ height: "60px" }} />
          {component}
        </div>
      )}
    </div>
  );
}

class Menu extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  render() {
    return (
      <div className="Menu">
        {Pages.map(page => (
          <span className="Menu-item">
            <a href={page.route}>{page.displayName}</a>
          </span>
        ))}
      </div>
    );
  }
}

export default Menu;
