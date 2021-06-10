package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {

	log.Println("Hello world")
	log.Println("Load the kubeconfig client")
	clientKate := connectKubernetes()
	nodes, err := clientKate.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	var allCpu float64
	var allMemory float64
	var allocatedCpu float64
	var allocatedMemory float64
	var maxLimitCpu float64
	var maxLimitMemory float64

	for _, node := range nodes.Items {

		convertCpu, err := strconv.ParseFloat(node.Status.Allocatable.Cpu().AsDec().String(), 64)
		PanicIfError(err)
		convertMemory, err := strconv.ParseFloat(node.Status.Allocatable.Memory().AsDec().String(), 64)
		PanicIfError(err)
		allCpu = allCpu + convertCpu
		allMemory = allMemory + convertMemory
	}

	// Get Namespace

	ns, err := clientKate.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	PanicIfError(err)
	for _, namespace := range ns.Items {
		rqs, err := clientKate.CoreV1().ResourceQuotas(namespace.GetName()).List(context.TODO(), metav1.ListOptions{})
		log.Println("Parsing : " + namespace.GetName() + " namespace in progress")
		PanicIfError(err)
		for _, rq := range rqs.Items {
			log.Println("Quotas found : " + rq.GetName())
			itemsCpu := rq.Spec.Hard.DeepCopy()["limits.cpu"]
			itemsMemory := rq.Spec.Hard.DeepCopy()["limits.memory"]

			// Convert to json and get the valid with trim
			jsonitems, err := json.MarshalIndent(itemsCpu, "", "    ")
			PanicIfError(err)
			convertCpuRqs, err := strconv.ParseFloat(strings.Trim(string(jsonitems), "\""), 64)
			PanicIfError(err)

			// Convert to json and get the valid with trim
			jsonitemsmemory, err := json.MarshalIndent(itemsMemory.AsApproximateFloat64(), "", "    ")
			PanicIfError(err)
			convertMemoryRqs, err := strconv.ParseFloat(strings.Trim(string(jsonitemsmemory), "\""), 64)
			PanicIfError(err)

			// Get the used memory and used Cpu
			usedMemory := rq.Status.Used.DeepCopy()["limits.memory"]
			usedCpu := rq.Status.Used.DeepCopy()["limits.cpu"]

			jsonUsedMemory, err := json.MarshalIndent(usedMemory.AsApproximateFloat64(), "", "    ")
			PanicIfError(err)
			jsonUsedCpu, err := json.MarshalIndent(usedCpu.AsApproximateFloat64(), "", "    ")
			PanicIfError(err)
			convertAllocatedCpu, err := strconv.ParseFloat(strings.Trim(string(jsonUsedCpu), "\""), 64)
			PanicIfError(err)

			convertAllocatedMemory, err := strconv.ParseFloat(strings.Trim(string(jsonUsedMemory), "\""), 64)

			// Create the sum
			maxLimitCpu = maxLimitCpu + convertCpuRqs
			maxLimitMemory = maxLimitMemory + convertMemoryRqs
			allocatedCpu = allocatedCpu + convertAllocatedCpu
			allocatedMemory = allocatedMemory + convertAllocatedMemory

		}
	}

	usageRate := math.Round(((allocatedCpu / maxLimitCpu) + (allocatedMemory / maxLimitMemory)) * 100)

	var responseObject = clusterCapacity{
		CPUCapacity:     allCpu,
		MemoryCapacity:  allMemory,
		MaxLimitCpu:     maxLimitCpu,
		MaxLimitMemory:  maxLimitMemory,
		RequestedCpu:    allocatedCpu,
		RequestedMemory: allocatedMemory,
		UsageRate:       FloatToString(usageRate) + "%",
	}
	prettyJSON, err := json.MarshalIndent(responseObject, "", "    ")
	fmt.Println(string(prettyJSON))

}
