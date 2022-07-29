// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package machine

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/google/cadvisor/fs"
	info "github.com/google/cadvisor/info/v1"
	"github.com/google/cadvisor/nvm"
	"github.com/google/cadvisor/utils/cloudinfo"
	"github.com/google/cadvisor/utils/sysfs"
	"github.com/google/cadvisor/utils/sysinfo"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"

	"k8s.io/klog/v2"
)

const hugepagesDirectory = "/sys/kernel/mm/hugepages/"
const memoryControllerPath = "/sys/devices/system/edac/mc/"

var machineIDFilePath = flag.String("machine_id_file", "/etc/machine-id,/var/lib/dbus/machine-id", "Comma-separated list of files to check for machine-id. Use the first one that exists.")
var bootIDFilePath = flag.String("boot_id_file", "/proc/sys/kernel/random/boot_id", "Comma-separated list of files to check for boot-id. Use the first one that exists.")

func getInfoFromFiles(filePaths string) string {
	if len(filePaths) == 0 {
		return ""
	}
	for _, file := range strings.Split(filePaths, ",") {
		id, err := ioutil.ReadFile(file)
		if err == nil {
			return strings.TrimSpace(string(id))
		}
	}
	klog.Warningf("Couldn't collect info from any of the files in %q", filePaths)
	return ""
}

func GetIsolCPUs(str1 string) []int {

	// rm last \n if is pressent
	if str1[len(str1)-1] == '\n' {
		str1 = str1[0 : len(str1)-1]
	}

	// strings that contains "="
	res1 := strings.SplitAfter(str1, "=")
	isolStr := res1[len(res1)-1]

	// strings that contains ","
	res2 := strings.SplitAfter(isolStr, ",")

	var isolCores []int

	for index, element := range res2 {
		if index != len(res2)-1 {
			element = element[0 : len(element)-1]
		}

		// is not a range
		if !strings.Contains(element, "-") {

			// fmt.Println(element, "at index ", index, " a single")
			y, e := strconv.Atoi(element)
			if e != nil {
				fmt.Println("Something went wrong1")
			}

			isolCores = append(isolCores, y)

			// is a range
		} else {
			coreStart, e := strconv.Atoi(strings.Split(element, "-")[0])

			if e != nil {
				fmt.Println("Something went wrong2")
				return []int{}
			}

			coreStop, e := strconv.Atoi(strings.SplitAfter(element, "-")[1])

			if e != nil {
				fmt.Println("Something went wrong3")
			}

			for i := coreStart; i <= coreStop; i++ {
				isolCores = append(isolCores, i)
			}
		}
	}

	return isolCores
}

func GetCPUsInfo(numCores int) info.CPUsInfo {

	path := "/proc/cmdline"

	content, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	str1 := string(content)

	isolCPUs := GetIsolCPUs(str1)

	cset_builder := cpuset.NewBuilder()
	for _, core := range isolCPUs {
		cset_builder.Add(core)
	}

	isoledSet := cset_builder.Result()

	cset_all_builder := cpuset.NewBuilder()
	for core := 0; core < numCores; core++ {
		cset_all_builder.Add(core)
	}

	cset_all := cset_all_builder.Result()

	shared := cset_all.Difference(isoledSet)

	return info.CPUsInfo{
		ExlusiveCPUs: isoledSet,
		SharedCPUs:   shared,
	}
}

func Info(sysFs sysfs.SysFs, fsInfo fs.FsInfo, inHostNamespace bool) (*info.MachineInfo, error) {
	rootFs := "/"
	if !inHostNamespace {
		rootFs = "/rootfs"
	}

	cpuinfo, err := ioutil.ReadFile(filepath.Join(rootFs, "/proc/cpuinfo"))
	if err != nil {
		return nil, err
	}
	clockSpeed, err := GetClockSpeed(cpuinfo)
	if err != nil {
		return nil, err
	}

	memoryCapacity, err := GetMachineMemoryCapacity()
	if err != nil {
		return nil, err
	}

	memoryByType, err := GetMachineMemoryByType(memoryControllerPath)
	if err != nil {
		return nil, err
	}

	nvmInfo, err := nvm.GetInfo()
	if err != nil {
		return nil, err
	}

	hugePagesInfo, err := sysinfo.GetHugePagesInfo(sysFs, hugepagesDirectory)
	if err != nil {
		return nil, err
	}

	filesystems, err := fsInfo.GetGlobalFsInfo()
	if err != nil {
		klog.Errorf("Failed to get global filesystem information: %v", err)
	}

	diskMap, err := sysinfo.GetBlockDeviceInfo(sysFs)
	if err != nil {
		klog.Errorf("Failed to get disk map: %v", err)
	}

	netDevices, err := sysinfo.GetNetworkDevices(sysFs)
	if err != nil {
		klog.Errorf("Failed to get network devices: %v", err)
	}

	topology, numCores, err := GetTopology(sysFs)
	if err != nil {
		klog.Errorf("Failed to get topology information: %v", err)
	}

	systemUUID, err := sysinfo.GetSystemUUID(sysFs)
	if err != nil {
		klog.Errorf("Failed to get system UUID: %v", err)
	}

	realCloudInfo := cloudinfo.NewRealCloudInfo()
	cloudProvider := realCloudInfo.GetCloudProvider()
	instanceType := realCloudInfo.GetInstanceType()
	instanceID := realCloudInfo.GetInstanceID()

	machineInfo := &info.MachineInfo{
		CPUsInfo:         GetCPUsInfo(numCores),
		Timestamp:        time.Now(),
		CPUVendorID:      GetCPUVendorID(cpuinfo),
		NumCores:         numCores,
		NumPhysicalCores: GetPhysicalCores(cpuinfo),
		NumSockets:       GetSockets(cpuinfo),
		CpuFrequency:     clockSpeed,
		MemoryCapacity:   memoryCapacity,
		MemoryByType:     memoryByType,
		NVMInfo:          nvmInfo,
		HugePages:        hugePagesInfo,
		DiskMap:          diskMap,
		NetworkDevices:   netDevices,
		Topology:         topology,
		MachineID:        getInfoFromFiles(filepath.Join(rootFs, *machineIDFilePath)),
		SystemUUID:       systemUUID,
		BootID:           getInfoFromFiles(filepath.Join(rootFs, *bootIDFilePath)),
		CloudProvider:    cloudProvider,
		InstanceType:     instanceType,
		InstanceID:       instanceID,
	}

	for i := range filesystems {
		fs := filesystems[i]
		inodes := uint64(0)
		if fs.Inodes != nil {
			inodes = *fs.Inodes
		}
		machineInfo.Filesystems = append(machineInfo.Filesystems, info.FsInfo{Device: fs.Device, DeviceMajor: uint64(fs.Major), DeviceMinor: uint64(fs.Minor), Type: fs.Type.String(), Capacity: fs.Capacity, Inodes: inodes, HasInodes: fs.Inodes != nil})
	}

	return machineInfo, nil
}

func ContainerOsVersion() string {
	os, err := getOperatingSystem()
	if err != nil {
		os = "Unknown"
	}
	return os
}

func KernelVersion() string {
	uname := &unix.Utsname{}

	if err := unix.Uname(uname); err != nil {
		return "Unknown"
	}

	return string(uname.Release[:bytes.IndexByte(uname.Release[:], 0)])
}
