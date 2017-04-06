package shell

import "errors"

// ManifestAttr defines a structure which stores a series of
// data pertaining to a specific resource.
type ManifestAttr struct {
	Size          int               `json:"size"`
	Remote        bool              `json:"remote"`
	Init          bool              `json:"init"`
	IsGlobal      bool              `json:"is_global"`
	Localize      bool              `json:"localize"`
	B64Encode     bool              `json:"b64_encode"`
	ContentB64    bool              `json:"content_b64"`
	Base64Padding bool              `json:"base64_padding"`
	ID            string            `json:"appmanifest_id,omitempty"`
	Name          string            `json:"name"`
	Path          string            `json:"path"`
	Content       string            `json:"content"`
	Meta          map[string]string `json:"meta"`
	HookName      string            `json:"hook_name,omitempty"`
}

// AppManifest defines a structure which holds a series of
// manifests data related to specific resources.
type AppManifest struct {
	GlobalScope bool               `json:"global_scope"`
	Name        string             `json:"name"`
	Depends     []string           `json:"depends"`
	Manifests   []ManifestAttr     `json:"manifests"`
	Relation    *ComponentRelation `json:"relation"`
}

// NewAppManifest returns a instance of the AppManifest type.
func NewAppManifest(name string) *AppManifest {
	return &AppManifest{
		Name: name,
	}
}

// ComponentRelation defines a structure which stores specific
// data about the relation of a giving component regarding the
// manifests that correlate to it.
type ComponentRelation struct {
	Name       string   `json:"name"`
	Package    string   `json:"package"`
	Composites []string `json:"composites,omitempty"`
	FieldTypes []string `json:"fieldtypes,omitempty"`
}

// FindByRelation provides a convenient method to search a giving manifests for a specific
// ComponentRelation.
func FindByRelation(apps []AppManifest, relationName string) (AppManifest, error) {
	for _, app := range apps {
		if app.Relation.Name != relationName {
			continue
		}

		return app, nil
	}

	return AppManifest{}, errors.New("Not Found")
}
