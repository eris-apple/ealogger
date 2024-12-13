package adapters

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/eris-apple/ealogger/ealogger/shared"
	"os"
	"strings"
	"time"
)

type ConsoleColorConfig struct {
	TimestampColor *string
	MessageColor   *string

	LevelColors map[shared.Level]string
}

type ConsoleConfig struct {
	Enable bool

	Level  shared.Level
	Colors *ConsoleColorConfig
}

type ConsoleAdapter struct {
	writer *log.Logger
	cfg    *ConsoleConfig
}

func (a *ConsoleAdapter) Log(log shared.Log) {
	if !a.cfg.Enable || !a.cfg.Level.IsEnabled(log.Level) {
		return
	}

	a.Format(&log)

	switch log.Level.String() {
	case shared.DebugLevel.String():
		a.writer.Debug(log.Data.TraceName + log.Message)
	case shared.InfoLevel.String():
		a.writer.Info(log.Data.TraceName + log.Message)
	case shared.WarnLevel.String():
		a.writer.Warn(log.Data.TraceName + log.Message)
	case shared.ErrorLevel.String():
		a.writer.Error(log.Data.TraceName + log.Message)
	case shared.FatalLevel.String():
		a.writer.Fatal(log.Data.TraceName + log.Message)
	case shared.UnselectedLevel.String():
		a.writer.Print(log.Data.TraceName + log.Message)
	default:
		a.writer.Info(log.Data.TraceName + log.Message)
	}
}

func (a *ConsoleAdapter) Format(log *shared.Log) {
	if log.Data.TraceName != "" {
		log.Data.TraceName = lipgloss.
			NewStyle().
			SetString(fmt.Sprintf("[%s]: ", log.Data.TraceName)).
			Foreground(lipgloss.Color(a.cfg.Colors.LevelColors[log.Level])).
			String()
	}

	log.Message = lipgloss.
		NewStyle().
		SetString(fmt.Sprintf("%s", log.Message)).
		Foreground(lipgloss.Color(*a.cfg.Colors.MessageColor)).
		String()

	if log.Data != nil {
		if log.Data.Error != nil {
			formattedError := lipgloss.
				NewStyle().
				SetString(fmt.Sprintf("err=%s", log.Data.Error)).
				Foreground(lipgloss.Color(a.cfg.Colors.LevelColors[shared.ErrorLevel])).
				String()

			log.Message = fmt.Sprintf("%s %s", log.Message, formattedError)
		}

		if len(log.Data.Fields) > 0 {
			formattedFields := ""
			for key, field := range log.Data.Fields {
				formattedField := lipgloss.
					NewStyle().
					SetString(fmt.Sprintf("%s=%v ", key, field)).
					Foreground(lipgloss.Color(a.cfg.Colors.LevelColors[log.Level])).
					String()

				formattedFields += formattedField
			}

			log.Message = fmt.Sprintf("%s %s", log.Message, formattedFields)
		}
	}
}

func NewConsoleAdapter(cfg *ConsoleConfig) *ConsoleAdapter {
	return &ConsoleAdapter{
		cfg:    cfg,
		writer: newConsoleLogger(cfg),
	}
}

func NewDefaultConsoleAdapter() *ConsoleAdapter {
	cfg := defaultConsoleConfig()

	return &ConsoleAdapter{
		cfg:    cfg,
		writer: newConsoleLogger(cfg),
	}
}

func NewDefaultConsoleAdapterWithLevel(level shared.Level) *ConsoleAdapter {
	cfg := defaultConsoleConfig()
	cfg.Level = level

	return &ConsoleAdapter{
		cfg:    cfg,
		writer: newConsoleLogger(cfg),
	}
}

func newConsoleLogger(cfg *ConsoleConfig) *log.Logger {
	if !cfg.Enable {
		return nil
	}

	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           cfg.Level.ToCharmbracelet(),
	})

	setupDefaultLoggerColors(cfg)
	logger.SetStyles(setupStyles(cfg.Colors))
	return logger
}

func defaultConsoleConfig() *ConsoleConfig {
	return &ConsoleConfig{
		Enable: true,
		Level:  shared.DebugLevel,
		Colors: &ConsoleColorConfig{},
	}
}

func setupStyles(cfg *ConsoleColorConfig) *log.Styles {
	styles := log.DefaultStyles()
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.Color(*cfg.MessageColor))
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color(*cfg.TimestampColor))

	for key, value := range cfg.LevelColors {
		styles.Levels[key.ToCharmbracelet()] = lipgloss.NewStyle().SetString(strings.ToUpper(key.String())).Foreground(lipgloss.Color(value))
	}

	return styles
}

func setupDefaultLoggerColors(cfg *ConsoleConfig) {
	timestampColor := "#8a8a8a"
	messageColor := "#e3e3e3"

	if cfg.Colors.TimestampColor == nil {
		cfg.Colors.TimestampColor = &timestampColor
	}
	if cfg.Colors.MessageColor == nil {
		cfg.Colors.MessageColor = &messageColor
	}

	if cfg.Colors.LevelColors == nil {
		cfg.Colors.LevelColors = make(map[shared.Level]string)
	}

	if cfg.Colors.LevelColors[shared.InfoLevel] == "" {
		cfg.Colors.LevelColors[shared.InfoLevel] = "#afd7ff"
	}

	if cfg.Colors.LevelColors[shared.DebugLevel] == "" {
		cfg.Colors.LevelColors[shared.DebugLevel] = "#969696"
	}

	if cfg.Colors.LevelColors[shared.WarnLevel] == "" {
		cfg.Colors.LevelColors[shared.WarnLevel] = "#ffff18"
	}

	if cfg.Colors.LevelColors[shared.ErrorLevel] == "" {
		cfg.Colors.LevelColors[shared.ErrorLevel] = "#af0000"
	}

	if cfg.Colors.LevelColors[shared.FatalLevel] == "" {
		cfg.Colors.LevelColors[shared.FatalLevel] = "#ff0000"
	}
}
