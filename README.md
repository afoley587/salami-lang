# Salami-Lang

My first attempt at creating a language with:

* Tokens
* Lexer
* AST
* Parser
* Interpreter

It only works for simple calculations and the only supported 
keyword is `var`. The output of the following salami-code:

```
var x = 5; 
var y = 10; 
var z = x + y; 
var t = z * z;
t;
```

would be `225`.

For a live demo using the [functions.salami](./examples/functions.salami) file, 
see the below video:

![Demo](./imgs/demo.gif)