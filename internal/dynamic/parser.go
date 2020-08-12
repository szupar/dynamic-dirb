package dynamic

import (
	"dynamic-dirb/internal/http_mng"
	service "dynamic-dirb/internal/services"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// response(ResponseInfo): struct containing the response information
// return([]string): javascript urls found in the body in src tag
func getJavascriptUrl(jsChannel chan<- []string, response http_mng.ResponseInfo) {
	var body = response.ResponseBodyString
	var return_matched []string
	//var regex = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	var regex = regexp.MustCompile(`src="(.*\.js)`)
	//result --> [][] string
	var result = regex.FindAllStringSubmatch(body, -1)
	//fmt.Println(regex.FindAllStringSubmatch(body, -1))
	for _, s := range result {
		//fmt.Println(s)
		for j := 1; j < len(s); j++ {
			return_matched = append(return_matched, s[j])
			service.GetParameters().PrintDebug("[+] Regex getting js url: " + s[j])
			//fmt.Println(s[j])
		}
	}
	jsChannel <- return_matched
}

// response(ResponseInfo): struct containing the response information
// return([]string): javascript urls found in the body in src tag
func getHtmlUrl(htmlChannel chan<- []string, response http_mng.ResponseInfo) {
	var body = response.ResponseBodyString
	var return_matched []string
	//var regex = regexp.MustCompile(`^[a-z]+\[[0-9]+\]$`)
	var regex = regexp.MustCompile(`href="(.*?)"`)     //non-greedy regex
	var result = regex.FindAllStringSubmatch(body, -1) //result --> [][] string
	for _, s := range result {
		for j := 1; j < len(s); j++ {

			return_matched = append(return_matched, s[j])
			service.GetParameters().PrintDebug("[+] Regex getting html url: " + s[j])
		}
	}
	htmlChannel <- return_matched
}

// response(ResponseInfo): struct containing the response information
// return([]string): javascript urls found in the body in src tag
func getGenericUrl(urlChannel chan<- []string, response http_mng.ResponseInfo) {
	var body = response.ResponseBodyString
	var return_matched []string
	//var regex = regexp.MustCompile(`"([a-z\.\/]*\.?[a-z]+)(?:.*)"`)
	var regex = regexp.MustCompile(`"([a-z\.\/]*\.?[a-z]+)(?:(?:\?[a-z]+=|"){1})`)
	var result = regex.FindAllStringSubmatch(body, -1) //result --> [][] string
	for _, s := range result {
		for j := 1; j < len(s); j++ {

			return_matched = append(return_matched, s[j])
			service.GetParameters().PrintDebug("[+] Regex getting generic url: " + s[j])
		}
	}
	urlChannel <- return_matched
}

// target(targetInfo): target struct with target information
// tovalidate([]string): urls array to check
// return([]string): array with only UNIQUE url containing domain target or start with / ./ [a-z] exept for http or https
func matchDomainUrl(toValidate []string) []string {

	var return_matched []string
	// creo una map di string in modo tale da non avere valori doppi dell'array che ritorno
	var uniqueList = map[string]bool{}

	for _, url := range toValidate {
		if len(url) >= 2 {
			//prendo gli url che iniziano con il nome del dominio di test e quelli che iniziano con / ./ [a-z]
			if service.GetParameters().GetDomain() == http_mng.GetDomainName(url) || checkUrlSyntax(url) {
				service.GetParameters().PrintDebug("[+] Matched domain or Url start with / or ./ or ../ or [a-z] for: " + url)
				//return_matched = append(return_matched, url)
				if !uniqueList[url] {
					uniqueList[url] = true
				}
			}
		}
	}
	for key, _ := range uniqueList {
		return_matched = append(return_matched, key)
	}
	return return_matched
}

//							[DO NOT DELETE]
// NOT USED ANYMORE SINGLE THREAD ONLY FOR TESTING PURPOSE [DO NOT DELETE]
// target(targetInfo): target struct with target information
// tovalidate([]string): urls array to check
// return([]string): url array containing url with response code = 200
/*func getExistsUrlOld(toValidate []string) []string {
	var return_matched []string

	for _, url := range toValidate {

		var completeUrl = ""
		var regex = regexp.MustCompile(`^http.*`)
		//se l'url inizia con https/http faccio la get
		if regex.MatchString(url) {
			completeUrl = url
		} else { //altrimenti ho url che iniziano con / ./ [a-z]

			if url[0] == '/' {
				completeUrl = service.GetParameters().GetUrlTarget() + url
			} else if url[0] == '.' {
				completeUrl = service.GetParameters().GetUrlTarget() + url[1:]
			} else {
				completeUrl = service.GetParameters().GetUrlTarget() + "/" + url
			}
		}
		fmt.Println("Worker...")
		service.GetParameters().PrintDebug("Checking if: " + url + " exist")
		response, code := http_mng.RequestResourceExist(completeUrl)
		// if code == true --> response code == 200
		if code {
			service.GetParameters().PrintDebug("[+] Exists: " + url)
			return_matched = append(return_matched, completeUrl)
		} else {
			service.GetParameters().PrintDebug("[-] Not exists: " + url + " Status code: " + strconv.Itoa(response.ResponseCode))
		}

	}

	return return_matched
}*/

// target(targetInfo): target struct with target information
// tovalidate([]string): urls array to check
// result (map[string]bool): urls already discovered in previous interaction
// return([]string): url array containing url with response code = 200
func getExistsUrl(toValidate []string, result map[string]bool) []string {
	var return_matched []string
	var thread = service.GetParameters().GetThreads()
	var divided [][]string

	// split the input array in n(thread) chunk
	chunkSize := (len(toValidate) + thread - 1) / thread
	for i := 0; i < len(toValidate); i += chunkSize {
		end := i + chunkSize

		if end > len(toValidate) {
			end = len(toValidate)
		}
		//[[thread11][thread2][...]]
		divided = append(divided, toValidate[i:end])
	}

	// spawn n thread one for each divided split
	var wg sync.WaitGroup
	var m sync.Mutex

	for _, chunk := range divided {
		wg.Add(1)
		go func(chunk []string, wg *sync.WaitGroup) {
			for _, url := range chunk {

				// Wait n millisecond based on user option
				//time.Sleep(service.GetParameters().GetDelay() * time.Millisecond)
				time.Sleep(service.GetParameters().GetDelay())
				if result[url] {
					continue
				}
				//Composing url string
				//if url start with https/http no problem
				var completeUrl = ""
				var regex = regexp.MustCompile(`^http.*`)
				if regex.MatchString(url) {
					//service.GetParameters().PrintDebug("Url starts with http")
					completeUrl = url
				} else {
					// if url start with /
					if url[0] == '/' {
						//service.GetParameters().PrintDebug("Url starts with /")
						completeUrl = service.GetParameters().GetUrlTarget() + url
					} else if url[0] == '.' && url[1] == '/' {
						//service.GetParameters().PrintDebug("Url starts with ./")
						completeUrl = service.GetParameters().GetUrlTarget() + url[1:]
					} else {
						//service.GetParameters().PrintDebug("Url starts with char")
						completeUrl = service.GetParameters().GetUrlTarget() + "/" + url
					}
				}

				service.GetParameters().PrintDebug("Check if " + completeUrl + " exist")
				response, code := http_mng.RequestResourceExist(completeUrl)
				// if code == true --> response code == 200
				if code {
					service.GetParameters().PrintDebug("[+] Exists: " + url)
					m.Lock()
					return_matched = append(return_matched, completeUrl)
					m.Unlock()
				} else {
					service.GetParameters().PrintDebug("[-] Not exists: " + url + " Status code: " + strconv.Itoa(response.ResponseCode))
				}

			}
			wg.Done()
		}(chunk, &wg)

	}
	wg.Wait()
	return return_matched
}

//target (*helper.ParamValidator): target details
//response (http_mng.ResponseInfo): response information
//return ([]string): url array found in source code of response. With response 200 and in target scope
func GetFinalUrls(response http_mng.ResponseInfo, result map[string]bool) []string {

	// creating channel
	jsChannel := make(chan []string, 1)
	hmtlChannel := make(chan []string, 1)
	urlChannel := make(chan []string, 1)

	// starting parsing js and html (multi-threads)
	go getJavascriptUrl(jsChannel, response)
	go getHtmlUrl(hmtlChannel, response)
	go getGenericUrl(urlChannel, response)
	// waiting end
	discoveredJs := <-jsChannel
	discoveredHtml := <-hmtlChannel
	discoveredGenericUrl := <-urlChannel

	discoveredUrl := append(discoveredJs, discoveredHtml...)
	discoveredUrl = append(discoveredUrl, discoveredGenericUrl...)
	discoveredUrl = matchDomainUrl(discoveredUrl)
	service.GetParameters().PrintDebug("Getting existing url")
	discoveredUrl = getExistsUrl(discoveredUrl, result)

	return discoveredUrl

}

// url(string): url string
// return(bool): true if url start with / ./ [a-z] exept for ^http or ^https else false
func checkUrlSyntax(url string) bool {
	// start with ./
	var regexDotSlash = regexp.MustCompile(`(?:^\.\/.*)`)
	// start with /
	var regexSlash = regexp.MustCompile(`^\/{1}[^\/][a-z|A-z]*`)
	// start with [a-z]
	var regexAZ = regexp.MustCompile(`^[a-z|A-z].*`)
	// start with https:// or http://
	var regexHTTP = regexp.MustCompile(`^(https:\/\/|http:\/\/)`)

	if (regexDotSlash.MatchString(url) || regexSlash.MatchString(url) || regexAZ.MatchString(url)) && !regexHTTP.MatchString(url) {
		return true
	}

	return false
}
