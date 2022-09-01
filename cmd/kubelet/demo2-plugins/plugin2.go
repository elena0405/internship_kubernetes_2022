package main

import (
	"fmt"
	"log"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func init() {
	log.Println("plugin1 init")
}

var V cpuset.CPUSet

func F_modif_V() cpuset.CPUSet {
	fmt.Println("plugin1: cpuSet variable V= ", V)

	return cpuset.NewCPUSet(10).Union(V)
}

type foo struct{}

func (foo) M1() {
	fmt.Println("plugin1: invoke foo.M1")
}

func M() {
	fmt.Printf("Hello world, from Ioan go plugin!\n")
}

func GetPluginName() string {

	return "plugin2"
}

var Foo foo

// type staticPolicy struct {
// 	// cpu socket topology
// 	topology *topology.CPUTopology
// 	// set of CPUs that is not available for exclusive assignment
// 	reserved cpuset.CPUSet
// 	// topology manager reference to get container Topology affinity
// 	affinity topologymanager.Store
// 	// set of CPUs to reuse across allocations in a pod
// 	cpusToReuse map[string]cpuset.CPUSet
// 	// options allow to fine-tune the behaviour of the policy
// 	// options StaticPolicyOptions
// }

// func (p *staticPolicy) GetAllocatableCPUs(s state.State) cpuset.CPUSet {
// 	return s.GetDefaultCPUSet().Difference(p.reserved)
// }

// func NewStaticPolicy(topology *topology.CPUTopology, numReservedCPUs int, reservedCPUs cpuset.CPUSet, affinity topologymanager.Store, cpuPolicyOptions map[string]string) (Policy, error) {
// 	policy := &staticPolicy{
// 		topology:    topology,
// 		affinity:    affinity,
// 		cpusToReuse: make(map[string]cpuset.CPUSet),
// 	}

// 	return policy, nil
// }

// func (p *staticPolicy) Allocate(s state.State, pod *v1.Pod, container *v1.Container) error {
// 	return nil
// }

// func (p *staticPolicy) RemoveContainer(s state.State, podUID string, containerName string) error {
// 	return nil
// }

// func (p *staticPolicy) GetTopologyHints(s state.State, pod *v1.Pod, container *v1.Container) map[string][]topologymanager.TopologyHint {
// 	return nil
// }

// func (p *staticPolicy) GetPodTopologyHints(s state.State, pod *v1.Pod) map[string][]topologymanager.TopologyHint {
// 	return nil
// }
