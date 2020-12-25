package landscape

import (
	"encoding/json"
	"log"
	"os"
)

var (
	jsonStr string
)

// Landscapes data structure
type Landscapes map[string]struct {
	CloudController string   `json:"cloudcontroller"`
	Uaa             string   `json:"uaa"`
	Labels          []string `json:"labels"`
}

// Get returns a lookup data structure
func Get() Landscapes {
	str := os.Getenv("LANDSCAPES")

	if str == "" {
		log.Println("LANDSCAPES environment variable not set")
		return Landscapes{}
	}

	data := Landscapes{}

	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		log.Printf("Error: %v", err)
		log.Printf("env was: %v", "LANDSCAPES")
		log.Printf("json object was: %v", str)
	}

	return data
}
