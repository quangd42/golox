# Lox Interpreter
Interpreter for the Lox language from [Crafting Interpreters](https://craftinginterpreters.com), implemented in Go.

---

## Features

- [x] Dynamic typing
- [x] Data types: boolean, numbers, string, nil
- [x] Expressions:
  - [x] Arithmetics 
  - [x] **Concatenate string and number with '+'
  - [x] Comparison and equality
  - [x] Logical operators: and/or
  - [x] **Ternary ( ? : )
- [x] Statements
  - [x] Print statement
  - [x] Expression statement
  - [x] Block statement
- [x] Control flows: if/else, while and for loop
  - [x] **`continue` and `break` with optional label
- [x] Variables
- [x] Functions
   - [x] Closures
   - [ ] Anonymous functions
- [x] Classes
   - [x] Inheritance
   - [ ] Getters & Setters
- [ ] Standard Library

**: Addition to features covered in the book.

Differences in syntax from proposed in the book:
- Go like syntax for if/else: parentheses not required for condition expression, thenBranch and elseBranch must be blocks (requires braces).
- Go like syntax for loops ('while' and 'for'): parentheses not required; loop body must be a block (requires braces).
- Keyword to define a function is `fn`.

Read more about [the Lox Language](https://craftinginterpreters.com/the-lox-language.html).

## Try it out!

Create a fib.lox script file:

```
fn fib(n) {
  if (n <= 1) { return n; }
  return fib(n - 2) + fib(n - 1);
}

print fib;

var start = clock();
print "started at " + start;
for var i = 0; i < 30; i = i + 1 {
	print fib(i);
}
var end = clock();
print "ended after " + (end - start) + " s";
```

Run it with:

```sh
go run github.com/quangd42/golox /path/to/fib.lox
```

Or run the interpreter by itself for a REPL:
```sh
go run github.com/quangd42/golox
```

See `/scripts/tests/` for more script examples.
