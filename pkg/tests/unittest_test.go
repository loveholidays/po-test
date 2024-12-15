package tests_test

import (
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sahma19/po-test/pkg/tests"
)

func TestPoTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Po-test Suite")
}

var _ = Describe("Po-test", func() {
	Context("Success", func() {
		It("Should mutate files and run unit tests", func() {
			testFilename := "prometheus-operator-unittest.yml"
			ruleFilename := "prometheus-operator-rules.yml"
			Expect(tests.RunUnitTests([]string{testFilename})).To(Succeed())

			file, err := os.ReadFile(ruleFilename)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(file)).To(ContainSubstring("PrometheusRule"))
		})

		It("Should run tests in relative paths", func() {
			testFilename := "subdir/prometheus-operator-unittest.yml"
			ruleFilename := "subdir/prometheus-operator-rules-subdir.yml"
			Expect(tests.RunUnitTests([]string{testFilename})).To(Succeed())

			file, err := os.ReadFile(ruleFilename)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(file)).To(ContainSubstring("PrometheusRule"))
		})
	})

	Context("Failure", func() {
		It("Should report error when tests fail", func() {
			testFilename := "bad-rules-error-test.yml"
			ruleFilename := "bad-rules-error.yml"
			err := tests.RunUnitTests([]string{testFilename})
			Expect(err).To(HaveOccurred())

			file, err := os.ReadFile(ruleFilename)
			Expect(err).ToNot(HaveOccurred())

			Expect(string(file)).To(ContainSubstring("PrometheusRule"))

		})
	})
})
