## 🧠 `brainfk`

A [brainfuck](https://brainfuck.org/) interpreter written in Go in a couple of hours because I wanted to learn Go.

### Usage

```shell
$ go run . <path-to-file>
```

### Implementation

There are two main parts, the [lexer](./lexer.go) and the [instruction generator](./instruction.go).

The lexer takes the input file as bytes, and when `Next()` is called, it returns the next available brainfuck
character (`+-<>[].,`), or returns 0 if the EOF has been reached.

The instruction generator takes the input from the lexer and generates instruction values for them. While doing that, 
it optimises repeating instructions like (`+++++`) into a single instruction of `{ instructionType: Inc, value: 5 }`. 
Once all instructions have been parsed, it then does one last iteration to resolve any loop references, i.e. a 
`JumpIfZero` must have a matching `JumpIfNotZero`, and vice versa.

Then, these instructions can be trivially interpreted, as seen in the [main](./main.go) file.

### Performance

Using the [bsort.b](./examples/bsort.b) example from [brainfuck.org](https://brainfuck.org/):

```shell
$ hyperfine --warmup 100 -N ./brainfk < ./examples/bsort.input
Benchmark 1: ./brainfk
  Time (mean ± σ):       1.7 ms ±   0.3 ms    [User: 0.7 ms, System: 0.7 ms]
  Range (min … max):     1.5 ms …   9.6 ms    1678 runs
```

Compared with a [bsort implementation that I threw together in 3 minutes in Go](./examples/bsort.go):

```shell
$ hyperfine --warmup 100 -N ./bsort
Benchmark 1: ./bsort
  Time (mean ± σ):       1.8 ms ±   0.2 ms    [User: 0.8 ms, System: 0.7 ms]
  Range (min … max):     1.6 ms …   4.7 ms    1766 runs
```

Both of these programs were ran with the same [input](./examples/bsort.input) and they both produced the same output.

These numbers don't really mean anything, but it at least shows that my interpreter is not stupidly slow I guess.

### License

This project is licensed under the [MIT](https://choosealicense.com/licenses/mit/) license.

Do note that the files under `examples` ending in `.b` are not mine, and are not covered by this license.
