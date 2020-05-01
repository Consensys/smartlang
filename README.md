# Smart: An Experimental Ethereum Smart Contract Language

**STATUS:** This was a short-term experiment and is currently not being developed/maintained. Many parts of the code are unfinished.

Smart is an experiment in designing a new smart contract language for Ethereum. You can read about the language design and motivation [on the ConsenSys Diligence blog](https://diligence.consensys.net/blog/2020/05/an-experiment-in-designing-a-new-smart-contract-language/).

## Code layout

* Code for lexical analysis is in the `parser/`, `scanner/`, and `token/` directories.
* The abstract syntax tree is implemented in `ast/`.
* `check.go`, `check2.go`, and `checkTypes.go` in `ast/` implement semantic analysis.
* `ast/codegen.go` and `assembler/` implement code generation.
* `compiler/` executes the full compilation stack.
* `grammar.txt` and `lexical.txt` were used to bootstrap the language syntax, but they're not directly used and aren't necessarily fully up-to-date.

## Running tests

`go test`, from the `parser/`, `scanner/`, `assembler/`, or `compiler/` directories will run some automated tests.

## Authors

* [Steve Marx](https://github.com/smarx)
* [Todd Proebsting](https://github.com/proebsting)
