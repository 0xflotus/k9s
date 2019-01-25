package resource

import (
	"log"
	"strconv"

	"github.com/k8sland/k9s/resource/k8s"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
)

// Secret tracks a kubernetes resource.
type Secret struct {
	*Base
	instance *v1.Secret
}

// NewSecretList returns a new resource list.
func NewSecretList(ns string) List {
	return NewSecretListWithArgs(ns, NewSecret())
}

// NewSecretListWithArgs returns a new resource list.
func NewSecretListWithArgs(ns string, res Resource) List {
	return newList(ns, "secret", res, AllVerbsAccess)
}

// NewSecret instantiates a new Secret.
func NewSecret() *Secret {
	return NewSecretWithArgs(k8s.NewSecret())
}

// NewSecretWithArgs instantiates a new Secret.
func NewSecretWithArgs(r k8s.Res) *Secret {
	cm := &Secret{
		Base: &Base{
			caller: r,
		},
	}
	cm.creator = cm
	return cm
}

// NewInstance builds a new Secret instance from a k8s resource.
func (*Secret) NewInstance(i interface{}) Columnar {
	cm := NewSecret()
	switch i.(type) {
	case *v1.Secret:
		cm.instance = i.(*v1.Secret)
	case v1.Secret:
		ii := i.(v1.Secret)
		cm.instance = &ii
	default:
		log.Fatalf("Unknown %#v", i)
	}
	cm.path = cm.namespacedName(cm.instance.ObjectMeta)
	return cm
}

// Marshal resource to yaml.
func (r *Secret) Marshal(path string) (string, error) {
	ns, n := namespaced(path)
	i, err := r.caller.Get(ns, n)
	if err != nil {
		return "", err
	}

	cm := i.(*v1.Secret)
	cm.TypeMeta.APIVersion = "v1"
	cm.TypeMeta.Kind = "Secret"
	raw, err := yaml.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

// Header return resource header.
func (*Secret) Header(ns string) Row {
	hh := Row{}
	if ns == AllNamespaces {
		hh = append(hh, "NAMESPACE")
	}
	return append(hh, "NAME", "TYPE", "DATA", "AGE")
}

// Fields retrieves displayable fields.
func (r *Secret) Fields(ns string) Row {
	ff := make(Row, 0, len(r.Header(ns)))
	i := r.instance
	if ns == AllNamespaces {
		ff = append(ff, i.Namespace)
	}

	return append(ff,
		i.Name,
		string(i.Type),
		strconv.Itoa(len(i.Data)),
		toAge(i.ObjectMeta.CreationTimestamp),
	)
}

// ExtFields returns extended fields in relation to headers.
func (*Secret) ExtFields() Properties {
	return Properties{}
}
