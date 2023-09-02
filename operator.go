package vib

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type OperatorStrategy string

const (
	FileSystemOperatorStrategy OperatorStrategy = "filesystem"
	GitOperatorStrategy        OperatorStrategy = "git"
)

type Operator interface {
	Get(name *string) ([]ResourceDefinition, error)
	Create(*ResourceDefinition) error
	Update(name *string, v *ResourceDefinition) error
	Delete(name string) error
}

//----------------------------------------------------------------------------------------------------------------------
// NewOperator
//----------------------------------------------------------------------------------------------------------------------

func NewOperator(strategy OperatorStrategy, options ...any) (Operator, error) {
	switch strategy {
	case FileSystemOperatorStrategy:
		if len(options) != 4 {
			fmt.Printf("%#v", options)
			return nil, NewErrAndLog(
				ErrType,
				fmt.Sprintf("wrong number of argument to construct FilesystemOperator; got: %d", len(options)),
			)
		}

		apiVersion, ok := options[0].(APIVersion)
		if !ok {
			return nil, NewErrAndLog(ErrType, "apiVersion must be of type string")
		}

		kind, ok := options[1].(Kind)
		if !ok {
			return nil, NewErrAndLog(ErrType, "kind must be of type string")
		}

		resourceDir, ok := options[2].(string)
		if !ok {
			return nil, NewErrAndLog(ErrType, "resourceDir must be of type string")
		}

		encoding, ok := options[3].(Encoding)
		if !ok {
			return nil, NewErrAndLog(ErrType, "encoding must be of type Encoding")
		}

		return NewFilesystemOperator(apiVersion, kind, resourceDir, encoding)
	case GitOperatorStrategy:
		// TODO implement me
		panic("not implemented yet")
	default:
		err := fmt.Errorf("%w: operator strategy %q is not supported", ErrType, strategy)
		logger.Error(err)

		return nil, err
	}
}

//----------------------------------------------------------------------------------------------------------------------
// FilesystemOperator
//----------------------------------------------------------------------------------------------------------------------

// FilesystemOperator operates T through the filesystem.
// Resources are stored on the filesystem using the following convention:
// - Filename: {{ T.APIVersion() }}.{{ T.Kind() }}.{{ T.IMetadata().Name }}. {{ s.encoder.Encoding() }}
type FilesystemOperator struct {
	apiVersion  APIVersion
	kind        Kind
	resourceDir string
	encoder     Encoder
}

func (s *FilesystemOperator) Get(name *string) ([]ResourceDefinition, error) {
	res := make([]ResourceDefinition, 0)
	// if name is specified T with specified name exist, then we return a list of length one containing T
	if name != nil {
		// read can return a nil pointer
		v, err := s.readWithResourceName(*name)
		if err != nil {
			return nil, err
		}

		// read can return a nil pointer. If nil pointer, we directly return a nil array
		if v == nil {
			return nil, nil
		}

		// If not nil then we can dereference and add the struct to the list
		res = append(res, *v)

		return res, nil
	}

	// Get all instance of T.
	list, err := s.list()
	if err != nil {
		return nil, err
	}

	for _, filename := range list {
		v, err := s.readWithFilename(filename)
		if err != nil {
			return nil, err
		}

		// read can return a nil pointer. If nil, continue the loop
		if v == nil {
			continue
		}

		// we can safely dereference the pointer
		res = append(res, *v)
	}

	return res, nil
}

// Create should create only if file does not already exist.
func (s *FilesystemOperator) Create(t *ResourceDefinition) error {
	exist, err := resourceExist(s, t.Metadata.Name)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("%w: cannot create resource; apiVersion: %q, kind: %q, name: %q",
			ErrAlreadyExist, t.APIVersion, t.Kind, t.Metadata.Name)
	}

	return s.write(t)
}

func (s *FilesystemOperator) Update(name *string, v *ResourceDefinition) error {
	// This operation rename the object
	if name != nil {
		name := *name
		if v.Metadata.Name != name {
			// Delete existing object with name `name *string` and rewrite with new name.
			if err := s.Delete(name); err != nil {
				return err
			}
		}
	}

	// Name is the same or former one was deleted, we can write T
	return s.write(v)
}

func (s *FilesystemOperator) Delete(name string) error {
	return os.Remove(s.filepathFromResourceName(name))
}

// list returns a list of reference to T.
func (s *FilesystemOperator) list() ([]string, error) {
	dirEntries, err := os.ReadDir(s.resourceDir)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	r := strings.ToLower(fmt.Sprintf(
		"%s\\.%s\\..*\\.%s",
		cleanAPIVersionForFilesystem(s.apiVersion),
		s.kind,
		s.encoder.Encoding(),
	))

	regex, err := regexp.Compile(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	refs := make([]string, 0)
	for _, entry := range dirEntries {
		if entry.IsDir() || !regex.MatchString(entry.Name()) {
			continue
		}

		refs = append(refs, entry.Name())
	}

	return refs, nil
}

// read tries to read file corresponding to the specified object's name.
// Returns a pointer to an unmarshalled T and an error.
// Warn: read can return a nil pointer if resource wasn't find
// The reason read returns a possible nil pointer is to avoid raising an error if the file we're trying to read does not
// exist.
func (s *FilesystemOperator) read(path string) (*ResourceDefinition, error) {
	if ok, err := fileExist(path); !ok {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	resource, err := ReadEncodedFile(path)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// readWithFilename makes a call to read
func (s *FilesystemOperator) readWithFilename(name string) (*ResourceDefinition, error) {
	return s.read(s.filepathFromFilename(name))
}

// readWithResourceName makes a call to read
func (s *FilesystemOperator) readWithResourceName(name string) (*ResourceDefinition, error) {
	return s.read(s.filepathFromResourceName(name))
}

// write tries to write T to filesystem.
func (s *FilesystemOperator) write(resource *ResourceDefinition) error {
	path := s.filepathFromResourceName(resource.Metadata.Name)
	logger.Debug(path)
	// Ensures the base dir exist before trying to write into it.
	if err := mkBaseDir(path); err != nil {
		return err
	}

	return WriteEncodedFile(path, resource)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (s *FilesystemOperator) filepathFromFilename(name string) string {
	return filepath.Join(s.resourceDir, name)
}

// filepathFromResourceName computes the resourceDir to the corresponding resource name, based on the naming convention
func (s *FilesystemOperator) filepathFromResourceName(name string) string {
	return s.filepathFromFilename(s.filename(name))
}

// filepathFromResourceName computes the resource filename, based on the naming convention
func (s *FilesystemOperator) filename(name string) string {
	return strings.ToLower(fmt.Sprintf(
		"%s.%s.%s.%s",
		cleanAPIVersionForFilesystem(s.apiVersion),
		s.kind,
		name,
		s.encoder.Encoding(),
	))
}

// NewFilesystemOperator instantiate a new strategy
func NewFilesystemOperator(apiVersion APIVersion, kind Kind, resourceDir string, encoding Encoding) (*FilesystemOperator, error) { //nolint:lll
	encoder, err := NewEncoder(encoding)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	apiVersion, err = apiVersion.Validate()
	if err != nil {
		return nil, err
	}

	kind, err = kind.Validate()
	if err != nil {
		return nil, err
	}

	return &FilesystemOperator{
		apiVersion:  apiVersion,
		kind:        kind,
		resourceDir: resourceDir,
		encoder:     encoder,
	}, nil
}

//----------------------------------------------------------------------------------------------------------------------
// GitOperator
//----------------------------------------------------------------------------------------------------------------------

// GitOperator uses FilesystemStrategy as a backend, and leverages Git for version control.
type GitOperator struct {
	innerStrategy FilesystemOperator
}

//----------------------------------------------------------------------------------------------------------------------
// Operator Utils
//----------------------------------------------------------------------------------------------------------------------

// resourceExist checks if a named resource already exist
func resourceExist(operator Operator, name string) (bool, error) {
	arr, err := operator.Get(ToPointer(name))
	if err != nil {
		return false, err
	}

	if len(arr) == 0 {
		return false, nil
	}

	return true, nil
}

func defaultOperatorStrategy() OperatorStrategy {
	return FileSystemOperatorStrategy
}

// cleanAPIVersionForFilesystem transforms `vib/v1alpha1` into `vib_v1alpha1`
func cleanAPIVersionForFilesystem(s APIVersion) string {
	return strings.ReplaceAll(string(s), "/", "_")
}
