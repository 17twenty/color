package color

var csi = "\033["

const (
	preFg256 = ";38;5;"
	preBg256 = ";48;5;"
)

// maps attributes names to their values.
var attrs = map[string]string{
	"reset":           ";0",
	"bold":            ";1",
	"faint":           ";2",
	"italic":          ";3",
	"underline":       ";4",
	"blink":           ";5",
	"inverse":         ";7",
	"invisible":       ";8",
	"crossedOut":      ";9",
	"doubleUnderline": ";21",
	"normal":          ";22",
	"notItalic":       ";23",
	"notUnderlined":   ";24",
	"steady":          ";25",
	"positive":        ";27",
	"visible":         ";28",
	"notCrossedOut":   ";29",
	"fgBlack":         ";30",
	"fgRed":           ";31",
	"fgGreen":         ";32",
	"fgYellow":        ";33",
	"fgBlue":          ";34",
	"fgMagenta":       ";35",
	"fgCyan":          ";36",
	"fgWhite":         ";37",
	"fgDefault":       ";39",
	"bgBlack":         ";40",
	"bgRed":           ";41",
	"bgGreen":         ";42",
	"bgYellow":        ";43",
	"bgBlue":          ";44",
	"bgMagenta":       ";45",
	"bgCyan":          ";46",
	"bgWhite":         ";47",
	"bgDefault":       ";49",
	"fgBrightBlack":   ";90",
	"fgBrightRed":     ";91",
	"fgBrightGreen":   ";92",
	"fgBrightYellow":  ";93",
	"fgBrightBlue":    ";94",
	"fgBrightMagenta": ";95",
	"fgBrightCyan":    ";96",
	"fgBrightWhite":   ";97",
	"bgBrightBlack":   ";100",
	"bgBrightRed":     ";101",
	"bgBrightGreen":   ";102",
	"bgBrightYellow":  ";103;",
	"bgBrightBlue":    ";104",
	"bgBrightMagenta": ";105",
	"bgBrightCyan":    ";106",
	"bgBrightWhite":   ";107",
}
