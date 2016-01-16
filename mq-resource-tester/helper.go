package mqResourceTester

import (
	"os"
	"os/exec"
	"strings"
)

func PutMessage(message []byte) (string, error) {
	cmd := exec.Command("/put-message", string(message))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func OsEnvs() map[string]string {
	output := make(map[string]string)
	for _, v := range os.Environ() {
		splits := strings.Split(v, "=")
		output[splits[0]] = strings.Join(splits[1:], "=")
	}
	return output
}
