package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Encoding string

const (
	JSONEncoding Encoding = "json"
	YAMLEncoding Encoding = "yaml"
)

type Encoder interface {
	Marshal(v *ResourceDefinition) ([]byte, error)
	Unmarshal([]byte) (*ResourceDefinition, error)
	Encoding() Encoding
}

func NewEncoder(encoding Encoding) (Encoder, error) {
	switch encoding {
	case JSONEncoding:
		return &JSONStrategy{}, nil
	case YAMLEncoding:
		return &YAMLStrategy{}, nil
	default:
		return nil, fmt.Errorf("%w: %q", logger.ErrEncoding, encoding)
	}
}

func NewEncoderFromFilepath(path string) (Encoder, error) {
	split := strings.Split(path, ".")
	extension := split[len(split)-1]

	encoder, err := NewEncoder(Encoding(extension))
	if err != nil {
		if errors.As(err, &logger.ErrEncoding) {
			err = fmt.Errorf("%w; %w: %q", err, logger.ErrFileExtension, path)
			logger.Error(err)
			return nil, err
		}

		logger.Error(err)
		return nil, err
	}

	return encoder, nil
}

type JSONStrategy struct{}

func (s *JSONStrategy) Marshal(v *ResourceDefinition) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *JSONStrategy) Unmarshal(data []byte) (*ResourceDefinition, error) {
	v := new(ResourceDefinition)
	if err := json.Unmarshal(data, v); err != nil {
		logger.Error(err)
		return nil, err
	}

	return v, nil
}

func (s *JSONStrategy) Encoding() Encoding {
	return JSONEncoding
}

type YAMLStrategy struct{}

func (s *YAMLStrategy) Marshal(v *ResourceDefinition) ([]byte, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil

}

func (s *YAMLStrategy) Unmarshal(data []byte) (*ResourceDefinition, error) {
	v := new(ResourceDefinition)
	if err := yaml.Unmarshal(data, v); err != nil {
		logger.Error(err)
		return nil, err
	}

	return v, nil

}

func (s *YAMLStrategy) Encoding() Encoding {
	return YAMLEncoding
}

func ReadEncodedFile(path string) (*ResourceDefinition, error) {
	encoder, err := NewEncoderFromFilepath(path)
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return encoder.Unmarshal(b)
}

func WriteEncodedFile(path string, v *ResourceDefinition) error {
	encoder, err := NewEncoderFromFilepath(path)
	if err != nil {
		return err
	}

	b, err := encoder.Marshal(v)
	if err != nil {
		logger.Error(err)
		return err
	}

	err = os.WriteFile(path, b, 0666)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
