import React from 'react';

import { connectGraphQL } from 'thunder-react';
import { GraphiQLWithFetcher } from './graphiql';

import './App.css';

// TODO
// type State interface {
// }

type Props = any

class CreateSnippetPage extends React.Component<Props, {}> {
  render() {
    return (<div className="App">
      {this.props.data.value.allSnippets.map((s: any) => <div>
        <h1>{s.id}</h1>
          {s.generatedText}
        </div>)}
    </div>)
  }
}

const ConnectedCreateSnippet = connectGraphQL(CreateSnippetPage, () => ({
  query: `
  {
    allSnippets {
      id
      generatedText
    }
  }`,
  variables: {
    // id: props.id
  },
  onlyValidData: true,
}));

const App: React.FC = () => {
  const pathname = window.location.pathname;
  if (pathname === "/graphiql") {
    return <GraphiQLWithFetcher />;
  } else {
    return <ConnectedCreateSnippet />
  }
}

export default App;
