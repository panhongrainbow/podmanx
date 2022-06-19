package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/containers/podman/v3/pkg/specgen"
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
	_, err = images.Pull(conn, "quay.io/libpod/alpine_nginx:latest", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 创建一个容器
	// create a container
	s := specgen.NewSpecGenerator("quay.io/libpod/alpine_nginx:latest", false)
	s.Name = "foobar"
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
	_, errs := images.Remove(ctx, []string{"quay.io/libpod/alpine_nginx:latest"}, nil)
	if err != nil {
		fmt.Println(errs)
		os.Exit(1)
	}
}
