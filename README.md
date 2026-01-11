# tiny-docker
> `tiny-docker` is a **toy container engine** that mocks Docker, built from scratch to replicate the core functionality of Docker in a lightweight form factor. 
>
> It is built on Go 1.23.4, allowing us to leverage more modern language features that simplify this lightweight project.


https://github.com/user-attachments/assets/180a38ca-d723-49d7-bd02-bd6fb5052f72


# description


```mermaid
---
config:
  kanban:
    ticketBaseUrl: 'https://github.com/mermaid-js/mermaid/issues/#TICKET#'
---
kanban
  Todo
    [Manage network for containers]
  [Not Support]
    [Pull/Commit an actual image file from docker-hub]
    [containerd/runc/containerd-shim]
    [cgroup setting for exec]
```

# Usage

## 1. Prepare Image

We have to prepare a runnable image because `tiny-docker` does not support pulling images from Docker Hub yet. Luckily, `docker export` is a convenient option for creating a image.

```bash
docker run -itd ubuntu /bin/bash # pull image from docker hub
# bdd68ffc0bee1f72367474f95e93a98a8fcd318033ec6982e696e26623c25a17
docker export -o /tmp/linux.tar bdd68ffc0bee1f72367474f95e93a98a8fcd318033ec6982e696e26623c25a17
# now we have a linux.tar that can be used as our image

docker run -d busybox top -b
# 553accaf523abbedd50195bd85173b714e10bb31c5e2995d4e28b7e28d5b88b8
docker export -o /tmp/busybox.tar 553accaf523abbedd50195bd85173b714e10bb31c5e2995d4e28b7e28d5b88b8
```

## 2. Running mini-dockerd

`mini-dockerd` is responsible for managing, deploying containers, and many other operations. So we should start it first.

```bash
git clone https://github.com/0x822a5b87/tiny-docker.git /tmp/tiny-docker && cd /tmp/tiny-docker

cd src

go build -o mini-docker . && ./restart_daemon.sh
```

![mini-docker-daemon-start-logs](./resources/mini-docker-daemon-start-logs.png)

## 3. Play with mini-docker

### Start Containers

Start a Linux container from `/tmp/linux.tar`, which we exported a few minutes ago.

```bash
./mini-docker run -d /tmp/linux.tar -- /bin/sh -c "while true; do sleep 1; done"
```

Start another Linux container with `environment`, `cgroup limit`.

```bash
./mini-docker run -d \
	-e PATH=/bin/ \
	-e name=mini-docker-linux \
	-m 2000m \
	-c '10000 100000' \
	/tmp/linux.tar \
	-- /bin/sh -c "while true; do sleep 1; done"
```

Start a busybox container from `/tmp/busybox.tar`.

```bash
./mini-docker run -d \
	-e PATH=/bin/ 
	-e name=mini-docker-busybox \
	-m 128m \
	-c '100000 100000' \
	/tmp/busybox.tar \ 
	-- /bin/ash -c "while true; do sleep 1; done"
```

### Manage containers

#### ps

```bash
./mini-docker ps
```

This command will show all running containers:

```
CONTAINER ID                      IMAGE    COMMAND                                       CREATED        STATUS            NAMES
26992886d8d94ab6bc3b5a9668afd46f  linux    "/bin/sh -c while true; do sleep 1; done"     8 minutes ago  Up 8 minutes ago  linux
9195b42c32f54b39a682b3294d188e12  busybox  "/bin/ash -c while true; do sleep 1; done"    6 minutes ago  Up 6 minutes ago  busybox
adc1dd03f37d4c8ba003b356e168d048  linux    "-- /bin/sh -c while true; do sleep 1; done"  2 minutes ago  Up 2 minutes ago  linux
```

#### stop

```bash
./mini-docker stop adc1dd03f37d4c8ba003b356e168d048

./mini-docker ps
#CONTAINER ID                      IMAGE    COMMAND                                     CREATED         STATUS             NAMES
#9195b42c32f54b39a682b3294d188e12  busybox  "/bin/ash -c while true; do sleep 1; done"  9 minutes ago   Up 9 minutes ago   busybox
#26992886d8d94ab6bc3b5a9668afd46f  linux    "/bin/sh -c while true; do sleep 1; done"   11 minutes ago  Up 11 minutes ago  linux

./mini-docker ps -a
#CONTAINER ID                      IMAGE    COMMAND                                       CREATED         STATUS                        NAMES
#26992886d8d94ab6bc3b5a9668afd46f  linux    "/bin/sh -c while true; do sleep 1; done"     11 minutes ago  Up 11 minutes ago             linux
#9195b42c32f54b39a682b3294d188e12  busybox  "/bin/ash -c while true; do sleep 1; done"    10 minutes ago  Up 10 minutes ago             busybox
#adc1dd03f37d4c8ba003b356e168d048  linux    "-- /bin/sh -c while true; do sleep 1; done"  5 minutes ago   Exited (0) a few seconds ago  linux
```

#### exec

```bash
./mini-docker exec -it 26992886d8d94ab6bc3b5a9668afd46f /bin/bash
```

This command allow us to enter a running container by `nsenter`.

```
root@VM-0-10-opencloudos:/# ls
bin  boot  dev	etc  home  lib	lib64  media  mnt  opt	proc  root  run  sbin  srv  sys  tmp  usr  var
```

#### others

Additionally, some core features of Docker are also supported:

-   `commit`
-   `logs`

# structure

## cgroup

```mermaid
flowchart TB

linux-cgroup("linux-cgroup"):::pink

linux-cgroup --> cgroup("cgroup"):::purple
linux-cgroup --> memory("memory"):::purple
linux-cgroup --> cpu("cpu"):::purple
linux-cgroup --> cpuset("cpuset"):::purple
linux-cgroup --> io("io"):::purple
linux-cgroup --> pids("pids"):::purple
linux-cgroup --> others("..."):::purple

cgroup --> procs("cgroup.procs"):::green
cgroup --> subtree_control("cgroup.subtree_control"):::green

memory --> memory.max("memory.max"):::green
memory --> memory.swap.max("memory.swap.max"):::green

cpu --> cpu.max("cpu.max"):::green


cpuset --> cpuset.cpus("cpuset.cpus"):::green

io --> io.max("io.max"):::green
io --> io.latency("io.latency"):::green

pids --> pids.max("pids.max"):::green
pids --> pids.current("pids.current"):::green


classDef pink   fill:#FFCCCC,stroke:#333,ont-weight: bold;
classDef green  fill:#696,color: #fff,font-weight: bold;
classDef purple fill:#969,stroke:#333;
classDef dotted fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5

```

## cgroup abstraction

```mermaid
classDiagram
    %% 严格按要求定义样式

    %% 核心抽象层（接口/泛型）
    class Item:::typeAlias {
        <<type alias>>
        any
    }
    class Value:::interface {
        <<interface>>
    }
    class BaseSubsystem:::interface {
        <<interface>>
    }
    class Subsystem:::genericInterface {
        <<generic interface>>
        BaseSubsystem
    }

    %% 基础设施层（具体实现）
    class CgroupFileSystem:::concrete {
        Path string
        AutoCreate bool
    }
    class CgroupManager:::concrete {
        fs *CgroupFileSystem
        procsSubsystem *ProcsValueSubsystem
        cpuMaxSubsystem *cpu.MaxValueSubsystem
        memoryMaxSubsystem *memory.MaxValueSubsystem
    }

    %% 外部依赖子系统（CPU/Procs/Memory）
    class cgroup.ProcsValueSubsystem:::external {
        <<external>>
        Procs 子系统
    }
    class cpu.MaxValueSubsystem:::external {
        <<external>>
        CPU 子系统
    }
    class memory.MaxValueSubsystem:::external {
        <<external>>
        Memory 子系统
    }

    %% 接口继承/实现关系
    BaseSubsystem <|-- Subsystem : 继承
    BaseSubsystem <|-- cgroup.ProcsValueSubsystem : 实现
    BaseSubsystem <|-- cpu.MaxValueSubsystem : 实现
    BaseSubsystem <|-- memory.MaxValueSubsystem : 实现
    Subsystem <|-- cgroup.ProcsValueSubsystem : 实现
    Subsystem <|-- cpu.MaxValueSubsystem : 实现
    Subsystem <|-- memory.MaxValueSubsystem : 实现
    Value <|-- cgroup.ProcsValue : 实现（隐含）
    Value <|-- cpu.MaxValue : 实现（隐含）
    Value <|-- memory.MaxValue : 实现（隐含）

    %% 组合/依赖关系
    CgroupManager o-- CgroupFileSystem : 包含
    CgroupManager o-- cgroup.ProcsValueSubsystem : 包含
    CgroupManager o-- cpu.MaxValueSubsystem : 包含
    CgroupManager o-- memory.MaxValueSubsystem : 包含
    CgroupFileSystem <-- newSubsystem : 读写依赖
    Subsystem <-- newSubsystem : 实例化依赖

    classDef typeAlias fill:#f0f8ff,stroke:#2196f3,stroke-width:1px,rounded:8px,font-style:italic;
    classDef interface fill:#fef7fb,stroke:#9c27b0,stroke-width:1.5px,rounded:8px,font-weight:600;
    classDef genericInterface fill:#e8f5e8,stroke:#4caf50,stroke-width:1.5px,rounded:8px,font-weight:600;
    classDef concrete fill:#fff8e1,stroke:#ff9800,stroke-width:1px,rounded:8px;
    classDef external fill:#f5f5f5,stroke:#607d8b,stroke-width:1px,rounded:8px,dashed:true;

```

## UnionFS

```mermaid
flowchart BT
    %% 样式定义：区分不同层级和操作
    classDef writeLayer fill:#FFE0B2,stroke:#E65100,stroke-width:2px,rounded:8px,font-weight:600;
    classDef readLayer fill:#E1F5FE,stroke:#0288D1,stroke-width:2px,rounded:8px,font-weight:600;
    classDef operation fill:#F3E5F5,stroke:#7B1FA2,stroke-width:2px,rounded:8px,font-weight:600;
    classDef step fill:#E8F5E8,stroke:#2E7D32,stroke-width:1px,rounded:6px;
    classDef arrow stroke:#555555,stroke-width:1.5px;

    %% 1. 分层结构
    subgraph LayerStructure["1. 分层结构"]
        direction BT
        WLayer(可写层 / Write Layer):::writeLayer
        RLayer(只读层 / Read Layer):::readLayer
        RLayer --> WLayer
    end
    class LayerStructure operation

    %% 2. 读操作
    subgraph ReadOperation["2. 读操作 (自上而下)"]
        direction LR
        RStart(读请求):::step
        RCheckWrite{检查可写层?}:::step
        RHitWrite(找到数据<br/>直接返回):::step
        RMissWrite(未找到):::step
        RCheckRead{检查只读层?}:::step
        RHitRead(找到数据<br/>直接返回):::step
        RMissRead(未找到<br/>返回不存在):::step

        RStart --> RCheckWrite
        RCheckWrite -->|"是"| RHitWrite
        RCheckWrite -->|"否"| RMissWrite
        RMissWrite --> RCheckRead
        RCheckRead -->|"是"| RHitRead
        RCheckRead -->|"否"| RMissRead
    end
    class ReadOperation operation

    %% 3. 写操作 (Copy-on-Write)
    subgraph WriteOperation["3. 写操作 (Copy-on-Write)"]
        direction LR
        WStart(写请求):::step
        WCheckWrite{数据在可写层?}:::step
        WUpdate(直接更新<br/>可写层数据):::step
        WMissWrite(数据不存在):::step
        WCheckRead{数据在只读层?}:::step
        WCopy(复制数据到<br/>可写层):::step
        WUpdateNew(在可写层<br/>创建新数据):::step

        WStart --> WCheckWrite
        WCheckWrite -->|"是"| WUpdate
        WCheckWrite -->|"否"| WMissWrite
        WMissWrite --> WCheckRead
        WCheckRead -->|"是"| WCopy
        WCheckRead -->|"否"| WUpdateNew
        WCopy --> WUpdate
    end
    class WriteOperation operation

    %% 整体布局连接
    LayerStructure -.->|为以下操作提供基础| ReadOperation:::arrow
    LayerStructure -.->|为以下操作提供基础| WriteOperation:::arrow
```

## mini-dockerd

```mermaid
---
title: mini-dockerd
---

flowchart TB
mini-docker("mini-docker"):::animate
dockerd("mini-dockerd"):::green
state("状态存储<br/>(内存+文件)"):::pale_pink
pid_sock("标识文件<br/>(dockerd.pid + dockerd.sock)"):::purple

ubuntu1("container-1(ubuntu)"):::yellow
ubuntu2("container-2(ubuntu)"):::yellow
busybox("container-3(busybox)"):::yellow
nginx("container-4(nginx)"):::yellow
other("..."):::yellow

mini-docker -->|1.读取标识文件| pid_sock
mini-docker <-.->|2.UDS通信| dockerd
dockerd <-->|读写状态| state
dockerd -->|"fork()/管控"| ubuntu1
dockerd -->|"fork()/管控"| ubuntu2
dockerd -->|"fork()/管控"| busybox
dockerd -->|"fork()/管控"| nginx
dockerd -->|"fork()/管控"| other
dockerd -->|写入PID/UDS地址| pid_sock

classDef pink 1,fill:#FFCCCC,stroke:#333, color: #fff, font-weight:bold;
classDef pale_pink fill:#E1BEE7,color:#000000;
classDef green fill: #696,color: #fff,font-weight: bold;
classDef purple fill:#969,stroke:#333, font-weight: bold;
classDef error fill:#bbf,stroke:#f66,stroke-width:2px,color:#fff,stroke-dasharray: 5 5
classDef coral fill:#f9f,stroke:#333,stroke-width:4px;
classDef animate stroke-dasharray: 9,5,stroke-dashoffset: 900,animation: dash 25s linear infinite;
classDef yellow fill:#FFF9C4,color:#000000;
```

## exec

```mermaid
flowchart TB
    docker("mini-docker exec <br/>&lt;容器ID&gt; &lt;命令&gt;"):::green

    subgraph proc["/proc/&lt;容器PID&gt;/"]
        direction TB
        ns["ns/ <br/>(mnt/pid/net/uts/ipc)"]:::pale_pink          
        environ("environ <br/>(环境变量)"):::pale_pink
        cwd("cwd <br/>(工作目录)"):::pale_pink
        cgroup("cgroup <br/>(可选)"):::pale_pink
    end

    step1("1.解析容器ID → 获取容器PID"):::green
    step2("2.读取上下文 <br/>(cwd/environ)"):::green
    step3("3.构造nsenter命令 <br/>nsenter <br/> -t &lt;PID&gt; -a <br/>/bin/bash <br/>-w &lt;cwd&gt;<br/> -e &lt;env> <br/>&lt;命令&gt;"):::yellow
    exec("4.syscall.Exec() <br/>执行nsenter，传递环境变量"):::yellow

    docker --> step1
    step1 <-.->|"读取PID"| proc
    step1 --> step2
    step2 <-.->|"读取上下文"| environ & cwd
    step2 --> step3
    step3 --> exec
    exec <-.->|"依赖命名空间文件"| ns

    classDef pale_pink fill:#E1BEE7,color:#000000;
    classDef green fill: #696,color: #fff,font-weight: bold;
    classDef yellow fill:#FFF9C4,color:#000000;
```

# references

- [A workshop on Linux containers: Rebuild Docker from Scratch](https://github.com/Fewbytes/rubber-docker/tree/master)
- [Linux containers in 500 lines of code](https://blog.lizzie.io/linux-containers-in-500-loc.html)
- [自己动手写docker](https://github.com/xianlubird/mydocker/tree/master)

