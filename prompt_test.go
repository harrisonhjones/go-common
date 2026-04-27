package common

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPrompter(input string) (*Prompter, *bytes.Buffer) {
	writer := &bytes.Buffer{}
	return NewPrompter(bytes.NewBufferString(input), writer), writer
}

func TestPrompterPrompt(t *testing.T) {
	tests := map[string]struct {
		input       string
		wantOutput  string
		wantWritten string
		wantErr     string
	}{
		"returns input": {
			input:       "bar\n",
			wantOutput:  "bar",
			wantWritten: "name: ",
		},
		"trims whitespace": {
			input:       "  bar  \n",
			wantOutput:  "bar",
			wantWritten: "name: ",
		},
		"empty input returns empty string": {
			input:       "\n",
			wantOutput:  "",
			wantWritten: "name: ",
		},
		"bad reader returns error": {
			input:       "",
			wantErr:     "failed to read line",
			wantWritten: "name: ",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p, writer := newTestPrompter(tc.input)
			output, err := p.Prompt("name")
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantOutput, output)
			}
			assert.Equal(t, tc.wantWritten, writer.String())
		})
	}
}

func TestPrompterPromptWithDefault(t *testing.T) {
	tests := map[string]struct {
		input       string
		def         string
		wantOutput  string
		wantWritten string
		wantErr     string
	}{
		"returns input over default": {
			input:       "bar\n",
			def:         "foo",
			wantOutput:  "bar",
			wantWritten: "name [Default: foo]: ",
		},
		"empty input returns default": {
			input:       "\n",
			def:         "foo",
			wantOutput:  "foo",
			wantWritten: "name [Default: foo]: ",
		},
		"blank input returns default": {
			input:       "   \n",
			def:         "foo",
			wantOutput:  "foo",
			wantWritten: "name [Default: foo]: ",
		},
		"empty default with empty input returns empty string": {
			input:       "\n",
			def:         "",
			wantOutput:  "",
			wantWritten: "name [Default: ]: ",
		},
		"bad reader returns error": {
			input:       "",
			def:         "foo",
			wantErr:     "failed to read line",
			wantWritten: "name [Default: foo]: ",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p, writer := newTestPrompter(tc.input)
			output, err := p.PromptWithDefault("name", tc.def)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantOutput, output)
			}
			assert.Equal(t, tc.wantWritten, writer.String())
		})
	}
}

func TestPrompterPromptContinue(t *testing.T) {
	tests := map[string]struct {
		input       string
		defaultYes  bool
		wantOutput  bool
		wantWritten string
		wantErr     string
	}{
		// Default no
		"default-no empty input": {
			input:       "\n",
			wantOutput:  false,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		"default-no input y": {
			input:       "y\n",
			wantOutput:  true,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		"default-no input Y": {
			input:       "Y\n",
			wantOutput:  true,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		"default-no input yes": {
			input:       "yes\n",
			wantOutput:  true,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		"default-no input YES": {
			input:       "YES\n",
			wantOutput:  true,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		"default-no input yuppers": {
			input:       "yuppers\n",
			wantOutput:  false,
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
		// Default yes
		"default-yes empty input": {
			input:       "\n",
			defaultYes:  true,
			wantOutput:  true,
			wantWritten: "Continue (Y/N) [Default: Y]: ",
		},
		// Bad read
		"bad reader returns error": {
			input:       "",
			wantErr:     "failed to read line",
			wantWritten: "Continue (Y/N) [Default: N]: ",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p, writer := newTestPrompter(tc.input)
			output, err := p.PromptContinue(tc.defaultYes)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantOutput, output)
			}
			assert.Equal(t, tc.wantWritten, writer.String())
		})
	}
}
