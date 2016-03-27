# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs for terminal color highlighting.  
All it does is replace the verbs with the appropiate terminal escape sequence.

## Install

`go get github.com/nhooyr/color`

## Why?

The API of similar packages requires calling a function everytime to color some text in a different color. E.g. once for red, then for yellow, and so on.  
That approach gets very verbose, I prefer the succinctness of using verbs.

## Usage

```
%h[attrs]		uses the attrs to highlight the following text
%r			an abbreviation for %h[reset]
```

attrs is a `+`  separated list of Colors (e.g. `fgRed`) or Attributes (e.g. `bold`).

Multiple highlight verbs do not reset preceeding verbs, they add onto them.  
For example, if you set the foreground to green in the first verb, then set the background to red in the second, any text following the second will have a green foreground and a red background.

## Examples:
### Standard Colors
```go
// red "panic:" and then normal "rip"
color.Printf("%h[fgRed]panic:%r rip\n")

// "panic:" with brightRed background and then normal "rip"
color.Printf("%h[bgBrightRed]panic:%r rip\n")
```

### Attributes
```go
// bold "panic:" and then normal "rip"
color.Printf("%h[bold]panic:%r rip\n")

// underlined "panic:" with and then normal "rip"
color.Printf("%h[underline]panic:%r rip\n")
```

### 256 Colors
```go
// red "panic:" and then normal "rip"
color.Printf("%h[fg1]panic:%r rip\n")

// "panic:" with a bright red background and then normal "rip"
color.Printf("%h[bg9]panic:%r rip\n")
```

### Mixing Colors and Attributes
```go
// bolded green "panic:" and then normal "rip"
color.Printf("%h[fgGreen+bold]panic:%r rip\n")

// underlined "panic:" with a bright black background and then normal "rip"
color.Printf("%h[bg8+underline]panic:%r rip\n")
```

### How does reset behave?
```go
// bolded green "panic:" and then bolded green "rip" with bright black background
color.Printf("%h[fgGreen+bold]panic: %h[bg8]rip\n")

// bolded green "hi" with bright black background
fmt.Printf("hi")

// finally resets the highlighting
color.Printf("%r")
```

## Reference
A syntax reference is [included](REFERENCE.md)

TODO:
-----
- [ ] True color support
- [ ] Windows support

