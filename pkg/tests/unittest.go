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
		b, err := os.ReadFile(testFile)
		if err != nil {
			return err
		}

		var unitTestInp unitTestFile
		unmarshalErr := yaml.Unmarshal(b, &unitTestInp)
		if unmarshalErr != nil {
			return unmarshalErr
		}

		for _, rulesFile := range unitTestInp.RuleFiles {
			relativeRulesFile := fmt.Sprintf("%s/%s", filepath.Dir(testFile), rulesFile)

			yamlFile, err := os.ReadFile(relativeRulesFile)
			if err != nil {
				return err
			}

			unstructured := make(map[interface{}]interface{})
			unmarshalErr := yaml.Unmarshal(yamlFile, &unstructured)
			if unmarshalErr != nil {
				return unmarshalErr
			}

			if spec, found := unstructured["spec"]; found {
				ruleFileContentWithoutMetadata, marshalErr := yaml.Marshal(spec)
				if marshalErr != nil {
					return marshalErr
				}

				originalRules = append(originalRules, &filenameAndData{relativeRulesFile, yamlFile})

				writeErr := os.WriteFile(relativeRulesFile, ruleFileContentWithoutMetadata, 0o600)
				if writeErr != nil {
					return writeErr
				}
			} else {
				log.Printf("No spec found in file %s", rulesFile)
			}
		}
	}

	promtoolArgs := append([]string{"test", "rules"}, testFiles...)
	command := exec.Command("promtool", promtoolArgs...)
	output, err := command.CombinedOutput()
	if err != nil {
		log.Printf("%s", output)
		restoreOriginalFiles(originalRules)
		return err
	}
	log.Printf("%s", output)
	restoreOriginalFiles(originalRules)
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
