package shell

import (
	"encoding/base64"
	"errors"
)

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

// UnwrapBody returns the response body as plain text if it has been base64
// encode else if not, returns the body as expected.
func (m ManifestAttr) UnwrapBody() ([]byte, error) {
	if m.ContentB64 {
		if m.Base64Padding {
			mo, err := base64.StdEncoding.DecodeString(m.Content)
			if err != nil {
				return nil, err
			}

			return mo, nil
		}

		mo, err := base64.RawStdEncoding.DecodeString(m.Content)
		if err != nil {
			return nil, err
		}

		return mo, nil
	}

	return []byte(m.Content), nil
}

// EncodeBase64Content encodes the value and sets the content which was encoded to base64.
func (m *ManifestAttr) EncodeBase64Content(content string) error {
	if m.Base64Padding {
		m.Content = base64.StdEncoding.EncodeToString([]byte(content))
	} else {
		m.Content = base64.RawStdEncoding.EncodeToString([]byte(content))
	}
	return nil
}

// EncodeContentBase64 returns the content converted from the base64 value.
func (m ManifestAttr) EncodeContentBase64() (string, error) {
	if m.Base64Padding {
		return string(base64.StdEncoding.EncodeToString([]byte(m.Content))), nil
	}

	return string(base64.RawStdEncoding.EncodeToString([]byte(m.Content))), nil
}

// IsBase64Encode returns true/false if the content is base64 or should be
// base64 encoded.
func (m ManifestAttr) IsBase64Encode() bool {
	var b64 bool

	if m.Content != "" {
		if m.ContentB64 {
			b64 = true
		}
	} else {
		b64 = m.B64Encode
	}

	return b64
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
