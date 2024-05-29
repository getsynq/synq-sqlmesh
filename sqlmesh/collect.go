package sqlmesh

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

func NewSqlMeshMetadata() *SqlMeshMetadata {
	return &SqlMeshMetadata{
		ModelDetails: make(map[string]json.RawMessage),
		ModelLineage: make(map[string]json.RawMessage),
	}
}

type SqlMeshMetadata struct {
	Meta         json.RawMessage            `json:"meta"`
	Models       json.RawMessage            `json:"models"`
	ModelDetails map[string]json.RawMessage `json:"model_details"`
	ModelLineage map[string]json.RawMessage `json:"model_lineage"`
	Files        json.RawMessage            `json:"files"`
	Environments json.RawMessage            `json:"environments"`
}

func CollectMetadata(url url.URL) (*SqlMeshMetadata, error) {

	api := NewAPIClient(url)

	res := NewSqlMeshMetadata()
	var err error
	res.Meta, err = api.GetMeta()
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
