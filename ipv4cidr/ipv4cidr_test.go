// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package ipv4cidr

import (
	"testing"

	"go-cidr-manager/ipv4cidr/consts"

	"github.com/stretchr/testify/assert"
)

// TestValidCIDRWithoutStandardization tests CIDR blocks where the IP is the first IP of the CIDR block
// Success Metric: Create the correct IPv4CIDR block from the input string
func TestValidCIDRWithoutStandardization(t *testing.T) {

	CIDR, err := NewIPv4CIDR("10.10.0.0/26", false)

	assert.Nil(t, err, "10.10.0.0/26 is a valid CIDR block, object should be created.")

	assert.Equal(t, "10.10.0.0/26", CIDR.ToString())
	assert.Equal(t, "10.10.0.0", CIDR.GetIP(), "IP in object should match expected IP.")
	assert.Equal(t, "255.255.255.192", CIDR.GetNetmask(), "Netmask in object should match expected netmask.")
	assert.Equal(t, uint32(64), CIDR.GetCIDRRangeLength(), "Range length in object should match expected range.")
	assert.Equal(t, uint8(26), CIDR.GetMask(), "Mask in object should match expected mask.")

}

// TestInvalidCIDRWithoutStandardization tests CIDR blocks where the IP is NOT the first IP of the CIDR block
// Success Metric: Throw an error pointing out that is isn't the standard notation
func TestInvalidCIDRWithoutStandardization(t *testing.T) {

	_, err := NewIPv4CIDR("10.10.0.1/26", false)

	if assert.Error(t, err, "10.10.0.1/26 is not standard because the IP isn't the first IP in range. An error should be thrown.") {

		assert.Equal(t, consts.NonStandardizedIPError, err.Error(), "Error thrown should be: \"%s\"", consts.NonStandardizedIPError)

	}

}

// TestInvalidRegexInput checks if the input matches the Regex pattern for a correct CIDR block
// Success Metric: Throw an error because all inputs are invalid
func TestInvalidRegexInput(t *testing.T) {

	testInputs := []string{
		"200.200.200.256/23",
		"10.10.0.0/33",
		"10.10.0.0/123",
		"10.10.0.0/",
		"10.20.3.",
		"10.20.3",
		"10./3",
	}

	for _, input := range testInputs {

		_, err := NewIPv4CIDR(input, false)
		if assert.Error(t, err, "%s is an invalid CIDR block. An error should be thrown.", input) {

			assert.Equal(t, consts.InvalidIPv4CIDRError, err.Error(), "For input %s, Error thrown should be: \"%s\"", input, consts.InvalidIPv4CIDRError)

		}

	}

}

// TestInvalidCIDRWithStandardization takes a non-standard IP/CIDR and converts the IP to the first IP in CIDR range
// Success Metric: Create the correct IPv4CIDR block from the input string
func TestInvalidCIDRWithStandardization(t *testing.T) {

	CIDR, err := NewIPv4CIDR("10.10.0.1/26", true)

	assert.Nil(t, err, "An IPv4CIDR object should be created for 10.10.0.0/26, as standardize flag is set to true")

	assert.Equal(t, "10.10.0.0", CIDR.GetIP(), "IP in object should match expected IP")
	assert.Equal(t, "255.255.255.192", CIDR.GetNetmask(), "Netmask in object should match expected netmask")
	assert.Equal(t, uint32(64), CIDR.GetCIDRRangeLength(), "Range length in object should match expected range")
	assert.Equal(t, uint8(26), CIDR.GetMask(), "Mask in object should match expected mask")

}

// TestSplittableCIDRRange takes a CIDR range with size > 1 and splits it into two equal ranges
// Success Metric: Create two CIDR blocks of half the length from the parent CIDR block
func TestSplittableCIDRRange(t *testing.T) {

	CIDR, _ := NewIPv4CIDR("10.10.0.0/26", false)
	subCIDR1, subCIDR2, err := CIDR.Split()
	subCIDRMask := CIDR.GetMask() + 1
	subCIDRRange := CIDR.GetCIDRRangeLength() / 2

	assert.Nil(t, err, "Successfully created an IPv4CIDR object for 10.10.0.0/26")

	assert.Equal(t, "10.10.0.0", subCIDR1.GetIP(), "IP in object should match expected IP")
	assert.Equal(t, "255.255.255.224", subCIDR1.GetNetmask(), "Netmask in object should match expected netmask")
	assert.Equal(t, subCIDRRange, subCIDR1.GetCIDRRangeLength(), "Range length in object should match expected range")
	assert.Equal(t, subCIDRMask, subCIDR1.GetMask(), "Mask in object should match expected mask")

	assert.Equal(t, "10.10.0.32", subCIDR2.GetIP(), "IP in object should match expected IP")
	assert.Equal(t, "255.255.255.224", subCIDR2.GetNetmask(), "Netmask in object should match expected netmask")
	assert.Equal(t, subCIDRRange, subCIDR2.GetCIDRRangeLength(), "Range length in object should match expected range")
	assert.Equal(t, subCIDRMask, subCIDR2.GetMask(), "Mask in object should match expected mask")

}

// TestUnsplittableCIDRRange takes a CIDR range with size = 1 and attempts to splits it into two equal ranges
// Success Metric: Throw an error saying this CIDR range cannot be split further
func TestUnsplittableCIDRRange(t *testing.T) {

	CIDR, _ := NewIPv4CIDR("10.10.0.0/32", false)
	_, _, err := CIDR.Split()

	if assert.Error(t, err, "%s cannot be split further. An error should be thrown.", "10.10.0.0/32") {

		assert.Equal(t, consts.NoMoreSplittingPossibleError, err.Error(), "Error thrown should be: \"%s\"", consts.NoMoreSplittingPossibleError)

	}

}

// TestSingleIPInput takes an IP address as valid CIDR input
// Success Metric: Create an IPv4CIDR object with mask = 32
func TestSingleIPInput(t *testing.T) {

	CIDR, err := NewIPv4CIDR("10.2.3.4", false)

	assert.Nil(t, err, "10.2.3.4 is a valid CIDR block, object should be created.")

	assert.Equal(t, "10.2.3.4", CIDR.GetIP(), "IP in object should match expected IP")
	assert.Equal(t, "255.255.255.255", CIDR.GetNetmask(), "Netmask in object should match expected netmask")
	assert.Equal(t, uint32(1), CIDR.GetCIDRRangeLength(), "Range length in object should match expected range")
	assert.Equal(t, uint8(32), CIDR.GetMask(), "Mask in object should match expected mask")

}

// TestNthIPInRange gets the nth IP address within the range of CIDR block
// Success Metric: Return string corresponding to the nth IP address
func TestNthIPInRange(t *testing.T) {

	CIDR, _ := NewIPv4CIDR("10.10.0.0/26", false)

	NthIP, err := CIDR.GetIPInRange(10, false)
	assert.Nil(t, err, "IP to get is within range, should be generated.")
	assert.Equal(t, "10.10.0.9", NthIP)

	NthIPWithCIDR, err := CIDR.GetIPInRange(10, true)
	assert.Nil(t, err, "IP to get is within range, should be generated.")
	assert.Equal(t, "10.10.0.9/26", NthIPWithCIDR)

}

// TestNthIPNotInRange tries to get the nth IP address exceeding the range of CIDR block
// Success Metric: Throw an error saying n is out of range
func TestNthIPNotInRange(t *testing.T) {

	CIDR, _ := NewIPv4CIDR("10.10.0.0/30", false)

	_, err := CIDR.GetIPInRange(10, false)
	if assert.Error(t, err, "IP to get is out range. An error should be thrown.") {

		assert.Equal(t, consts.RequestedIPExceedsCIDRRangeError, err.Error(), "Error thrown should be: \"%s\"", consts.RequestedIPExceedsCIDRRangeError)

	}

}
