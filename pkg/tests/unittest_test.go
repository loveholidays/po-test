/*
po-test
Copyright (C) 2023 loveholidays

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU Lesser General Public
License as published by the Free Software Foundation; either
version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program; if not, write to the Free Software Foundation,
Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

package tests_test

import (
	"github.com/loveholidays/po-test/pkg/tests"
	"os"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestPoTest(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Po-test Suite")
}

var _ = ginkgo.Describe("Po-test", func() {
	ginkgo.Context("Success", func() {
		ginkgo.It("Should mutate files and run unit tests", func() {
			testFilename := "prometheus-operator-unittest.yml"
			ruleFilename := "prometheus-operator-rules.yml"
			gomega.Expect(tests.RunUnitTests([]string{testFilename})).To(gomega.Succeed())

			file, err := os.ReadFile(ruleFilename)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			gomega.Expect(string(file)).To(gomega.ContainSubstring("PrometheusRule"))
		})

		ginkgo.It("Should run tests in relative paths", func() {
			testFilename := "subdir/prometheus-operator-unittest.yml"
			ruleFilename := "subdir/prometheus-operator-rules-subdir.yml"
			gomega.Expect(tests.RunUnitTests([]string{testFilename})).To(gomega.Succeed())

			file, err := os.ReadFile(ruleFilename)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			gomega.Expect(string(file)).To(gomega.ContainSubstring("PrometheusRule"))
		})
	})

	ginkgo.Context("Failure", func() {
		ginkgo.It("Should report error when tests fail", func() {
			testFilename := "bad-rules-error-test.yml"
			ruleFilename := "bad-rules-error.yml"
			err := tests.RunUnitTests([]string{testFilename})
			gomega.Expect(err).To(gomega.HaveOccurred())

			file, err := os.ReadFile(ruleFilename)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			gomega.Expect(string(file)).To(gomega.ContainSubstring("PrometheusRule"))

		})
	})
})
