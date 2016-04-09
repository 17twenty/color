# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color wraps the `fmt.Printf` functions with verbs for producing colored output.

## Install
```
go get github.com/nhooyr/color
```

## Usage
```
%h[attr...]	replaced with a SGR code that sets all of the attributes in []
			multiple attributes are + separated
%r			an abbreviation for %h[reset]
```

See [godoc](https://godoc.org/github.com/nhooyr/color) for more information.

## Examples
### 16 Colors
```go
// "panic:" with a red foreground then normal "rip".
color.Printf("%h[fgRed]panic:%r rip\n")

// "panic:" with a brightRed background then normal "rip".
color.Printf("%h[bgBrightRed]panic:%r rip\n")
```

### 256 Colors
```go
// "panic:" with a green foreground then normal "rip".
color.Printf("%h[fg2]panic:%r rip\n")

// "panic:" with a bright green background then normal "rip".
color.Printf("%h[bg10]panic:%r rip\n")
```

### Other Attributes
```go
// Bold "panic:" then normal "rip".
color.Printf("%h[bold]panic:%r rip\n")

// Underlined "panic:" with then normal "rip".
color.Printf("%h[underline]panic:%r rip\n")
```

### Mixing Attributes
```go
// Bolded "panic:" with a green foreground then normal "rip".
color.Printf("%h[fgGreen+bold]panic:%r rip\n")

// Underlined "panic:" with a bright black background then normal "rip".
color.Printf("%h[bg8+underline]panic:%r rip\n")
```

### Reset's Behavior
```go
// Bolded "panic:" with a blue foreground then
// bolded "rip" with a blue foreground and bright black background.
color.Printf("%h[fgBlue+bold]panic: %h[bg8]rip\n")

// Bolded "hi" with a blue foreground and bright black background because
// we did not reset the highlighting above.
fmt.Printf("hi")

// Resets the highlighting and then prints "hello" normally.
color.Printf("%rhello")
```

### Printer
```go
p := color.NewPrinter(os.Stderr)

// Prints "hi" with red foreground.
p.Printf("%h[fgRed]hi%r\n")

p.DisableColor()

// Prints "hi" normally.
p.Printf("%h[fgRed]hi%r\n")

p.EnableColor()

// Prints "hi" with red foreground.
p.Printf("%h[fgRed]hi%r\n")
```

### `*log.Logger` wrapper
```go
l := color.NewLogger(os.Stderr, "%h[bold]color:%r ", 0)

// Prints bold "color:" and then "hi" with red foreground.
l.Printf("%h[fgRed]hi%r")

l.DisableColor()

// Prints "color: hi" normally.
l.Printf("%h[fgRed]hi%r")

l.EnableColor()

// Prints bold "color:" and then "hi" with red foreground and
// then exits with status code 1.
l.Fatalf("%h[fgRed]hi%r")
```

## TODO
- [ ] True color support
- [ ] Windows support
- [ ] Respect $TERM
- [ ] Fully wrap \*log.Logger, perhaps a format string that defines the prefix, date, content etc. Perhaps another package?
