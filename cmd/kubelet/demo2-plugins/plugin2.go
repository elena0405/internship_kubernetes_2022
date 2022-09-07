package main

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/state"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpumanager/topology"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

func init() {
	log.Println("plugin2 init")
}

var V cpuset.CPUSet

func F_modif_V() cpuset.CPUSet {
	fmt.Println("plugin2: cpuSet variable V= ", V)

	return cpuset.NewCPUSet(9).Union(V)
}

type foo struct{}

func (foo) M1() {
	fmt.Println("plugin2: invoke foo.M1")
}

func M() {
	fmt.Printf("Hello world, from Ioan go plugin!\n")
}

func GetPluginName() string {

	return "policy2"
}

var Foo foo

type Policy interface {
	Name() string
	Start(s state.State) error
	// Allocate call is idempotent
	Allocate(s state.State, pod *v1.Pod, container *v1.Container) error
	// RemoveContainer call is idempotent
	RemoveContainer(s state.State, podUID string, containerName string) error
	// GetTopologyHints implements the topologymanager.HintProvider Interface
	// and is consulted to achieve NUMA aware resource alignment among this
	// and other resource controllers.
	GetTopologyHints(s state.State, pod *v1.Pod, container *v1.Container) map[string][]topologymanager.TopologyHint
	// GetPodTopologyHints implements the topologymanager.HintProvider Interface
	// and is consulted to achieve NUMA aware resource alignment per Pod
	// among this and other resource controllers.
	GetPodTopologyHints(s state.State, pod *v1.Pod) map[string][]topologymanager.TopologyHint
	// GetAllocatableCPUs returns the assignable (not allocated) CPUs
	GetAllocatableCPUs(m state.State) cpuset.CPUSet
}

type staticPolicy struct {
	// cpu socket topology
	topology *topology.CPUTopology
	// set of CPUs that is not available for exclusive assignment
	reserved cpuset.CPUSet
	// topology manager reference to get container Topology affinity
	affinity topologymanager.Store
	// set of CPUs to reuse across allocations in a pod
	cpusToReuse map[string]cpuset.CPUSet
	// options allow to fine-tune the behaviour of the policy
	// options StaticPolicyOptions
}

// var _ Policy = &staticPolicy{}

// func (p *staticPolicy) GetAllocatableCPUs(s state.State) cpuset.CPUSet {
func (staticPolicy) GetAllocatableCPUs(s state.State) {
	fmt.Println("[from plugin2]: GetAllocatableCPUs")
	// return s.GetDefaultCPUSet().Difference(p.reserved)
}

func NewPolicy(topology *topology.CPUTopology, numReservedCPUs int, reservedCPUs cpuset.CPUSet, affinity topologymanager.Store, cpuPolicyOptions map[string]string) (Policy, error) {
	fmt.Println("[from plugin2]: NewStaticPolicy")
	// policy := &staticPolicy{
	// 	topology:    topology,
	// 	affinity:    affinity,
	// 	cpusToReuse: make(map[string]cpuset.CPUSet),
	// }

	return nil, nil
}

// func (p *staticPolicy) Allocate(s state.State, pod *v1.Pod, container *v1.Container) error {
func (staticPolicy) Allocate(s state.State, pod *v1.Pod, container *v1.Container) {
	fmt.Println("[from plugin2]: Allocate")
	// return nil
}

// func (p *staticPolicy) RemoveContainer(s state.State, podUID string, containerName string) error {
func (staticPolicy) RemoveContainer(s state.State, podUID string, containerName string) {
	fmt.Println("[from plugin2]: RemoveContainer")
	// return nil
}

// func (p *staticPolicy) GetTopologyHints(s state.State, pod *v1.Pod, container *v1.Container) map[string][]topologymanager.TopologyHint {
func (staticPolicy) GetTopologyHints(s state.State, pod *v1.Pod, container *v1.Container) {
	fmt.Println("[from plugin2]: GetTopologyHints")
	// return nil
}

// func (p *staticPolicy) GetPodTopologyHints(s state.State, pod *v1.Pod) map[string][]topologymanager.TopologyHint {
func (staticPolicy) GetPodTopologyHints(s state.State, pod *v1.Pod) {
	fmt.Println("[from plugin2]: GetPodTopologyHints")
	// return nil
}

// func (p *staticPolicy) GetAllocatableCPUs(s state.State) cpuset.CPUSet {
// 	return nil
// }

var MyPolicy staticPolicy
