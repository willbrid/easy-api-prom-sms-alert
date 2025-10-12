package httpparam

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	contentTypeForm = "application/x-www-form-urlencoded"
	contentTypeJson = "application/json"

	PostType  = "post"
	QueryType = "query"
)

type Param struct {
	Post  map[string]string
	Query map[string]string
}

func NewParam() *Param {
	return &Param{
		Post:  make(map[string]string, 0),
		Query: make(map[string]string, 0),
	}
}

func (p *Param) AddPostParam(key, value string) *Param {
	p.Post[key] = value

	return p
}

func (p *Param) AddQueryParam(key, value string) *Param {
	p.Query[key] = value

	return p
}

func (p *Param) AddParam(paramType, key, value string) *Param {
	switch strings.ToLower(paramType) {
	case PostType:
		p.AddPostParam(key, value)
	case QueryType:
		p.AddQueryParam(key, value)
	default:
		p.AddPostParam(key, value)
	}

	return p
}

func (p *Param) EncodeQuery() string {
	return encodeUrlParams(p.Query)
}

func (p *Param) EncodePost(contentType string) (string, error) {
	var encodedPost string

	switch contentType {
	case contentTypeForm:
		encodedPost = encodeUrlParams(p.Post)

	case contentTypeJson:
		postParamStr, err := json.Marshal(p.Post)
		if err != nil {
			return "", err
		}
		encodedPost = string(postParamStr)
	}

	return encodedPost, nil
}

// encodeUrlParams encodes URL parameters from a map to a query string format.
func encodeUrlParams(params map[string]string) string {
	var encodedUrlParams string

	for key, value := range params {
		encodedUrlParams += fmt.Sprintf("&%s=%s", key, value)
	}
	encodedUrlParams = encodedUrlParams[1:]

	return encodedUrlParams
}
