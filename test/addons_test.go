package test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/blang/semver"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/repositories/addons/revisions"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const testNamespace = "default"

// TestAddons tests deployment of all addons in this repository
func TestAddons(t *testing.T) {
	addons, err := temp.Addons("../addons/")
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Kafka needs special attention right now to ensure it can find its dependency ZK
	kafkaFilters(addons)

	// TODO: for speed, prometheus options and requirements are disabled for Jenkins for now
	if err := jenkinsFilters(addons); err != nil {
		t.Fatal(err)
	}

	testAddons := []v1beta1.AddonInterface{}
	for _, v := range addons {
		// TODO - for right now, we're only testing the latest revision.
		// We're waiting on additional features from the test harness to
		// expand this, see https://jira.mesosphere.com/browse/DCOS-61266
		testAddons = append(testAddons, v[0])
	}

	cluster, err := kind.NewCluster(semver.MustParse("1.16.3"))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster, "kind"); err != nil {
		t.Fatal(err)
	}

	if err := deployKudo(cluster); err != nil {
		t.Fatal(err)
	}

	for _, testAddon := range testAddons {
		testAddon.SetNamespace(testNamespace)
	}

	th, err := test.NewBasicTestHarness(t, cluster, testAddons...)
	if err != nil {
		t.Fatal(err)
	}
	defer th.Cleanup()

	th.Validate()
	th.Deploy()

}

func kafkaFilters(addons map[string]revisions.AddonRevisions) {
	if revisions, ok := addons["kafka"]; ok {
		for _, revision := range revisions {
			zkuri := fmt.Sprintf("ZOOKEEPER_URI: zookeeper-cs.%s.svc", addons["zookeeper"][0].GetNamespace())
			revision.GetAddonSpec().KudoReference.Parameters = &zkuri
		}
	}
}

// TODO - due to issues with how CRDs are handled by the helm chart for Jenkins
// we're going to use this to disable options that would utilize CRDs that wont
// actually be deployed by the chart. This is something that will need to be improved upstream.
var jenkinsOverrides = `---
master:
  useSecurity: false
  installPlugins:
    - kubernetes:1.18.2
    - workflow-job:2.33
    - workflow-aggregator:2.6
    - credentials-binding:1.19
    - git:3.11.0
  csrf:
    defaultCrumbIssuer:
      enabled: false
      proxyCompatability: false
  serviceType: "ClusterIP"
  jenkinsUriPrefix: "/jenkins"
  path: /jenkins
  ingress:
    enabled: false
  prometheus:
    enabled: false
`

func jenkinsFilters(addons map[string]revisions.AddonRevisions) error {
	if revisions, ok := addons["jenkins"]; ok {
		for _, revision := range revisions {
			// TODO: for now we're going to remove deps for speed, jenkins can deploy usably without traefik.
			revision.GetAddonSpec().Requires = make([]v1.LabelSelector, 0)
			revision.GetAddonSpec().ChartReference.Values = &jenkinsOverrides
		}
	}
	return nil
}

// TODO - this is a temporary lockin to a specific revision, later we should switch to using
// the repository library instead and selecting the latest `v0.7.x` release by revision semver
const kudoAddonURL = "https://raw.githubusercontent.com/mesosphere/kubernetes-base-addons/master/addons/kudo/0.7.x/kudo-1.yaml"

func deployKudo(cluster test.Cluster) error {
	cmd := exec.Command("kubectl", "--kubeconfig", cluster.ConfigPath(), "apply", "-f", kudoAddonURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
