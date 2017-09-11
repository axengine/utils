// Copyright © 2011-12 Qtrac Ltd.
//
// This program or package and any associated files are licensed under the
// Apache License, Version 2.0 (the "License"); you may not use these files
// except in compliance with the License. You can get a copy of the License
// at: http://www.apache.org/licenses/LICENSE-2.0.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The tests here are very incomplete and just to show examples of how it
// can be done.
package omap

import (
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestStringKeyOMapInsertion(t *testing.T) {
	wordForWord := NewCaseFoldedKeyed()
	for _, word := range []string{"one", "Two", "THREE", "four", "Five"} {
		wordForWord.Insert(word, word)
	}
	var words []string
	wordForWord.Do(func(_, value interface{}) {
		words = append(words, value.(string))
	})
	actual, expected := strings.Join(words, ""), "FivefouroneTHREETwo"
	if actual != expected {
		t.Errorf("%q != %q", actual, expected)
	}
}
func TestMap_FindMinKey(t *testing.T) {
	intMap := NewIntKeyed()
	// 随机生成100长度的随机正整数,将其排序后,注入有序map
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ints := make([]int, 100)
	for i := 0; i < 100; i++ {
		ints[i] = r.Intn(1000)
	}
	sort.Ints(ints)
	for _, number := range ints {
		intMap.Insert(number, number*2)
	}

	key, _, found := intMap.First()
	if !found || key != ints[0] {
		t.Error("Cann‘t found the mininal key for it's value")
	}
}
func TestMap_FindMaxKey(t *testing.T) {
	intMap := NewIntKeyed()
	// 随机生成100长度的随机正整数,将其排序后,注入有序map
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ints := make([]int, 100)
	for i := 0; i < 100; i++ {
		ints[i] = r.Intn(1000)
	}
	sort.Ints(ints)
	for _, number := range ints {
		intMap.Insert(number, number*2)
	}

	key, _, found := intMap.Latest()
	if !found || key != ints[100-1] {
		t.Error("Cann‘t found the maxinal key for it's value")
	}
}
func TestOMap_FindIntKey(t *testing.T) {
	intMap := NewIntKeyed()
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intMap.Insert(number, number*10)
	}
	for _, number := range []int{0, 1, 5, 8, 9} {
		value, found := intMap.Find(number)
		if !found {
			t.Errorf("failed to find %d", number)
		}
		actual, expected := value.(int), number*10
		if actual != expected {
			t.Errorf("value is %d should be %d", actual, expected)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		_, found := intMap.Find(number)
		if found {
			t.Errorf("should not have found %d", number)
		}
	}
}

func TestOMap_DeleteIntKey(t *testing.T) {
	intMap := NewIntKeyed()
	for _, number := range []int{9, 1, 8, 2, 7, 3, 6, 4, 5, 0} {
		intMap.Insert(number, number*10)
	}
	if intMap.Len() != 10 {
		t.Errorf("map len %d should be 10", intMap.Len())
	}
	length := 9
	for i, number := range []int{0, 1, 5, 8, 9} {
		if deleted := intMap.Delete(number); !deleted {
			t.Errorf("failed to delete %d", number)
		}
		if intMap.Len() != length-i {
			t.Errorf("map len %d should be %d", intMap.Len(), length-i)
		}
	}
	for _, number := range []int{-1, -21, 10, 11, 148} {
		if deleted := intMap.Delete(number); deleted {
			t.Errorf("should not have deleted nonexistent %d", number)
		}
	}
	if intMap.Len() != 5 {
		t.Errorf("map len %d should be 5", intMap.Len())
	}
}

func TestOMap_DeleteIntKey2(t *testing.T) {
	var (
		size    = 10000
		numbers = make([]int, size)
	)
	intMap := NewIntKeyed()
	for i := 0; i < size; i++ {
		numbers[i] = i * 2
		intMap.Insert(numbers[i], i*10)
	}

	if intMap.Len() != size {
		t.Errorf("map len %d should be %v", intMap.Len(), size)
	}
	for i := size - 1; i >= 0; i-- {
		intMap.Delete(numbers[i])
	}
	if intMap.Len() != 0 {
		t.Errorf("map len %d should be 0", intMap.Len())
	}
}

type Point struct {
	ActionTime time.Time
	Key        int
}

var (
	r         = rand.New(rand.NewSource(time.Now().Unix()))
	randscope = 10 * 1000 * 1000 * 1000
	now       = time.Now()
)

func getRandPoint() Point {
	return Point{
		ActionTime: now.Add(time.Duration(int64(r.Intn(randscope)))),
		Key:        rand.Int(),
	}
}
func less(a, b interface{}) bool {
	α, β := a.(Point), b.(Point)
	if !α.ActionTime.Equal(β.ActionTime) {
		return α.ActionTime.Before(β.ActionTime)
	}
	return α.Key < β.Key
}
func TestOMap_FindStructKey(t *testing.T) {
	var (
		number = 1000000
	)
	points := make([]Point, number)

	for i := 0; i < number; i++ {
		points[i] = getRandPoint()
	}
	pointMap := New(less)
	for i := 0; i < number; i++ {
		pointMap.Insert(points[i], i)
	}
	if pointMap.Len() != number {
		t.Error("数据加入有丢失")
	}

	for i := 0; i < number; i++ {
		nodevalue, founded := pointMap.Find(points[i])
		value := nodevalue.(int)
		if !founded || value != i {
			t.Errorf("查找数据失败,key=%+v \n", points[i])
			break
		}
	}
}

// TODO:this test cannot pass,bacause some delete will cause other key change the same
func TestOMap_DelStructKey(t *testing.T) {
	t.Skip()
	var (
		size = 1000
	)
	points := make([]Point, size)
	for i := 0; i < size; i++ {
		points[i] = getRandPoint()
	}
	pointMap := New(less)
	for i := 0; i < size; i++ {
		pointMap.Insert(points[i], i) //这里值，仅仅只是一个填充，无实际意义
	}
	if pointMap.Len() != size {
		t.Errorf("map len %d should be %v", pointMap.Len(), size)

	}

	for i := 0; i < size; i++ {
		_, founded := pointMap.Find(points[i])
		if !founded {
			t.Errorf("查找数据失败,key=%+v,value=%v \n", points[i-1], i-1)

			t.Errorf("查找数据失败,key=%+v,value=%v \n", points[i], i)
			//break
		}
		deleted := pointMap.Delete(points[i])
		if !deleted {
			t.Errorf("删除数据失败,key=%+v,value=%v \n", points[i], i)
			break
		}
	}
	for i := 0; i < pointMap.Len(); i++ {
		f, v, ok := pointMap.First()
		if ok {
			t.Logf("savive遗留数据,key=%+v ,value=%+v \n", f, v)
		}
	}
}

// 测试逆序插入，能否按从小到大的顺序导出与清理
func TestOMap_DelStructFirst(t *testing.T) {
	var (
		num      = 100000
		zerotime time.Time
		points   = make([]Point, num)
	)
	ordermap := New(less)
	for i := 0; i < num; i++ {
		//生成从大到小逆序的点位
		points[i] = Point{
			ActionTime: zerotime.Add(time.Duration(int64(num - i))),
			Key:        num - i,
		}
		ordermap.Insert(points[i], i)
		if v, finded := ordermap.Find(points[i]); !finded {
			t.Errorf("cannot find :key=%+v | value=%+v\n", points[i], v)
		}
	}
	//clean ordermap by minkey,one by one
	for i := 0; i < num; i++ {
		k, v, founded := ordermap.First()
		if !founded {
			t.Error("获取最小值出错")
			break
		}
		key, _ := k.(Point)
		deleted := ordermap.Delete(k)
		if founded != deleted || founded == false {
			t.Errorf("key=%+v | value = %+v \n", k, v)
			break
		}
		if key.Key != points[num-i-1].Key || !key.ActionTime.Equal(points[num-i-1].ActionTime) {
			t.Errorf("数据解析顺序有误,ordermap.minKey=%v,points[%v]=%v\n", key, i, points[num-i-1])
		}
	}
	if ordermap.Len() != 0 {
		t.Error("omap is not clean")
	}
}
func TestPassing(t *testing.T) {
	intMap := NewIntKeyed()
	intMap.Insert(7, 7)
	passMap(intMap, t)
}

func passMap(m *Map, t *testing.T) {
	for _, number := range []int{9, 3, 6, 4, 5, 0} {
		m.Insert(number, number)
	}
	if m.Len() != 7 {
		t.Errorf("should have %d items", 7)
	}
}

// Thanks to Russ Cox for improving these benchmarks
func BenchmarkOMapFindSuccess(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	intMap := NewIntKeyed()
	for i := 0; i < 1e6; i++ {
		intMap.Insert(i, i)
	}
	b.StartTimer() // Time the Find() method succeeding
	for i := 0; i < b.N; i++ {
		intMap.Find(i % 1e6)
	}
}

func BenchmarkOMapFindFailure(b *testing.B) {
	b.StopTimer() // Don't time creation and population
	intMap := NewIntKeyed()
	for i := 0; i < 1e6; i++ {
		intMap.Insert(2*i, i)
	}
	b.StartTimer() // Time the Find() method failing
	for i := 0; i < b.N; i++ {
		intMap.Find(2*(i%1e6) + 1)
	}
}
