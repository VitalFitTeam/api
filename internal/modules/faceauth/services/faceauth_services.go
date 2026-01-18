package faceauthservices

import (
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/vitalfit/api/internal/store"
)

type FaceAuthServices struct {
	store       store.Storage
	collection  string
	rekognition *rekognition.Client
}

func NewFaceAuthServices(store store.Storage, rekognition *rekognition.Client) *FaceAuthServices {
	return &FaceAuthServices{
		store:       store,
		collection:  "vitalfit-faces",
		rekognition: rekognition,
	}
}
