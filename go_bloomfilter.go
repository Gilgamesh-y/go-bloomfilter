package go_bloomfilter

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"strconv"
)

type Bitmap struct {
	set []uint64
	k uint
	m uint64
}

func New() *Bitmap {
	return &Bitmap{k: 4}
}

func (bitmap *Bitmap) getIndexPos(h [4]uint64, k uint) (index, position uint64) {
	return (h[k] >> 6) % 2, h[k] % 64
}

func hash(data []byte) [4]uint64 {
	h := murmur3.New128()
	h.Write(data)
	v1, v2 := h.Sum128()
	h.Write([]byte{byte(v1)})
	v3, v4 := h.Sum128()
	return [4]uint64{v1, v2, v3, v4}
}

func (bitmap *Bitmap) Remove(bytes []byte) bool {
	for i := uint(0); i < bitmap.k; i++ {
		h := hash(bytes)
		index, position  := bitmap.getIndexPos(h, i)
		if (index >= uint64(len(bitmap.set))) {
			return false
		}
		if bitmap.set[index]&(1<<position) != 0 {
			bitmap.set[index] &^= 1<<position
			bitmap.m--
		}
	}

	return true
}

func (bitmap *Bitmap) RemoveInt(num uint) bool {
	if num < 0 {
		return false
	}

	return bitmap.Remove([]byte{byte(num)})
}

func (bitmap *Bitmap) RemoveString(s string) bool {
	if s == "" {
		return false
	}

	return bitmap.Remove([]byte(s))
}

func (bitmap *Bitmap) Has(index, position uint64) bool {
	return index < uint64(len(bitmap.set)) && (bitmap.set[index]&(1<<position)) != 0
}

func (bitmap *Bitmap) HasInt(num int) bool {
	for i := uint(0); i < bitmap.k; i++ {
		h := hash([]byte{byte(num)})
		index, position  := bitmap.getIndexPos(h, i)
		if (!bitmap.Has(index, position)) {
			return false
		}
	}
	return true
}

func (bitmap *Bitmap) HasString(s string) bool {
	for i := uint(0); i < bitmap.k; i++ {
		h := hash([]byte(s))
		index, position  := bitmap.getIndexPos(h, i)
		if (!bitmap.Has(index, position)) {
			return false
		}
	}
	return true
}

func (bitmap *Bitmap) Add(i, position uint64) {
	for i >= uint64(len(bitmap.set)) {
		bitmap.set = append(bitmap.set, 0)
	}
	// 若num不存在于bitmap中
	if bitmap.set[i]&(1<<position) == 0 {
		bitmap.set[i] |= 1 << position
		bitmap.m++
	}
}

func (bitmap *Bitmap) AddInt(num int) {
	for i := uint(0); i < bitmap.k; i++ {
		h := hash([]byte{byte(num)})
		index, position  := bitmap.getIndexPos(h, i)
		bitmap.Add(index, position)
	}
}

func (bitmap *Bitmap) AddString(s string) {
	for i := uint(0); i < bitmap.k; i++ {
		h := hash([]byte(s))
		index, position  := bitmap.getIndexPos(h, i)
		bitmap.Add(index, position)
	}
}

func (bitmap *Bitmap) String() string {
	str := "{"
	for i, v := range bitmap.set {
		if v == 0 {
			continue
		}
		for j := uint(0); j < 64; j++ {
			if v&(1<<j) != 0 {
				if len(str) > len("{") {
					str += ","
				}
				str += strconv.Itoa(int(64*uint(i)+j))
			}
		}
	}
	str += "}"
	fmt.Printf("\nLength: %d\n", len(bitmap.set))
	return str
}