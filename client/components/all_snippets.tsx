import React from "react";

import Snippet from "./snippet";

import { connectGraphQL } from "thunder-react";

type Props = any; // TODO: generate this or specify

class AllSnippets extends React.Component<Props, {}> {
  render() {
    return (
      <div>
        {this.props.data.value.allSnippets.map((s: any) => (
          <div>
            <h1>{s.id}</h1>
            <Snippet generatedText={s.generatedText} seedText={s.seedText} />
          </div>
        ))}
      </div>
    );
  }
}

export default connectGraphQL(AllSnippets, () => ({
  query: `
  {
    allSnippets {
      id
      seedText
      generatedText
    }
  }`,
  variables: {},
  onlyValidData: true
}));
