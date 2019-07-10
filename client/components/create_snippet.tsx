import React from "react";

import { mutate } from "thunder-react";

import "./create_snippet.css";

type Props = any; // TODO: generate this or specify

interface State {
  inputText?: string;
  finalTokenLength?: number;
  numTokensInputError?: boolean;
}

// const numTokensValueIsValid = (rawValue: string): boolean => {
//   if (!rawValue) {
//     return true // Empty value is valid.
//   }
//   const parsedValue = parseInt(rawValue)
//   if (!parsedValue || parsedValue < 0 )
// }

class CreateSnippet extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  handleTextChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    this.setState({ inputText: (e.target as HTMLTextAreaElement).value });
  };

  handleNumTokensInputChange = (e: React.FormEvent<HTMLTextAreaElement>) => {
    this.setState({ inputText: (e.target as HTMLTextAreaElement).value });
  };

  handleSubmit = () => {
    mutate({
      query:
        "{ createSnippet(text: $text, finalTokenLength: $finalTokenLength) }",
      variables: {
        text: this.state.inputText,
        finalTokenLength: this.state.finalTokenLength
      }
    });
  };

  handleKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
    // KeyCode for Enter is 13
    if (event.keyCode === 13 && (event.ctrlKey || event.metaKey)) {
      this.handleSubmit();
    }
  };

  render() {
    return (
      <div className="CreateSnippet">
        <textarea
          className="CreateSnippet-textInput"
          maxLength={850}
          autoFocus={true}
          placeholder="Type something and a neural network will guess what comes next.&#10;Press Cmd+Enter or Ctrl+Enter to submit."
          rows={4}
          onChange={this.handleTextChange}
          onKeyDown={this.handleKeyDown}
        />
      </div>
    );
  }
}

export default CreateSnippet;
