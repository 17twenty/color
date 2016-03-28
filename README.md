# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs that produce colored output.

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
// "panic:" with a red foreground then normal "rip"
color.Printf("%h[fgRed]panic:%r rip\n")

// "panic:" with a brightRed background then normal "rip"
color.Printf("%h[bgBrightRed]panic:%r rip\n")
```

### 256 Colors
```go
// "panic:" with a green foreground then normal "rip"
color.Printf("%h[fg2]panic:%r rip\n")

// "panic:" with a bright green background then normal "rip"
color.Printf("%h[bg10]panic:%r rip\n")
```

### Other Attributes
```go
// bold "panic:" then normal "rip"
color.Printf("%h[bold]panic:%r rip\n")

// underlined "panic:" with then normal "rip"
color.Printf("%h[underline]panic:%r rip\n")
```

### Mixing Attributes
```go
// bolded "panic:" with a green foreground then normal "rip"
color.Printf("%h[fgGreen+bold]panic:%r rip\n")

// underlined "panic:" with a bright black background then normal "rip"
color.Printf("%h[bg8+underline]panic:%r rip\n")
```

### How does reset behave?
```go
// bolded "panic:" with a blue foreground
// then bolded "rip" with a green foreground and bright black background
color.Printf("%h[fgBlue+bold]panic: %h[bg8]rip\n")

// bolded "hi" with a green foreground and bright black background
fmt.Printf("hi")

// finally resets the highlighting
color.Printf("%r")
```

### `*log.Logger` wrapper
```go
logger := color.NewLogger(os.Stderr, "", 0)

// prints hi in red
logger.Printf("%h[fgRed]hi%r")

// prints hi in red and then exits
logger.Fatalf("%h[fgRed]hi%r")
```

## TODO
- [ ] True color support
- [ ] Windows support
