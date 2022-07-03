package echo

import (
	"context"
	"fmt"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestEcho(t *testing.T) {

	// >>>>> mimic "podman --version"
	// >>>>> 相对于 "podman --version"

	// create a client
	// 创建一个客户端
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	conn, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// >>>>> mimic "podman image exists k8s.gcr.io/echoserver:1.4"
	// >>>>> 相对于 "podman image exists k8s.gcr.io/echoserver:1.4"

	// check the image exists
	// 检查镜像是否存在
	exists, err := images.Exists(conn, "k8s.gcr.io/echoserver:1.4", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("image exists ?", exists)

	// pull an image if it doesn't exist
	// 如果不存在，重新下载镜像
	if !exists {

		// >>>>> mimic "podman image pull k8s.gcr.io/echoserver:1.4"
		// >>>>> 相对于 "podman image pull k8s.gcr.io/echoserver:1.4"

		// 下载镜像
		// pull an image
		_, err = images.Pull(conn, "k8s.gcr.io/echoserver:1.4", nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check the image exists again
		// 再一次检查镜像是否存在
		exists, err := images.Exists(conn, "k8s.gcr.io/echoserver:1.4", nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// the image must exit
		// 镜像必须存在
		require.Equal(t, true, exists)

		fmt.Println("image exists after pull ?", exists)
	}

	// >>>>> mimic "podman ps --filter label=io.podman.compose.project=echo -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'"
	// >>>>> 相对于 "podman ps --filter label=io.podman.compose.project=echo -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'"

	// 启动计划是否已经存在
	// check if the plan exists

	listOptions := containers.ListOptions{
		Filters: map[string][]string{
			"label": {"io.podman.compose.project=echo"},
		},
	}

	listContainer, err := containers.List(conn, &listOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(listContainer) > 0 {
		fmt.Println("plan exists ?", exists)
		exists = true
		for _, container := range listContainer {
			if value, ok := container.Labels["io.podman.compose.config-hash"]; ok == true {
				fmt.Println("hash: ", value)
			}
		}
		// 存在就中断
		// if it exists, stop all
		os.Exit(1)
	} else {
		fmt.Println("plan exists ?", exists)
		exists = false
	}

	// >>>>> mimic "podman pod create --name=pod_echo --infra=false --share="
	// >>>>> 相对于 "podman pod create --name=pod_echo --infra=false --share="

	// 建立夹子
	// create a pod

	pspec := entities.PodSpec{
		PodSpecGen: specgen.PodSpecGenerator{
			PodBasicConfig: specgen.PodBasicConfig{
				Name:    "pod_echo",
				NoInfra: true,
			},
			PodCgroupConfig:   specgen.PodCgroupConfig{},
			PodResourceConfig: specgen.PodResourceConfig{},
			InfraContainerSpec: &specgen.SpecGenerator{
				specgen.ContainerBasicConfig{},
				specgen.ContainerStorageConfig{},
				specgen.ContainerSecurityConfig{},
				specgen.ContainerCgroupConfig{},
				specgen.ContainerNetworkConfig{},
				specgen.ContainerResourceConfig{},
				specgen.ContainerHealthCheckConfig{},
			},
		},
	}

	preport, _ := pods.CreatePodFromSpec(conn, &pspec)

	_, _ = pods.Start(conn, preport.Id, nil)

}
