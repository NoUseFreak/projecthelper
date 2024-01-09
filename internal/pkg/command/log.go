package command

import (
	"bytes"
	"os"
	"strings"

	"github.com/nousefreak/projecthelper/internal/pkg/color"
	"github.com/sirupsen/logrus"
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
		b.WriteString(color.Color(color.FgGreen, "* "))
	case logrus.WarnLevel:
		b.WriteString(color.Color(color.FgYellow, "W "))
	case logrus.ErrorLevel:
		b.WriteString(color.Color(color.FgRed, "E "))
	case logrus.FatalLevel:
		b.WriteString(color.Color(color.FgRed, "F "))
	default:
		b.WriteString("  ")
	}

	entry.Message = strings.TrimSuffix(entry.Message, "\n")
	b.WriteString(entry.Message)

	b.WriteString("\n")
	return b.Bytes(), nil
}

