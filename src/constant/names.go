package constant

const (
	CgroupBasePath       = "/sys/fs/cgroup/system.slice"
	DefaultContainerName = "tiny-docker"

	CgroupProcs = "cgroup.procs"
	CpuMax      = "cpu.max"
	MemoryMax   = "memory.max"
)
