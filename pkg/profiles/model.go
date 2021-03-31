package profiles

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ProfileSpec defines the desired state of Profile
type ProfileSpec struct {
	// Description is some text to allow a user to identify what this profile installs.
	Description string `json:"description,omitempty"`
	// Artifacts is a list of Profile artifacts
	// can be one of HelmChart, TODO
	Artifacts []Artifact `json:"artifacts,omitempty"`
}

type Artifact struct {
	// Name is the name of the Artifact
	Name string `json:"name,omitempty"`
	// Path is the local path to the Artifact in the Profile repo
	Path string `json:"path,omitempty"`
}

// Profile is the Schema for the profiles API
type Profile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ProfileSpec `json:"spec,omitempty"`
}
