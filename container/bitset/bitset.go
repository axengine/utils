// github.com/alex023/basekit/tree/master/container/bitset
package bitset

import "math"

const (
	max_value = math.MaxUint32
	shift     = 5
	mask      = 0x1F
)

//var defaultBitSet = NewBitSet()

// BitSet  是位图索引的基本数据结构。用一个bit位来标记某个元素对应的value，而key即是这个元素。由于采用bit为单位来存储数据，因此在可以大大的节省存储空间
// 通过该结构提供的方法,可用于排序、查重。对于包括所有uint32的数字的key（约40亿）,存储所有的整形,只需要其1/32的容量,即500M即可。
type BitSet struct {
	//保存实际的bit数据
	data []uint32
	//指示该BitSet的bit容量
	bitsize uint32
	//该bitmap被设置1的最大位数
	maxpos uint32
}

// 默认的位图索引创建方法,按照uint32的最大范围开辟500M内存空间以作计算使用。
func NewBitSet() *BitSet {
	return &BitSet{
		data:    make([]uint32, max_value>>shift),
		bitsize: max_value,
	}
}

// 根据指定的值创建位图索引,分配适合的空间,但输入值不得大于uint32的最大值。。
func NewBitSetByMax(maxnum uint32) *BitSet {
	if maxnum == 0 || maxnum > max_value {
		maxnum = max_value
	}
	return &BitSet{
		data:    make([]uint32, maxnum>>shift),
		bitsize: maxnum,
	}
}

// 添加给定key
func (bitset *BitSet) Add(key uint32) {
	index, pos := key>>shift, key&mask
	bitset.data[index] |= (1 << (pos))
	if bitset.maxpos < key {
		bitset.maxpos = key
	}
}

// 删除给定key
func (bitset *BitSet) Del(key uint32) {
	index, pos := key>>shift, key&mask
	bitset.data[index] &^= (1 << (pos))
}

// 判断给定key是否存在
func (BitSet *BitSet) IsExist(key uint32) bool {
	index, pos := key>>shift, key&mask
	return (BitSet.data[index]>>pos)&0x01 == 1
}
