package services

import (
	"dynamic-dirb/internal/file_mng"
	"strings"
)

// Local storage of the graph (used before file persist. es. a->b;a->c)
var graph string
var graphName string

// Node queue to analyze (used to retore graph)
var queue []string

// Check if dynamic analysis is finished
var isFish bool

func CreateGraphFromString(graphInput string) {
	graph = graphInput
}

// Create a node and append it to the local variable
func CreateGraphNode(root string, node string) {
	graph = graph + "\t\"" + root + "\"" + " -> " + "\"" + node + "\";\n"
}

// Return the variable completed with header and footer
func GetGraph() string {
	var localGraph = "digraph \"" + GetGraphName() + "\" {\n" + graph
	localGraph = localGraph + "}"
	return localGraph
}

func SetQueueGraph(newQueue []string) {
	queue = newQueue
}

func SetGraphName(name string) {
	graphName = name
}

func GetGraphName() string {
	return graphName
}

func GetQueueGraph() []string {
	return queue
}

func SetIsFinish(isFinish bool) {
	isFish = isFinish
}

func GetIsFinish() bool {
	return isFish
}

// Persist the graph in the filesystem if user set -graph flag
func PersistFileCompleteGraph(fileName string, graph string) bool {
	return file_mng.Append(fileName, graph)
}

func PersistFileNodeGraph(fileName string, node string) bool {
	return file_mng.Append(fileName, node)
}

// NOT USED YET
// Write the ending symbol for the graph file
func writeGraphEnding(fileName string) bool {
	if !file_mng.FileExists(fileName) {
		return false
	}
	return file_mng.Append(fileName, "}")
}

// NOT USED YET
func writeGraphNodes(fileName string, root string, node []string) bool {
	for _, childNode := range node {
		data := "\"" + root + "\"" + " -> " + "\"" + childNode + "\""
		if !file_mng.Append(fileName, data) {
			return false
		}
	}
	return true
}

func PersistFileRestoreGraph(fileName string, graph string, queue []string) bool {
	file_mng.EmptyFileOverride(fileName)
	content := graph + "\n-------------\n" + strings.Join(queue, "\n")
	return file_mng.Append(fileName, content)
}
