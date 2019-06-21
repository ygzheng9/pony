import React from "react";
import ReactDOM from "react-dom"
// "presets": ["@babel/preset-env", "@babel/preset-react"]


import { HookUI }  from "../src/reactDemo";

import Game  from "../src/board";
ReactDOM.render(<Game />, document.getElementById("game"));


ReactDOM.render(<HookUI />, document.getElementById("idHookui"));

