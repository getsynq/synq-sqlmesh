package git

import (
	ingestgitv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/git/v1"
	"context"
	"os/exec"
	"strings"
	"time"
)

func CollectGitContext(ctx context.Context, dir string) *ingestgitv1.GitContext {
	if !commandExists("git") {
		return nil
	}

	return &ingestgitv1.GitContext{
		CloneUrl:  getCloneUrl(ctx, dir),
		Branch:    getCurrentBranch(ctx, dir),
		CommitSha: getCommitSHA(ctx, dir),
	}
}

func getCurrentBranch(ctx context.Context, dir string) string {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getCommitSHA(ctx context.Context, dir string) string {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getCloneUrl(ctx context.Context, dir string) string {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
