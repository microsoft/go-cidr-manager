// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package consts

// This set of constants defines strings corresponding to the new errors introduced in this package
const (
	InvalidIPv4CIDRError             string = "IP address is invalid, it should be of the format a.b.c.d or a.b.c.d/e, where 0 <= a, b, c, d < 256 and 0 <= e <= 32"
	NonStandardizedIPError           string = "IP address is not standardized, the IP part of IP/CIDR should be the first IP in the range"
	NoMoreSplittingPossibleError     string = "There is only one IP address in this CIDR range, further splitting is not possible"
	RequestedIPExceedsCIDRRangeError string = "Requested IP exceeds the CIDR range"
)
