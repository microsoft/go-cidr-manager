// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package ipv4cidr

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/microsoft/go-cidr-manager/ipv4cidr/consts"
	"github.com/microsoft/go-cidr-manager/ipv4cidr/utils"
)

// IPv4CIDR models an IPv4 CIDR range.
// @field ip uint32: Holds the IP address
// @field mask uint8: Holds the CIDR mask
// @field netmask uint32: Holds the netmask for the subnet
// @field rangeLength uint32: Holds the number of IP addresses in the CIDR range
type IPv4CIDR struct {
	ip          uint32
	mask        uint8
	netmask     uint32
	rangeLength uint32
}

// NewIPv4CIDR instantiates a new IPv4CIDR object and returns it
// @param IP string: A string representation of CIDR range in the format a.b.c.d/e or a.b.c.d
// @param standardize bool: If the IP part of the CIDR range is not the first IP in range, then setting this value to "true" will automatically convert it to the first IP in range. If set to "false", a non-standard CIDR will give an error
// @returns *IPv4CIDR: If the input parameters are valid, returns a pointer to a new IPv4CIDR object
// @returns error: If the input parameters are invalid, or any processing errors occur, returns the appropriate error back to caller.
func NewIPv4CIDR(IP string, standardize bool) (*IPv4CIDR, error) {

	// Use regex to check if the input string is valid
	isValid, err := regexp.Match(consts.IPv4CIDRRegex, []byte(IP))
	if err != nil {
		return nil, err
	}
	if !isValid {
		err := errors.New(consts.InvalidIPv4CIDRError)
		return nil, err
	}

	// Create an IPv4CIDR object
	ip := IPv4CIDR{}

	// Parse the input string into the IPv4CIDR object
	err = ip.parse(IP, standardize)
	if err != nil {
		return nil, err
	}

	return &ip, nil

}

// parse takes as input the IP string and standardize flag, and parses it
// @input ipString string: A valid IP/CIDR string
// @input standardize bool: Flag for whether to standardize non-standard IP string or throw an error
// @returns error: If there is any processing error, the appropriate error is returned to caller.
func (i *IPv4CIDR) parse(ipString string, standardize bool) error {

	// Instantiate IP as a 32-bit 0.0.0.0
	ip := uint32(0)

	// Instantiate mask with a default value of 32
	mask := uint8(32)

	// Split the IP string into the IP part (ipSections[0]) and optional CIDR part (ipSections[1])
	ipSections := strings.Split(ipString, "/")

	// If there are 2 sections, a CIDR part was provided, use that to set the mask. Else, let mask have default value of 32
	if len(ipSections) == 2 {
		tempMask, err := strconv.Atoi(ipSections[1])
		if err != nil {
			return err
		}
		mask = uint8(tempMask)
	}

	// Split the IP part into 4 sections (a.b.c.d => [a,b,c,d])
	ipNumbers := strings.Split(ipSections[0], ".")

	// Convert each 8-bit section into its integer representation, and set the corresponding 8 bits of the IP's integer representation
	for i := 0; i < 4; i++ {

		tempIP, err := strconv.Atoi(ipNumbers[i])
		if err != nil {
			return err
		}

		ip = ip << consts.GroupSize
		ip = ip | uint32(tempIP)

	}

	netmask := utils.GetNetmask(mask)
	rangeLength := utils.GetCIDRRangeLength(mask)

	// If standardize is true, then standardize the IP part of the object
	// If standardize is false, check if the representation is correct. If not, return an error
	if standardize {
		ip = utils.Standardize(ip, netmask)
	} else {
		err := utils.CheckStandardized(ip, netmask)
		if err != nil {
			return err
		}
	}

	// Set values in the IP object
	i.ip = ip
	i.mask = mask
	i.rangeLength = rangeLength
	i.netmask = netmask

	return nil

}

// Split splits the IPv4CIDR into two IPv4CIDRs of half the size (mask + 1)
// @returns *IPv4CIDR: The first (lower) block
// @returns *IPv4CIDR: The second (higher) block
// @returns error: If CIDR cannot be split further, the appropriate error is returned.
func (i *IPv4CIDR) Split() (*IPv4CIDR, *IPv4CIDR, error) {

	// If we are already at a single-IP CIDR block, further splitting is not possible. Hence return an error
	if i.rangeLength == 1 {
		return nil, nil, errors.New(consts.NoMoreSplittingPossibleError)
	}

	// The new mask becomes the old mask + 1
	newMask := i.mask + 1

	// The new range is half of old range
	newRange := i.rangeLength / 2

	// The new netmask has the leftmost 0 of the old netmask also set
	// In other words, shift right and set the highest bit
	newNetmask := (i.netmask >> 1) | consts.HighestBitSet

	// The lower CIDR block has the same IP
	newIP1 := i.ip

	// The higher CIDR block has the leftmost 0 of the rightmost block of 0s also set.
	// The XOR of the old and new netmasks gives us the bit that needs to be set, which can be done by bitwise OR
	newIP2 := (i.ip | (newNetmask ^ i.netmask))

	// Create the two new IPv4CIDR objects
	IP1 := IPv4CIDR{
		ip:          newIP1,
		mask:        newMask,
		rangeLength: newRange,
		netmask:     newNetmask,
	}

	IP2 := IPv4CIDR{
		ip:          newIP2,
		mask:        newMask,
		rangeLength: newRange,
		netmask:     newNetmask,
	}

	return &IP1, &IP2, nil

}

// GetIPInRange returns the nth IP address in the CIDR block
// @input n uint32: The value of n, representing the nth IP to return
// @input withCIDR bool: Flag corresponding to whether to append the CIDR mask with the returned IP or not
// @returns string: The nth IP address
// @returns error: If nth IP is out of range of the CIDR block, an error is returned
func (i *IPv4CIDR) GetIPInRange(n uint32, withCIDR bool) (string, error) {

	// Check if range exceeded, return error if yes
	if i.rangeLength < n {
		return "", errors.New(consts.RequestedIPExceedsCIDRRangeError)
	}

	// The nth IP is obtained by simply adding n-1 to the 1st IP in CIDR range
	nthIP := i.ip + n - 1

	// Convert the IP to string
	nthIPstr := utils.ConvertIPToString(nthIP)

	// If withCIDR is set, append the CIDR mask to string
	if withCIDR {
		mask := strconv.Itoa(int(i.mask))
		nthIPstr = strings.Join([]string{nthIPstr, mask}, "/")
	}

	return nthIPstr, nil

}

// ToString converts the IP into its string representation
// @returns string: String corresponding to the IP address in format a.b.c.d
func (i *IPv4CIDR) ToString() string {

	ip := utils.ConvertIPToString(i.ip)
	mask := strconv.Itoa(int(i.mask))

	return strings.Join([]string{ip, mask}, "/")

}

// GetIP returns the IP part of the CIDR range
// @returns string: String corresponding to the first IP address in CIDR range in format a.b.c.d
func (i *IPv4CIDR) GetIP() string {

	return utils.ConvertIPToString(i.ip)

}

// GetCIDRRangeLength returns the number of IP addresses contained in the CIDR range
// @returns uint32: Length of the CIDR range
func (i *IPv4CIDR) GetCIDRRangeLength() uint32 {

	return i.rangeLength

}

// GetMask returns the mask part of the CIDR range (0-32)
// @returns uint8: Mask of the CIDR range
func (i *IPv4CIDR) GetMask() uint8 {

	return i.mask

}

// GetNetmask returns the netmask for the CIDR range
// @returns string: Netmask of the CIDR range
func (i *IPv4CIDR) GetNetmask() string {

	return utils.ConvertIPToString(i.netmask)

}
