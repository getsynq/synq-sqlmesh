package sqlmesh

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/valyala/fasthttp"
)

type Api interface {
	GetMeta() (json.RawMessage, error)
	GetModels() (json.RawMessage, error)
	GetModel(modelName string) (json.RawMessage, error)
	GetLineage(modelName string) (json.RawMessage, error)
	GetEnvironments() (json.RawMessage, error)
	GetFiles() (json.RawMessage, error)
	GetFileContent(filePath string) (json.RawMessage, error)
	Health() (json.RawMessage, error)
}

type Directory struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Directories []Directory `json:"directories,omitempty"`
	Files       []File      `json:"files,omitempty"`
}

type File struct {
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	Extension *string `json:"extension,omitempty"`
	Content   *string `json:"content,omitempty"`
}

func NewAPIClient(url url.URL) Api {
	c := &fasthttp.Client{}
	return &ApiImpl{
		c:       c,
		baseUrl: url,
	}
}

type ApiImpl struct {
	c       *fasthttp.Client
	baseUrl url.URL
}

func (a ApiImpl) Health() (json.RawMessage, error) {
	urlPath := a.buildUrlPath("health")
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetMeta() (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "meta")
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetModels() (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "models")
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetModel(modelName string) (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "models", modelName)
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetLineage(modelName string) (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "lineage", modelName)
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetEnvironments() (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "environments")
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetFiles() (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "files")
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) GetFileContent(filePath string) (json.RawMessage, error) {
	urlPath := a.buildUrlPath("api", "files", filePath)
	statusCode, body, err := a.c.Get(nil, urlPath)
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, a.createStatusError(urlPath, statusCode, body)
	}
	return body, nil
}

func (a ApiImpl) buildUrlPath(path ...string) string {
	return a.baseUrl.JoinPath(path...).String()
}

type SQLMeshApiError struct {
	UrlPath string `json:"path"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *SQLMeshApiError) Error() string {
	return fmt.Sprintf("SQLMesh UI API error at %s (%d): %s ", s.UrlPath, s.Code, s.Message)
}

func (a ApiImpl) createStatusError(urlPath string, code int, body []byte) error {
	return &SQLMeshApiError{
		UrlPath: urlPath,
		Code:    code,
		Message: string(body),
	}
}
