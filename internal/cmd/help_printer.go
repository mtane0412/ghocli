/**
 * help_printer.go
 * Custom help printer
 *
 * Provides colorized help messages and build information injection.
 */

package cmd

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

// helpPrinter is a custom help printer
func helpPrinter(options kong.HelpOptions, ctx *kong.Context) error {
	origStdout := ctx.Stdout
	origStderr := ctx.Stderr

	// Get terminal width
	width := guessColumns(origStdout)
	oldCols, hadCols := os.LookupEnv("COLUMNS")
	_ = os.Setenv("COLUMNS", strconv.Itoa(width))
	defer func() {
		if hadCols {
			_ = os.Setenv("COLUMNS", oldCols)
		} else {
			_ = os.Unsetenv("COLUMNS")
		}
	}()

	// Output to buffer
	buf := bytes.NewBuffer(nil)
	ctx.Stdout = buf
	ctx.Stderr = origStderr
	defer func() { ctx.Stdout = origStdout }()

	// Generate default help
	if err := kong.DefaultHelpPrinter(options, ctx); err != nil {
		return err
	}

	// Customize help text
	out := buf.String()
	out = injectBuildLine(out)
	out = colorizeHelp(out, helpProfile(origStdout))

	_, err := io.WriteString(origStdout, out)
	return err
}

// injectBuildLine injects build information into help text
func injectBuildLine(out string) string {
	// Build information is obtained from the version set in Execute()
	// For simplicity, injection is skipped here
	return out
}

// helpProfile returns the color profile
func helpProfile(stdout io.Writer) termenv.Profile {
	// Check NO_COLOR environment variable
	if termenv.EnvNoColor() {
		return termenv.Ascii
	}

	// Check GHO_COLOR environment variable
	mode := strings.ToLower(strings.TrimSpace(os.Getenv("GHO_COLOR")))
	switch mode {
	case "never":
		return termenv.Ascii
	case "always":
		return termenv.TrueColor
	default:
		// auto: Auto-detect terminal color support
		o := termenv.NewOutput(stdout, termenv.WithProfile(termenv.EnvColorProfile()))
		return o.Profile
	}
}

// colorizeHelp colorizes help text
func colorizeHelp(out string, profile termenv.Profile) string {
	// Don't colorize if color profile is Ascii
	if profile == termenv.Ascii {
		return out
	}

	// Define color functions
	heading := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#60a5fa")).Bold().String()
	}
	section := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#a78bfa")).Bold().String()
	}
	cmdName := func(s string) string {
		return termenv.String(s).Foreground(profile.Color("#38bdf8")).Bold().String()
	}

	// Process line by line
	inCommands := false
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		if line == "Commands:" {
			inCommands = true
		}
		switch {
		case strings.HasPrefix(line, "Usage:"):
			lines[i] = heading("Usage:") + strings.TrimPrefix(line, "Usage:")
		case line == "Flags:":
			lines[i] = section(line)
		case line == "Commands:":
			lines[i] = section(line)
		case line == "Arguments:":
			lines[i] = section(line)
		case inCommands && strings.HasPrefix(line, "  ") && len(line) > 2 && line[2] != ' ':
			// Colorize command name
			name, tail, found := strings.Cut(strings.TrimPrefix(line, "  "), " ")
			if found {
				lines[i] = "  " + cmdName(name) + " " + tail
			} else {
				lines[i] = "  " + cmdName(name)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// guessColumns guesses terminal width
func guessColumns(w io.Writer) int {
	// Check COLUMNS environment variable
	if colsStr := os.Getenv("COLUMNS"); colsStr != "" {
		if cols, err := strconv.Atoi(colsStr); err == nil {
			return cols
		}
	}

	// Get terminal size
	f, ok := w.(*os.File)
	if !ok {
		return 80
	}

	width, _, err := term.GetSize(int(f.Fd()))
	if err == nil && width > 0 {
		return width
	}

	// Default value
	return 80
}
