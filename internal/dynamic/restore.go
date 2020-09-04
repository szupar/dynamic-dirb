package dynamic

import (
	"dynamic-dirb/internal/file_mng"
	"dynamic-dirb/internal/helper"
	service "dynamic-dirb/internal/services"
	"dynamic-dirb/internal/string_mng"
	"os"
	"strings"

	"github.com/awalterschulze/gographviz"
)

func RestoreDynamicExecution() {
	service.GetParameters().PrintDebug("Validating parameters...")
	if service.GetParameters().ValidateRestoreDynamic() {
		service.GetParameters().PrettyPrintParameters()
	} else {
		string_mng.PrintError("[-] Invalid parameters")
		helper.PrintResumeDynamicUsage()
		os.Exit(0)
	}

	service.GetParameters().PrintDebug("...Parameters validated")

	service.GetParameters().PrintDebug("Starting BFS...")
	// Restoring queue and result for BFS
	queueString, startResult := restoreGraph()
	// Restoring output file
	_ = restoreOutputFile(queueString, startResult)

	_ = BFS(queueString, startResult, make(map[string]bool))
	service.GetParameters().PrintDebug("...BFS is done!")
	string_mng.PrintNotice("\nFinish\n")
	service.GetParameters().PrintDebug("...Worker is done!")
}

// Restoring queue and result from .restore file
func restoreGraph() ([]string, map[string]bool) {
	g := gographviz.NewGraph()
	fileLine := file_mng.ReadFileByLine(service.GetParameters().GetInputRestoreFile())

	separator := 0
	for i, line := range fileLine {
		if line == "-------------" {
			separator = i
		}
	}

	// Get array of lines containing only the graph (first line: diagraph <name> { , last line: })
	graphArray := fileLine[0:separator]
	// Graph array to string
	graphString := strings.Join(graphArray, "")

	// Get queue array from file
	queueString := fileLine[separator+1 : len(fileLine)]

	// Import the graph string
	g, _ = gographviz.Read([]byte(graphString))

	// Create the result start for BFS visit (by getting all graph node)
	startResult := make(map[string]bool)
	// Getting all graph node
	for _, value := range g.Nodes.Nodes {
		// Remove " char from node name
		node := value.Name[1 : len(value.Name)-1]
		startResult[node] = true
	}
	// Set the graph name
	service.SetGraphName(g.Name[1 : len(g.Name)-1])
	service.SetIsFinish(false)
	service.GetParameters().SetUrlDomain(g.Name[1 : len(g.Name)-1])
	// Getting the edges from graphArray (by removing graph name and } char)
	graphEdgeString := graphArray[1 : len(graphArray)-1]
	service.CreateGraphFromString(strings.Join(graphEdgeString, "\n"))

	return queueString, startResult
}

// Restoring new ouput file using queue and result
func restoreOutputFile(queueString []string, startResult map[string]bool) string {
	outputFileMap := startResult
	// Merging queue and result in a unique map
	for _, value := range queueString {
		if !outputFileMap[value] {
			outputFileMap[value] = true
		}
	}

	// Converting map to array
	var outputFileArray []string
	for key, _ := range outputFileMap {
		outputFileArray = append(outputFileArray, key)
	}

	outputFileString := strings.Join(outputFileArray, "\n")
	file_mng.Append(service.GetParameters().GetOutputFile(), outputFileString)

	// Retoring output service
	service.SetOutput(outputFileArray)
	return outputFileString

}
