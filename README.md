# Lox Interpreter
Interpreter for the Lox language [Crafting Interpreters](https://craftinginterpreters.com) written in Go.

---

Differences from the book:
- Go like syntax for if/else: no required parentheses for condition, but thenBranch and elseBranch must be blocks (requires braces).
- Go like syntax for 'for' loop: no required parentheses; loop body must be a block.
