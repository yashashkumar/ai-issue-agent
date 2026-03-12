package spawner

import (
	"context"
	"fmt"
	"os/exec"
)

// RunGit runs a raw git command and returns output + error
func RunGit(ctx context.Context, dir string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %v failed: %w, output: %s", args, err, string(out))
	}
	return string(out), nil
}

func CreateWorktreeBranch(ctx context.Context, baseDir string, branchName string, worktreePath string) error {
	// e.g. baseDir is the project root (where the main repo is)
	// create branch and directly create a linked worktree
	_, err := RunGit(ctx, baseDir, "worktree", "add", "-b", branchName, worktreePath)
	return err
}

func CommitAndPush(ctx context.Context, worktreePath string, message string, branchName string) error {
	_, err := RunGit(ctx, worktreePath, "add", "-A")
	if err != nil {
		return err
	}
	_, err = RunGit(ctx, worktreePath, "commit", "-m", message)
	if err != nil {
		return err // Might be empty commit if no changes, handle safely
	}
	_, err = RunGit(ctx, worktreePath, "push", "origin", branchName)
	return err
}

func CreatePR(ctx context.Context, dir string, title string, body string, baseBranch string, headBranch string, repo string) error {
	cmd := exec.CommandContext(ctx, "gh", "pr", "create",
		"--base", baseBranch,
		"--head", headBranch,
		"--title", title,
		"--body", body,
		"--repo", repo,
	)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh pr create failed: %w, output: %s", err, string(out))
	}
	return nil
}
