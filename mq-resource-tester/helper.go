package mqResourceTester

import (
	"os"
	"os/exec"
	"strings"
)

func PutMessage(message []byte) error {
	cmd := exec.Command("/put-message", string(message))
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	return err
}

func OsEnvs() map[string]string {
	output := make(map[string]string)
	for _, v := range os.Environ() {
		splits := strings.Split(v, "=")
		output[splits[0]] = strings.Join(splits[1:], "=")
	}
	return output
}
