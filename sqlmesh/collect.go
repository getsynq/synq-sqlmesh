package sqlmesh

import (
	ingestsqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsynq/synq-sqlmesh/build"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/url"
	"path/filepath"
	"strings"
)

func NewSQLMeshMetadata() *ingestsqlmeshv1.IngestMetadataRequest {
	return &ingestsqlmeshv1.IngestMetadataRequest{
		ModelDetails: make(map[string][]byte),
		ModelLineage: make(map[string][]byte),
		FileContent:  make(map[string][]byte),
		StateAt:      timestamppb.Now(),
	}
}

func CollectMetadata(url url.URL, fileContentGlobFilter GlobFilter) (*ingestsqlmeshv1.IngestMetadataRequest, error) {

	api := NewAPIClient(url)

	res := NewSQLMeshMetadata()
	res.UploaderVersion = strings.TrimSpace(fmt.Sprintf("synq-sqlmesh/%s", build.Version))
	res.UploaderBuildTime = strings.TrimSpace(build.Time)

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

	if len(res.Files) > 0 {
		dir := &Directory{}
		err = json.Unmarshal(res.Files, dir)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal /files/ response")
		} else {
			filesToProcess, err := collectFilesForProcessing(res.Files, fileContentGlobFilter)
			if err != nil {
				logError(err, "Failed to collect files for processing")
			} else {
				for _, fileToProcess := range filesToProcess {
					fileContent, err := api.GetFileContent(fileToProcess)
					logError(err, "Failed to get file content")
					if err == nil {
						res.FileContent[fileToProcess] = fileContent
					}
				}
			}
		}
	}

	res.Environments, err = api.GetEnvironments()
	logError(err, "Failed to get environments information")

	return res, nil
}

func collectFilesForProcessing(files []byte, fileContentGlobFilter GlobFilter) ([]string, error) {
	var filesToGetContent []string
	dir := Directory{}
	err := json.Unmarshal(files, &dir)
	if err != nil {
		logrus.WithError(err).Error("Failed to unmarshal /files/ response")
		return nil, err
	} else {

		dirsToProcess := []Directory{dir}
		for len(dirsToProcess) > 0 {
			dir := dirsToProcess[0]
			dirsToProcess = dirsToProcess[1:]
			for _, file := range dir.Files {
				accepted, err := fileContentGlobFilter.Match(file.Path)
				if err != nil {
					if errors.Is(err, filepath.ErrBadPattern) {
						return nil, err
					}
					logrus.WithError(err).Error("Failed to match file path")
					continue
				}
				if accepted {
					filesToGetContent = append(filesToGetContent, file.Path)
				}
			}
			for _, subDir := range dir.Directories {
				dirsToProcess = append(dirsToProcess, subDir)
			}
		}
	}
	return filesToGetContent, nil
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
