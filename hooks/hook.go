package hooks

import (
	"errors"

	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/trees"
)

// Hook defines an interface which handles the retrieval and installation of
// a shell.ManifestAttr. It expects to return two values, the markup to be installed
// into the page and a boolean indicating if it should be added into the head.
type Hook interface {
	Fetch(*router.Router, shell.ManifestAttr) (res *trees.Markup, addToHeader bool, err error)
}

// hooks provides a global registry for registering hooks.
var hooks = newRegistry()

// Register adds the Hook into the registry, passing the giving struct with the
// provider for usage.
func Register(name string, hook Hook) {
	hooks.Create(name, hook)
}

// RegisterManifest adds a manifest into the global register for access by the name.
func RegisterManifest(attr shell.ManifestAttr) {
	hooks.AddManifest(attr)
}

// GetManifest retrieves the manifest with the given name.
func GetManifest(name string) (shell.ManifestAttr, error) {
	return hooks.GetManifest(name)
}

// Get retrieves the giving provider by the name.
func Get(name string) (Hook, error) {
	return hooks.Get(name)
}

// registry defines a structure for registery for Providers where all providers
// handles the lifecycle process for a asset from installation to removal.
type registry struct {
	hooks     map[string]Hook
	manifests map[string]shell.ManifestAttr
}

// newRegistry returns a new instance of Registry.
func newRegistry() *registry {
	return &registry{
		hooks:     make(map[string]Hook),
		manifests: make(map[string]shell.ManifestAttr),
	}
}

// Create returns a provider associated with the given key if found and if not, then
// it is created and added into the hook lists.
func (r *registry) Create(name string, hook Hook) {
	r.hooks[name] = hook
}

// AddManifest adds the giving manifest into the global manifest file.
// It avoids manifests of same name once added.
func (r *registry) AddManifest(attr shell.ManifestAttr) {
	r.manifests[attr.Name] = attr
}

// ErrHookNotFound is returned when the desired hook is not on the registry.
var ErrHookNotFound = errors.New("Hook not found")

// Get returns a giving provider with the provided key.
func (r *registry) Get(name string) (Hook, error) {
	hl, ok := r.hooks[name]
	if !ok {
		return nil, ErrHookNotFound
	}

	return hl, nil
}

// ErrManifestNotFound is returned when the desired manifest is not on the registry.
var ErrManifestNotFound = errors.New("Manifest not found")

// GetManifest returns a giving provider with the provided key.
func (r *registry) GetManifest(name string) (shell.ManifestAttr, error) {
	hl, ok := r.manifests[name]
	if !ok {
		return shell.ManifestAttr{}, ErrManifestNotFound
	}

	return hl, nil
}
