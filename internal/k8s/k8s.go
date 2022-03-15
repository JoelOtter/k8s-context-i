package k8s

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Context struct {
	Name    string
	Current bool
}

func GetContexts() ([]Context, error) {
	k8sCmd := exec.Command("kubectl", "config", "get-contexts", "--no-headers")
	k8sCmd.Stderr = os.Stderr
	output, err := k8sCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s contexts: %w", err)
	}
	var contexts []Context
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		ctxText := scanner.Text()
		fields := strings.Fields(ctxText)
		if strings.HasPrefix(fields[0], "*") {
			contexts = append(contexts, Context{
				Name:    strings.TrimSpace(fields[1]),
				Current: true,
			})
		} else {
			contexts = append(contexts, Context{
				Name:    strings.TrimSpace(fields[0]),
				Current: false,
			})
		}
	}
	return contexts, nil
}

func ChangeContext(ctx string, output io.Writer) error {
	k8sCmd := exec.Command("kubectl", "config", "use-context", ctx)
	k8sCmd.Stdout = output
	k8sCmd.Stderr = output
	if err := k8sCmd.Run(); err != nil {
		return fmt.Errorf("failed to switch to context %s: %w", ctx, err)
	}
	return nil
}
