// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package consts

// This set of constants contains the numeric constants used throughout this package
const (
	MaxUInt32     uint32 = ^uint32(0)
	MaxBits       uint8  = 32
	EightBits     uint32 = 255
	GroupSize     uint8  = 8
	HighestBitSet uint32 = uint32(1) << (MaxBits - 1)
)
