package main_test

import (
	"testing"

	"github.com/aminjam/mq-resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var source resource.Source
var inPath string

var _ = BeforeSuite(func() {
	var err error
	inPath, err = gexec.Build("github.com/aminjam/mq-resource/in/cmd")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}
