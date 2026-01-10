package constant

const (
	CgroupBasePath       = "/sys/fs/cgroup/system.slice"
	CgroupServiceName    = "tiny-docker.service"
	CgroupSubtreeControl = "cgroup.subtree_control"

	CgroupProcs = "cgroup.procs"
	CpuMax      = "cpu.max"
	MemoryMax   = "memory.max"
)
