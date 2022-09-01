package main

import (
	"fmt"
	"log"

	//"github.com/MarianNoaghea/plugin-go/common"
	//v2 "k8s.io/klog/v2"
	k8sCpuSet "k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func init() {
	log.Println("plugin1 init")
}

var V k8sCpuSet.CPUSet

func F_print_V() k8sCpuSet.CPUSet {
	fmt.Println("plugin1: cpuSet variable V= ", V)

	return k8sCpuSet.NewCPUSet(9)
}

type foo struct{}

func (foo) M1() {
	fmt.Println("plugin1: invoke foo.M1")
}

func M() {
	fmt.Printf("Hello world, from Elena go plugin!\n")
}

var Foo foo
