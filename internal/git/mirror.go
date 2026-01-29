package git

import (
	"fmt"
	"os"
	"os/exec"
)

func MirrorClone(url, path string) error {
	cmd := exec.Command("git", "clone", "--mirror", url, path)

	// Force SSH + disable password prompts
	cmd.Env = append(os.Environ(),
		"GIT_SSH_COMMAND=ssh -o BatchMode=yes",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("clone failed: %v\n%s", err, out)
	}
	return nil
}

func MirrorPush(repoPath, remote string) error {
	cmd := exec.Command("git", "--git-dir", repoPath, "push", "--mirror", remote)

	cmd.Env = append(os.Environ(),
		"GIT_SSH_COMMAND=ssh -o BatchMode=yes",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("push failed: %v\n%s", err, out)
	}
	return nil
}
