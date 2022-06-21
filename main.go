package main

import (
	"context"
	"fmt"
	"github.com/containers/podman/v3/libpod/define"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"path/filepath"
)

const (
	// imageUrl 为去拉取的镜像地址
	// imageUrl is the URL of the image to be pulled.
	imageUrl = "docker.io/library/httpd:latest"
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

	// 挂载目录
	// mount a directory
	mountFlag := define.TypeBind
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
	}

	// 创建一个容器
	// create a container
	s := specgen.NewSpecGenerator(imageUrl, false)
	s.Name = "web01"
	s.Mounts = finalMounts
	createResponse, err := containers.CreateWithSpec(conn, s, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 启动容器
	// start the container
	if err := containers.Start(conn, createResponse.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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

}
