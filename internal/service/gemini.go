package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/beriholic/geminic/internal/config"
	"github.com/google/generative-ai-go/genai"
)

type GeminiService struct {
	SystemInstruction string
	Prompts           []string
}

var (
	geminiServerOnce sync.Once
	geminiServer     *GeminiService = nil
)

func GetGeminiService() *GeminiService {
	geminiServerOnce.Do(func() {
		geminiServer = &GeminiService{
			Prompts: make([]string, 0),
		}
	})
	return geminiServer
}

func (g *GeminiService) BuildGitCommitStle() *GeminiService {
	emoji := config.Get().Emoji
	rule := `
You now need to help the user generate the message for the git commit please follow the rules 
<rule>
- Write in first-person singular present tense
- Be concise and direct
- Output only the commit message without any explanations
- Commit message should starts with lowercase letter.
- Commit message must be a maximum of 72 characters.
- Exclude anything unnecessary such as translation. Your entire response will be passed directly into git commit.
</rule>
`

	gitCommitStyle := `
<git-commit-style>
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
</git-commit-style>
`
	gitCommitEmoji := `
<git-commit-emoji>
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
</git-commit-emoji>
`
	g.Prompts = append(g.Prompts, gitCommitStyle)
	if emoji {
		g.Prompts = append(g.Prompts, gitCommitEmoji)
	}
	g.Prompts = append(g.Prompts, rule)
	return g
}

func (g *GeminiService) BuildCommitInfoPrompt(
	commit string,
	diff string,
	files []string,
) *GeminiService {
	if commit != "" {
		const pmt = `<user-commit> %s (write on this basis) </user-commit>`
		g.Prompts = append(g.Prompts, fmt.Sprintf(pmt, commit))
	}

	fileChanged := fmt.Sprintf(`<files-changed> %s </files-changed>`, strings.Join(files, ", "))
	codeDiff := fmt.Sprintf(`<code-diff> %s </code-diff>`, diff)

	g.Prompts = append(g.Prompts, fileChanged)
	g.Prompts = append(g.Prompts, codeDiff)
	return g
}
func (g *GeminiService) BuildCot() *GeminiService {
	emoji := config.Get().Emoji
	prompt := ""
	if emoji {
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
	g.Prompts = append(g.Prompts, prompt)
	return g
}

func (g *GeminiService) BuildResponseStructure() *GeminiService {
	prompt := `return git commit message using this JSON schema:
	           Return {
			     "typ": string,
				 "emoji": string?,
				 "scope": string?,
				 "msg":string
			   }`
	g.Prompts = append(g.Prompts, prompt)
	return g
}

func (g *GeminiService) GetPrompt(commit string, diff string, files []string) string {
	g.BuildCot().
		BuildGitCommitStle().
		BuildCommitInfoPrompt(commit, diff, files).
		BuildResponseStructure()

	return strings.Join(g.Prompts, "\n")
}

func (g *GeminiService) AnalyzeChanges(
	userCommit string,
	geminiClient *genai.Client,
	ctx context.Context,
	diff string,
	relatedFiles *map[string]string,
	modelName *string,
) (string, error) {
	relatedFilesArray := make([]string, 0, len(*relatedFiles))
	for dir, ls := range *relatedFiles {
		relatedFilesArray = append(relatedFilesArray, fmt.Sprintf("%s/%s", dir, ls))
	}

	model := geminiClient.GenerativeModel(*modelName)
	model.ResponseMIMEType = "application/json"
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
	}

	model.SetTemperature(1.02)
	// model.SetTopK(40)
	// model.SetTopP(0.92)

	prompt := g.GetPrompt(userCommit, diff, relatedFilesArray)

	resp, err := model.GenerateContent(
		ctx,
		genai.Text(prompt),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
