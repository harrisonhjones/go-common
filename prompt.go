package common

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// DefaultPrompter is the package-level Prompter using stdin/stdout.
var DefaultPrompter = &Prompter{
	Reader: os.Stdin,
	Writer: os.Stdout,
}

// Prompter reads from Reader and writes to Writer.
// The zero value is not usable; use NewPrompter or DefaultPrompter.
type Prompter struct {
	Reader io.Reader
	Writer io.Writer
}

// NewPrompter returns a Prompter using the given reader and writer.
func NewPrompter(r io.Reader, w io.Writer) *Prompter {
	return &Prompter{Reader: r, Writer: w}
}

// Prompt writes "name: " to the writer and returns the trimmed input from the reader.
// If the input is empty, an empty string is returned.
func (p *Prompter) Prompt(name string) (string, error) {
	fmt.Fprint(p.Writer, name+": ")
	return p.readLine()
}

// PromptWithDefault writes "name [Default: def]: " to the writer and returns
// the trimmed input from the reader. If the input is empty, def is returned.
func (p *Prompter) PromptWithDefault(name, def string) (string, error) {
	fmt.Fprint(p.Writer, name+" [Default: "+def+"]: ")
	line, err := p.readLine()
	if err != nil {
		return "", err
	}
	if line == "" {
		return def, nil
	}
	return line, nil
}

// PromptContinue prompts for a yes/no confirmation and returns true if the
// user enters "y" or "yes" (case-insensitive), false otherwise.
// If defaultYes is true the prompt defaults to "Y", otherwise "N".
// An empty or whitespace-only input accepts the default.
func (p *Prompter) PromptContinue(defaultYes bool) (bool, error) {
	def := "N"
	if defaultYes {
		def = "Y"
	}
	line, err := p.PromptWithDefault("Continue (Y/N)", def)
	if err != nil {
		return false, err
	}
	switch strings.ToLower(line) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

func (p *Prompter) readLine() (string, error) {
	line, err := bufio.NewReader(p.Reader).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read line: %w", err)
	}
	return strings.TrimSpace(line), nil
}

// Prompt calls DefaultPrompter.Prompt.
func Prompt(name string) (string, error) {
	return DefaultPrompter.Prompt(name)
}

// PromptWithDefault calls DefaultPrompter.PromptWithDefault.
func PromptWithDefault(name, def string) (string, error) {
	return DefaultPrompter.PromptWithDefault(name, def)
}

// PromptContinue calls DefaultPrompter.PromptContinue.
func PromptContinue(defaultYes bool) (bool, error) {
	return DefaultPrompter.PromptContinue(defaultYes)
}
