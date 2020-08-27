package main

import (
	"dynamic-dirb/internal/dynamic"
	"dynamic-dirb/internal/helper"
	service "dynamic-dirb/internal/services"
	"dynamic-dirb/internal/static"
	"dynamic-dirb/internal/string_mng"
	"dynamic-dirb/internal/web_mng"
	"flag"
	_ "os"
	"sync"
)

// set $GOPATH to ddirb then go build src/main/main.go then ./main
// Compiling for raspberrys env GOOS=linux GOARCH=arm go build -o dynamic-dirb_linux ./dynamic-dirb.go
// Compilinf for kali env GOOS=linux GOARCH=amd64 go build -o dynamic-dirb_linux ./dynamic-dirb.go

func main() {

	// flag.String returns *string (string var that stores the value of the flag)
	url := flag.String("url", "", "Target website")
	outputFile := flag.String("output", "/tmp/ddirb-output", "Output file path")
	overrideFile := flag.Bool("override", false, "Override output file path")
	restoreFile := flag.String("restoreFile", "/tmp/ddirb-output.restore", "Input file [Dynamic mode]")
	threadsCount := flag.Int("threads", 1, "Threads count")
	delay := flag.String("delay", "250", "Delay between requests in milliseconds")
	methods := flag.String("methods", "GET", "HTTP methods to use (HEAD/GET) [Not used yet]")
	recursive := flag.Bool("recursive", false, "Recursive search [Static mode]")
	debugFlag := flag.Bool("debug", false, "Debug information enabled")
	graph := flag.Bool("graph", false, "Print graph in dot lang [Dynamic mode]")
	scanType := flag.String("type", "", "Scan type dynamic/static/resumeDynamic")
	wordlist := flag.String("wordlist", "", "Wordlist file [Static mode]")
	headers := flag.String("headers", "", "Header (es. Header1:value1;value2,Header2:value1)")
	exclusionPath := flag.String("exclusionPath", "", "Path to exclude es. logout")
	//port := flag.Int("port", 8081, "Http listener port default: 8081")
	port := 5689
	flag.Usage = helper.PrintUsage //override flag.Usage
	string_mng.ColorizeErrorStart()
	flag.Parse()

	GlobalParameters := new(helper.ParamValidator)
	GlobalParameters.Init(url, outputFile, overrideFile, threadsCount, delay, methods, recursive, debugFlag, graph, scanType,
		wordlist, headers, restoreFile, &port, exclusionPath)
	service.SetParameters(GlobalParameters)

	switch service.GetParameters().GetScanType() {
	case "dynamic":
		service.GetParameters().PrintDebug("-dynamic case")
		// set up a ctrl+c handler
		dynamic.SetupCloseHandler()
		var wg sync.WaitGroup
		wg.Add(1)
		//spawn web server
		go web_mng.StartWebServer()
		// compute dynamic web scraping
		dynamic.Dynamic()
		wg.Wait()
		break
	case "static":
		service.GetParameters().PrintDebug("-static case")
		static.Static()
		break
	case "resumeDynamic":
		service.GetParameters().PrintDebug("-resumeDynamic case")
		// set up a ctrl+c handler
		dynamic.SetupCloseHandler()
		var wg sync.WaitGroup
		wg.Add(1)
		//spawn web server
		go web_mng.StartWebServer()
		dynamic.RestoreDynamicExecution()
		wg.Wait()
		break
	default:
		string_mng.PrintError("[!] Invalid option detected:")
		string_mng.PrintError("[*]\tscan type")
		helper.PrintUsage()
		// os.Exit(0)
	}
}
