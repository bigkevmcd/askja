package operations

import (
	"fmt"
	"path"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func filenameFrom(base string, o runtime.Object) (string, error) {
	oa, err := meta.Accessor(o)
	if err != nil {
		return "", fmt.Errorf("failed to get the object meta for object %#v: %w", o, err)
	}
	ta, err := meta.TypeAccessor(o)
	if err != nil {
		return "", fmt.Errorf("failed to get the type meta for object %#v: %w", o, err)
	}
	filename := strings.Join([]string{strings.ToLower(ta.GetKind()), oa.GetName()}, "_") + ".yaml"
	return path.Join(base, filename), nil
}
