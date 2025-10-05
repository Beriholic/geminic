package prompt

import (
	"fmt"
	"strings"

	"github.com/Beriholic/geminic/internal/config"
	"github.com/Beriholic/geminic/internal/model/dto"
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

func (p *Prompt) Build(commitDTO *dto.CommitDTO) string {
	if commitDTO == nil {
		return ""
	}

	p.
		AddRule().
		AddCommitType().
		AddCommitEmoji().
		AddCommitInfo(commitDTO.Commit, commitDTO.Diff, commitDTO.Files).
		AddI18n().
		AddOutputTemplateStruct()

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

func (p *Prompt) AddI18n() *Prompt {
	prompt := fmt.Sprintf("You need to write it in %s language", config.Get().I18n)

	p.AddStructStart("I18n")
	p.AddStruct(prompt)
	p.AddStructEnd("I18n")
	return p
}

func (p *Prompt) AddOutputTemplateStruct() *Prompt {
	p.AddStructStart("OutputTempalte")
	p.AddStruct(`
		Output only the following JSON structure, without any additional content
		{
			"typ": "(required)The type of git commit",
			"msg": "(required)The subject of git commit"
			"scope": "(optinal)The scope of git commit",
			"emoji": "(optinal)The emoji of git commit" 
		}`)
	p.AddStructEnd("OutputTempalte")
	return p
}
