package main

type nodeResponse struct {
	Nodes []nodeItem `json:"nodes"`
}

type nodeItem struct {
	Name              string   `json:"name"`
	UUID              string   `json:"uuid"`
	CapacityCPU       string   `json:"capacityCpu"`
	CapacityMemory    string   `json:"capacityMemory"`
	AllocatableCPU    string   `json:"allocableCpu"`
	AllocatableMemory string   `json:"allocableMemory"`
	Machineid         string   `json:"machineId"`
	SystemUUID        string   `json:"systemUUID"`
	KernelVersion     string   `json:"kernelVersion"`
	KubeletVersion    string   `json:"kubeletVersion"`
	CRIVersion        string   `json:"criVersion"`
	OperatingSystem   string   `json:"operatingSystem"`
	Architecture      string   `json:"architecture"`
	PodCidrs          []string `json:"podCidr"`
}

type clusterCapacity struct {
	CPUCapacity     float64 `json:"cpuCapacity"`
	MemoryCapacity  float64 `json:"memoryCapacity"`
	RequestedCpu    float64 `json:"requestdCpu"`
	RequestedMemory float64 `json:"requestdMemory"`
	MaxLimitCpu     float64 `json:"maxLimitCpu"`
	MaxLimitMemory  float64 `json:"maxLimitMemory"`
	UsageRate       string  `json:"usageRate"`
}
