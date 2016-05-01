# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color adds verbs to `fmt.Printf` for producing colored output.

__note: if you have any ideas for a better prepare API, please let me know__

## Install
```
go get github.com/nhooyr/color
```

## Examples
See [godoc](https://godoc.org/github.com/nhooyr/color) for more information.

### Setting Attributes
```go
// "panic:" with a maroon foreground then normal "rip".
color.Printfh("%h[fgMaroon]panic:%r %s\n", "rip")

// "panic:" with a red background then normal "rip".
color.Printfh("%h[bgRed]panic:%r %s\n", "rip")

// Bold "panic:" then normal "rip".
color.Printfh("%h[bold]panic:%r %s\n", "rip")

// Underlined "panic:" with then normal "rip".
color.Printfh("%h[underline]panic:%r %s\n", "rip")

// "panic:" using color 83 as the foreground then normal "rip".
color.Printfh("%h[fg83]panic:%r %s\n", "rip")

// "panic:" using color 158 as the background then normal "rip".
color.Printfh("%h[bg158]panic:%r %s\n", "rip")
```

### Mixing Attributes
```go
// Bolded "panic:" with a green foreground then normal "rip".
color.Printfh("%h[fgGreen+bold]panic:%r %s\n", "rip")

// Underlined "panic:" with a gray background then normal "rip".
color.Printfh("%h[bg8+underline]panic:%r %s\n", "rip")
```

### Preparing Strings
```go
// Prepare only processes the highlight verbs in the string,
// letting you print it repeatedly with performance.
panicFormat := color.Prepare("%h[fgMaroon+bold]panic:%r %s\n")

// Each prints a bolded "panic:" in red foreground and some normal text after.
// Notice that fmt.Printf is used, this works because only the highlight verbs
// were processed above, the %s verb was not.
fmt.Printf(panicFormat, "rip")
fmt.Printf(panicFormat, "yippie")
fmt.Printf(panicFormat, "dsda").Printf(panicFormat, "dsda")
```

### Printer
A `Printer` wraps around an `io.Writer`, but unlike `color.Fprintf`, it gives you full control over whether color output is enabled.

```go
// "hi" with red foreground.
p := color.NewPrinter(os.Stderr, true)
// See the Prepare example for an explanation of this.
redFormat := p.Prepare("%h[fgMaroon]%s%r\n")
p.Printf(redFormat, "hi")

// normal "hi", the highlight verbs are ignored.
p = color.NewPrinter(os.Stderr, false)
p.Printfh("%h[fgMaroon]%s%r\n", "hi")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi".
p = color.NewPrinter(os.Stderr, color.IsTerminal(os.Stderr))
p.Printfh("%h[fgMaroon]%s%r\n", "hi")
```

### `*log.Logger` wrapper
```go
// "hi" with a red foreground.
l := color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, true)
// See the Prepare example for an explanation of this.
redFormat := l.Prepare("%h[fgMaroon]%s%r\n")
l.Printf(redFormat, "hi")

// normal "hi", the highlight verbs are ignored.
l = color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, false)
l.Printfh("%h[fgMaroon]%s%r", "hi")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi".
l = color.NewLogger(os.Stderr, "%h[bold]color:%r ", log.LstdFlags, color.IsTerminal(os.Stderr))
l.Fatalf("%h[fgMaroon]%s%r", "hi")
```

### How does reset behave?
```go
// "rip" will be printed with a navy foreground and gray background
// because we never reset the highlighting after "panic:". The navy foreground is
// carried on from "panic:".
color.Printfh("%h[fgNavy+bgGray]panic: %h[bg8]%s\n", "rip")

// The attributes carry onto anything written to the terminal until reset.
// This prints "rip" in the same attributes as above.
fmt.Println("rip")

// Resets the highlighting and then prints "hello" normally.
color.Printfh("%r%s", "hello")
```

## TODO
- [ ] True color support
- [ ] Windows support
- [x] Respect $TERM
- [ ] Fully wrap \*log.Logger, perhaps a format string that defines the prefix, date, content etc. Perhaps another package?
