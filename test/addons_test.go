package test

import (
	"io/ioutil"
	"testing"

	"github.com/blang/semver"

	"github.com/mesosphere/kubeaddons/api/v1beta1"
	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

// TestAddons tests deployment of all addons in this repository
func TestAddons(t *testing.T) {
	cluster, err := kind.NewCluster(semver.MustParse("1.16.3"))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster); err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("../addons/zookeeper/0.x/zookeeper.yaml")
	if err != nil {
		t.Fatal(err)
	}
	ro, err := test.DecodeObjectFromManifest(b)
	if err != nil {
		t.Fatal(err)
	}
	addon, ok := ro.(v1beta1.AddonInterface)
	if !ok {
		t.Fatalf("invalid addon provided: %+v", ro)
	}
	addon.SetNamespace("default")

	th, err := test.NewBasicTestHarness(t, cluster, addon)
	if err != nil {
		t.Fatal(err)
	}

	th.Validate()
	th.Deploy()
	th.Cleanup()
}
