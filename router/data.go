package router

type JsonMessage struct {
	Service string `json:"service"`
	Action  string `json:"action"`
	Params  string `json:"params"`
}
