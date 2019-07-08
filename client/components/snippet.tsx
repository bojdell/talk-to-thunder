import React from "react";
import classNames from "classnames";

import { mutate } from "thunder-react";

import "./snippet.css";

// interface Highlight {
//   startIdx: number;
//   endIdx: number;
// }

interface Props {
  seedText: string;
  generatedText: string;
}

interface State {
  hoverIdx?: number;
  highlightStartIdx?: number;
  highlightEndIdx?: number;
}

class Snippet extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {};
  }

  handleMouseEnter = (idx: number) => () => {
    this.setState({ hoverIdx: idx });
  };

  handleMouseLeave = () => {
    this.setState({ hoverIdx: undefined });
  };

  handleClick = (idx: number) => () => {
    if (!this.state.highlightStartIdx) {
      this.setState({ highlightStartIdx: idx });
    } else {
      if (!this.state.highlightEndIdx) {
        this.setState({ highlightEndIdx: idx });
        // TODO: mutation?
      } else {
        this.setState({
          highlightStartIdx: undefined,
          highlightEndIdx: undefined
        });
      }
    }
  };

  // handleClick = () => {
  //   mutate({
  //     query: '{ createSnippet(text: $text) }',
  //     variables: { text: this.state.inputText },
  //   });
  // }

  render() {
    const text = this.props.generatedText || this.props.seedText;
    const tokenElems: JSX.Element[] = [];
    const lines = text.split("\n");
    let highlightedText = "";
    let idx = 0;
    lines.forEach(line => {
      const tokens = line.split(" ");
      const numSeedTokens = this.props.seedText.split(" ").length;
      let continues = false;
      tokens.forEach(token => {
        const { hoverIdx, highlightStartIdx, highlightEndIdx } = this.state;
        // TODO: add a test for this.
        const endIdx = highlightEndIdx || hoverIdx;
        let shouldHighlight =
          idx === hoverIdx ||
          idx === highlightStartIdx ||
          Boolean(
            highlightStartIdx &&
              endIdx &&
              Math.min(highlightStartIdx, endIdx) <= idx &&
              Math.max(highlightStartIdx, endIdx) >= idx
          );
        const highlightTrailingSpace =
          shouldHighlight &&
          highlightStartIdx &&
          endIdx &&
          idx !== Math.max(highlightStartIdx, endIdx);
        if (shouldHighlight) {
          highlightedText += token;
        }
        if (highlightTrailingSpace) {
          highlightedText += " ";
        }
        const roundLeft =
          (idx === hoverIdx && !highlightStartIdx) ||
          (highlightStartIdx &&
            (!endIdx || idx === Math.min(highlightStartIdx, endIdx)));
        const roundRight =
          (idx === hoverIdx && !highlightStartIdx) ||
          (highlightStartIdx &&
            (!endIdx || idx === Math.max(highlightStartIdx, endIdx)));
        continues = shouldHighlight && !roundRight;
        tokenElems.push(
          <span>
            {/* <span
              className="Snippet-highlightedToken"
              style={{
                position: "relative"
              }}
            /> */}
            <span
              key={idx}
              onMouseEnter={this.handleMouseEnter(idx)}
              onMouseLeave={this.handleMouseLeave}
              onClick={this.handleClick(idx)}
              className={classNames({
                "Snippet-highlightedToken": shouldHighlight,
                "Snippet-token": !shouldHighlight,
                "Snippet-borderRadiusLeft": roundLeft,
                "Snippet-borderRadiusRight": roundRight,
                "u-textBold": idx < numSeedTokens
              })}
            >
              {token}
              {highlightTrailingSpace ? " " : ""}
            </span>
            {highlightTrailingSpace ? "" : " "}
          </span>
        );
        idx++;
      });
      if (continues) {
        highlightedText += "\n";
      }
      tokenElems.push(<br />);
    });
    highlightedText = `"${highlightedText}" #talktothunder`;
    const encodedText = encodeURIComponent(highlightedText);
    return (
      <div>
        {tokenElems}
        <a
          href={`https://twitter.com/intent/tweet?text=${encodedText}`}
          className="twitter-share-button"
          data-hashtags="talktothunder"
          data-show-count="false"
        >
          Tweet
        </a>
        <script async src="https://platform.twitter.com/widgets.js" />
      </div>
    );
  }
}

export default Snippet;
