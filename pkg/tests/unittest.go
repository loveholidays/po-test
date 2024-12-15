package tests

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func RunUnitTests(testFiles []string) error {
	var originalRules []*filenameAndData

	for _, testFile := range testFiles {
		unitTestInp, err := parseUnitTestFile(testFile)
		if err != nil {
			return err
		}

		for _, rulesFile := range unitTestInp.RuleFiles {
			err := processRulesFile(testFile, rulesFile, &originalRules)
			if err != nil {
				return err
			}
		}
	}

	if err := runPromtoolTests(testFiles); err != nil {
		restoreOriginalFiles(originalRules)
		return err
	}

	restoreOriginalFiles(originalRules)
	return nil
}

func parseUnitTestFile(testFile string) (*unitTestFile, error) {
	b, err := os.ReadFile(testFile)
	if err != nil {
		return nil, err
	}

	var unitTestInp unitTestFile
	if err := yaml.Unmarshal(b, &unitTestInp); err != nil {
		return nil, err
	}

	return &unitTestInp, nil
}

func processRulesFile(testFile, rulesFile string, originalRules *[]*filenameAndData) error {
	relativeRulesFile := fmt.Sprintf("%s/%s", filepath.Dir(testFile), rulesFile)

	yamlFile, err := os.ReadFile(relativeRulesFile)
	if err != nil {
		return err
	}

	unstructured := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(yamlFile, &unstructured); err != nil {
		return err
	}

	if spec, found := unstructured["spec"]; found {
		ruleFileContentWithoutMetadata, err := yaml.Marshal(spec)
		if err != nil {
			return err
		}

		*originalRules = append(*originalRules, &filenameAndData{relativeRulesFile, yamlFile})

		if err := os.WriteFile(relativeRulesFile, ruleFileContentWithoutMetadata, 0o600); err != nil {
			return err
		}
	} else {
		log.Printf("No spec found in file %s", rulesFile)
	}

	return nil
}

func runPromtoolTests(testFiles []string) error {
	promtoolArgs := append([]string{"test", "rules"}, testFiles...)
	command := exec.Command("promtool", promtoolArgs...)
	output, err := command.CombinedOutput()
	log.Printf("%s", output)
	if err != nil {
		return err
	}
	return nil
}

func restoreOriginalFiles(rules []*filenameAndData) {
	for _, nameAndData := range rules {
		err := os.WriteFile(nameAndData.filename, nameAndData.data, 0o600)
		if err != nil {
			log.Fatalf("Failed to write file: %v", err)
		}
	}
}

type filenameAndData struct {
	filename string
	data     []byte
}

type unitTestFile struct {
	RuleFiles []string `yaml:"rule_files"`
}
