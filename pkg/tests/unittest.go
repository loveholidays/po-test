package tests

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func RunUnitTests(testFiles []string) error {
	var originalRules []*filenameAndData
	for _, testFile := range testFiles {
		b, err := os.ReadFile(testFile)
		if err != nil {
			return err
		}

		var unitTestInp unitTestFile
		if err := yaml.Unmarshal(b, &unitTestInp); err != nil {
			return err
		}

		for _, rulesFile := range unitTestInp.RuleFiles {
			relativeRulesFile := fmt.Sprintf("%s/%s", filepath.Dir(testFile), rulesFile)

			yamlFile, err := os.ReadFile(relativeRulesFile)
			if err != nil {
				return err
			}

			unstructured := make(map[interface{}]interface{})
			err = yaml.Unmarshal(yamlFile, &unstructured)
			if err != nil {
				return err
			}

			if spec, found := unstructured["spec"]; found {
				ruleFileContentWithoutMetadata, err := yaml.Marshal(spec)
				if err != nil {
					return err
				}

				originalRules = append(originalRules, &filenameAndData{relativeRulesFile, yamlFile})

				err = os.WriteFile(relativeRulesFile, ruleFileContentWithoutMetadata, 0o600)
				if err != nil {
					return err
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
