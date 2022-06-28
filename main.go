package main

import (
	"context"
	"fmt"
	nettypes "github.com/containers/podman/v3/libpod/network/types"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/bindings/pods"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/specgen"
	"os"
)

const (
	// imageUrl 为去拉取的镜像地址
	// imageUrl is the URL of the image to be pulled.
	imageUrl = "docker.io/library/httpd:latest"

	// 测试用的命令
	// test command
	// podman run --detach --name some-mariadb --env MARIADB_USER=xiaomi --env MARIADB_PASSWORD=12345 --env MARIADB_ROOT_PASSWORD=12345  docker.io/library/mariadb:latest
	// imageUrl = "docker.io/library/mariadb:latest"
)

// 先测试再封装
// test first and encapsulate later
func main() {
	// 创建一个客户端
	// create a client
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	conn, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 下载镜像
	// pull an image
	_, err = images.Pull(conn, imageUrl, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	portMapping := make([]nettypes.PortMapping, 1, 1)
	portMapping[0].HostPort = 8080
	portMapping[0].ContainerPort = 80 // or 3306

	// >>>>> >>>>> >>>>> 建立夹子
	// create a pod
	pspec := entities.PodSpec{
		PodSpecGen: specgen.PodSpecGenerator{
			PodBasicConfig: specgen.PodBasicConfig{
				Name: "some-pod",
			},
			PodNetworkConfig: specgen.PodNetworkConfig{
				// StaticIP:     &net.IP{127, 0, 0, 1},
				PortMappings: portMapping,
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

	// >>>>> >>>>> >>>>> 挂载目录
	// mount a directory
	/*mountFlag := define.TypeBind
	volumeMounts, _, _, _ := specgen.GenVolumeMounts([]string{"/home/panhong/web01/:/usr/local/apache2/htdocs/:ro"})
	unifiedMounts := make(map[string]spec.Mount)
	for dst, mount := range volumeMounts {
		fmt.Println(mount.Type == mountFlag)
		mount.Type = mountFlag
		unifiedMounts[dst] = mount
	}

	finalMounts := make([]spec.Mount, 0, 1)
	for _, mount := range unifiedMounts {
		if mount.Type == define.TypeBind {
			absSrc, _ := filepath.Abs(mount.Source)
			mount.Source = absSrc
		}
		finalMounts = append(finalMounts, mount)
	}*/

	// >>>>> >>>>> >>>>> 挂载卷
	// mount a volume
	/*_, volumeVolumes, _, _ := specgen.GenVolumeMounts([]string{"web01:/usr/local/apache2/htdocs/"})
	unifiedVolumes := make(map[string]*specgen.NamedVolume)
	for dest, volume := range volumeVolumes {
		unifiedVolumes[dest] = volume
	}

	finalVolumes := make([]*specgen.NamedVolume, 0, 1)
	for _, volume := range unifiedVolumes {
		finalVolumes = append(finalVolumes, volume)
	}*/

	// 设置 port mapping
	// make port mapping
	// 使用 pod 时，关闭网路设置
	/*
		portMapping := make([]nettypes.PortMapping, 1, 1)
		portMapping[0].HostPort = 8080
		portMapping[0].ContainerPort = 80 // or 3306
	*/

	// 创建一个容器
	// create a container
	s := specgen.NewSpecGenerator(imageUrl, false)
	s.Name = "web01"
	// s.Mounts = finalMounts
	// s.Volumes = finalVolumes
	// s.PortMappings = portMapping // 使用 pod 时，关闭网路设置

	// 容器加入 pod
	// add a container to a pod
	s.Pod = preport.Id

	// 设定帐号密码
	// set account and password
	s.Env = map[string]string{
		"MARIADB_USER":          "xiaomi",
		"MARIADB_PASSWORD":      "12345",
		"MARIADB_ROOT_PASSWORD": "12345",
	}

	// 创建容器的 spec
	// create a container's spec
	createResponse, err := containers.CreateWithSpec(conn, s, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 分离容器的命令
	// detach container command
	detachKeys := "ctrl-a"
	startOptions := containers.StartOptions{DetachKeys: &detachKeys}

	// 启动容器
	// start the container
	if err := containers.Start(conn, createResponse.ID, &startOptions); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 先中断
	// interrupt
	return

	// 停止容器
	// stop the container
	if err := containers.Stop(conn, createResponse.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 删除容器
	// remove the container
	if err := containers.Remove(conn, createResponse.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 删除镜像
	// reomve an image
	ctx := context.Background()
	_, errs := images.Remove(ctx, []string{imageUrl}, nil)
	if err != nil {
		fmt.Println(errs)
		os.Exit(1)
	}

	// 创建一个卷
	// create a volume
	// volumeName := "web01"
	/*ctx := context.Background()
	volumeOptions := entities.VolumeCreateOptions{
		Name: volumeName,
	}
	_, err = volumes.Create(conn, volumeOptions, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/

	// 删除卷
	// remove a volume
	/*err = volumes.Remove(conn, volumeName, nil)
	if err != nil {
		fmt.Println(errs)
		os.Exit(1)
	}*/
}
