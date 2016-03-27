# color

This package adds color verbs to fmt.Printf. All three printf functions are wrapped in this package.

##Examples:
```go
package main

import (
	"github.com/nhooyr/color"
)

func main() {
	// bolded red "panic:" and then normal "rip"
	color.Printf("%h#bold+fgRed#panic:%r rip\n")

	// bolded red with bright black background "panic:" and then cyan "rip"
	color.Printf("%h#bold+fgRed+bgBrightBlack#panic:%r %h#fgCyan#rip%r\n")

	// underlined red with bright black background "panic:" and then normal "rip"
	color.Printf("%h#underline+fgRed+bgBrightBlack#panic:%r rip\n")

	// red with black background "panic:" and then normal "rip"
	color.Printf("%h#underline+bold+fg1+bg0#panic:%r rip\n")
}
```

##Reference:
```
16 Colors:
%h#fgBlack#
%h#fgRed#
%h#fgGreen#
%h#fgYellow#
%h#fgBlue#
%h#fgMagenta#
%h#fgCyan#
%h#fgWhite#
%h#fgDefault#
%h#bgBlack#
%h#bgRed#
%h#bgGreen#
%h#bgYellow#
%h#bgBlue#
%h#bgMagenta#
%h#bgCyan#
%h#bgWhite#
%h#bgDefault#
%h#fgBrighBlack#
%h#fgBrightRed#
%h#fgBrightGreen#
%h#fgBrightYellow#
%h#fgBrightBlue#
%h#fgBrightMagenta#
%h#fgBrightCyan#
%h#fgBrightWhite#
%h#bgBrighBlack#
%h#bgBrightRed#
%h#bgBrightGreen#
%h#bgBrightYellow#
%h#bgBrightBlue#
%h#bgBrightMagenta#
%h#bgBrightCyan#
%h#bgBrightWhite#

256 Colors:
%h#fg144#
%h#bg144#

Attributes:
%h#reset# or just %r
%h#bold#
%h#faint#
%h#italic#
%h#underline#
%h#blink#
%h#inverse#
%h#invisible#
%h#crossedOut#
%h#doubleUnderline#
%h#normal#
%h#notItalic#
%h#notUnderlined#
%h#steady#
%h#positive#
%h#visible#
%h#notCrossedOut#

To combine:
%h#fgBlue+bgRed+bold#
```

TODO:
-----
- [ ] True color support, just not sure on the schema
