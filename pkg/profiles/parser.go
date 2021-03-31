package profiles

import (
	"fmt"

	"sigs.k8s.io/yaml"
)

// ParseBytes takes a slice of bytes, parses it as YAML and returns the
// resulting Profile.
func ParseBytes(b []byte) (*Profile, error) {
	p := &Profile{}
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}
	return p, nil
}
