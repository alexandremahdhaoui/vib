package api

import (
	"fmt"
	"github.com/alexandremahdhaoui/vib/pkg/logger"
)

type APIServer interface {
	//-----------------
	// Manage the API Server

	// Register a new APIKind to the APIServer
	Register(apiKind APIKind) error

	// Unregister an APIKind from the APIServer
	Unregister(apiVersion APIVersion, kind Kind) error

	//-----------------
	// Basic Operations

	// Get existing resources by Kind and optionally by APIVersion and/or name
	// Please set APIVersion & name to nil if you don't need to filter
	Get(apiVersion *APIVersion, kind Kind, name *string) ([]ResourceDefinition, error)

	// Create a new resource definition
	Create(resource *ResourceDefinition) error

	// Update an existing resource definition or creates it
	Update(apiVersion *APIVersion, kind Kind, name string, resource *ResourceDefinition) error

	// Delete by kind & name. Optional apiVersion.
	Delete(apiVersion *APIVersion, kind Kind, name string) error
}

func NewAPIServer() APIServer {
	return defaultAPIServer()
}

func defaultAPIServer() *localAPIServer {
	return newLocalAPIServer()
}

type localAPIServer struct {
	// apiKinds is a list of registered APIKind
	apiKinds map[Kind]map[APIVersion]APIKind
}

func (l *localAPIServer) Register(apiKind APIKind) error {
	apiVersion, err := apiKind.APIVersion().Validate()
	if err != nil {
		return err
	}

	kind, err := apiKind.Kind().Validate()
	if err != nil {
		return err
	}

	// get apiVersions' map if it exists
	if apiVersions, ok := l.apiKinds[kind]; ok {
		// We know there is an entry for kind in APIKinds map

		// Check if it's a nil entry
		if apiVersions == nil {
			// Create the map if we got a nil map
			l.apiKinds[kind] = make(map[APIVersion]APIKind)

			// Register the APIKind and return
			l.registerUnsafe(apiKind)
			return nil
		}

		// apiVersions map exist
		// apiVersions map is not nil.
		// -> check if there is already an entry for the specified kind-apiVersion
		if _, ok = apiVersions[apiVersion]; ok {
			err := fmt.Errorf("%w: cannot register APIKind %#v", logger.ErrAlreadyExist, apiKind)
			logger.Error(err)
			return err
		}

		// apiVersions map exist
		// apiVersion map is not nil
		// There is no already existing entry for apiVersion
		// -> We register APIKind and return
		l.registerUnsafe(apiKind)
		return nil
	}

	// No entry for kind in APIKinds map

	// -> Create an apiVersion Map for this kind
	l.apiKinds[kind] = make(map[APIVersion]APIKind)

	// -> Register APIKind and return
	l.registerUnsafe(apiKind)
	return nil
}

func (l *localAPIServer) registerUnsafe(apiKind APIKind) {
	apiVersion, err := apiKind.APIVersion().Validate()
	if err != nil {
		logger.Fatal(err)
	}

	kind, err := apiKind.Kind().Validate()
	if err != nil {
		logger.Fatal(err)
	}

	l.apiKinds[kind][apiVersion] = apiKind
}

func (l *localAPIServer) Unregister(apiVersion APIVersion, kind Kind) error {
	var err error

	apiVersion, err = apiVersion.Validate()
	if err != nil {
		return err
	}

	kind, err = kind.Validate()
	if err != nil {
		return err
	}

	// If Kind exist, get the available apiVersions map
	if versions, ok := l.apiKinds[kind]; ok {
		// delete the APIKind from the APIVersion if it exists
		// nothing happens if APIKind is not present or if the map is nil
		delete(versions, apiVersion)

		// if there is no APIVersions left for this Kind, then delete the Kind map
		if len(versions) == 0 {
			delete(l.apiKinds, kind)
		}
	}

	return nil
}

func (l *localAPIServer) Get(apiVersion *APIVersion, kind Kind, name *string) ([]ResourceDefinition, error) {
	err := ValidateAPIVersionPtr(apiVersion)
	if err != nil {
		return nil, err
	}

	kind, err = kind.Validate()
	if err != nil {
		return nil, err
	}

	// check if kind is registered
	if _, ok := l.apiKinds[kind]; !ok {
		// kind is not registered, return error
		err := ErrKind(kind)
		logger.Error(err)

		return nil, err
	}

	results := make([]ResourceDefinition, 0)
	apiKinds := l.apiKinds[kind]

	// We handle the condition where user specified an apiVersion
	if apiVersion != nil {
		dereferencedAPIVersion := *apiVersion
		// if corresponding APIKind exist, then we construct an "apiKinds" map made only of the specified apiVersion
		if apiKind, ok := apiKinds[dereferencedAPIVersion]; ok {
			apiKinds = make(map[APIVersion]APIKind)
			apiKinds[dereferencedAPIVersion] = apiKind
		} else {
			// specified APIVersion is not registered, return error
			err := ErrApiVersion(dereferencedAPIVersion, kind)
			logger.Error(err)

			return nil, err
		}
	}

	for _, apiKind := range apiKinds {
		// let's validate name before propagating it to Operator.Get
		err = ValidateResourceNamePtr(name)
		if err != nil {
			return nil, err
		}

		// Get(name *string) handles the condition where user specified a resource name.
		resource, err := apiKind.Operator().Get(name)
		if err != nil {
			return nil, err
		}

		results = append(results, resource...)
	}

	return results, nil
}

func (l *localAPIServer) Create(resource *ResourceDefinition) error {
	apiVersion, err := resource.APIVersion.Validate()
	if err != nil {
		return err
	}

	kind, err := resource.Kind.Validate()
	if err != nil {
		return err
	}

	if err = ValidateResourceName(resource.Metadata.Name); err != nil {
		return err
	}

	apiKinds, err := l.queryAPIKinds(&apiVersion, kind)
	if err != nil {
		return err
	}

	// Create resource
	// queryAPIKinds safely returns a non nil slice of size > 0
	return apiKinds[0].Operator().Create(resource) //nolint:wrapcheck
}

func (l *localAPIServer) Update(apiVersion *APIVersion, kind Kind, name string, resource *ResourceDefinition) error {
	err := ValidateAPIVersionPtr(apiVersion)
	if err != nil {
		return err
	}

	kind, err = resource.Kind.Validate()
	if err != nil {
		return err
	}

	if apiVersion == nil {
		if version, err := resource.APIVersion.Validate(); err != nil {
			return err
		} else {
			apiVersion = &version
		}
	}

	if err = ValidateResourceName(name); err != nil {
		return err
	}

	apiKinds, err := l.queryAPIKinds(apiVersion, kind)
	if err != nil {
		return err
	}

	// Update resource
	// queryAPIKinds safely returns a non nil slice of size > 0
	return apiKinds[0].Operator().Update(&name, resource)
}

func (l *localAPIServer) Delete(apiVersion *APIVersion, kind Kind, name string) error {
	err := ValidateAPIVersionPtr(apiVersion)
	if err != nil {
		return err
	}

	kind, err = kind.Validate()
	if err != nil {
		return err
	}

	apiKinds, err := l.queryAPIKinds(apiVersion, kind)
	if err != nil {
		return err
	}

	return apiKinds[0].Operator().Delete(name) //nolint:wrapcheck
}

// queryAPIKinds returns a list of APIKinds, querying by Kind and optionally by APIVersion
func (l *localAPIServer) queryAPIKinds(apiVersion *APIVersion, kind Kind) ([]APIKind, error) {
	err := ValidateAPIVersionPtr(apiVersion)
	if err != nil {
		return nil, err
	}

	kind, err = kind.Validate()
	if err != nil {
		return nil, err
	}

	// check if kind is registered
	if _, ok := l.apiKinds[kind]; !ok {
		// kind is not registered, return error
		err = ErrKind(kind)
		logger.Error(err)

		return nil, err
	}

	results := make([]APIKind, 0)

	apiKinds := l.apiKinds[kind]
	if len(apiKinds) == 0 {
		err = fmt.Errorf("%w: kind %q cannot be found", logger.ErrNotFound, kind)
		logger.Error(err)

		return nil, err
	}

	// We handle the condition where user specified an apiVersion
	if apiVersion != nil {
		version := *apiVersion

		// if corresponding APIKind exist, then we construct an "apiKinds" map made only of the specified apiVersion
		if _, ok := apiKinds[version]; !ok {
			// specified APIVersion is not registered, return error
			err = ErrApiVersion(version, kind)
			logger.Error(err)

			return nil, err
		}

		results = append(results, apiKinds[version])
		return results, nil
	}

	for _, apiKind := range apiKinds {
		results = append(results, apiKind)
	}

	return results, nil
}

func newLocalAPIServer() *localAPIServer {
	return &localAPIServer{
		apiKinds: make(map[Kind]map[APIVersion]APIKind),
	}
}
