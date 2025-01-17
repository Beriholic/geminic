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
			SystemInstruction: `You now need to help the user generate the message for the git commit please follow the rules 
<rule>
- Write in first-person singular present tense
- Be concise and direct
- Output only the commit message without any explanations
- Follow the format: <type>(<optional scope>): <commit message>
- Commit message should starts with lowercase letter.
- Commit message must be a maximum of 72 characters.
- Exclude anything unnecessary such as translation. Your entire response will be passed directly into git commit.
</rule>
<git-commit-specification>
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
</git-commit-specificatio>
<emoji>
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
</emoji>
`,
			Prompts: []string{},
		}
	})
	return geminiServer
}

func (g *GeminiService) BuildCommitInfoPrompt(
	commit string,
	diff string,
	files []string,
) *GeminiService {
	prompt := fmt.Sprintf(`
<user-commit>
Reference commit: %s
</user-commit>
<files>
%s
</files>
<code-diff>
%s
<code-diff>
`,
		commit,
		strings.Join(files, ", "),
		diff,
	)

	g.Prompts = append(g.Prompts, prompt)
	return g
}
func (g *GeminiService) BuildCot() *GeminiService {
	prompt := `
Use the following format to output the chain of thought before each response
<thinking>
1. what the code changed
2. what the purpose of the change was
3. do you use emoji?
</thinking>
	`
	g.Prompts = append(g.Prompts, prompt)
	return g
}

func (g *GeminiService) BuildCorePrompt() *GeminiService {
	emoji := config.Get().Emoji
	prompt := fmt.Sprintf(`
Creative Requirements:
Write a git commit based on the changes made to the user's git repository, strictly following the format below

current emoji usage: %v

Formatting Demonstration:
<content>
if use emoji:
<type>: <emoji>(<optional scope>): <commit message>
else:
<type>(<optional scope>): <commit message>
</content>
`, emoji)
	g.Prompts = append(g.Prompts, prompt)
	return g
}

func (g *GeminiService) GetPrompt(commit string, diff string, files []string) string {
	g.BuildCot().BuildCommitInfoPrompt(commit, diff, files).BuildCorePrompt()

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
	safetySettings := []*genai.SafetySetting{
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

	model.SafetySettings = safetySettings

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(g.SystemInstruction)},
		Role:  "system",
	}

	model.SetTemperature(1.02)
	model.SetTopK(40)
	model.SetTopP(0.92)

	userPrompt := g.GetPrompt(userCommit, diff, relatedFilesArray)

	resp, err := model.GenerateContent(
		ctx,
		genai.Text(userPrompt),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return "", nil
	}

	return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
}
