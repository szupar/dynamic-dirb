package helper

import (
	"dynamic-dirb/internal/file_mng"
	"dynamic-dirb/internal/string_mng"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ParamValidator struct {
	url               string
	urlTarget         string
	outputFile        string
	outputFileGraph   string
	outputFileRestore string
	threads           int
	delay             time.Duration
	methods           string
	recursive         bool
	domain            string
	debugFlag         bool
	graph             bool
	scanType          string
	wordlist          string
	headers           string
	restoreFile       string
	port              int
}

// Init method to initialize the struct
func (pv *ParamValidator) Init(url *string,
	outputFile *string,
	threads *int,
	delay *int,
	methods *string,
	recursive *bool,
	debugFlag *bool,
	graph *bool,
	scanType *string,
	wordlist *string,
	headers *string,
	restoreFile *string,
	port *int) {

	pv.url = *url
	if len(pv.url) != 0 {
		if pv.url[len(pv.url)-1] == '/' {
			pv.url = pv.url[:len(pv.url)-1]
		}
	}
	pv.urlTarget = getInitialUrl(*url)
	pv.outputFile = *outputFile
	pv.threads = *threads
	pv.delay = time.Duration(*delay)
	pv.methods = *methods
	pv.recursive = *recursive
	pv.domain = getDomainName(*url)
	pv.debugFlag = *debugFlag
	pv.graph = *graph
	pv.scanType = *scanType
	pv.wordlist = *wordlist
	pv.headers = *headers
	if *graph {
		pv.outputFileGraph = *outputFile + ".dot"
	}
	pv.outputFileRestore = *outputFile + ".restore"
	pv.restoreFile = *restoreFile
	pv.port = *port

}

func PrintUsage() {
	string_mng.PrintNormal("-type:\t\t\tSet -type dynamic/static/resumeDynamic\n")
	//fmt.Fprintf(os.Stderr, "-type:\t\t\tSet -type dynamic/static/resumeDynamic\n\n")
	PrintDynamicUsage()
	PrintStaticUsage()
	PrintResumeDynamicUsage()
}

func PrintDynamicUsage() {
	string_mng.PrintNormal("dynamic Usage:")
	string_mng.PrintNormal("\t-url:\t\tTarget url in http/https format (http://example.com)")
	string_mng.PrintNormal("\t-output:\tOutput file (default /tmp/ddirb-output)")
	//string_mng.PrintNormal("\t-port:\tHttp listener port default: 8080")
	string_mng.PrintNormal("\t-threads:\tNumber of threads to use")
	string_mng.PrintNormal("\t-delay:\t\tDelay in seconds between each thread")
	string_mng.PrintNormal("\t-debug:\t\tFlag to print in verbose mode")
	string_mng.PrintNormal("\t-graph:\t\tFlag to save the graph in .dot language")
	string_mng.PrintNormal("\t-headers:\tSet headers (es. Header1:value1;value2,Header2:value1)")
}

func PrintResumeDynamicUsage() {
	string_mng.PrintNormal("resumeDynamic Usage:")
	string_mng.PrintNormal("\t-restoreFile:\tInput file (default /tmp/ddirb-output.restore)")
	//string_mng.PrintNormal("\t-port:\tHttp listener port default: 8080")
	string_mng.PrintNormal("\t-output:\tOutput file (default /tmp/ddirb-output)")
	string_mng.PrintNormal("\t-threads:\tNumber of threads to use")
	string_mng.PrintNormal("\t-delay:\t\tDelay in seconds between each thread")
	string_mng.PrintNormal("\t-debug:\t\tFlag to print in verbose mode")
	string_mng.PrintNormal("\t-graph:\t\tFlag to save the graph in .dot language")
	string_mng.PrintNormal("\t-headers:\tSet headers (es. Header1:value1;value2,Header2:value1)")
}

func PrintStaticUsage() {
	string_mng.PrintNormal("[NOT IMPLEMENTED YET] static Usage:")
	string_mng.PrintNormal("\t-url:\t\tTarget url in http/https format (http://example.com)")
	string_mng.PrintNormal("\t-output:\tOutput file (default /tmp/ddirb-output)")
	string_mng.PrintNormal("\t-threads:\tNumber of threads to use")
	string_mng.PrintNormal("\t-delay:\t\tDelay in seconds between each thread")
	string_mng.PrintNormal("\t-debug:\t\tFlag to print in verbose mode")
	string_mng.PrintNormal("\t-headers:\tSet headers (es. Header1:value1;value2,Header2:value1)")
}

// From url set url and domain
func (pv *ParamValidator) SetUrlDomain(url string) {
	pv.url = url
	pv.domain = getDomainName(url)
	pv.urlTarget = getInitialUrl(url)
}

// Various getter

func (pv *ParamValidator) GetUrlTarget() string {
	return pv.urlTarget
}

func (pv *ParamValidator) GetUrl() string {
	return pv.url
}
func (pv *ParamValidator) GetHeaders() string {
	return pv.headers
}
func (pv *ParamValidator) GetDomain() string {
	return pv.domain
}
func (pv *ParamValidator) GetOutputFile() string {
	return pv.outputFile
}
func (pv *ParamValidator) GetOutputFileGraph() string {
	return pv.outputFileGraph
}
func (pv *ParamValidator) GetOutputFileRestore() string {
	return pv.outputFileRestore
}
func (pv *ParamValidator) GetInputRestoreFile() string {
	return pv.restoreFile
}
func (pv *ParamValidator) GetThreads() int {
	return pv.threads
}
func (pv *ParamValidator) GetPort() int {
	return pv.port
}
func (pv *ParamValidator) GetDelay() time.Duration {
	return pv.delay
}
func (pv *ParamValidator) GetMethods() string {
	return pv.methods
}
func (pv *ParamValidator) GetRethods() bool {
	return pv.recursive
}
func (pv *ParamValidator) isDebug() bool {
	return pv.debugFlag
}

func (pv *ParamValidator) GetScanType() string {
	return pv.scanType
}

func (pv *ParamValidator) GetWordlist() string {
	return pv.wordlist
}

func (pv *ParamValidator) IsGraph() bool {
	return pv.graph
}

// PrintStructure method prints the structure
func (pv *ParamValidator) PrintStructure() {
	fmt.Printf("%+v\n", pv)
}

func (pv *ParamValidator) PrintDebug(message string) {
	if pv.debugFlag {
		string_mng.PrintDebug(message)
	}
}

func (pv *ParamValidator) PrettyPrintParameters() {
	string_mng.PrintInfo("---------- Execution Arguments ----------")
	string_mng.PrintInfo("[*] Scan Type: " + pv.scanType)
	string_mng.PrintInfo("[*] URL: " + pv.url)
	string_mng.PrintInfo("[*] URL Target: " + pv.urlTarget)
	string_mng.PrintInfo("[*] Domain: " + pv.domain)
	string_mng.PrintInfo("[*] Output file: " + pv.outputFile)
	if pv.graph {
		string_mng.PrintInfo("[*] Output file: " + pv.outputFileGraph)
	}
	string_mng.PrintInfo("[*] Output restore: " + pv.outputFileRestore)
	string_mng.PrintInfo("[*] Threads: " + strconv.Itoa(pv.threads))
	string_mng.PrintInfo("[*] Requests delay: " + strconv.Itoa(int(pv.delay)))
	string_mng.PrintInfo("[*] Methods: " + pv.methods)
	string_mng.PrintInfo("[*] Recursive flag: " + strconv.FormatBool(pv.recursive))
	string_mng.PrintInfo("[*] Restore file: " + pv.restoreFile)
	string_mng.PrintInfo("[*] Debug flag: " + strconv.FormatBool(pv.debugFlag))
	if pv.headers != "" {
		headerSplit := strings.Split(pv.headers, ",")
		for i, singleHeader := range headerSplit {
			headerDetail := strings.Split(singleHeader, ":")
			string_mng.PrintInfo("[*] Header " + strconv.Itoa(i) + " Name: " + headerDetail[0] + " Value: " + headerDetail[1])
		}
	}

	string_mng.PrintInfo("-----------------------------------------")
	string_mng.PrintWarning("[*] Listening on port: " + strconv.Itoa(pv.port) + " Browsing http://127.0.0.1:" + strconv.Itoa(pv.port) + "/graphView/graphView.html and http://127.0.0.1:" + strconv.Itoa(pv.port) + "/graphView/outputView.html")
}

// PrettyPrintParameters

// ValidateUrl method verifies if url is valid i.e. non empty
func (pv *ParamValidator) ValidateUrl() bool {
	if pv.url == "" {
		// empty string verification
		string_mng.PrintWarning("[?] URL is empty (example: http://www.google.com)")
		return false
	} else if !(strings.HasPrefix(pv.url, "http://") ||
		strings.HasPrefix(pv.url, "https://")) {
		// invalid schema verification
		string_mng.PrintWarning("[?] Invalid URL format (example: http://www.google.com)")
		return false
	} else {
		// everything ok
		return true
	}
}

// ValidateOutputFile verifies if path is valid, already exists
func (pv *ParamValidator) ValidateOutputFile() bool {
	if file_mng.FileExists(pv.outputFile) {
		if file_mng.RequireConfirmation("Output file already exists, do you want to override it? y/n") {
			file_mng.EmptyFileOverride(pv.outputFile)
			// if user need .dot graph override also .dot output file
			if pv.IsGraph() {
				file_mng.EmptyFileOverride(pv.outputFile + ".dot")
			}

			return true
		}
		return false
	}

	file_mng.EmptyFileOverride(pv.outputFile)
	if pv.IsGraph() {
		file_mng.EmptyFileOverride(pv.outputFile + ".dot")
	}
	return true
}

// ValidateThreads verifies the bounds range for threads field
func (pv *ParamValidator) ValidateThreads() bool {
	if pv.threads < 1 || pv.threads > 10 {
		// threads bounds verification (0-10 hardcoded limit)
		return false
	}
	return true
}

func (pv *ParamValidator) validateWordlist() bool {

	if !file_mng.FileExists(pv.outputFile) {
		string_mng.PrintError("[-] Wordlist not found")
		return false
	}
	return true
}

// ValidateDelay  must be implemented
func (pv *ParamValidator) ValidateDelay() bool {
	return true
}

// ValidateMethods validates methods: at the momend GET/HEAD
func (pv *ParamValidator) ValidateMethods() bool {
	return false
}

func (pv *ParamValidator) ValidateDynamic() bool {
	return (pv.ValidateUrl() &&
		pv.ValidateOutputFile() &&
		pv.ValidateThreads() &&
		pv.ValidateDelay())
}

func (pv *ParamValidator) ValidateRestoreDynamic() bool {
	return (file_mng.FileExists(pv.restoreFile) && pv.ValidateDelay() && pv.ValidateThreads() && pv.ValidateOutputFile())
}

func (pv *ParamValidator) ValidateStatic() bool {
	return (pv.ValidateUrl() &&
		pv.ValidateOutputFile() &&
		pv.ValidateThreads() &&
		pv.ValidateDelay() &&
		pv.validateWordlist())
}

func getInitialUrl(completeUrl string) string {
	//https://admin:admin@www.test.com
	var regex = regexp.MustCompile(`^(https?:\/\/)?([^@\n]+@)?(www\.)?([^:\/\n?]+)`)
	var result = regex.FindAllStringSubmatch(completeUrl, -1)
	//[[https://example.com example.com]]
	if len(result) == 0 {
		return ""
	}
	return result[0][0]
}

// url(string): url string
// result(string): return only the domain name starting from a url
func getDomainName(url string) string {
	//https://admin:admin@www.test.com
	var regex = regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+)`)
	var result = regex.FindAllStringSubmatch(url, -1)
	//[[https://example.com example.com]]
	if len(result) == 0 {
		return ""
	}
	return result[0][1]
}
