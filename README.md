# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs for producing colored output.

__note: things may change but it looks pretty stable. If you have any new ideas, let me know ASAP__

## Install
```
go get github.com/nhooyr/color
```

## Examples
See [godoc](https://godoc.org/github.com/nhooyr/color) for more information, especially for a better understanding of `color.Prepare`.

### Setting Attributes
```go
// "panic:" with a red foreground then normal "foo".
f := color.Prepare("%h[fgRed]panic:%r %s\n")
color.Printf(f, "foo")

// "panic:" with a red background then normal "bar".
f = color.Prepare("%h[bgRed]panic:%r %s\n")
color.Printf(f, "bar")

// Bold "panic:" then normal "foo".
f = color.Prepare("%h[bold]panic:%r %s\n")
color.Printf(f, "foo")

// Underlined "panic:" with then normal "bar".
f = color.Prepare("%h[underline]panic:%r %s\n")
color.Printf(f, "bar")

// "panic:" using color 83 as the foreground then normal "foo".
f = color.Prepare("%h[fg83]panic:%r %s\n")
color.Printf(f, "foo")

// "panic:" using color 158 as the background then normal "bar".
f = color.Prepare("%h[bg158]panic:%r %s\n")
color.Printf(f, "bar")
```

### Mixing Attributes
```go
// Bolded "panic:" with a green foreground then normal "foo".
f := color.Prepare("%h[fgGreen+bold]panic:%r %s\n")
color.Printf(f, "foo")

// Underlined "panic:" with a bright black background then normal "bar".
f = color.Prepare("%h[bg8+underline]panic:%r %s\n")
color.Printf(f, "bar")
```

### Printing
```go
f := color.Prepare("%h[fgBrightMagenta+underline]panic:%r %s\n")

// There are two methods of printing a color.Format, either as the
// format to a Printf like function, or as one of the variadic
// arguments to any Print like function.

// f is printed and used as the format.
color.Printf(f, "foo")

// f is simply printed.
color.Print(f)
```

### Printer
A `Printer` writes to an `io.Writer`.

```go
f := color.Prepare("%h[fgRed]%s%r\n")

// If standard error is a terminal, this will print in color.
// Otherwise it will print a normal "bar".
p := color.New(os.Stderr, color.IsTerminal(os.Stderr))
p.Printf(f, "bar")

// "foo" with red foreground.
p = color.New(os.Stderr, true)
p.Printf(f, "foo")

// Normal "bar", the highlight verbs are ignored.
p = color.New(os.Stderr, false)
p.Printf(f, "bar")
```

### `github.com/nhooyr/color/log`
```go
f := color.Prepare("%h[fgRed]%s%r\n")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "foo".
log.Printf(f, "foo")

// Normal "bar", the highlight verbs are ignored.
log.SetColor(false)
log.Printf(f, "bar")

// "foo" with a red foreground.
log.SetColor(true)
log.Fatalf(f, "foo")
```

### How does reset behave?
```go
// "hello" will be printed with a black foreground and bright green background
// because we never reset the highlighting after "panic:". The black foreground is
// carried on from "panic:".
f := color.Prepare("%h[fgBlack+bgBrightRed]panic: %h[bgBrightGreen]%s")
color.Printf(f, "hello")

// The attributes carry onto anything written to the terminal until reset.
// This prints "world" in the same attributes as above.
fmt.Println("world")

// Resets the highlighting and then prints "hello" normally.
f = color.Prepare("%r%s")
color.Printf(f, "foo")
```

## TODO
- [ ] True color support (needs work in my terminfo package)
- [ ] Windows support
- [x] Respect $TERM
- [x] Seperate log package
- [x] color.Format docs
- [x] A better way to combine `color.Format`s
- [ ] Logging/Printing/Format tests
