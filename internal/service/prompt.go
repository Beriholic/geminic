package service

import (
	"fmt"
	"strings"

	"github.com/beriholic/geminic/internal/config"
)

type Prompt struct {
	Basic  string
	Struct []string
}

func NewPrompt() *Prompt {
	prompt := Prompt{
		Basic:  "You now need to help the user generate the message for the git commit please follow the rules",
		Struct: []string{},
	}

	return &prompt
}

func (p *Prompt) AddStructStart(name string) *Prompt {
	prompt := fmt.Sprintf("<%s>", name)
	return p.AddStruct(prompt)
}
func (p *Prompt) AddStructEnd(name string) *Prompt {
	prompt := fmt.Sprintf("</%s>", name)
	return p.AddStruct(prompt)
}

func (p *Prompt) Build(
	commit string,
	diff string,
	files []string,
) string {
	p.
		AddRule().
		AddCommitType().
		AddCommitEmoji().
		AddCommitInfo(commit, diff, files).
		AddCot().
		AddOutputStruct()
	return p.Basic + "\n" + strings.Join(p.Struct, "\n")
}
func (p *Prompt) AddStruct(prompt string) *Prompt {
	p.Struct = append(p.Struct, prompt)
	return p
}

func (p *Prompt) AddRule() *Prompt {
	prompt := `
<Rule>
- Write in first-person singular present tense
- Be concise and direct
- Output only the commit message without any explanations
- Commit message should starts with lowercase letter.
- Commit message must be a maximum of 72 characters.
- Exclude anything unnecessary such as translation. Your entire response will be passed directly into git commit.
- Commit Message without subject
</Rule>
`
	return p.AddStruct(prompt)
}

func (p *Prompt) AddCommitType() *Prompt {
	propmt := `
<GitCommitType>
"feat":     "A new feature"
"fix":      "A bug fix"
"docs":     "Documentation only changes"
"style":    "Changes that do not affect the meaning of the code (white-space formatting missing semi-colons etc)"
"refactor": "A code change that neither fixes a bug nor adds a feature"
"perf":     "A code change that improves performance"
"test":     "Adding missing tests or correcting existing tests"
"build":    "Changes that affect the build system or external dependencies"
"ci":       "Changes to our CI configuration files and scripts"
"chore":    "Other changes that don't modify src or test files"
"revert":   "Reverts a previous commit"
</GitCommitType>
`
	return p.AddStruct(propmt)
}

func (p *Prompt) AddCommitEmoji() *Prompt {
	if !config.Get().Emoji {
		return p
	}

	prompt := `<GitCommitEmoji>
"feat": ":sparkles:"
"fix": ":bug:"
"docs": ":memo:":
"style": ":lipstick:":
"refactor": ":recycle:":
"perf": ":zap:"
"test: ":white_check_mark:":
"build: ":package:":
"ci: ":ferris_wheel:":
"chore: ":hammer:":
"revert: ":rewind:":
</GitCommitEmoji>
`
	return p.AddStruct(prompt)
}

func (p *Prompt) AddOutputStruct() *Prompt {
	prompt := ""
	if config.Get().Emoji {
		prompt = `return git commit message using this JSON schema:
Return {
  "typ": string,
  "emoji": string,
  "scope": string?,
  "msg":string
}`
	} else {
		prompt = `return git commit message using this JSON schema:
Return {
  "typ": string,
  "scope": string?,
  "msg":string
}`
	}
	return p.AddStruct(prompt)
}

func (p *Prompt) AddCot() *Prompt {
	prompt := ""
	if config.Get().Emoji {
		prompt = `
Use the following format to output the chain of thought before each response
<thinking>
1. what the code changed?
2. what the purpose of the change was?
3. what the type of the change was?
4. if you could use emoji, what would you use?
</thinking>
`
	} else {
		prompt = `
Use the following format to output the chain of thought before each response
<thinking>
1. what the code changed?
2. what the purpose of the change was?
3. what the type of the change was?
</thinking>
`
	}

	return p.AddStruct(prompt)
}

func (p *Prompt) AddCommitInfo(
	commit string,
	diff string,
	files []string,
) *Prompt {
	userInput := ""
	if commit != "" {
		userInput = fmt.Sprintf(`<UserInput> %s (write on this basis) </UserInput>`, commit)
	}
	fileChanged := fmt.Sprintf(`<FilesChanged> %s </-changed>`, strings.Join(files, ", "))
	codeDiff := fmt.Sprintf(`<CodeDiff> %s </CodeDiff>`, diff)

	p.AddStructStart("CommitInfo")
	if userInput != "" {
		p.AddStruct(userInput)
	}
	p.AddStruct(fileChanged)
	p.AddStruct(codeDiff)
	p.AddStructEnd("CommitInfo")
	return p
}
