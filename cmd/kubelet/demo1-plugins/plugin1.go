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

func F_print_V() cpuset.CPUSet {
	fmt.Println("plugin1: cpuSet variable V= ", V)

	return cpuset.NewCPUSet(9)
}

type foo struct{}

func (foo) M1() {
	fmt.Println("plugin1: invoke foo.M1")
}

func M() {
	fmt.Printf("Hello world, from Ioan go plugin\n")
}

var Foo foo
