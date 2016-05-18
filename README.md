# color [![GoDoc](https://godoc.org/github.com/nhooyr/color?status.svg)](https://godoc.org/github.com/nhooyr/color)

Color extends `fmt.Printf` with verbs for producing colored output.

__note: things may change but it looks pretty stable. Email me if you have any new ideas.__

## Install
```
go get github.com/nhooyr/color
```

## Examples
See [godoc](https://godoc.org/github.com/nhooyr/color) for more information.

### Setting Attributes
```go
// "panic:" with a red foreground then normal "foo".
color.Printf("%h[fgRed]panic:%r %s\n", "foo")

// "panic:" with a red background then normal "bar".
color.Printf("%h[bgRed]panic:%r %s\n", "bar")

// Bold "panic:" then normal "foo".
color.Printf("%h[bold]panic:%r %s\n", "foo")

// Underlined "panic:" with then normal "bar".
color.Printf("%h[underline]panic:%r %s\n", "bar")

// "panic:" using color 83 as the foreground then normal "foo".
color.Printf("%h[fg83]panic:%r %s\n", "foo")

// "panic:" using color 158 as the background then normal "bar".
color.Printf("%h[bg158]panic:%r %s\n", "bar")
```

### Mixing Attributes
```go
// Bolded "panic:" with a green foreground then normal "foo".
color.Printf("%h[fgGreen+bold]panic:%r %s\n", "foo")

// Underlined "panic:" with a bright black background then normal "bar".
color.Printf("%h[bg8+underline]panic:%r %s\n", "bar")
```

### Preparing
```go
// Prepare only processes the highlight verbs in the string,
// letting you print it repeatedly with performance.
panicFormat := color.Prepare("%h[fgRed+bold]panic:%r %s\n")

// Each will print "panic:" and some normal text after, but if standard output
// is a terminal, "panic:" will be printed in bold with a red foreground.
color.Printfp(panicFormat, "foo")
color.Printfp(panicFormat, "bar")
color.Printfp(panicFormat, "foo")

hello := color.Prepare("%h[fgBlue+bold]HELLO%r")

// Each will print "HELLO", but if standard output is a terminal,
// "HELLO" will be printed in bold with a blue foreground.
color.Println(hello)
color.Println(hello)
color.Println(hello)
```

### Printer
A `Printer` writes to an `io.Writer`.

```go
redFormat := color.Prepare("%h[fgRed]%s%r\n")

// If standard error is a terminal, this will print in color.
// Otherwise it will print a normal "bar".
p := color.New(os.Stderr, color.IsTerminal(os.Stderr))
p.Printfp(redFormat, "bar")

// "foo" with red foreground.
p = color.New(os.Stderr, true)
p.Printfp(redFormat, "foo")

// Normal "bar", the highlight verbs are ignored.
p = color.New(os.Stderr, false)
p.Printfp(redFormat, "bar")
```

### `github.com/nhooyr/color/log`
```go
redFormat := color.Prepare("%h[fgRed]%s%r\n")

// If os.Stderr is a terminal, this will print in color.
// Otherwise it will be a normal "foo".
log.Printfp(redFormat, "foo")

// Normal "bar", the highlight verbs are ignored.
log.SetColor(false)
log.Printfp(redFormat, "bar")

// "foo" with a red foreground.
log.SetColor(true)
log.Fatalfp(redFormat, "foo")
```

## Vim syntax highlighting
Add the following to `after/syntax/go.vim` to highlight the highlight verbs within strings.
```vim
syn match goFormatSpecifier /%[-#0 +]*\%(\*\|\d\+\)\=\%(\.\%(\*\|\d\+\)\)*\%([vTtbcdoqxXUeEfgGspr]\|h\[[a-zA-Z+0-9]\+\]\)/ contained containedin=goString
```
![syntax highlighting](https://i.imgur.com/yFZnJXg.png)

## TODO
- [ ] True color support
- [ ] Windows support
- [x] Respect $TERM
- [x] Seperate log package
- [x] color.Format docs
- [x] A better way to combine `color.Format`s
- [x] Logging/Printing/Format tests
- [ ] cleanup above tests
