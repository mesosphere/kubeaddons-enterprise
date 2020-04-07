package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/blang/semver"
	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/kind"
	"github.com/mesosphere/ksphere-testing-framework/pkg/experimental"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/repositories"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	addontesters "github.com/mesosphere/kubeaddons/test/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha3"
)

const (
	controllerBundle         = "https://mesosphere.github.io/kubeaddons/bundle.yaml"
	defaultKubernetesVersion = "1.16.4"
)

var addonTestingGroups = make(map[string][]AddonTestConfiguration)

type AddonTestConfiguration struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Override the values for helm chart and parameters for kudo operators
	Override string `json:"override,omitempty" yaml:"override,omitempty"`
	// List of requirements to be removed from the addon
	RemoveDependencies []string `json:"removeDependencies,omitempty" yaml:"removeDependencies,omitempty"`
}

var (
	cat       catalog.Catalog
	localRepo repositories.Repository
	groups    map[string][]v1beta1.AddonInterface
)

func init() {
	var err error
	localRepo, err = local.NewRepository("local", "../addons/")
	if err != nil {
		panic(err)
	}

	cat, err = catalog.NewCatalog(localRepo)
	if err != nil {
		panic(err)
	}

	groups, err = experimental.AddonsForGroupsFile("groups.yaml", cat)
	if err != nil {
		panic(err)
	}

	for group, addons := range groups {
		for _, addon := range addons {
			overrides := overridesForAddon(addon.GetName())
			removeDeps := removeDepsForAddon(addon.GetName())
			cfg := AddonTestConfiguration{
				Name:               addon.GetName(),
				Override:           overrides,
				RemoveDependencies: removeDeps,
			}
			addonTestingGroups[group] = append(addonTestingGroups[group], cfg)
		}
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

func TestSparkGroup(t *testing.T) {
	if err := testgroup(t, "spark"); err != nil {
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

func cleanupNodeVolumes(numberVolumes int, nodePrefix string) error {
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
		if err := cleanupNodeVolumes(3, u.String()); err != nil {
			t.Logf("error: %s", err)
		}
	}()

	cluster, err := kind.NewCluster(version)
	if err != nil {
		return err
	}
	defer cluster.Cleanup()

	if err := kubectl("apply", "-f", controllerBundle); err != nil {
		return err
	}

	addons, err := addons(addonTestingGroups[groupname]...)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})
	go experimental.LoggingHook(t, cluster, wg, stop)

	deployplan, err := addontesters.DeployAddons(t, cluster, addons...)
	if err != nil {
		return err
	}

	defaultplan, err := addontesters.WaitForAddons(t, cluster, addons...)
	if err != nil {
		return err
	}

	cleanupplan, err := addontesters.CleanupAddons(t, cluster, addons...)
	if err != nil {
		return err
	}

	th := testharness.NewSimpleTestHarness(t)
	th.Load(addontesters.ValidateAddons(addons...), deployplan, defaultplan, cleanupplan)

	defer th.Cleanup()
	th.Validate()
	th.Deploy()
	th.Default()

	close(stop)
	wg.Wait()

	return nil
}

func addons(addonConfigs ...AddonTestConfiguration) ([]v1beta1.AddonInterface, error) {
	var testAddons []v1beta1.AddonInterface

	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return testAddons, err
	}
	for _, addonConfig := range addonConfigs {
		addon, err := repo.GetAddon(addonConfig.Name)
		if err != nil {
			return testAddons, err
		}
		overrides(addon[0], addonConfig)
		if addon[0].GetNamespace() == "" {
			addon[0].SetNamespace("default")
		}
		// TODO - we need to re-org where these filters are done (see: https://jira.mesosphere.com/browse/DCOS-63260)
		testAddons = append(testAddons, addon[0])
	}

	if len(testAddons) != len(addonConfigs) {
		return testAddons, fmt.Errorf("got %d addons, expected %d", len(testAddons), len(addonConfigs))
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
			for _, addonConfig := range v {
				if addonConfig.Name == addon.GetName() {
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

// -----------------------------------------------------------------------------
// Private - CI Values Overrides
// -----------------------------------------------------------------------------

func overrides(addon v1beta1.AddonInterface, config AddonTestConfiguration) {
	if config.Override != "" {
		// override helm chart values
		if addon.GetAddonSpec().ChartReference != nil {
			addon.GetAddonSpec().ChartReference.Values = &config.Override
		}
		//override kudo operator default values
		if addon.GetAddonSpec().KudoReference != nil {
			addon.GetAddonSpec().KudoReference.Parameters = &config.Override
		}
	}

	for _, toRemove := range config.RemoveDependencies {
		removeDependencyFromAddon(addon, toRemove)
	}

}

func removeDependencyFromAddon(addon v1beta1.AddonInterface, toRemove string) {
	for index, labelSelector := range addon.GetAddonSpec().Requires {
		for label, value := range labelSelector.MatchLabels {
			if value == toRemove {
				delete(labelSelector.MatchLabels, label)
				if len(labelSelector.MatchLabels) == 0 {
					addon.GetAddonSpec().Requires = removeLabelsIndex(addon.GetAddonSpec().Requires, index)
				}
				return
			}

		}
	}
}

func removeLabelsIndex(s []metav1.LabelSelector, index int) []metav1.LabelSelector {
	return append(s[:index], s[index+1:]...)
}

// TODO currently overrides are hardcoded but will be promoted into Groups in the future
// see D2IQ-64898
func overridesForAddon(name string) string {
	switch name {
	case "kafka":
		return `ZOOKEEPER_URI: zookeeper-cs
BROKER_MEM: 32Mi
BROKER_CPUS: 20m
BROKER_COUNT: 1
ADD_SERVICE_MONITOR: false
`
	case "cassandra":
		return `NODE_COUNT: 1
NODE_DISK_SIZE_GIB: 1
NODE_MEM_MIB: 128
PROMETHEUS_EXPORTER_ENABLED: "false"
`
	case "spark":
		return `enableMetrics: "false"`
	case "zookeeper":
		return `MEMORY: "32Mi"
CPUS: 50m
NODE_COUNT: 1
`
	}

	return ""
}

// TODO currently depremovals are hardcoded but will be promoted into Groups in the future
// see D2IQ-64898
func removeDepsForAddon(name string) []string {
	switch name {
	// remove promethues dependency for CI
	// https://jira.d2iq.com/browse/D2IQ-63819
	case "spark":
		return []string{"prometheus"}
	case "kafka":
		return []string{"prometheus"}
	}
	return []string{}
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
