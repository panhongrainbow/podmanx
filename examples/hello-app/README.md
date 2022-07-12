# GCR Hello App

## create by using podman-compose

A small ~2MB image, type

```
podman-compose up
```

then open your browser on [http://localhost:8080/](http://localhost:8080/)

## create by using podman

### establish

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman ps --filter label=io.podman.compose.project=hello-app -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'
# ef191a8214ebad4a7d3c0c6981f2437bde31ed9a951a6db0b44a6aabe1e76d3d

$ podman pod create --name=pod_hello-app --infra=false --share=
# b3a90166fb3297c5edd88e7b517f9118979c3d2675b54ed528aff4479bd936d4

$ podman network exists hello-app_default

$ podman network create --label io.podman.compose.project=hello-app --label com.docker.compose.project=hello-app hello-app_default
# /home/panhong/.config/cni/net.d/hello-app_default.conflist

$ podman create --name=hello-app_web_1 --pod=pod_hello-app --label io.podman.compose.config-hash=ef191a8214ebad4a7d3c0c6981f2437bde31ed9a951a6db0b44a6aabe1e76d3d --label io.podman.compose.project=hello-app --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=hello-app --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/hello-app --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=web --net hello-app_default --network-alias web -p 8080:8080 gcr.io/google-samples/hello-app:1.0
# a7744b16fe76ce11a2dfff870ca1f1a575b5322ac994bba3858e4f0adf2ba9f0

$ podman start -a hello-app_web_1
# 2022/07/12 15:52:05 Server listening on port 8080
```

### Tear down

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman stop -t 10 hello-app_web_1
# hello-app_web_1

$ podman rm hello-app_web_1
# a7744b16fe76ce11a2dfff870ca1f1a575b5322ac994bba3858e4f0adf2ba9f0

$ podman pod rm pod_hello-app
# b3a90166fb3297c5edd88e7b517f9118979c3d2675b54ed528aff4479bd936d4

$ podman network rm hello-app_default
# hello-app_default
```

## create by using golang



















