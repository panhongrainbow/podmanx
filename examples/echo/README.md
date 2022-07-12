# Echo Service example

## create by using podman-compose

```
podman-compose up
```

Test the service with `curl like this`

```
$ curl -X POST -d "foobar" http://localhost:8080/; echo

CLIENT VALUES:
client_address=10.89.31.2
command=POST
real path=/
query=nil
request_version=1.1
request_uri=http://localhost:8080/

SERVER VALUES:
server_version=nginx: 1.10.0 - lua: 10001

HEADERS RECEIVED:
accept=*/*
content-length=6
content-type=application/x-www-form-urlencoded
host=localhost:8080
user-agent=curl/7.76.1
BODY:
foobar
```

## create by using podman

### establish

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman ps --filter label=io.podman.compose.project=echo -a --format '{{ index .Labels "io.podman.compose.config-hash"}}'
# 57c9635d928f88954e01491e81ca0f9049014d0d205f0fb03c951e2fe09d582a

$ podman pod create --name=pod_echo --infra=false --share=
# a5a083c2c43ada30fcc97dd2876055cb5abcd7b89bbd3e1a321a0643903c2973

$ podman network exists echo_default
# /home/panhong/.config/cni/net.d/echo_default.conflist

$ podman network create --label io.podman.compose.project=echo --label com.docker.compose.project=echo echo_default
# /home/panhong/.config/cni/net.d/echo_default.conflist

$ podman create --name=echo_web_1 --pod=pod_echo --label io.podman.compose.config-hash=57c9635d928f88954e01491e81ca0f9049014d0d205f0fb03c951e2fe09d582a --label io.podman.compose.project=echo --label io.podman.compose.version=1.0.4 --label com.docker.compose.project=echo --label com.docker.compose.project.working_dir=/home/panhong/go/src/github.com/panhongrainbow/podmanx/examples/echo --label com.docker.compose.project.config_files=docker-compose.yaml --label com.docker.compose.container-number=1 --label com.docker.compose.service=web --net echo_default --network-alias web -p 8080:8080 k8s.gcr.io/echoserver:1.4
# 47bac49dd23fb5a0ba636df7bec36f68db1c24495f214b15621d997ac4d7fde9

$ podman start -a echo_web_1
```

### tear down

```bash
$ podman-compose --version
# podman-compose version: 1.0.4

$ podman stop -t 10 echo_web_1
# echo_web_1

$ podman rm echo_web_1
# 47bac49dd23fb5a0ba636df7bec36f68db1c24495f214b15621d997ac4d7fde9

$ podman pod rm pod_echo
# 6c88ae1e4428117a359d6c50e26c2eba8092ab1810082cd3d84186bdb26f65fe

$ podman network rm echo_default
# echo_default
```

## create by using golang

```bash
$ cd podmanx/examples/echo

$ go test -v
# === RUN   TestEcho
# image exists ? true
# plan exists ? false
# pod exists after creating ? true
# pod:  pod_echo
# network exists ? false
# network path  /home/panhong/.config/cni/net.d/echo_default.conflist
# network exists after creating ? true
# INFO[0000] Going to start container "5e762c8b424624a5607fafe137980593d74372a1b96fc65fdbcb55f78efc93de" 
# connection successful
# remove container successful
# remove pods successful
# remove network successful
# --- PASS: TestEcho (2.76s)
# PASS
# ok      podmanx/examples/echo   2.771s
```

