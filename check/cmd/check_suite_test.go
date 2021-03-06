package main_test

import (
	"encoding/json"

	"testing"

	"github.com/aminjam/mq-resource"
	"github.com/aminjam/mq-resource/check"
	"github.com/aminjam/mq-resource/mq-resource-tester"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var source resource.Source
var checkPath string

var _ = BeforeSuite(func() {
	sourceStr := mqResourceTester.OsEnvs()["SOURCE"]
	var req check.Request
	err := json.Unmarshal([]byte(sourceStr), &req)
	Expect(err).ToNot(HaveOccurred())
	source = req.Source

	checkPath, err = gexec.Build("github.com/aminjam/mq-resource/check/cmd")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Suite")
}
