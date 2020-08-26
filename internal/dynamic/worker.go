package dynamic

import (
	"dynamic-dirb/internal/helper"
	"dynamic-dirb/internal/http_mng"
	service "dynamic-dirb/internal/services"
	"dynamic-dirb/internal/string_mng"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func Dynamic() {
	service.GetParameters().PrintDebug("Validating parameters...")
	if service.GetParameters().ValidateDynamic() {
		service.GetParameters().PrettyPrintParameters()
	} else {
		string_mng.PrintError("[-] Invalid parameters")
		helper.PrintDynamicUsage()
		os.Exit(0)
	}
	service.GetParameters().PrintDebug("...Parameters validated")

	worker()
}

// url(string): url target
func worker() {
	service.GetParameters().PrintDebug("Worker instantiated...")
	url := service.GetParameters().GetUrl()

	response, code := http_mng.RequestResourceExist(url)
	//string_mng.PrintNotice(response)
	service.GetParameters().PrintDebug("Testing if target reachable...")
	if !code {
		string_mng.PrintError("Get request response with code: " + strconv.Itoa(response.ResponseCode) + "\n")
		for key, _ := range response.ResponseHeader {
			string_mng.PrintWarning(key + " " + response.ResponseHeader[key][0])
		}
		os.Exit(1)
	}
	/*parameters.PrintDebug("...Target reachable!")
	parameters.PrintDebug("Starting DFS...")
	result := DFS(parameters, parameters.Url(), map[string]bool{})
	parameters.PrintDebug("...DFS is done!")
	string_mng.PrintNotice("Finish\n")
	fmt.Println(result)*/

	service.GetParameters().PrintDebug("...Target reachable!")
	var startUrl []string
	startUrl = append(startUrl, service.GetParameters().GetUrl())
	startResult := map[string]bool{service.GetParameters().GetUrl(): true}
	service.GetParameters().PrintDebug("Starting BFS...")
	service.SetGraphName(service.GetParameters().GetUrl())
	_ = BFS(startUrl, startResult)
	service.GetParameters().PrintDebug("...BFS is done!")
	string_mng.PrintNotice("\nFinish\n")
	service.GetParameters().PrintDebug("...Worker is done!")

}

// target(targetInfo): target information
// result(map[string]bool): array containing the found urls for recursion
// return(map[string]bool): array containing the found urls
func DFS(target *helper.ParamValidator, url string, result map[string]bool) map[string]bool {

	// request to the target url
	response := http_mng.RequestGet(url)
	if response.ResponseBodyString == "" {
		fmt.Println("Null response for: " + url)
		return result
	}

	// discovering js script and href url and validating it
	discoveredUrl := GetFinalUrls(response, result)
	//string_mng.UpdateLine(">" + url + "")

	// if url is already present not use it
	if result[url] {
		//string_mng.PrintInfo("[-] Already present: " + url + "\n")
		return result
	}

	// if url is not present add it to result and do the recursion
	result[url] = true
	//string_mng.PrintInfo("[+] Added: " + url + "\n")
	for _, newUrl := range discoveredUrl {

		DFS(target, newUrl, result)
	}
	return result
}

// target(helper.ParamValidator): target information
// queue([]string): array containing the queue
// result(map[string]bool): array containing the found urls
// return(map[string]bool): array containing the found urls
func BFS(queue []string, result map[string]bool) map[string]bool {
	service.SetIsFinish(false)

	for (len(queue)) != 0 {

		url := queue[0]
		queue = queue[1:]

		service.GetParameters().PrintDebug("Analizing url: " + url)

		// request to the target url
		response := http_mng.RequestGet(url)
		if response.ResponseBodyString == "" {
			fmt.Println("Null response for: " + url)
		}

		// discovering js script and href url and validating it
		discoveredUrl := GetFinalUrls(response, result)

		string_mng.ClearProgress()
		string_mng.PrintDynamicScanStatus(len(queue), url, len(discoveredUrl))

		for _, newUrl := range discoveredUrl {
			// if node is node visited yet
			if !result[newUrl] {
				result[newUrl] = true
				queue = append(queue, newUrl)
				// Add node to the .dot notation graph
				service.CreateGraphNode(url, newUrl)
				// Back-up queue in case of restore
				service.SetQueueGraph(queue)
				// Print to file the new urls
				service.PersistFileNodeGraph(service.GetParameters().GetOutputFile(), newUrl)
				// Update output service
				service.SetNewNodeOutput(newUrl)
			}
		}
	}

	service.SetIsFinish(true)
	// Check if this return cause problem to recursion
	return result
}

// Handle the SIGTERM Signal
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		parameters := service.GetParameters()
		// if user want to save the .dot output the program will save it
		if parameters.IsGraph() {
			service.PersistFileCompleteGraph(parameters.GetOutputFileGraph(), service.GetGraph())
		}

		// if dynamic analysis not finished save restore file
		if !service.GetIsFinish() {
			service.PersistFileRestoreGraph(parameters.GetOutputFileRestore(), service.GetGraph(), service.GetQueueGraph())
		}

		os.Exit(0)
	}()
}
