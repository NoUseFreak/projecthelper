package color

import (
    "fmt"
)

const (
    Reset int = iota
    Bold
    Faint
    Italic
    Underline
    BlinkSlow
    BlinkRapid
    ReverseVideo
    Concealed
    CrossedOut
)

const (
    FgBlack int = iota + 30
    FgRed
    FgGreen
    FgYellow
    FgBlue
    FgMagenta
    FgCyan
    FgWhite
)

const (
    Escape = "\x1b"
)


func Color(a int, msg string) string {
    return fmt.Sprintf("%s[%dm%s%s[%dm", Escape, a, msg, Escape, Reset)
}
