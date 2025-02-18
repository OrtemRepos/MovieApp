package model

import (
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
)

func MovieDetailsToProto(m *MovieDetails) *gen.MovieDetails {
	return &gen.MovieDetails{
		Rating:   float32(m.Rating),
		Metadata: model.MetadataToProto(&m.Metadata),
	}
}

func MovieDetailsFromProto(m *gen.MovieDetails) *MovieDetails {
	return &MovieDetails{
		Rating:   float64(m.Rating),
		Metadata: *model.MetadataFromProto(m.Metadata),
	}
}
