package main

import (
	"encoding/json"
	"log"

	"github.com/rodpadev/aws-profiles/lib"
)

func main() {
	data, err := lib.LoadAWSProfileData()
	if err != nil {
		log.Printf("ERR@LoadAWSProfileData")
		log.Fatal(err)
	}

	parsed := lib.ParseAWSProfileData(data)

	profileMap, err := lib.BuildProfileMap(parsed)
	if err != nil {
		log.Printf("ERR@BuildProfileMap")
		log.Fatal(err)
	}

	if jsonBytes, err := json.MarshalIndent(profileMap, "", "  "); err == nil {
		log.Println(string(jsonBytes))
	}

}
