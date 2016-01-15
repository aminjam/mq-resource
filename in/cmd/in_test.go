package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/aminjam/mq-resource"
	"github.com/aminjam/mq-resource/in"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("In", func() {
	var inCmd *exec.Cmd
	var tmpdir string
	var destination string

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).ToNot(HaveOccurred())
		destination = path.Join(tmpdir, "in-dir")
		inCmd = exec.Command(inPath, destination)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request in.Request
		var response in.Response

		BeforeEach(func() {
			request = in.Request{
				Version: resource.StringMap{
					"type": "DELIVERED_PRE_REQ",
					"sha":  "123", "repo": "github.com/user/repo"},
				Source: resource.Source{},
			}
			response = in.Response{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).ToNot(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).ToNot(HaveOccurred())

			//wailt for 10 second
			Eventually(session, 10*time.Second).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).ToNot(HaveOccurred())
		})

		It("writes the version to the destination", func() {
			contents, err := ioutil.ReadFile(path.Join(destination, "message.json"))
			Expect(err).ToNot(HaveOccurred())
			Expect(string(contents)).To(ContainSubstring(request.Version["type"]))
			Expect(string(contents)).To(ContainSubstring(request.Version["sha"]))
			Expect(string(contents)).To(ContainSubstring(request.Version["repo"]))
		})
		Context("when writing to a custom file", func() {
			BeforeEach(func() {
				request = in.Request{
					Source:  resource.Source{},
					Version: resource.StringMap{"sha": "431"},
					Params:  in.Params{File: "custom.json"},
				}
			})
			It("writes the version to the destination", func() {
				contents, err := ioutil.ReadFile(path.Join(destination, "custom.json"))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(contents)).To(ContainSubstring(request.Version["sha"]))
			})
		})
	})
})
