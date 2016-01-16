package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aminjam/mq-resource"
	"github.com/aminjam/mq-resource/out"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Out", func() {
	var outCmd *exec.Cmd
	var destination string

	BeforeEach(func() {
		var err error

		destination, err = ioutil.TempDir("", "out-dest")
		Expect(err).ToNot(HaveOccurred())

		outCmd = exec.Command(outPath, destination)
	})

	AfterEach(func() {
		os.RemoveAll(destination)
	})
	Context("with a plugin", func() {
		var request out.Request
		var response out.Response

		BeforeEach(func() {
			request = out.Request{
				Version: resource.StringMap{},
				Source:  source,
				Params:  out.Params{},
			}
			response = out.Response{}
		})

		JustBeforeEach(func() {
			stdin, err := outCmd.StdinPipe()
			Expect(err).ToNot(HaveOccurred())

			session, err := gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).ToNot(HaveOccurred())

			//wailt for resource.WaitFor + 5 second
			Eventually(session, (resource.WaitFor+5)*time.Second).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("when setting the version", func() {
			BeforeEach(func() {
				request.Params.File = "output.json"
			})

			Context("when a valid version is in the file", func() {
				BeforeEach(func() {
					err := ioutil.WriteFile(filepath.Join(destination, request.Params.File), []byte("{\"name\":\"jane\"}"), 0644)
					Expect(err).ToNot(HaveOccurred())
				})

				It("reports the version as the resource's version", func() {
					Expect(response["name"]).To(Equal("jane"))
				})
			})
		})
	})
})
