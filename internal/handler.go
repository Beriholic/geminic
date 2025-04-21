package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beriholic/geminic/internal/config"
	"github.com/beriholic/geminic/internal/model"
	"github.com/beriholic/geminic/internal/service"
	"github.com/beriholic/geminic/internal/ui"
	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GeneratorCommit(ctx context.Context, userCommit string) error {
	cfg := config.Get()
	if err := config.Verify(); err != nil {
		return err
	}

	client, err := genai.NewClient(
		ctx,
		option.WithAPIKey(cfg.Key),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	gitService := service.GetGitService()

	if err := gitService.VerifyGitInstallation(); err != nil {
		return err
	}
	if err := gitService.VerifyGitRepository(); err != nil {
		return err
	}

	files, diff, err := gitService.DetectDiffChanges()

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf(
			"no staged changes found. stage your changes manually",
		)
	} else {
		fmt.Printf("Detected %v staged file:\n", len(files))
	}

	relatedFiles := getRelatedFiles(files)

	for {
		geminiService := service.GetGeminiService()

		errChan := make(chan error, 1)
		genCommitChan := make(chan string, 1)

		err := ui.RenderSpinner("Generating commit message...", func() {
			genCommit, err := geminiService.AnalyzeChanges(
				userCommit,
				client,
				ctx,
				diff,
				&relatedFiles,
				&cfg.Model,
			)
			errChan <- err
			genCommitChan <- genCommit
		})
		if err != nil {
			return err
		}

		rawGenCommit, err := <-genCommitChan, <-errChan
		if err != nil {
			return err
		}

		var gitCommit model.GitCommit
		json.Unmarshal([]byte(rawGenCommit), &gitCommit)

		genCommit := gitCommit.String()

		fmt.Println(ui.FormatText("Generated commit message", genCommit))

		action, err := ui.RenderActionForm()
		if err != nil {
			return err
		}

		switch action {
		case ui.CONFIRM:
			fmt.Println("committed")
			return gitService.CommitChanges(genCommit)
		case ui.REGENERATE:
			continue
		case ui.EDIT_COMMIT:
			_, err := ui.RenderEditorForm(genCommit)
			if err != nil {
				return err
			}
			return nil
		case ui.CANCEL:
			fmt.Println("cancelled")
			return nil
		default:
			fmt.Println("invalid action")
			return nil
		}
	}
}

func getRelatedFiles(files []string) map[string]string {
	relatedFiles := make(map[string]string)
	visitedDirs := make(map[string]bool)

	for idx, file := range files {
		color.New(color.Bold).Printf("\t%d. %s\n", idx+1, file)

		dir := filepath.Dir(file)
		if !visitedDirs[dir] {
			lsEntry, err := os.ReadDir(dir)
			if err == nil {
				var ls []string
				for _, entry := range lsEntry {
					ls = append(ls, entry.Name())
				}
				relatedFiles[dir] = strings.Join(ls, ", ")
				visitedDirs[dir] = true
			}
		}
	}

	return relatedFiles
}
