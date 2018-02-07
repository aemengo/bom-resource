package cmd

import (
	"github.com/pivotalservices/bom-resource/features"
)

type FeatureFileCommand struct {
	InputFiles      []string `long:"input-file" description:"input yaml feature file, can be specified multiple times" required:"true"`
	ValidateKeys    bool     `long:"validate-keys" description:"will enforce that you have provided a list of expected keys"`
	ExpectedKeyPath string   `long:"expected-keys-file" description:"path to expected keys yaml file" required:"false"`
	OutputFile      string   `long:"output-file" description:"output file path for yaml feature file" required:"true"`
}

//Execute - produces a feature file
func (c *FeatureFileCommand) Execute([]string) error {
	return features.Compile(c.ValidateKeys, c.InputFiles, c.ExpectedKeyPath, c.OutputFile)
}
