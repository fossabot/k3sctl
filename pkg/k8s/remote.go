package k8s

import (
	"fmt"
)

// RemoteKubeconfigReader can fetch a file's contents from a remote host.
type RemoteKubeconfigReader interface {
	ExecOutput(cmd string) (string, error)
}

// remoteReaderAdapter wraps an object with ExecOutput(string, ...opts) to
// satisfy RemoteKubeconfigReader.
type remoteReaderAdapter struct {
	fn func(string) (string, error)
}

func (a *remoteReaderAdapter) ExecOutput(cmd string) (string, error) {
	return a.fn(cmd)
}

// NewFromRemoteHost fetches the kubeconfig from a remote K3s server and
// builds a Client. The execOutput function should execute a command on the
// remote host and return stdout (e.g. leader.Sudo().ExecOutput).
func NewFromRemoteHost(execOutput func(string) (string, error), kubeconfigPath, apiAddress string) (Client, error) {
	output, err := execOutput(fmt.Sprintf("cat '%s'", kubeconfigPath))
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig from remote: %w", err)
	}

	rewritten, err := RewriteKubeconfig([]byte(output), apiAddress, "", "")
	if err != nil {
		return nil, fmt.Errorf("rewrite kubeconfig: %w", err)
	}

	return NewFromKubeconfig(rewritten)
}
