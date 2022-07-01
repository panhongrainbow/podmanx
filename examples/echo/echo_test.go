package echo

import (
	"context"
	"fmt"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestEcho(t *testing.T) {
	// create a client
	// 创建一个客户端
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	conn, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check the image exists
	// 检查镜像是否存在
	exists, err := images.Exists(conn, "k8s.gcr.io/echoserver:1.4", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("image exists ?", exists)

	// 下载镜像
	// pull an image
	if !exists {
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
}
