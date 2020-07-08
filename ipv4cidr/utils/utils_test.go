// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package utils

import (
	"testing"

	"github.com/microsoft/go-cidr-manager/ipv4cidr/consts"

	"github.com/stretchr/testify/assert"
)

// TestGetNetmask generates the netmask for /32 to /0 block sizes
// Success Metric: The correct netmask is generated for each value.
// For /32 the netmask is a 32-bit number with all bits set, and left-shifted for each lower mask number
func TestGetNetmask(t *testing.T) {

	netMask := consts.MaxUInt32
	var mask uint8

	for mask = 32; mask <= 32; mask-- {

		assert.Equal(t, netMask, GetNetmask(mask), "Netmask for %d should be %d", mask, netMask)
		netMask <<= 1

	}

}

// TestGetCIDRRangeLength calculates the range size for /32 to /0 block sizes
// Success Metric: The correct range size is calculated for each value.
// For /32, the range is 1 and multiplied by 2 for each lower mask number
func TestGetCIDRRangeLength(t *testing.T) {

	rangeLength := uint32(1)
	var mask uint8

	for mask = 32; mask <= 32; mask-- {

		assert.Equal(t, rangeLength, GetCIDRRangeLength(mask), "Range length for %d should be %d", mask, rangeLength)
		rangeLength *= 2

	}

}

// TestStandardize uses the IP address and CIDR block number to calculate the first IP address in the CIDR block
// Success Metric: The first IP address of the CIDR block is returned
func TestStandardize(t *testing.T) {

	standardIP := uint32(168427520)    // 10.10.0.0
	nonStandardIP := uint32(168427620) // 10.10.0.100
	mask := uint8(20)                  // IPs = 4096
	netMask := GetNetmask(mask)

	assert.Equal(t, standardIP, Standardize(standardIP, netMask), "The standardized form of 10.10.0.0/20 is 10.10.0.0/20")
	assert.Equal(t, standardIP, Standardize(nonStandardIP, netMask), "The standardized form of 10.10.0.100/20 is 10.10.0.0/20")

}

// TestCheckStandardized checks if the IP of the IP/Mask pair is the first IP in CIDR block
// Success Metric: If 1st IP, error is nil. Else, error is thrown saying IP is not standard
func TestCheckStandarized(t *testing.T) {

	standardIP := uint32(168427520)    // 10.10.0.0
	nonStandardIP := uint32(168427620) // 10.10.0.100
	mask := uint8(20)                  // IPs = 4096
	netMask := GetNetmask(mask)

	err := CheckStandardized(standardIP, netMask)
	assert.Nil(t, err, "10.10.0.0/20 is a standard CIDR representation. No error should be thrown")

	err = CheckStandardized(nonStandardIP, netMask)
	if assert.Error(t, err, "IP from a non-standard IP/CIDR was passed. An error should be thrown.") {
		assert.Equal(t, consts.NonStandardizedIPError, err.Error(), "Error thrown should be: \"%s\"", consts.NonStandardizedIPError)
	}

}

// TestConvertIPToString converts an IP in integer format to string format
// Success Metric: IP is successfully converted to its string representation
func TestConvertIPToString(t *testing.T) {

	IP1 := uint32(168427520) // 10.10.0.0
	IP2 := uint32(168427620) // 10.10.0.100
	assert.Equal(t, "10.10.0.0", ConvertIPToString(IP1))
	assert.Equal(t, "10.10.0.100", ConvertIPToString(IP2))

}
