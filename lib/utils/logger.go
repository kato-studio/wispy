package utils

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var logger = log.New(os.Stderr)

func init() {
	var styles = log.DefaultStyles()

	// Override the default fatal level style.
	styles.Levels[log.FatalLevel] = lipgloss.NewStyle().
		SetString("FATAL!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("196")).
		Foreground(lipgloss.Color("15"))

	// Override the default error level style.
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("0"))

	// Override the default warn level style.
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().
		SetString("WARN!!").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("220")).
		Foreground(lipgloss.Color("0"))

	// Override the default info level style.
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("238")).
		Foreground(lipgloss.Color("0"))

	// Override the default debug level style.
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("245")).
		Foreground(lipgloss.Color("0"))
}

var serverLogStyle = lipgloss.NewStyle().
	SetString("SERVER").
	Padding(0, 1, 0, 1).
	Background(lipgloss.Color("46,139,87")).
	Foreground(lipgloss.Color("255"))

var debugLogStyle = lipgloss.NewStyle().
	Padding(0, 1, 0, 1).
	Background(lipgloss.Color("20")).
	Foreground(lipgloss.Color("255"))

// example of title
// ╭─────────────────────────────────────────────────╮
// │ DEBUG Loading components from ./view/components │
// ╰─────────────────────────────────────────────────╯
var TitleStyle = lipgloss.NewStyle().
	SetString("").
	Padding(0, 1, 0, 1).
	Background(lipgloss.Color("100")).
	Foreground(lipgloss.Color("255")).
	BorderStyle(lipgloss.RoundedBorder())

func Fatal(str string) {
	logger.Fatal(str)
}

func Error(str string) {
	logger.Error(str)
}

func Warn(str string) {
	logger.Warn(str)
}

func Info(str string) {
	logger.Info(str)
}

func Debug(str string) {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		logger.Print(debugLogStyle.Render("DEBUG:") + "" + str)
	}
}

func ServerPrint(str string) {
	logger.Print("----------------")
	logger.Print(serverLogStyle.Render(str))
	logger.Print("----------------")
}

func Print(str interface{}) {
	logger.Print(str)
}
