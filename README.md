# Building A Programming Language From Scratch
## The Quentin Tarantino Inspired Salami-Lang

### Overview

Over the past few weeks, I have been trying to find some good documentation
on how to build an interpreted language. There's a lot out there and it
was a bit overwhelming. There's a lot of in-the-weed docs that were amazing
once you were ready for them, but I found it hard to get from the high-level
to the low-level. This will be the main point of this project - to talk about
what components are critical for an interpreted language, then guide ourselves
into the nitty-gritty.

Also - as an FYI - I won't be going through each line of code in detail in this
project. However, I empower you to go to my [GitHub Repo]() and check out all of
the code that you would like. Feel free to open issues, ask me questions, add
functionality, etc.

So, what are we building? We'll we are going to build an interpreted language.
It's called salami, but it resembles a lot of code you've seen in other 
languages (with a twist). Let's look at the below code snippet:

```shell
gorlami sq(a) {
    dicocco a * a;
}

gorlami add(a, b) {
    if (a < b) {
        dicocco sq(a) + b;
    }
    dicocco a + b;
}
    

var x = 5;
var y = 10;
exit add(x, y);
```

If we ran the above, we would expect to see a return value of 35.

We use the `gorlami` keyword to define a function, much like golang's `fn` 
keyword or python's `function` keyword. We then return values from functions
with `dicocco`. If you aren't familiar with Gorlami or Dicocco, do yourself a 
favor and watch [this clip](https://www.youtube.com/watch?v=krtnt191Drg).

We will define:
* variables with the `var` keyword
* functions with the `gorlami` keyword
* return values with the `dicocco` keyword
* exit codes with the `exit` keyword

We will also support greater than, less than, assignment, and mathematical 
operations.

For a live demo using the [functions.salami](./examples/functions.salami) file, 
see the below video:

![Demo](./imgs/demo.gif)


### Components Of An Interpreted Programming Language

There are a few components of an interpreted programming language.

1. Source Code: This is the human-readable code written by the programmer 
    in the interpreted language.
2. Tokens: Tokens are the indiviual characters or set of characters from the
    source code. They comprise the operations, keywords, identifiers, etc. that
    your language is going to support.
3. Lexer: The lexer takes the source code as input and breaks it down into tokens. 
    For example, the lexer would take something like `var x = 3;` and break it into 
    the tokens: `var`, `x`, `=`, `3`, `;`.
4. Abstract Syntax Tree:
5. Parser:
6. Interpreter:
7. Environment:
8. Runtime:


### Example Of An Interpreted Program

### Tokens

### Lexer

### AST

### Parser

### Interpreter

### Environment

### Runtime