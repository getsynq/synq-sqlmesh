package sqlmesh

import (
	ingestsqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"strings"
)

func NewSQLMeshMetadata() *ingestsqlmeshv1.IngestMetadataRequest {
	return &ingestsqlmeshv1.IngestMetadataRequest{
		ModelDetails: make(map[string][]byte),
		ModelLineage: make(map[string][]byte),
		StateAt:      timestamppb.Now(),
	}
}

func CollectMetadata(url url.URL) (*ingestsqlmeshv1.IngestMetadataRequest, error) {

	api := NewAPIClient(url)

	res := NewSQLMeshMetadata()
	var err error
	res.ApiMeta, err = api.GetMeta()
	logError(err, "Failed to get meta information")
	res.Models, err = api.GetModels()
	logError(err, "Failed to get models information")
	modelNames, err := ModelNames(res.Models)
	logError(err, "Failed to get model names")
	for _, modelName := range modelNames {
		res.ModelDetails[modelName], err = api.GetModel(modelName)
		logError(err, "Failed to get model details")
		res.ModelLineage[modelName], err = api.GetLineage(modelName)
		logError(err, "Failed to get model lineage")
	}
	res.Files, err = api.GetFiles()
	logError(err, "Failed to get files information")

	res.Environments, err = api.GetEnvironments()
	logError(err, "Failed to get environments information")

	return res, nil
}

func logError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Error(msg)
	}
}

func ModelNames(models json.RawMessage) ([]string, error) {
	type model struct {
		Name string `json:"name"`
	}

	var decodedModels []*model
	err := json.Unmarshal(models, &decodedModels)
	if err != nil {
		return nil, err
	}

	var modelNames []string
	for _, m := range decodedModels {
		modelName := strings.TrimSpace(m.Name)
		if len(modelName) > 0 {
			modelNames = append(modelNames, modelName)
		}
	}
	return modelNames, nil
}
