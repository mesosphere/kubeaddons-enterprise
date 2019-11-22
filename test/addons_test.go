package test

import (
	"fmt"
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

	addons, err := temp.Addons("../addons/")
	if err != nil {
		t.Fatal(err)
	}

	// Kafka needs special attention right now to ensure it can find its dependency ZK
	if revisions, ok := addons["kafka"]; ok {
		for _, revision := range revisions {
			zkuri := fmt.Sprintf("ZOOKEEPER_URI: zookeeper-cs.%s.svc", addons["zookeeper"][0].GetNamespace())
			revision.GetAddonSpec().KudoReference.Parameters = &zkuri
		}
	}

	testAddons := []v1beta1.AddonInterface{}
	for _, v := range addons {
		// TODO - for right now, we're only testing the latest revision.
		// We're waiting on additional features from the test harness to
		// expand this, see https://jira.mesosphere.com/browse/DCOS-61266
		testAddons = append(testAddons, v[0])
	}

	th, err := test.NewBasicTestHarness(t, cluster, testAddons...)
	if err != nil {
		t.Fatal(err)
	}
	defer th.Cleanup()

	th.Validate()
	th.Deploy()

}
