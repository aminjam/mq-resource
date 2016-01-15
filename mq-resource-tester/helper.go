package mqResourceTester

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func PutMessage(message []byte) error {
	endpoint := OsEnvs()["PUT_MESSAGE_ENDPOINT"]
	reader := bytes.NewReader(message)
	res, err := http.Post(endpoint, "application/json", reader)
	if err != nil {
		return err
	}
	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
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
