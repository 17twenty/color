# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs for terminal color highlighting. All it does is replace the verbs with the appropiate terminal escape sequence.

It also provides a wrapper around `*log.Logger`'s `Printf` functions to make use of the new verbs. See `color.NewLogger`.

## Install

`go get github.com/nhooyr/color`

## Why?

The API of similar packages requires calling a function everytime to color some text in a different color. E.g. once for red, then for yellow, and so on.  
That approach gets very verbose, I prefer the succinctness of using verbs.

## Usage

`%h[colors+attributes]text`
Highlights text with the `colors+attributes` in `[]`.
Subsquent highlighting verbs will not reset the highlighting first, they will just add onto it.

You can also use the `%r` verb as an abbreviation for `%h[reset]`

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
// bold "panic:" and then underlined "rip"
color.Printf("%h[bold]panic:%r %h[underline]rip%r\n")
```

### 256 Colors
```go
// red "panic:" and then "rip" with a cyan background
color.Printf("%h[fg1]panic:%r %h[bg6]rip\n")
```

### Mixing Colors and Attributes
```go
// bolded green "panic:" and then underlined "rip" with bright black background
color.Printf("%h[fgGreen+bold]panic:%r %h[bg8+underline]rip%r\n")
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
See [this](http://invisible-island.net/xterm/ctlseqs/ctlseqs.html) for explanations
```
16 Colors:
%h[fgBlack]
%h[fgRed]
%h[fgGreen]
%h[fgYellow]
%h[fgBlue]
%h[fgMagenta]
%h[fgCyan]
%h[fgWhite]
%h[fgDefault]
%h[bgBlack]
%h[bgRed]
%h[bgGreen]
%h[bgYellow]
%h[bgBlue]
%h[bgMagenta]
%h[bgCyan]
%h[bgWhite]
%h[bgDefault]
%h[fgBrighBlack]
%h[fgBrightRed]
%h[fgBrightGreen]
%h[fgBrightYellow]
%h[fgBrightBlue]
%h[fgBrightMagenta]
%h[fgBrightCyan]
%h[fgBrightWhite]
%h[bgBrighBlack]
%h[bgBrightRed]
%h[bgBrightGreen]
%h[bgBrightYellow]
%h[bgBrightBlue]
%h[bgBrightMagenta]
%h[bgBrightCyan]
%h[bgBrightWhite]

256 Colors:
%h[fg144]
%h[bg144]

Attributes:
%h[reset] or the %r verb
%h[bold]
%h[faint]
%h[italic]
%h[underline]
%h[blink]
%h[inverse]
%h[invisible]
%h[crossedOut]
%h[doubleUnderline]
%h[normal]
%h[notItalic]
%h[notUnderlined]
%h[steady]
%h[positive]
%h[visible]
%h[notCrossedOut]

To combine:
%h[fgBlue+bgRed+bold]
```

TODO:
-----
- [ ] True color support
- [ ] Windows support

