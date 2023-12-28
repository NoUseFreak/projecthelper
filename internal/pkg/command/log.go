package command

import (
    "bytes"
    "fmt"
    "os"
    "strings"

    "github.com/sirupsen/logrus"
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

func init() {
    logrus.SetOutput(os.Stderr)
    logrus.SetFormatter(&cliFormatter{})
}

var CmdOutput = os.Stdout

type cliFormatter struct {
}

func (f cliFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    var b *bytes.Buffer
    if entry.Buffer != nil {
        b = entry.Buffer
    } else {
        b = &bytes.Buffer{}
    }

    switch entry.Level {
    case logrus.InfoLevel:
        b.WriteString(f.color(FgGreen, "* "))
    case logrus.WarnLevel:
        b.WriteString(f.color(FgYellow, "W "))
    case logrus.ErrorLevel:
        b.WriteString(f.color(FgRed, "E "))
    case logrus.FatalLevel:
        b.WriteString(f.color(FgRed, "F "))
    default:
        b.WriteString("  ")
    }

    entry.Message = strings.TrimSuffix(entry.Message, "\n")
    b.WriteString(entry.Message)

    b.WriteString("\n")
    return b.Bytes(), nil
}

func (f cliFormatter) color(a int, msg string) string {
    return fmt.Sprintf("%s[%dm%s%s[%dm", Escape, a, msg, Escape, Reset)
}
