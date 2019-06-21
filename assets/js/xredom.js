import {mount, el, list} from 'redom';


// define Login component
class Login {
    constructor() {
        this.el = el('form#login',
            this.email = el('input.email', {type: 'email'}),
            this.pass = el('input.pass', {type: 'password'}),
            this.submit = el('button', {type: 'submit'},
                'Sign in'
            )
        );
        this.el.onsubmit = e => {
            e.preventDefault();

            console.log({
                email: this.email.value,
                pass: this.pass.value
            });
        }
    }
}

// list
class Li {
    constructor() {
        this.el = el('li');
    }

    update({_id, name}) {
        this.el.textContent = `Item ${name}`;
    }
}

// define <td> component
class Td {
    constructor () {
        this.el = el('td');
    }
    update (i) {
        this.el.textContent = i;
    }
}
const Tr = list.extend('tr', Td);


(() => {
    let a = el('h1', 'Hello RE:DOM!');
    mount(document.getElementById("title"), a);

    // create login
    const login = new Login();
    // mount to DOM
    mount(document.getElementById("login"), login);

    // create list
    const ul = list('ul', Li, '_id');

    // update with data
    ul.update([
        {_id: 1, name: "abc"},
        {_id: 2, name: "great"},
        {_id: 3, name: "wonderful"}
        ]);

    // mount to DOM
    mount(document.getElementById("list1"), ul);

    // schedule another update
    setTimeout(() => {
        ul.update([
            {_id: 2, name: "wall"},
            {_id:5, name: "new Item"}
        ]);
    }, 1000);


    // create table
    const table = list('table', Tr);

    // update with data
    table.update([
        [ 1, 2, 3 ],
        [ 4, 5, 6 ],
        [ 7, 8, 9 ]
    ]);

    // mount to DOM
    mount(document.getElementById("table1"), table);

    // schedule another update
    setTimeout(() => {
        table.update([
            [ 3, 2, 1 ],
            [ 6, 5, 4 ],
            [ 9, 8, 7 ]
        ]);
    }, 1000);
})();


