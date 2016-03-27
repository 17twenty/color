# color

This package adds color verbs to fmt.Printf. All three printf functions are wrapped in this package.

Use `color.NewLogger()` to wrap the `Printf` functions of `*log.Logger`

## Why?

The API of similar packages requires calling a function everytime to color some text differently.  
It gets very verbose, I much prefer the succinctness of using verbs.

##Usage

`%h[color]text` is the highlighting verb. The text between the `[]` is used to describe the highlighting  
Everything after it is highlighted.  
Note: that the next highlighting verb will not reset the highlighting first, it will just add onto the first.

`%r` is the reset verb. It resets all highlighting.

##Examples:
```go
package main

import (
	"github.com/nhooyr/color"
)

func main() {
	// bolded red "panic:" and then normal "rip"
	color.Printf("%h[fgRed+bold]panic:%r rip\n")

	// bolded red with bright black background "panic:" and then cyan "rip"
	color.Printf("%h[fgRed+bgBrightBlack+bold]panic:%r %h[fgCyan]rip%r\n")

	// underlined red with bright black background "panic:" and then normal "rip"
	color.Printf("%h[fgRed+bgBrightBlack+underline]panic:%r rip\n")

	// red with black background "panic:" and then normal "rip"
	color.Printf("%h[fg1+bg0+underline+bold]panic:%r rip\n")

	// green "panic:" and then green with bright black background "rip"
	color.Printf("%h[fg2]panic: %h[bg8]rip%r\n")
}
```

##Reference:
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
%h[reset] or use the %r verb
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
- [ ] True color support, just not sure on the schema
