package models

type Graph struct {
	Graph string `json:"graph"`
}

type Output struct {
	Name string   `json:"name"`
	Body []string `json:"output"`
}

type StaticNode struct{
	Page string
	Info []ResponseNodeInfo
	Nodes []StaticNode
}

type ResponseNodeInfo struct{
	RequestMethod string
	ResponseCode string
	ResponseServer string
	ResponseDate string
	ResponseHasForm bool
	ResponseFormCounter int
}
