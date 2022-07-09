# Podman Support

## Installation

### Debian 11:

> - Consider to install a single package from `Debian sid` to get better `Podman` experiences.
> - Install stable Podman package first, then slightly upgrade Podman to the cutting edge version. It makes Debian more stable.

Install necessary `stable` packages by using `Apt` command

```bash
$ apt-get update

# install all stable packages
$ apt-get install podman buildah skopeo fuse-overlayfs slirp4netns golang-github-containers-buildah-dev golang-github-containers-common golang-github-containers-common-dev golang-github-containers-image golang-github-containers-image-dev golang-github-containers-libpod-dev golang-github-containers-ocicrypt-dev golang-github-containers-psgo-dev golang-github-containers-storage-dev golang-github-containernetworking-plugin-dnsname golang-github-containernetworking-plugins-dev crun libc6 podman-compose
```

Check current Podman newest version

```bash
$ cat << EOF >> /etc/apt/sources.list
# [unstable repo]
deb http://deb.debian.org/debian unstable main contrib non-free
deb-src http://deb.debian.org/debian unstable main contrib non-free
EOF

$ apt-get update

$ apt search podman engine # newest vesion is 3.4.7
# podman/unstable 3.4.7+ds1-3+b1 amd64 [upgradable from: 3.0.1+dfsg1-3+deb11u1]
#   engine to run OCI-based containers in Pods
```

Gather newest Podman package

```bash
$ mkdir -p /var/lib/apt/podman/3.4.7

$ cd /var/lib/apt/podman/3.4.7/
$ apt-get --download-only --only-upgrade --no-install-recommends -o dir::cache=`pwd` install podman buildah skopeo fuse-overlayfs slirp4netns golang-github-containers-buildah-dev golang-github-containers-common golang-github-containers-common-dev golang-github-containers-image golang-github-containers-image-dev golang-github-containers-libpod-dev golang-github-containers-ocicrypt-dev golang-github-containers-psgo-dev golang-github-containers-storage-dev golang-github-containernetworking-plugin-dnsname golang-github-containernetworking-plugins-dev podman-compose crun libc6 podman-compose

$ cd /var/lib/apt/podman/3.4.7/
$ dpkg-scanpackages --multiversion . /dev/null | gzip -9c > Packages.gz
```

Close unstable repository and add specific downloaded repository

```bash
$ vim /etc/apt/sources.list
# correct /etc/apt/sources.list in the following

# [unstable repo]
# deb http://deb.debian.org/debian/ unstable main
# deb-src http://deb.debian.org/debian/ unstable main

# Podman 3.4.7
deb [trusted=yes] file:/var/lib/apt/podman/3.4.7 ./ # add here !
```

Upgrade Podman

```bash
$ apt-get update

$ apt-get --only-upgrade --no-install-recommends install podman buildah skopeo fuse-overlayfs slirp4netns golang-github-containers-buildah-dev golang-github-containers-common golang-github-containers-common-dev golang-github-containers-image golang-github-containers-image-dev golang-github-containers-libpod-dev golang-github-containers-ocicrypt-dev golang-github-containers-psgo-dev golang-github-containers-storage-dev golang-github-containernetworking-plugin-dnsname golang-github-containernetworking-plugins-dev podman-compose crun libc6 podman-compose
```

Check Podman

```bash
# non rootless mode
$ sudo podman --log-level=debug info
# No error messages
```

Check Podman-compose

```bash
# non rootless mode
$ sudo podman-compose --version
['podman', '--version', '']
using podman version: 3.4.7
podman-composer version  1.0.3
podman --version 
podman version 3.4.7
exit code: 0
```

## Rootless mode

> Refer to [Basic Setup and Use of Podman in a Rootless environment](https://github.com/containers/podman/blob/main/docs/tutorials/rootless_tutorial.md). set up the `rootless mode` step by step.

### Debian 11:

Make rootless the config files

```bash
# config files for root in the following

$ mv /etc/containers/containers.conf /etc/containers/containers.conf.old
$ mv /etc/containers/libpod.conf /etc/containers/libpod.conf.old
$ mv /etc/containers/storage.conf /etc/containers/storage.conf.old

$ cat << EOF > /etc/containers/containers.conf 
[containers]
default_capabilities = [
  "CHOWN",
  "DAC_OVERRIDE",
  "FOWNER",
  "FSETID",
  "KILL",
  "NET_BIND_SERVICE",
  "SETFCAP",
  "SETGID",
  "SETPCAP",
  "SETUID",
  "SYS_CHROOT"
]
default_sysctls = [
  "net.ipv4.ping_group_range=0 0",
]
rootless_networking = "slirp4netns"
[secrets]
[secrets.opts]
[network]
[engine]
network_cmd_path = "/usr/bin/slirp4netns"
[engine.runtimes]
[engine.volume_plugins]
[machine]
EOF

$ cat << EOF > /etc/containers/libpod.conf 
image_default_transport = "docker://"
conmon_path = [
    "/usr/bin/conmon",
    "/usr/sbin/conmon",
    "/usr/libexec/podman/conmon",
    "/usr/local/libexec/crio/conmon",
    "/usr/lib/podman/bin/conmon",
    "/usr/libexec/crio/conmon",
    "/usr/lib/crio/bin/conmon"
]
conmon_env_vars = [
    "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
]
cgroup_manager = "systemd"
tmp_dir = "/var/run/libpod"
max_log_size = -1
no_pivot_root = false
cni_config_dir = "/etc/cni/net.d/"
cni_plugin_dir = [
    "/usr/lib/cni",
    "/usr/local/lib/cni",
    "/opt/cni/bin"
]
cni_default_network = "podman"
infra_image = "k8s.gcr.io/pause:3.1"
infra_command = "/pause"
lock_type = "shm"
num_locks = 2048
runtime = "crun"
runtime_supports_json = ["crun", "runc"]
runtime_supports_nocgroups = ["crun"]
[runtimes]
runc = [
    "/usr/sbin/runc",
]
crun = [
    "/usr/bin/crun"
]
EOF

$ cat << EOF > /etc/containers/storage.conf
[storage]
driver = "overlay"
runroot = "/run/containers/storage"
graphroot = "/var/lib/containers/storage"
[storage.options]
additionalimagestores = [
]
[storage.options.overlay]
mount_program = "/usr/bin/fuse-overlayfs"
mountopt = "nodev"
[storage.options.thinpool]
EOF

# config files for user in the following

$ mv ~/.config/containers/storage.conf ~/.config/containers/storage.conf.old

$ cat << EOF > ~/.config/containers/storage.conf
[storage]
driver = "overlay"
runroot = "/run/containers/storage"
graphroot = "/var/lib/containers/storage"
rootless_storage_path = "$HOME/.local/share/containers/storage"
[storage.options]
additionalimagestores = [
]
[storage.options.overlay]
mount_program = "/usr/bin/fuse-overlayfs"
mountopt = "nodev"
[storage.options.thinpool]
EOF
```

#### Adjust the permission of folders

```bash
$ mkdir ~/.podman

$ cat << EOF > ~/.podman/podman.sh
#!/bin/bash
PURPLE="\[$(tput setaf 12)\]"
RESET="\[$(tput sgr0)\]"
PS1="${PURPLE}Podman User >${RESET} "

# adjust the permission of folders
sudo chown -R "\$USER"."\$USER" /var/lib/containers
sudo chown -R "\$USER"."\$USER" /run/containers
sudo chown -R "\$USER"."\$USER" /run/libpod

# firewall forward
sudo iptables -P FORWARD ACCEPT
EOF

$ . ~/.podman/podman.sh
```

#### Check Podman information

Get Podman information

```bash
$ podman --log-level=debug info

# rootless mode configurations are in following !
INFO[0000] podman filtering at log level debug          
DEBU[0000] Called info.PersistentPreRunE(podman --log-level=debug info) 
DEBU[0000] Merged system config "/usr/share/containers/containers.conf" 
DEBU[0000] Merged system config "/etc/containers/containers.conf" 
DEBU[0000] Using conmon: "/usr/bin/conmon"              
DEBU[0000] Initializing boltdb state at /var/lib/containers/storage/libpod/bolt_state.db 
DEBU[0000] Overriding tmp dir "/run/user/1001/libpod/tmp" with "/run/libpod" from database 
DEBU[0000] Using graph driver overlay                   
DEBU[0000] Using graph root /var/lib/containers/storage 
DEBU[0000] Using run root /run/containers/storage       
DEBU[0000] Using static dir /var/lib/containers/storage/libpod 
DEBU[0000] Using tmp dir /run/libpod                    
DEBU[0000] Using volume path /var/lib/containers/storage/volumes 
DEBU[0000] Set libpod namespace to ""                   
DEBU[0000] Not configuring container store              
DEBU[0000] Initializing event backend journald          
DEBU[0000] configured OCI runtime kata initialization failed: no valid executable found for OCI runtime kata: invalid argument 
DEBU[0000] configured OCI runtime runsc initialization failed: no valid executable found for OCI runtime runsc: invalid argument 
DEBU[0000] Using OCI runtime "/usr/bin/crun"            
INFO[0000] Found CNI network jaeger-go-example_default (type=bridge) at /home/panhong/.config/cni/net.d/jaeger-go-example_default.conflist 
DEBU[0000] Default CNI network name podman is unchangeable 
INFO[0000] Setting parallel job count to 25             
INFO[0000] podman filtering at log level debug          
DEBU[0000] Called info.PersistentPreRunE(podman --log-level=debug info) 
DEBU[0000] Merged system config "/usr/share/containers/containers.conf" 
DEBU[0000] Merged system config "/etc/containers/containers.conf" 
DEBU[0000] Using conmon: "/usr/bin/conmon"              
DEBU[0000] Initializing boltdb state at /var/lib/containers/storage/libpod/bolt_state.db 
DEBU[0000] Overriding tmp dir "/run/user/1001/libpod/tmp" with "/run/libpod" from database 
DEBU[0000] Using graph driver overlay                   
DEBU[0000] Using graph root /var/lib/containers/storage 
DEBU[0000] Using run root /run/containers/storage       
DEBU[0000] Using static dir /var/lib/containers/storage/libpod 
DEBU[0000] Using tmp dir /run/libpod                    
DEBU[0000] Using volume path /var/lib/containers/storage/volumes 
DEBU[0000] Set libpod namespace to ""                   
DEBU[0000] [graphdriver] trying provided driver "overlay" 
DEBU[0000] overlay: mount_program=/usr/bin/fuse-overlayfs 
DEBU[0000] backingFs=extfs, projectQuotaSupported=false, useNativeDiff=false, usingMetacopy=false 
DEBU[0000] Initializing event backend journald          
DEBU[0000] configured OCI runtime runsc initialization failed: no valid executable found for OCI runtime runsc: invalid argument 
DEBU[0000] configured OCI runtime kata initialization failed: no valid executable found for OCI runtime kata: invalid argument 
DEBU[0000] Using OCI runtime "/usr/bin/crun"            
INFO[0000] Found CNI network jaeger-go-example_default (type=bridge) at /home/panhong/.config/cni/net.d/jaeger-go-example_default.conflist 
DEBU[0000] Default CNI network name podman is unchangeable 
INFO[0000] Setting parallel job count to 25             
DEBU[0000] Loading registries configuration "/etc/containers/registries.conf" 
DEBU[0000] Loading registries configuration "/etc/containers/registries.conf.d/shortnames.conf" 

(...)

DEBU[0000] Called info.PersistentPostRunE(podman --log-level=debug info) 
```

Get plugin information

```yaml
host:
  arch: amd64
  buildahVersion: 1.23.1
  cgroupControllers:
  - cpuset
  - cpu
  - io
  - memory
  - pids
  cgroupManager: systemd
  cgroupVersion: v2
  conmon:
    package: 'conmon: /usr/bin/conmon'
    path: /usr/bin/conmon
    version: 'conmon version 2.0.25, commit: unknown'
  cpus: 8
  distribution:
    codename: bullseye
    distribution: debian
    version: "11"
  eventLogger: journald
  hostname: debian5
  idMappings:
    gidmap:
    - container_id: 0
      host_id: 1001
      size: 1
    - container_id: 1
      host_id: 165536
      size: 65536
    uidmap:
    - container_id: 0
      host_id: 1001
      size: 1
    - container_id: 1
      host_id: 165536
      size: 65536
  kernel: 5.10.0-15-amd64
  linkmode: dynamic
  logDriver: journald
  memFree: 15494119424
  memTotal: 33419116544
  ociRuntime:
    name: crun
    package: 'crun: /usr/bin/crun'
    path: /usr/bin/crun
    version: |-
      crun version 0.17
      commit: 0e9229ae34caaebcb86f1fde18de3acaf18c6d9a
      spec: 1.0.0
      +SYSTEMD +SELINUX +APPARMOR +CAP +SECCOMP +EBPF +YAJL
  os: linux
  remoteSocket:
    exists: true
    path: /run/user/1001/podman/podman.sock
  security:
    apparmorEnabled: false
    capabilities: CAP_CHOWN,CAP_DAC_OVERRIDE,CAP_FOWNER,CAP_FSETID,CAP_KILL,CAP_NET_BIND_SERVICE,CAP_SETFCAP,CAP_SETGID,CAP_SETPCAP,CAP_SETUID,CAP_SYS_CHROOT
    rootless: true
    seccompEnabled: true
    seccompProfilePath: /usr/share/containers/seccomp.json
    selinuxEnabled: false
  serviceIsRemote: false
  slirp4netns:
    executable: /usr/bin/slirp4netns
    package: 'slirp4netns: /usr/bin/slirp4netns'
    version: |-
      slirp4netns version 1.0.1
      commit: 6a7b16babc95b6a3056b33fb45b74a6f62262dd4
      libslirp: 4.4.0
  swapFree: 9999216640
  swapTotal: 9999216640
  uptime: 3h 19m 41.57s (Approximately 0.12 days)
plugins:
  log:
  - k8s-file
  - none
  - journald
  network:
  - bridge
  - macvlan
  volume:
  - local
registries: {}
store:
  configFile: /home/panhong/.config/containers/storage.conf
  containerStore:
    number: 5
    paused: 0
    running: 3
    stopped: 2
  graphDriverName: overlay
  graphOptions:
    overlay.mount_program:
      Executable: /usr/bin/fuse-overlayfs
      Package: 'fuse-overlayfs: /usr/bin/fuse-overlayfs'
      Version: |-
        fusermount3 version: 3.10.3
        fuse-overlayfs: version 1.8.2
        FUSE library version 3.10.3
        using FUSE kernel interface version 7.31
    overlay.mountopt: nodev
  graphRoot: /var/lib/containers/storage
  graphStatus:
    Backing Filesystem: extfs
    Native Overlay Diff: "false"
    Supports d_type: "true"
    Using metacopy: "false"
  imageStore:
    number: 22
  runRoot: /run/containers/storage
  volumePath: /var/lib/containers/storage/volumes
version:
  APIVersion: 3.4.7
  Built: 0
  BuiltTime: Thu Jan  1 08:00:00 1970
  GitCommit: ""
  GoVersion: go1.18.1
  OsArch: linux/amd64
  Version: 3.4.7
```

# Podman Debug

## Could not upgrade podman buildah and skopeo packages ?

Solution:

```bash
# install libc6 package
$ apt-get install -y libc6_2.33-7_amd64.deb libc6_2.33-7_i386.deb
```

## Podman-compose could not mount container network ?

Solution:

```bash
# upgrage
$ apt-get install -y containernetworking-plugins_1.1.0+ds1-1+b1_amd64.deb
```

# Podman Tools

Podman desktop graphical interface tools

- [podman-desktop-companion](https://iongion.github.io/podman-desktop-companion/)
