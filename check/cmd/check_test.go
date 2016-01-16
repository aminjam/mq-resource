package main_test

import (
	"encoding/json"
	"os/exec"
	"time"

	//. "github.com/aminjam/mq-resource/check/cmd"

	"github.com/aminjam/mq-resource"
	"github.com/aminjam/mq-resource/check"
	"github.com/aminjam/mq-resource/mq-resource-tester"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check", func() {
	var checkCmd *exec.Cmd

	BeforeEach(func() {
		checkCmd = exec.Command(checkPath)
	})
	Context("with a plugin", func() {
		var request check.Request
		var response check.Response
		var messages []string

		BeforeEach(func() {
			messages = []string{}
			request = check.Request{
				Version: resource.StringMap{},
				Source:  source,
			}
			response = check.Response{}
		})

		JustBeforeEach(func() {
			stdin, err := checkCmd.StdinPipe()
			Expect(err).ToNot(HaveOccurred())

			session, err := gexec.Start(checkCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).ToNot(HaveOccurred())

			for _, v := range messages {
				_, err = mqResourceTester.PutMessage([]byte(v))
				Expect(err).ToNot(HaveOccurred())
			}
			//wailt for 10 second
			Eventually(session, 10*time.Second).Should(gexec.Exit(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("with no version", func() {
			BeforeEach(func() {
				request.Version = resource.StringMap{}
			})

			Context("when a version is present in the source", func() {
				BeforeEach(func() {
					messages = []string{"{\"name\":\"john\"}"}
				})

				It("returns the version present at the source", func() {
					Expect(response).To(HaveLen(1))
					Expect(response[0]["name"]).To(Equal("john"))
				})
			})
		})

		Context("with a version present", func() {
			BeforeEach(func() {
				request.Version = resource.StringMap{
					"name": "johnson",
				}
			})

			Context("when there is no current version", func() {
				It("outputs an empty list", func() {
					Expect(response).To(HaveLen(0))
				})
			})

			Context("when the source has a higher version", func() {
				BeforeEach(func() {
					messages = []string{"{\"name\":\"john\"}"}
				})

				It("returns the version present at the source", func() {
					Expect(response).To(HaveLen(1))
					Expect(response[0]["name"]).Should(Equal("john"))
				})
			})

			Context("when the source has multiple new versions", func() {
				BeforeEach(func() {
					messages = []string{"{\"name\":\"john\"}", "{\"name\":\"jacky\"}"}
				})

				It("returns the version present at the source", func() {
					Expect(response).To(HaveLen(2))
					Expect(response[0]["name"]).Should(Equal("john"))
					Expect(response[1]["name"]).Should(Equal("jacky"))
				})
			})

			Context("when it's the same as the current version", func() {
				BeforeEach(func() {
					messages = []string{"{\"name\":\"johnson\"}"}
				})

				It("outputs an empty list", func() {
					Expect(response).To(HaveLen(0))
				})
			})
		})
	})
})
