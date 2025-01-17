package service

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type GitService struct{}

var (
	gitServer     *GitService
	gitServerOnce sync.Once
)

func GetGitService() *GitService {
	gitServerOnce.Do(func() {
		gitServer = &GitService{}
	})

	return gitServer
}

func (g *GitService) VerifyGitInstallation() error {
	if err := exec.Command("git", "--version").Run(); err != nil {
		return fmt.Errorf("git is not installed. %v", err)
	}

	return nil
}

func (g *GitService) VerifyGitRepository() error {
	if err := exec.Command("git", "rev-parse", "--show-toplevel").Run(); err != nil {
		return fmt.Errorf("current directory is not a git repository. %v", err)
	}

	return nil
}

func (g *GitService) DetectDiffChanges() ([]string, string, error) {
	files, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal", "--name-only").
		Output()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, "", err
	}
	filesStr := strings.TrimSpace(string(files))

	if filesStr == "" {
		return nil, "", fmt.Errorf("no changes detected")
	}

	diff, err := exec.Command("git", "diff", "--cached", "--diff-algorithm=minimal").Output()

	if err != nil {
		return nil, "", err
	}

	return strings.Split(filesStr, "\n"), string(diff), nil
}

func (g *GitService) CommitChanges(message string) error {
	_, err := exec.Command("git", "commit", "-m", message).Output()
	if err != nil {
		return fmt.Errorf("failed to commit changes. %v", err)
	}

	return nil
}
