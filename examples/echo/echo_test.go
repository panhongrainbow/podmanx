package echo

import (
	"context"
	"fmt"
	nettypes "github.com/containers/podman/v3/libpod/network/types"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/bindings/network"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {

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

	// prepare for listing containers
	// 准备列出容器的过滤条件
	containerListOptions := containers.ListOptions{
		Filters: map[string][]string{
			"label": {"io.podman.compose.project=echo"},
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
			"label": {"io.podman.compose.project=echo"},
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

		// >>>>> mimic "podman pod create --label io.podman.compose.project=echo --name=pod_echo --infra=false --share=" (adjusted)
		// >>>>> 相对于 "podman pod create --label io.podman.compose.project=echo --name=pod_echo --infra=false --share=" (调整)
		// ( original one: "podman pod create --name=pod_echo --infra=false --share=" )

		// prepare data for creating a pod
		// 准备建立夹子的资料
		pspec := entities.PodSpec{
			PodSpecGen: specgen.PodSpecGenerator{
				PodBasicConfig: specgen.PodBasicConfig{
					Name:    "pod_echo",
					NoInfra: true,
					Labels: map[string]string{
						"io.podman.compose.project": "echo",
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

	// >>>>> mimic "podman network exists echo_default"
	// >>>>> 相对于 "podman network exists echo_default"

	// check if the network exists
	// 检查网络是否存在
	exists, err = network.Exists(conn, "echo_default", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("network exists ?", exists)

	// crate a network if it doesn't exist
	// 如果不存在，重新建立网路
	if !exists {

		// >>>>> mimic "podman network create --label io.podman.compose.project=echo --label com.docker.compose.project=echo echo_default"
		// >>>>> 相对于 "podman network create --label io.podman.compose.project=echo --label com.docker.compose.project=echo echo_default"

		// prepare data for creating a network
		// 准备创建网路的资料
		nw := "echo_default"
		createOptions := network.CreateOptions{
			Name: &nw,
			Labels: map[string]string{
				"io.podman.compose.project":  "echo",
				"com.docker.compose.project": "echo",
			},
		}

		// create a network
		// 创建网路
		path, err := network.Create(conn, &createOptions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("network path ", path.Filename)

		// check if the network exists after creating it
		// 在创建网路完成后，再检查网络是否存在
		exists, err = network.Exists(conn, "echo_default", nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("network exists after creating ?", exists)
	}

	// >>>>> mimic "podman create --name=echo_web_1 --pod=pod_echo --label io.podman.compose.config-hash=57c9635d928f88954e01491e81ca0f9049014d0d205f0fb03c951e2fe09d582a --label io.podman.compose.project=echo --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=echo --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/echo --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=web --net echo_default --network-alias web -p 8080:8080 k8s.gcr.io/echoserver:1.4"
	// >>>>> 相对于 "podman create --name=echo_web_1 --pod=pod_echo --label io.podman.compose.config-hash=57c9635d928f88954e01491e81ca0f9049014d0d205f0fb03c951e2fe09d582a --label io.podman.compose.project=echo --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=echo --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/echo --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=web --net echo_default --network-alias web -p 8080:8080 k8s.gcr.io/echoserver:1.4"

	// prepare data for creating a container
	// 准备创建容器的资料
	s := specgen.NewSpecGenerator("k8s.gcr.io/echoserver:1.4", false)
	s.Name = "echo_web_1"

	// set the pod's network
	// 设定容器网路配置
	portMapping := make([]nettypes.PortMapping, 1, 1)
	portMapping[0].HostPort = 8080
	portMapping[0].ContainerPort = 8080

	s.ContainerNetworkConfig = specgen.ContainerNetworkConfig{
		CNINetworks:    []string{"echo_default"},
		NetworkOptions: map[string][]string{},
		Aliases: map[string][]string{
			"echo_default": {"web"},
		},
		PortMappings: portMapping,
	}

	// add a container to a pod by using pod's id
	// 容器利用编号加入到夹子
	s.Pod = podID

	// set the container's tags
	// 设定容器标签
	s.Labels = map[string]string{
		"io.podman.compose.config-hash":           "57c9635d928f88954e01491e81ca0f9049014d0d205f0fb03c951e2fe09d582a",
		"io.podman.compose.project":               "echo",
		"io.podman.compose.version":               "1.0.4",
		"com.docker.compose.project":              "echo",
		"com.docker.compose.project.working_dir":  "/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/echo",
		"com.docker.compose.project.config_files": "docker-compose.yaml",
		"com.docker.compose.container-number":     "1",
		"com.docker.compose.service":              "web",
	}

	// create a container spec
	// 创建一个容器规格
	containerCreateResponse, err := containers.CreateWithSpec(conn, s, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start the container
	// 启动容器
	if err := containers.Start(conn, containerCreateResponse.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// >>>>> mimic "curl -X POST -d "foobar" http://localhost:8080/; echo"
	// >>>>> 相对于 "curl -X POST -d "foobar" http://localhost:8080/; echo"

	// connect to container
	// 进行连线
	params := url.Values{}
	params.Add("foobar", ``)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "http://localhost:8080/;", body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// >>>>> >>>>> >>>>> testing
	// >>>>> >>>>> >>>>> 测试部份

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, "200 OK", resp.Status)

	fmt.Println("connection successful")

	// >>>>> >>>>> >>>>> tear down
	// >>>>> >>>>> >>>>> 拆除部份

	// >>>>> mimic "podman stop -t 10 echo_web_1"
	// >>>>> 相对于 "podman stop -t 10 echo_web_1"

	// set timeout to 10 seconds
	// 设置超时时间为10秒
	timeOut := uint(10)
	stopOptions := containers.StopOptions{
		Timeout: &timeOut,
	}

	// stop the container
	// 停止容器
	containers.Stop(conn, containerCreateResponse.ID, &stopOptions)
}
