package web_mng

import (
	"dynamic-dirb/internal/models"
	service "dynamic-dirb/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/browser"
)

// Start the web server and listen for the requests
func StartWebServer() {
	// Getwd returns a rooted path name corresponding to the
	path, err := os.Getwd()

	if err != nil {
		log.Println(err)
		return
	}

	browser.OpenURL("http://127.0.0.1:" + strconv.Itoa(service.GetParameters().GetPort()) + "/graphView/graphView.html")
	browser.OpenURL("http://127.0.0.1:" + strconv.Itoa(service.GetParameters().GetPort()) + "/graphView/outputView.html")
	http.HandleFunc("/getGraph", handleHttpGraphResponse)
	http.HandleFunc("/getOutput", handleHttpOutputResponse)
	http.Handle("/graphView/", http.StripPrefix("/graphView/", http.FileServer(http.Dir(path+"/js-graph-parser"))))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(service.GetParameters().GetPort()), nil))
}

func handleHttpGraphResponse(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	response := models.Graph{Graph: service.GetGraph()}
	enc := json.NewEncoder(responseWriter)
	enc.SetEscapeHTML(false)
	enc.Encode(response)
}

func handleHttpOutputResponse(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Add("Content-Type", "application/json")
	responseWriter.Header().Add("Access-Control-Allow-Origin", "*")
	response := models.Output{Name: service.GetGraphName(), Body: service.GetOutput()}
	enc := json.NewEncoder(responseWriter)
	enc.SetEscapeHTML(false)
	enc.Encode(response)

}
