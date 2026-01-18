package faceauthservices

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/google/uuid"
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

func (s *FaceAuthServices) EnrollFace(ctx context.Context, userID uuid.UUID, imageBytes []byte) error {
	input := &rekognition.IndexFacesInput{
		CollectionId: aws.String(s.collection),
		Image: &types.Image{
			Bytes: imageBytes,
		},
		ExternalImageId: aws.String(userID.String()),
		MaxFaces:        aws.Int32(1),
		QualityFilter:   types.QualityFilterAuto,
	}

	output, err := s.rekognition.IndexFaces(ctx, input)
	if err != nil {
		return fmt.Errorf("error de comunicación con AWS Rekognition: %w", err)
	}

	if len(output.FaceRecords) == 0 {
		return fmt.Errorf("no se detectó ningún rostro válido en la imagen")
	}

	faceID := *output.FaceRecords[0].Face.FaceId

	if err := s.store.FaceAuth.UpdateUserFaceID(ctx, userID, faceID); err != nil {
		return fmt.Errorf("error guardando credenciales faciales: %w", err)
	}

	return nil
}

func (s *FaceAuthServices) AuthenticateUser(ctx context.Context, imageBytes []byte) (uuid.UUID, error) {
	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(s.collection),
		Image:              &types.Image{Bytes: imageBytes},
		FaceMatchThreshold: aws.Float32(95.0),
		MaxFaces:           aws.Int32(1),
	}

	output, err := s.rekognition.SearchFacesByImage(ctx, input)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error buscando rostro en AWS: %w", err)
	}

	if len(output.FaceMatches) == 0 {
		return uuid.Nil, fmt.Errorf("rostro no reconocido")
	}

	matchFaceID := *output.FaceMatches[0].Face.FaceId

	userID, err := s.store.FaceAuth.GetUserIDByFaceID(ctx, matchFaceID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("rostro reconocido en nube pero usuario no encontrado en sistema")
	}

	return userID, nil
}
