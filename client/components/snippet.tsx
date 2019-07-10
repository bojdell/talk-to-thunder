import React from "react";

import ShareButtons from "./share_buttons";

import { mutate } from "thunder-react";

import "./snippet.css";

// interface Highlight {
//   startIdx: number;
//   endIdx: number;
// }

interface Props {
  id: number;
  title: string;
  seedText: string;
  generatedText: string;
  isDeleted?: boolean;
}

interface State {
  selectedText?: string;
}

const newlinesToBr = (text?: string) =>
  text
    ? text.split("\n").map((item, idx) => {
        return (
          <span key={idx} className="Snippet-text">
            {item}
            <br />
          </span>
        );
      })
    : "";

class Snippet extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  handleDeleteClick = () => {
    mutate({
      query: "{ deleteSnippet(id: $id) }",
      variables: {
        id: this.props.id
      }
    });
  };

  renderDeleteButton = () => {
    return (
      <a href="/" onClick={this.handleDeleteClick} className="u-marginMd">
        Delete
      </a>
    );
  };

  handleUndeleteClick = () => {
    mutate({
      // Challenge 2: Implement the mutation to restore a deleted snippet.
      query: ``,
      variables: {
        id: this.props.id
      }
    });
  };

  renderUndeleteButton = () => {
    return (
      <a href="/" onClick={this.handleUndeleteClick} className="u-marginMd">
        Restore
      </a>
    );
  };

  handleMouseMove = () => {
    const selectedText = window.getSelection();
    if (selectedText) {
      this.setState({ selectedText: selectedText.toString() });
    } else {
      this.setState({ selectedText: undefined });
    }
  };

  renderButtons = (shareText: string) => {
    if (this.props.isDeleted) {
      return (
        <div className="u-flexAlignItemsCenter">
          <div className="u-flex1" />
          {this.renderUndeleteButton()}
        </div>
      );
    }

    // TODO: why aren't fragments working? :thinking:
    return (
      <div className="u-flexAlignItemsCenter">
        <div className="u-flex1" />
        {this.renderDeleteButton()}
        <ShareButtons
          textContent={shareText}
          isSelection={Boolean(this.state.selectedText)}
        />
      </div>
    );
  };

  render() {
    const snippetText = this.props.generatedText || this.props.seedText;
    const shareText = this.state.selectedText || snippetText;
    return (
      <div className="Snippet">
        <h1>{this.props.title}</h1>
        <div onMouseMove={this.handleMouseMove}>
          {newlinesToBr(snippetText)}
          {this.renderButtons(shareText)}
        </div>
      </div>
    );
  }
}

export default Snippet;
