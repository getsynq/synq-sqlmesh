package sqlmesh

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	ingestsqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"github.com/getsynq/synq-sqlmesh/build"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	processErr(res, err, "Failed to get meta information")
	res.Models, err = api.GetModels()
	processErr(res, err, "Failed to get models information")
	modelNames, err := ModelNames(res.Models)
	processErr(res, err, "Failed to get model names")
	for _, modelName := range modelNames {
		res.ModelDetails[modelName], err = api.GetModel(modelName)
		processErr(res, err, "Failed to get model details of %s", modelName)
		res.ModelLineage[modelName], err = api.GetLineage(modelName)
		processErr(res, err, "Failed to get model lineage of %s", modelName)
	}
	res.Files, err = api.GetFiles()
	processErr(res, err, "Failed to get files information")

	if len(res.Files) > 0 {
		dir := &Directory{}
		err = json.Unmarshal(res.Files, dir)
		if err != nil {
			logrus.WithError(err).Error("Failed to unmarshal /files/ response")
		} else {
			filesToProcess, err := collectFilesForProcessing(res.Files, fileContentGlobFilter)
			if err != nil {
				processErr(res, err, "Failed to collect files for processing")
			} else {
				for _, fileToProcess := range filesToProcess {
					fileContent, err := api.GetFileContent(fileToProcess)
					processErr(res, err, "Failed to get file content %s", fileToProcess)
					if err == nil {
						res.FileContent[fileToProcess] = fileContent
					}
				}
			}
		}
	}

	res.Environments, err = api.GetEnvironments()
	processErr(res, err, "Failed to get environments information")

	return res, nil
}

func processErr(res *ingestsqlmeshv1.IngestMetadataRequest, err error, msg string, args ...interface{}) {
	if err == nil {
		return
	}
	var sqlMeshApiErr *SQLMeshApiError
	if errors.As(err, &sqlMeshApiErr) {
		res.Errors = append(res.Errors, &ingestsqlmeshv1.IngestMetadataRequest_Error{
			Path:    lo.ToPtr(sqlMeshApiErr.UrlPath),
			Code:    lo.ToPtr(int64(sqlMeshApiErr.Code)),
			Message: sqlMeshApiErr.Message,
		})
	} else {
		res.Errors = append(res.Errors, &ingestsqlmeshv1.IngestMetadataRequest_Error{
			Message: err.Error(),
		})
	}

	logError(err, msg, args...)
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

func logError(err error, msg string, args ...interface{}) {
	if err != nil {
		logrus.Errorf(msg, args...)
		logrus.Error(err)
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
