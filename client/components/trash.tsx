import React from "react";

import Snippet from "./snippet";

import { connectGraphQL, mutate } from "thunder-react";

type Props = any; // TODO: generate this or specify

class Trash extends React.Component<Props, {}> {
  handleSubmit = () => {
    mutate({
      query: "{ emptyTrash }"
    });
  };

  renderEmptyTrashButton = () => {
    return (
      <a href="/" onClick={this.handleSubmit} className="u-marginMd">
        Empty Trash
      </a>
    );
  };

  render() {
    return (
      <div className="u-marginTopLg">
        <div className="u-flexJustifyContentCenter">
          {this.renderEmptyTrashButton()}
        </div>
        {this.props.data.value.deletedSnippets.map((s: any) => (
          <Snippet
            id={s.id}
            title={s.id}
            generatedText={s.generatedText}
            seedText={s.seedText}
            isDeleted
          />
        ))}
      </div>
    );
  }
}

export default connectGraphQL(Trash, () => ({
  query: `
  {
    deletedSnippets {
      id
      seedText
      generatedText
    }
  }`,
  variables: {},
  onlyValidData: true
}));
