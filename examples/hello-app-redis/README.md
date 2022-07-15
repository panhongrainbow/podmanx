# GCR Hello App Redis

A 6-node redis cluster using [Bitnami](https://github.com/bitnami/bitnami-docker-redis-cluster)
with a [simple hit counter](https://github.com/GoogleCloudPlatform/kubernetes-engine-samples/tree/main/hello-app-redis) that persists on that redis cluster

```
podman-compose up
```

then open your browser on [http://localhost:8080/](http://localhost:8080/)

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman ps --filter label=io.podman.compose.project=hello-app-redis -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'

$ podman pod create --name=pod_hello-app-redis --infra=false --share=
# ac12095db54cb04163bc631c98680a1a4aab9756a4f08b5b1b1ff7031f8eb273

$ podman network exists hello-app-redis_default

$ podman network create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_default

# >>>>> node1

$ podman volume inspect hello-app-redis_redis-node1-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-node1-data

$ podman create --name=hello-app-redis_redis-node1_1 --pod=pod_hello-app-redis --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-node1 -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -v hello-app-redis_redis-node1-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-node1 docker.io/bitnami/redis-cluster:6.2
# 985e58415766b83420175d17167982f169c8f70204eb0642d4ff51b78e9c580a

# >>>>> node2

$ podman volume inspect hello-app-redis_redis-node2-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-node2-data

$ podman create --name=hello-app-redis_redis-node2_1 --pod=pod_hello-app-redis --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-node2 -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -v hello-app-redis_redis-node2-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-node2 docker.io/bitnami/redis-cluster:6.2
# 5b23188ae6bc5f816b7c7dc1f13be036b9d2712714d5028dd0c865a853593f62

# >>>>> node3

$ podman volume inspect hello-app-redis_redis-node3-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-node3-data

$ podman create --name=hello-app-redis_redis-node3_1 --pod=pod_hello-app-redis --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-node3 -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -v hello-app-redis_redis-node3-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-node3 docker.io/bitnami/redis-cluster:6.2
# e92a57a8a0ff887eefb77e09fc9ecef6d89de142681eb64f2073ab7c2816064f

# >>>>> node4

$ podman volume inspect hello-app-redis_redis-node4-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-node4-data

$ podman create --name=hello-app-redis_redis-node4_1 --pod=pod_hello-app-redis --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-node4 -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -v hello-app-redis_redis-node4-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-node4 docker.io/bitnami/redis-cluster:6.2
# 315f030923b9108f51878cb61c176c2f41f2bae18e0ab8c27572bc9b350a5b8b

# >>>>> node5

$ podman volume inspect hello-app-redis_redis-node5-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-node5-data

$ podman create --name=hello-app-redis_redis-node5_1 --pod=pod_hello-app-redis --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-node5 -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -v hello-app-redis_redis-node5-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-node5 docker.io/bitnami/redis-cluster:6.2
# afd1b70fe07a8ac0c8e9107620377583af38c2ff901a783484ba0712877464b9

# >>>>> hello-app

$ podman volume inspect hello-app-redis_redis-data

$ podman volume create --label io.podman.compose.project=hello-app-redis --label com.docker.compose.project=hello-app-redis hello-app-redis_redis-data

$ podman create --name=hello-app-redis_redis-cluster_1 --pod=pod_hello-app-redis --requires=hello-app-redis_redis-node5_1,hello-app-redis_redis-node4_1,hello-app-redis_redis-node2_1,hello-app-redis_redis-node3_1,hello-app-redis_redis-node1_1 --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=redis-cluster -e ALLOW_EMPTY_PASSWORD=yes -e REDIS_NODES="redis-node1 redis-node2 redis-node3 redis-node4 redis-node5 redis-cluster" -e REDIS_CLUSTER_CREATOR=yes -v hello-app-redis_redis-data:/bitnami/redis/data --net hello-app-redis_default --network-alias redis-cluster docker.io/bitnami/redis-cluster:6.2
# 77ce7432c0697ca14c3430017db6a5be569682617913189c3b00ddbcf4c9580d

$ podman create --name=hello-app-redis_web_1 --pod=pod_hello-app-redis --requires=hello-app-redis_redis-node5_1,hello-app-redis_redis-node2_1,hello-app-redis_redis-node4_1,hello-app-redis_redis-cluster_1,hello-app-redis_redis-node3_1,hello-app-redis_redis-node1_1 --label io.podman.compose.config-hash=f8290d1648ed78d029eafeb3596d02a662735e61e00379b29e94a021e588d4c2 --label io.podman.compose.project=hello-app-redis --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app-redis --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app-redis --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=web --net hello-app-redis_default --network-alias web -p 8080:8080 gcr.io/google-samples/hello-app-redis:1.0
# 1579455f835f1ff16c607707464ec87bfaebf5bf8f3d42c5a97da5cf79f8a405

$ podman start -a hello-app-redis_redis-node1_1
```

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman stop -t 10 hello-app-redis_web_1
# hello-app-redis_web_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-cluster_1
# hello-app-redis_redis-cluster_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-node5_1
# hello-app-redis_redis-node5_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-node4_1
# hello-app-redis_redis-node4_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-node3_1
# hello-app-redis_redis-node3_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-node2_1
# hello-app-redis_redis-node2_1
# exit code: 0

$ podman stop -t 10 hello-app-redis_redis-node1_1
# hello-app-redis_redis-node1_1
# exit code: 0

$ podman rm hello-app-redis_web_1
# 8ba574193bc073eb8a372540d665ad8f6af08fef1161165c70c4256b8fb35d89
# exit code: 0

$ podman rm hello-app-redis_redis-cluster_1
# ed4c98e86854b77455401d8199a1e53bc6e863005b0017fa21169610431c4e7f
# exit code: 0

$ podman rm hello-app-redis_redis-node5_1
# bf74bda77b29c3b347cfa275d82ad2e93c9f1cdb64b2c18bcbd021ebb09578b1
# exit code: 0

$ podman rm hello-app-redis_redis-node4_1
# 4e7a2f0999786bcb6f0b33ddf036f9ef993ad4d74567640584b67f28a8db502c
# exit code: 0

$ podman rm hello-app-redis_redis-node3_1
# a29de9ebd67b82e249bd2e22a21326e66f1aa386ebf5e361c2902c835d868028
# exit code: 0

$ podman rm hello-app-redis_redis-node2_1
# 64705a80c34d379dc36ed5ac35100d6fe4a367a1a3d00b64e29966fc8a63749f
# exit code: 0

$ podman rm hello-app-redis_redis-node1_1
# 18b4298a89017e202e32ac048025e3e52932f6428aa6e2b8ff190d64005ea642
# exit code: 0

$ podman pod rm pod_hello-app-redis
# 58a355e31511a1cf1dcbcff984cd34dd84f8c9bf501e9e6f609cf2c9fbe9e01e
# exit code: 0

$ podman volume rm hello-app-redis_redis-data hello-app-redis_redis-node1-data hello-app-redis_redis-node2-data hello-app-redis_redis-node3-data hello-app-redis_redis-node4-data hello-app-redis_redis-node5-data

$ podman network exists hello-app-redis_default
```
