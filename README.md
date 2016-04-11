# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color wraps `fmt.Printf` with verbs for producing colored output.

__note: this is still a WIP and things may change__

## Install
```
go get github.com/nhooyr/color
```

## Examples
See [godoc](https://godoc.org/github.com/nhooyr/color) for more information.

### Setting Attributes
```go
// "panic:" with a red foreground then normal "rip".
color.Printf("%h[fgRed]panic:%r %s\n", "rip")

// "panic:" with a brightRed background then normal "rip".
color.Printf("%h[bgBrightRed]panic:%r %s\n", "rip")

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

// Underlined "panic:" with a bright black background then normal "rip".
color.Printf("%h[bg8+underline]panic:%r %s\n", "rip")
```

### Preparing Strings
```go
// Prepare processes the highlight verbs in the string only once,
// letting you print it repeatedly with performance.
f := color.Prepare("%h[fgRed+bold]panic:%r %s\n")

// Each prints bolded "panic:" with a red foreground and some normal text after.
color.Eprintf(f, "rip")
color.Eprintf(f, "yippie")
color.Eprintf(f, "dsda")
```

### Printer
A `Printer` wraps around an `io.Writer`, but unlike `color.Fprintf`, it gives you full control over whether color output is enabled.

```go
// "hi" with red foreground.
p := color.NewPrinter(os.Stderr, color.EnableColor)
p.Printf("%h[fgRed]%s%r\n", "hi")

// normal "hi", the highlight verbs are ignored.
p = color.NewPrinter(os.Stderr, color.DisableColor)
p.Printf("%h[fgRed]%s%r\n", "hi")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi".
p = color.NewPrinter(os.Stderr, color.PerformCheck)
p.Printf("%h[fgRed]%s%r\n", "hi")
```

### `*log.Logger` wrapper
```go
l := color.NewLogger(os.Stderr, "%h[bold]color:%r ", 0)

// "hi" with a red foreground.
l.Printf("%h[fgRed]%s%r", "hi")

// normal "hi", the highlight verbs are ignored.
l.DisableColor()
l.Printf("%h[fgRed]%s%r", "hi")

// "hi" with a red foreground.
l.EnableColor()
l.Fatalf("%h[fgRed]%s%r", "hi")
```

### How does reset behave?
```go
// "rip" will be printed with a blue foreground and bright black background
// because we never reset the highlighting after "panic:". The blue foreground is
// carried on from "panic:".
color.Printf("%h[fgBlue+bgBlack]panic: %h[bg8]%s\n", "rip")

// The attributes carry onto anything written to the terminal until reset.
// This prints "rip" in the same attributes as above.
fmt.Println("rip")

// Resets the highlighting and then prints "hello" normally.
color.Printf("%r%s", "hello")
```

## TODO
- [ ] True color support
- [ ] Windows support
- [ ] Respect $TERM
- [ ] Fully wrap \*log.Logger, perhaps a format string that defines the prefix, date, content etc. Perhaps another package?
