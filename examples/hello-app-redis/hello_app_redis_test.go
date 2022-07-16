package hello_app_redis

import (
	"context"
	"fmt"
	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/images"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	//
	podmanImages = [2]string{
		"docker.io/bitnami/redis-cluster:6.2",
		"gcr.io/google-samples/hello-app-redis:1.0",
	}
)

func TestHelloRedisApp(t *testing.T) {

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

	// >>>>> mimic "podman image exists docker.io/bitnami/redis-cluster:6.2; podman image exists gcr.io/google-samples/hello-app-redis:1.0"
	// >>>>> 相对于 "podman image exists docker.io/bitnami/redis-cluster:6.2; podman image exists gcr.io/google-samples/hello-app-redis:1.0"

	for i := 0; i < len(podmanImages); i++ {
		// check the image exists
		// 检查镜像是否存在
		exists, err := images.Exists(conn, podmanImages[i], nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("image [", podmanImages[i], "] exists ?", exists)

		// pull an image if it doesn't exist
		// 如果不存在，重新下载镜像
		if !exists {

			// >>>>> mimic "podman image pull docker.io/bitnami/redis-cluster:6.2 gcr.io/google-samples/hello-app-redis:1.0"
			// >>>>> 相对于 "podman image pull docker.io/bitnami/redis-cluster:6.2 gcr.io/google-samples/hello-app-redis:1.0"

			// 下载镜像
			// pull an image
			_, err = images.Pull(conn, podmanImages[i], nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// check the image exists again
			// 再一次检查镜像是否存在
			exists, err := images.Exists(conn, podmanImages[i], nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// the image must exit
			// 镜像必须存在
			require.Equal(t, true, exists)

			fmt.Println("image [", podmanImages[i], "] exists after pull ?", exists)
		}
	}
}
