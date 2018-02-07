package cmd

type Command struct {
	FeatureFileCommand FeatureFileCommand `command:"compile" description:"compiles a feature file"`
}

var TheCommand Command
