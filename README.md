# tiny-docker
> `tiny-docker` is a **toy container engine** that mocks Docker, built from scratch to replicate the core functionality of Docker in a lightweight form factor. 
>
> It is built on Go 1.23.4, allowing us to leverage more modern language features that simplify this lightweight project.

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

