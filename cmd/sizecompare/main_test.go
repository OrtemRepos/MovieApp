package main

import (
	"encoding/json"
	"encoding/xml"
	"testing"

	"google.golang.org/protobuf/proto"
)

func BenchmarkSerializeToJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(metadata)
	}
}

func BenchmarkSerializeToXML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = xml.Marshal(metadata)
	}
}

func BenchmarkSerializeToProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = proto.Marshal(genMetadata)
	}
}
