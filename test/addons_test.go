package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/blang/semver"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"

	"sigs.k8s.io/kind/pkg/apis/config/v1alpha3"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const (
	defaultKubernetesVersion = "1.16.4"
	patchStorageClass        = `{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}`
)

var addonTestingGroups = make(map[string][]string)

func init() {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, addonTestingGroups); err != nil {
		panic(err)
	}
}

func TestValidateUnhandledAddons(t *testing.T) {
	unhandled, err := findUnhandled()
	if err != nil {
		t.Fatal(err)
	}

	if len(unhandled) != 0 {
		names := make([]string, len(unhandled))
		for _, addon := range unhandled {
			names = append(names, addon.GetName())
		}
		t.Fatal(fmt.Errorf("the following addons are not handled as part of a testing group: %+v", names))
	}
}

func TestGeneralGroup(t *testing.T) {
	if err := testgroup(t, "general"); err != nil {
		t.Fatal(err)
	}
}

func TestKafkaGroup(t *testing.T) {
	if err := testgroup(t, "kafka"); err != nil {
		t.Fatal(err)
	}
}

func TestCassandraGroup(t *testing.T) {
	if err := testgroup(t, "cassandra"); err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

func createNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha3.Node) error {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	dockerClient.NegotiateAPIVersion(context.TODO())

	for index := 0; index < numberVolumes; index++ {
		volumeName := fmt.Sprintf("%s-%d", nodePrefix, index)

		volume, err := dockerClient.VolumeCreate(context.TODO(), volumetypes.VolumeCreateBody{
			Driver: "local",
			Name:   volumeName,
		})
		if err != nil {
			return fmt.Errorf("creating volume for node: %w", err)
		}

		node.ExtraMounts = append(node.ExtraMounts, v1alpha3.Mount{
			ContainerPath: fmt.Sprintf("/mnt/disks/%s", volumeName),
			HostPath:      volume.Mountpoint,
		})
	}

	return nil
}

func cleanupNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha3.Node) error {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	dockerClient.NegotiateAPIVersion(context.TODO())

	for index := 0; index < numberVolumes; index++ {
		volumeName := fmt.Sprintf("%s-%d", nodePrefix, index)

		if err := dockerClient.VolumeRemove(context.TODO(), volumeName, false); err != nil {
			return fmt.Errorf("removing volume for node: %w", err)
		}
	}

	return nil
}

func testgroup(t *testing.T, groupname string) error {
	t.Logf("testing group %s", groupname)

	version, err := semver.Parse(defaultKubernetesVersion)
	if err != nil {
		return err
	}

	u := uuid.New()

	node := v1alpha3.Node{}
	if err := createNodeVolumes(3, u.String(), &node); err != nil {
		return err
	}
	defer func() {
		if err := cleanupNodeVolumes(3, u.String(), &node); err != nil {
			t.Logf("error: %s", err)
		}
	}()

	cluster, err := kind.NewCluster(version)
	if err != nil {
		return err
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster, "kind"); err != nil {
		return err
	}

	addons, err := addons(addonTestingGroups[groupname]...)
	if err != nil {
		return err
	}

	ph, err := test.NewBasicTestHarness(t, cluster, addons...)
	if err != nil {
		return err
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()

	return nil
}

func addons(names ...string) ([]v1beta1.AddonInterface, error) {
	var testAddons []v1beta1.AddonInterface

	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return testAddons, err
	}
	addons, err := repo.ListAddons()
	if err != nil {
		return testAddons, err
	}

	for _, addon := range addons {
		for _, name := range names {
			overrides(addon[0])
			if addon[0].GetName() == name {
				if addon[0].GetNamespace() == "" {
					addon[0].SetNamespace("default")
				}
				// TODO - we need to re-org where these filters are done (see: https://jira.mesosphere.com/browse/DCOS-63260)
				filters(addon[0])
				testAddons = append(testAddons, addon[0])
			}
		}
	}

	if len(testAddons) != len(names) {
		return testAddons, fmt.Errorf("got %d addons, expected %d", len(testAddons), len(names))
	}

	return testAddons, nil
}

func findUnhandled() ([]v1beta1.AddonInterface, error) {
	var unhandled []v1beta1.AddonInterface
	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return unhandled, err
	}
	addons, err := repo.ListAddons()
	if err != nil {
		return unhandled, err
	}

	for _, revisions := range addons {
		addon := revisions[0]
		found := false
		for _, v := range addonTestingGroups {
			for _, name := range v {
				if name == addon.GetName() {
					found = true
				}
			}
		}
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func filters(addon v1beta1.AddonInterface) {
	switch addon.GetName() {
	case "cassandra":
		cassandraParameters := "NODE_COUNT: 1\nNODE_DISK_SIZE_GIB: 1\nPROMETHEUS_EXPORTER_ENABLED: \"false\""
		addon.GetAddonSpec().KudoReference.Parameters = &cassandraParameters
	case "kafka":
		zkuri := fmt.Sprintf("ZOOKEEPER_URI: zookeeper-cs.%s.svc", addon.GetNamespace())
		addon.GetAddonSpec().KudoReference.Parameters = &zkuri
	}
}

// -----------------------------------------------------------------------------
// Private - CI Values Overrides
// -----------------------------------------------------------------------------

// TODO: a temporary place to put configuration overrides for addons
// See: https://jira.mesosphere.com/browse/DCOS-62137
func overrides(addon v1beta1.AddonInterface) {
	if v, ok := addonOverrides[addon.GetName()]; ok {
		addon.GetAddonSpec().ChartReference.Values = &v
	}
}

var addonOverrides = map[string]string{}
