package services

var output []string

func SetOutput(newOutput []string) {
	output = newOutput
}

func SetNewNodeOutput(newOutput string) {
	output = append(output, newOutput)
}

func GetOutput() []string {
	return output
}
