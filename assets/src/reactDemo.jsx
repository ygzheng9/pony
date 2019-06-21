import React, { Component, useState, useEffect } from "react"

const ElemTitle = () => (<h1>Hello, Hook API, wowwowowo</h1>);

const App = (props) => {
    return (
        <div>
            <h5>params: {props.a} - {props.b}</h5>
        </div>
    )
};


class Button extends Component {
    state = { buttonText: "Click me, please" };

    handleClick = () => {
        this.setState(() => {
            return { buttonText: "Thanks, been clicked!" };
        });
    };

    render() {
        const { buttonText } = this.state;
        return <button onClick={this.handleClick}>{buttonText}</button>;
    }
}


function ButtonFunc() {
    const [buttonText, setButtonText] = useState("Click me, please");

    const btnClick = () => {
        setButtonText("Thanks, been clicked!")
    };

    return (
        <button onClick={btnClick}>
            {buttonText}
        </button>
    );
}

function Bubble() {
    const [items, setItems] = useState([{}]);

    useEffect(() => {
        fetch("/assets/data/bubble.json")
            .then(response => response.json())
            .then(data => setItems(data));

        // console.log(items);
    }, []);

    return (
        <div>
            <ul>
                {items.map(el => (
                    <li key={el.Country}>{el.Country}: {el.GDP} - {el.Population}</li>
                ))}
            </ul>
        </div>
    );
}

function Counter() {
    // Declare a new state variable, which we'll call "count"
    const [count, setCount] = useState(0);

    useEffect(() => {
        document.title = `You clicked ${count} times`;
    });

    return (
        <div>
            <p>You clicked {count} times</p>
            <button onClick={() => setCount(count + 1)}>
                Click me
            </button>
        </div>
    );
}

function HookUI() {
    const params = {
        "a": "aaa",
        "b": "bbb",
    };

    return (
        <div>
            <ElemTitle />
            <br />
            <App {...params} />
            <br />
            <Button />
            <ButtonFunc />
            <Counter />
            <br/>
            <Bubble />
        </div>
    );
}

export { HookUI };