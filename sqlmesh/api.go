package sqlmesh

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/url"
)

type Api interface {
	GetMeta() (json.RawMessage, error)
	GetModels() (json.RawMessage, error)
	GetModel(modelName string) (json.RawMessage, error)
	GetLineage(modelName string) (json.RawMessage, error)
	GetEnvironments() (json.RawMessage, error)
	GetFiles() (json.RawMessage, error)
	Health() (json.RawMessage, error)
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
	statusCode, body, err := a.c.Get(nil, a.url("health"))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetMeta() (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "meta"))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetModels() (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "models"))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetModel(modelName string) (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "models", modelName))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetLineage(modelName string) (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "lineage", modelName))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetEnvironments() (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "environments"))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) GetFiles() (json.RawMessage, error) {
	statusCode, body, err := a.c.Get(nil, a.url("api", "files"))
	if err != nil {
		return nil, err
	}
	if statusCode != fasthttp.StatusOK {
		return nil, err
	}
	return body, nil
}

func (a ApiImpl) url(path ...string) string {
	return a.baseUrl.JoinPath(path...).String()
}
