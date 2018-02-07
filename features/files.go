package features

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//Generate - creates a map of all the keys that can be overridden based on order
func Generate(validateKeys bool, filePaths []string, expectedKeys []string) (map[string]interface{}, error) {
	valueMap := make(map[string]interface{})
	for _, filePath := range filePaths {
		if _, err := os.Stat(filePath); err == nil {
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}
			yaml.Unmarshal(data, &valueMap)
		}
	}
	if validateKeys {
		var missingKeys []string
		expectedMap := make(map[string]string)
		for _, expectedKey := range expectedKeys {
			expectedMap[expectedKey] = expectedKey
			if _, ok := valueMap[expectedKey]; !ok {
				missingKeys = append(missingKeys, expectedKey)
			}
		}
		if len(missingKeys) > 0 {
			return nil, fmt.Errorf("Missing keys %s", missingKeys)
		}

		var unexpectedKeys []string
		for mapKey := range valueMap {
			if _, ok := expectedMap[mapKey]; !ok {
				unexpectedKeys = append(unexpectedKeys, mapKey)
			}
		}
		if len(unexpectedKeys) > 0 {
			return nil, fmt.Errorf("Keys %s were provided but not expected", unexpectedKeys)
		}
	}
	return valueMap, nil
}

//Compile - takes in ordered slice of file paths and produces a output feature file
func Compile(validateKeys bool, inputFilePaths []string, expectedKeyPath string, outputFilePath string) error {
	var expectedKeys ExpectedKeys
	if validateKeys {
		data, err := ioutil.ReadFile(expectedKeyPath)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(data, &expectedKeys)
		if err != nil {
			return err
		}
	}
	output, err := Generate(validateKeys, inputFilePaths, expectedKeys.Keys)
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(output)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(outputFilePath, bytes, 0644)
}

type ExpectedKeys struct {
	Keys []string `yaml:"keys"`
}
