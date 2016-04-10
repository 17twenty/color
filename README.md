# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color wraps `fmt.Printf` with verbs for producing colored output.

__note: this is still a WIP and things may change__

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

### Printer
```go
p := color.NewPrinter(os.Stderr, color.EnableColor)

// Prints "hi" with red foreground.
p.Printf("%h[fgRed]hi%r\n")

p = color.NewPrinter(os.Stderr, color.DisableColor)

// Prints "hi" normally.
p.Printf("%h[fgRed]hi%r\n")

p = color.NewPrinter(os.Stderr, color.PerformCheck)

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "hi"
p.Printf("%h[fgRed]hi%r\n")
```

### Preparing Strings
```go
// Prepare processes the string only once and lets you print it repeatedly
// without any additional processing overhead.
f := color.Prepare("%h[fgRed+bold]panic:%r %s\n")

// Each prints bolded "panic:" with a red foreground and some normal text after.
color.Aprintf(f, "rip")
color.Aprintf(f, "yippie")
color.Aprintf(f, "dsda")
```

### `*log.Logger` wrapper
```go
l := color.NewLogger(os.Stderr, "%h[bold]color:%r ", 0)

// Prints "hi" with a red foreground.
l.Printf("%h[fgRed]hi%r")

l.DisableColor()

// now a normal "hi", the highlight verbs are ignored.
l.Printf("%h[fgRed]hi%r")

l.EnableColor()

// Prints "hi" with a red foreground.
l.Fatalf("%h[fgRed]hi%r")
```

### How does reset behave?
```go
// "rip" will be printed with a blue foreground and bright black background
// because we never reset the highlighting after "panic:". The blue foreground is
// carried on from "panic:".
color.Printf("%h[fgBlue+bgBlack+bold]panic: %h[bg8]rip\n")

// The attributes carry onto anything written to the terminal until reset.
// This prints "rip" in the same attributes as above.
fmt.Println("rip")

// Resets the highlighting and then prints "hello" normally.
color.Printf("%rhello")
```

## TODO
- [ ] True color support
- [ ] Windows support
- [ ] Respect $TERM
- [ ] Fully wrap \*log.Logger, perhaps a format string that defines the prefix, date, content etc. Perhaps another package?
