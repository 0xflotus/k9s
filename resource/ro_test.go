package resource_test

import (
	"testing"

	"github.com/derailed/k9s/resource"
	"github.com/derailed/k9s/resource/k8s"
	m "github.com/petergtz/pegomock"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRoleMarshal(t *testing.T) {
	setup(t)

	ca := NewMockCaller()
	m.When(ca.Get("blee", "fred")).ThenReturn(k8sRole(), nil)

	cm := resource.NewRoleWithArgs(ca)
	ma, err := cm.Marshal("blee/fred")
	ca.VerifyWasCalledOnce().Get("blee", "fred")
	assert.Nil(t, err)
	assert.Equal(t, roleYaml(), ma)
}

func TestRoleListData(t *testing.T) {
	setup(t)

	ca := NewMockCaller()
	m.When(ca.List("blee")).ThenReturn(k8s.Collection{*k8sRole()}, nil)

	l := resource.NewRoleListWithArgs("blee", resource.NewRoleWithArgs(ca))
	// Make sure we can get deltas!
	for i := 0; i < 2; i++ {
		err := l.Reconcile()
		assert.Nil(t, err)
	}

	ca.VerifyWasCalled(m.Times(2)).List("blee")
	td := l.Data()
	assert.Equal(t, 1, len(td.Rows))
	assert.Equal(t, "blee", l.GetNamespace())
	assert.True(t, l.HasXRay())
	row := td.Rows["blee/fred"]
	assert.Equal(t, 2, len(row.Deltas))
	for _, d := range row.Deltas {
		assert.Equal(t, "", d)
	}
	assert.Equal(t, resource.Row{"fred"}, row.Fields[:1])
}

func TestRoleListDescribe(t *testing.T) {
	setup(t)

	ca := NewMockCaller()
	m.When(ca.Get("blee", "fred")).ThenReturn(k8sRole(), nil)
	l := resource.NewRoleListWithArgs("blee", resource.NewRoleWithArgs(ca))
	props, err := l.Describe("blee/fred")

	ca.VerifyWasCalledOnce().Get("blee", "fred")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(props))
}

// Helpers...

func k8sRole() *v1.Role {
	return &v1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:         "blee",
			Name:              "fred",
			CreationTimestamp: metav1.Time{Time: testTime()},
		},
	}
}

func newRole() resource.Columnar {
	return resource.NewRole().NewInstance(k8sRole())
}

func roleYaml() string {
	return `typemeta:
  kind: Role
  apiversion: rbac.authorization.k8s.io/v1
objectmeta:
  name: fred
  generatename: ""
  namespace: blee
  selflink: ""
  uid: ""
  resourceversion: ""
  generation: 0
  creationtimestamp: "2018-12-14T10:36:43.326972-07:00"
  deletiontimestamp: null
  deletiongraceperiodseconds: null
  labels: {}
  annotations: {}
  ownerreferences: []
  initializers: null
  finalizers: []
  clustername: ""
rules: []
`
}
