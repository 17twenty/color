# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs for producing colored output.

__note: things may change__

## Install
```
go get github.com/nhooyr/color
```

## Examples
See [godoc](https://godoc.org/github.com/nhooyr/color) for more information.

### Setting Attributes
```go
// "panic:" with a maroon foreground then normal "rip".
color.Printf("%h[fgMaroon]panic:%r %s\n", "rip")

// "panic:" with a red background then normal "rip".
color.Printf("%h[bgRed]panic:%r %s\n", "rip")

// Bold "panic:" then normal "rip".
color.Printf("%h[bold]panic:%r %s\n", "rip")

// Underlined "panic:" with then normal "rip".
color.Printf("%h[underline]panic:%r %s\n", "rip")

// "panic:" using color 83 as the foreground then normal "rip".
color.Printf("%h[fg83]panic:%r %s\n", "rip")

// "panic:" using color 158 as the background then normal "rip".
color.Printf("%h[bg158]panic:%r %s\n", "rip")
```

### Mixing Attributes
```go
// Bolded "panic:" with a green foreground then normal "rip".
color.Printf("%h[fgGreen+bold]panic:%r %s\n", "rip")

// Underlined "panic:" with a gray background then normal "rip".
color.Printf("%h[bg8+underline]panic:%r %s\n", "rip")
```

### Preparing Strings
```go
// Prepare only processes the highlight verbs in the string,
// letting you print it repeatedly with performance.
panicFormat := color.Prepare("%h[fgMaroon+bold]panic:%r %s\n")

// If os.Stdout is a terminal, each will print a bolded "panic:" in red foreground
// and some normal text after. Otherwise each will print normally.
color.Printfp(panicFormat, "rip")
color.Printfp(panicFormat, "yippie")
color.Printfp(panicFormat, "dsda")
```

### Printer
A `Printer` writes to an `io.Writer`.

```go
redFormat := color.Prepare("%h[fgMaroon]%s%r\n")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi".
p = color.New(os.Stderr, color.IsTerminal(os.Stderr))
p.Printfp(redFormat, "hi")

// "hi" with red foreground.
p := color.New(os.Stderr, true)
p.Printfp(redFormat, "hi")

// normal "hi", the highlight verbs are ignored.
p = color.New(os.Stderr, false)
p.Printfp(redFormat, "hi")
```

### Logger
Import `github.com/nhooyr/color/log` to use it.

```go
redFormat := color.Prepare("%h[fgMaroon]%s%r\n")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi".
log.Printfp(redFormat, "hi")

// normal "hi", the highlight verbs are ignored.
log.SetColor(false)
log.Printfp(redFormat, "hi")

// "hi" with a red foreground.
log.SetColor(true)
log.Fatalfp(redFormat, "hi")
```

### How does reset behave?
```go
// "rip" will be printed with a navy foreground and gray background
// because we never reset the highlighting after "panic:". The navy foreground is
// carried on from "panic:".
color.Printf("%h[fgNavy+bgGray]panic: %h[bg8]%s\n", "rip")

// The attributes carry onto anything written to the terminal until reset.
// This prints "rip" in the same attributes as above.
fmt.Println("rip")

// Resets the highlighting and then prints "hello" normally.
color.Printf("%r%s", "hello")
```

## TODO
- [ ] True color support (needs work in my terminfo package)
- [ ] Windows support
- [x] Respect $TERM
- [x] Seperate log package
- [x] color.Format docs
