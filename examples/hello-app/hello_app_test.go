package hello_app

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

func TestHelloApp(t *testing.T) {

	// >>>>> >>>>> >>>>> establish
	// >>>>> >>>>> >>>>> 创建部份

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

	// >>>>> mimic "podman image gcr.io/google-samples/hello-app:1.0"
	// >>>>> 相对于 "podman image gcr.io/google-samples/hello-app:1.0"

	// check the image exists
	// 检查镜像是否存在
	exists, err := images.Exists(conn, "gcr.io/google-samples/hello-app:1.0", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("image exists ?", exists)

	// pull an image if it doesn't exist
	// 如果不存在，重新下载镜像
	if !exists {

		// >>>>> mimic "podman image pull gcr.io/google-samples/hello-app:1.0"
		// >>>>> 相对于 "podman image pull gcr.io/google-samples/hello-app:1.0"

		// 下载镜像
		// pull an image
		_, err = images.Pull(conn, "gcr.io/google-samples/hello-app:1.0", nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// check the image exists again
		// 再一次检查镜像是否存在
		exists, err := images.Exists(conn, "gcr.io/google-samples/hello-app:1.0", nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// the image must exit
		// 镜像必须存在
		require.Equal(t, true, exists)

		fmt.Println("image exists after pull ?", exists)
	}

	// >>>>> mimic "podman ps --filter label=io.podman.compose.project=hello-app -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'"
	// >>>>> 相对于 "podman ps --filter label=io.podman.compose.project=hello-app -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'"

	// prepare for listing containers
	// 准备列出容器的过滤条件
	containerListOptions := containers.ListOptions{
		Filters: map[string][]string{
			"label": {"io.podman.compose.project=hello-app"},
		},
	}

	// check if the containers exists in the plan
	// 确认计划内的容器是否已经存在
	listContainer, err := containers.List(conn, &containerListOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// list all containers in the plan
	// 显示所有在计划中的容器
	if len(listContainer) > 0 {
		exists = true
		for _, container := range listContainer {
			if value, ok := container.Labels["io.podman.compose.config-hash"]; ok == true {
				fmt.Println("hash: ", value)
			}
		}
	} else {
		exists = false
	}
	fmt.Println("plan exists ?", exists)

	// 容器存在就中断
	// if it exists, stop all
	if exists == true {
		fmt.Println("exit !")
		os.Exit(1)
	}

	// >>>>> add "podman pod ls --format '{{ index .Labels "io.podman.compose.project"}}'"
	// >>>>> 新增 "podman pod ls --format '{{ index .Labels "io.podman.compose.project"}}'"

	// prepare for listing pods
	// 准备列出夹子的过滤条件
	podListOptions := pods.ListOptions{
		map[string][]string{
			"label": {"io.podman.compose.project=hello-app"},
		},
	}

	// check if the pods exists in the plan
	// 确认计划内的夹子是否已经存在
	listPods, err := pods.List(conn, &podListOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// list all pods in the plan
	// 显示所有在计划中的夹子
	exists = false
	podID := "" // pod's ID 夹子的编号
	for i := 0; i < len(listPods); i++ {
		exists = true
		fmt.Println("pod exists ?", exists, "(", listPods[i].Name, ")")
	}

	// create a pod if it doesn't exist
	// 如果不存在，重新建立夹子
	if !exists {

		// >>>>> mimic "podman pod create --label io.podman.compose.project=echo --name=hello-app --infra=false --share=" (adjusted)
		// >>>>> 相对于 "podman pod create --label io.podman.compose.project=echo --name=hello-app --infra=false --share=" (调整)
		// ( original one: "podman pod create --name=pod_hello-app --infra=false --share=" )

		// prepare data for creating a pod
		// 准备建立夹子的资料
		pspec := entities.PodSpec{
			PodSpecGen: specgen.PodSpecGenerator{
				PodBasicConfig: specgen.PodBasicConfig{
					Name:    "pod_hello-app",
					NoInfra: true,
					Labels: map[string]string{
						"io.podman.compose.project": "hello-app",
					},
				},
				PodCgroupConfig:    specgen.PodCgroupConfig{},
				PodResourceConfig:  specgen.PodResourceConfig{},
				InfraContainerSpec: &specgen.SpecGenerator{},
			},
		}

		// create a pod
		// 建立夹子
		preport, err := pods.CreatePodFromSpec(conn, &pspec)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// renew the pod's ID
		// 更新夹子的编号
		podID = preport.Id

		// the pod starts
		// 启动夹子
		_, err = pods.Start(conn, podID, nil)
		if err.Error() != "" {
			fmt.Println(err)
			os.Exit(1)
		}

		// check if the pod exists again
		// 再一次检查夹子是否存在
		listPods, err := pods.List(conn, &podListOptions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// crate a pod if it doesn't exist
		// 如果不存在，重新建立夹子
		if len(listPods) > 0 {
			exists = true
		} else {
			exists = false
		}
		fmt.Println("pod exists after creating ?", exists)

		// it must exist after creating a pod
		// 创建夹子后，检查夹子必须存在
		require.Equal(t, true, exists)

		// list the pod's name
		// 列出夹子的名称
		for i := 0; i < len(listPods); i++ {
			fmt.Println("pod: ", listPods[i].Name)
		}
	}
}
