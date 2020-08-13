package static

import (
	"dynamic-dirb/internal/string_mng"
	"dynamic-dirb/internal/file_mng"
	service "dynamic-dirb/internal/services"
)

func Static() {
	/* service.GetParameters().PrintDebug("Validating parameters...")
	if service.GetParameters().ValidateStatic() {
		service.GetParameters().PrettyPrintParameters()
	} else {
		string_mng.PrintError("[-] Invalid parameters")
		os.Exit(0)
	}
	service.GetParameters().PrintDebug("...Parameters validated") */
	/*
	// TODO:
	1) wordlist
	- from file
	- from url (download  in tmp + parse)
	2) threads

	3) requestQueue -> Response
	-> status code
	-> content length
	-> server response header

	*/
	service.GetParameters().PrintDebug("Validating static parameters...")
	if (service.GetParameters().ValidateStatic()) {
		service.GetParameters().PrettyPrintParameters()
		// print Wordlist content
		wordlistContent := file_mng.ReadFileByLine(service.GetParameters().GetWordlist())
		wordlistLen := len(wordlistContent)

		for i:=0; i< wordlistLen; i=i+1{
			service.GetParameters().PrintDebug(wordlistContent[i])
		}
		// check if target reachable with HEAD

	}else{
		string_mng.PrintWarning("[NOT IMPLEMENTED YET]")
	}
}
