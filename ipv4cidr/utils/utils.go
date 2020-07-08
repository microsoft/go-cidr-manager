// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package utils

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/microsoft/go-cidr-manager/ipv4cidr/consts"
)

// GetNetmask takes the mask number as input and creates the netmask from it
// @input mask uint8: The mask for the CIDR range
// @returns uint32: The integer representation of the netmask
func GetNetmask(mask uint8) uint32 {

	// Netmask = 32-bit number with all bits set, shifted left by (32-mask)
	return consts.MaxUInt32 << (consts.MaxBits - mask)

}

// GetCIDRRangeLength calculates the number of IP addresses in that CIDR range
// @input mask uint8: The mask for the CIDR range
// @returns uint32: The length of the CIDR range
func GetCIDRRangeLength(mask uint8) uint32 {

	// Length of CIDR range = 2^(32-mask)
	return uint32(math.Pow(float64(2), float64((consts.MaxBits - mask))))

}

// Standardize converts the IP to the first IP address of the CIDR range
// @input ip uint32: The IP address in integer representation
// @input netmask uint32: The netmask of the CIDR range
// @returns uint32: First IP in CIDR range
func Standardize(ip uint32, netmask uint32) uint32 {

	// A bitwise AND of the input IP and the netmask gives the first IP address in range
	return (ip & netmask)

}

// CheckStandardized checks if the IP stored in object is the first IP in range or not
// @input ip uint32: The IP address in integer representation
// @input netmask uint32: The netmask of the CIDR range
// @returns error: If not the first IP in range, an error is returned. Else, return value is nil
func CheckStandardized(ip uint32, netmask uint32) error {

	// If IP stored in object is same as the standardized representation, then the check passes
	if ip == Standardize(ip, netmask) {
		return nil
	}

	// If above check fails, return an error
	return errors.New(consts.NonStandardizedIPError)

}

// ConvertIPToString converts an integer IP address to its string representation
// @param ip uint32: IP address in integer representation
// @returns string: IP address in string representation
func ConvertIPToString(ip uint32) string {

	// IP addresses consist of 4 sections (a.b.c.d)
	ipSections := make([]string, 4)

	for i := 3; i >= 0; i-- {

		// To generate each section, we go in reverse order (right to left)
		// 1. Pull the least significant 8 bits into another var
		// 2. Convert to int and save in the corresponding section
		// 3. Shift the IP by 8 bits to the right
		sectionInt := int(ip & consts.EightBits)
		ipSections[i] = strconv.Itoa(sectionInt)
		ip = ip >> consts.GroupSize

	}

	return strings.Join(ipSections, ".")

}
