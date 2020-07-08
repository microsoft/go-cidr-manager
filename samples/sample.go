// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package main

import (
	"fmt"

	"github.com/microsoft/go-cidr-manager/ipv4cidr"
)

func main() {

	IP, err := ipv4cidr.NewIPv4CIDR("10.10.0.0/26", true)
	if err != nil {
		panic(err)
	}

	fmt.Println(IP.ToString())
	fmt.Println("Netmask: ", IP.GetNetmask())
	fmt.Println("Range: ", IP.GetCIDRRangeLength())

	IP8, err := IP.GetIPInRange(8, false)
	if err != nil {
		panic(err)
	}

	IP8CIDR, err := IP.GetIPInRange(8, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(IP8)
	fmt.Println(IP8CIDR)

	IP1, IP2, err := IP.Split()
	if err != nil {
		panic(err)
	}

	fmt.Println(IP1.ToString())
	fmt.Println("Netmask: ", IP1.GetNetmask())
	fmt.Println("Range: ", IP1.GetCIDRRangeLength())

	fmt.Println(IP2.ToString())
	fmt.Println("Netmask: ", IP2.GetNetmask())
	fmt.Println("Range: ", IP2.GetCIDRRangeLength())

}
