// import preact from 'preact';

// Tell Babel to transform JSX into preact.h() calls:
// /** @jsx preact.h */

// "plugins": [
//     ["@babel/plugin-transform-react-jsx", { "pragma":"h" }]
// ]

// please use react

import { h, render } from 'preact';

render((
    <div id="foo">
        <span>Hello, world!</span>
        <button onClick={ e => alert("hi!") }>Click Me</button>
    </div>
), document.getElementById("title"));