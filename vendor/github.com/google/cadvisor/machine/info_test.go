package machine

import (
	"testing"
)

type isolCoresTest struct {
	description string
	fileContent string
	expErr      error
	expCPUs     []int
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestIsolCores(t *testing.T) {
	testCases := []isolCoresTest{
		{
			description: "everything works fine",
			expCPUs:     []int{3, 4, 5},
			fileContent: "BOOT_IMAGE=/boot/vmlinuz-5.15.0-41-generic root=UUID=484f5b5d-8372-4855-84e5-3967bd55ad35 ro quiet splash isolcpus=3-5",
		},
		{
			description: "everything works fine 2",
			expCPUs:     []int{7, 8, 9, 10, 11},
			fileContent: "BOOT_IMAGE=/boot/vmlinuz-5.15.0-41-generic root=UUID=484f5b5d-8372-4855-84e5-3967bd55ad35 ro quiet splash isolcpus=7-11",
		},
		{
			description: "no isoled CPUs",
			expCPUs:     []int{},
			fileContent: "BOOT_IMAGE=/boot/vmlinuz-5.15.0-41-generic root=UUID=484f5b5d-8372-4855-84e5-3967bd55ad35 ro quiet splash",
		},
		{
			description: "many isoled CPUs",
			expCPUs:     []int{3, 4, 5, 6, 7, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
			fileContent: "BOOT_IMAGE=/boot/vmlinuz-5.15.0-41-generic root=UUID=484f5b5d-8372-4855-84e5-3967bd55ad35 ro quiet splash isolcpus=3,4,5-7,10-30",
		},
		{
			description: "isoled CPUs declared in between",
			expCPUs:     []int{3, 4, 5, 7},
			fileContent: "BOOT_IMAGE=/boot/vmlinuz-5.15.0-41-generic isolcpus=3-5,7 root=UUID=484f5b5d-8372-4855-84e5-3967bd55ad35 ro quiet splash ",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			cores := GetIsolCPUs(testCase.fileContent)

			for i := 0; i < len(cores); i++ {
				if !contains(testCase.expCPUs, cores[i]) {
					t.Errorf("something went wrong!")
				}
			}
		})
	}
}
