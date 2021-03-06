// Copyright 2021 go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/gballet/go-verkle"

	"github.com/holiman/uint256"
)

const (
	VersionLeafKey    = 0
	BalanceLeafKey    = 1
	NonceLeafKey      = 2
	CodeKeccakLeafKey = 3
	CodeSizeLeafKey   = 4
)

var (
	zero                = uint256.NewInt(0)
	HeaderStorageOffset = uint256.NewInt(64)
	CodeOffset          = uint256.NewInt(128)
	MainStorageOffset   = new(uint256.Int).Lsh(uint256.NewInt(256), 31)
	VerkleNodeWidth     = uint256.NewInt(256)
	codeStorageDelta    = uint256.NewInt(0).Sub(CodeOffset, HeaderStorageOffset)
)

// GetTreeKey performs both the work of the spec's get_tree_key function, and that
// of pedersen_hash: it builds the polynomial in pedersen_hash without having to
// create a mostly zero-filled buffer and "type cast" it to a 128-long 16-byte
// array. Since at most the first 5 coefficients of the polynomial will be non-zero,
// these 5 coefficients are created directly.
func GetTreeKey(address []byte, treeIndex *uint256.Int, subIndex byte) []byte {
	if len(address) < 32 {
		var aligned [32]byte
		address = append(aligned[:32-len(address)], address...)
	}
	var poly [5]fr.Element

	// (2 + 256 * length) little endian, length = 64 bytes
	verkle.FromLEBytes(&poly[0], []byte{2, 64})

	// 32-byte address, interpreted as two little endian
	// 16-byte numbers.
	verkle.FromLEBytes(&poly[1], address[:16])
	verkle.FromLEBytes(&poly[2], address[16:])

	// little-endian, 32-byte aligned treeIndex
	var index [32]byte
	for i, b := range treeIndex.Bytes() {
		index[len(treeIndex.Bytes())-1-i] = b
	}
	verkle.FromLEBytes(&poly[3], index[:16])
	verkle.FromLEBytes(&poly[4], index[16:])

	cfg, _ := verkle.GetConfig()
	ret := cfg.CommitToPoly(poly[:], 0)

	// The output of Byte() is big engian for banderwagon. This
	// introduces an inbalance in the tree, because hashes are
	// elements of a 253-bit field. This means more than half the
	// tree would be empty. To avoid this problem, use a little
	// endian commitment and chop the MSB.
	var retb [32]byte
	retb = ret.Bytes()
	for i := 0; i < 16; i++ {
		retb[31-i], retb[i] = retb[i], retb[31-i]
	}
	retb[31] = subIndex
	return retb[:]

}

func GetTreeKeyAccountLeaf(address []byte, leaf byte) []byte {
	return GetTreeKey(address, zero, leaf)
}

func GetTreeKeyVersion(address []byte) []byte {
	return GetTreeKey(address, zero, VersionLeafKey)
}

func GetTreeKeyBalance(address []byte) []byte {
	return GetTreeKey(address, zero, BalanceLeafKey)
}

func GetTreeKeyNonce(address []byte) []byte {
	return GetTreeKey(address, zero, NonceLeafKey)
}

func GetTreeKeyCodeKeccak(address []byte) []byte {
	return GetTreeKey(address, zero, CodeKeccakLeafKey)
}

func GetTreeKeyCodeSize(address []byte) []byte {
	return GetTreeKey(address, zero, CodeSizeLeafKey)
}

func GetTreeKeyCodeChunk(address []byte, chunk *uint256.Int) []byte {
	chunkOffset := new(uint256.Int).Add(CodeOffset, chunk)
	treeIndex := new(uint256.Int).Div(chunkOffset, VerkleNodeWidth)
	subIndexMod := new(uint256.Int).Mod(chunkOffset, VerkleNodeWidth).Bytes()
	var subIndex byte
	if len(subIndexMod) != 0 {
		subIndex = subIndexMod[0]
	}
	return GetTreeKey(address, treeIndex, subIndex)
}

func GetTreeKeyStorageSlot(address []byte, storageKey *uint256.Int) []byte {
	pos := storageKey.Clone()
	if storageKey.Cmp(codeStorageDelta) < 0 {
		pos.Add(HeaderStorageOffset, storageKey)
	} else {
		pos.Add(MainStorageOffset, storageKey)
	}
	treeIndex := new(uint256.Int).Div(pos, VerkleNodeWidth)

	// calculate the sub_index, i.e. the index in the stem tree.
	// Because the modulus is 256, it's the last byte of treeIndex
	subIndexMod := new(uint256.Int).Mod(pos, VerkleNodeWidth).Bytes()
	var subIndex byte
	if len(subIndexMod) != 0 {
		// uint256 is broken into 4 little-endian quads,
		// each with native endianness. Extract the least
		// significant byte.
		subIndex = byte(subIndexMod[0] & 0xFF)
	}
	return GetTreeKey(address, treeIndex, subIndex)
}
