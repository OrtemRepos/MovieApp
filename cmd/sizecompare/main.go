package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"google.golang.org/protobuf/proto"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
)

var metadata = &model.Metadata{
	ID:          "1",
	Title:       "The Incredible Hulk",
	Description: "The Incredible Hulk is a 2008 American superhero film based on the Marvel Comics character the Hulk.",
	Director:    "Louis Leterrier",
}

var genMetadata = &gen.Metadata{
	Id:          "1",
	Title:       "The Incredible Hulk",
	Description: "The Incredible Hulk is a 2008 American superhero film based on the Marvel Comics character the Hulk.",
	Director:    "Louis Leterrier",
}

func main() {
	jsonBytes, err := json.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := xml.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := proto.Marshal(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size: %d\n", len(jsonBytes))
	fmt.Printf("XML size: %d\n", len(xmlBytes))
	fmt.Printf("Proto size: %d\n", len(protoBytes))
}
