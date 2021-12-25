package ft_test

import (
	"context"
	"functools/ft"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	f := func(input, expected []int) {
		iter := ft.SliceIter(input)
		result := ft.Collect(iter)
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4}, []int{1, 2, 3, 4})
	f([]int{}, []int{})
	f([]int{1}, []int{1})
}

func TestFilter(t *testing.T) {
	f := func(input, expected []int) {
		iter := ft.SliceIter(input)
		result := ft.Collect(ft.Filter(iter, func(t int) bool {
			return t%2 == 0
		}))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4, 5}, []int{2, 4})
	f([]int{}, []int{})
	f([]int{1, 1, 1, 1, 1}, []int{})
	f([]int{2, 2, 2}, []int{2, 2, 2})
}

func TestMap(t *testing.T) {
	f := func(input []int, expected []string) {
		iter := ft.SliceIter(input)
		result := ft.Collect(ft.Map(iter, func(t int) string {
			return strconv.Itoa(t)
		}))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3}, []string{"1", "2", "3"})
}

func TestReduce(t *testing.T) {
	f := func(input []int, expected string, initial ...string) {
		var i string
		if len(initial) > 0 {
			i = initial[0]
		}
		iter := ft.SliceIter(input)
		result := ft.Reduce(iter, func(o string, t int) string {
			return o + strconv.Itoa(t)
		}, i)
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3}, "__123", "__")
	f([]int{1, 2, 3}, "123")
	f([]int{}, "__", "__")
}

func TestFilterMapReduce(t *testing.T) {
	f := func(input []int, expected string) {
		iter := ft.SliceIter(input)
		result := ft.Reduce(ft.Map(ft.Filter(iter, func(t int) bool {
			return t%2 != 0
		}), func(t int) int {
			return t * 2
		}), func(o string, t int) string {
			return o + strconv.Itoa(t)
		})
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4, 5}, "2610")
	f([]int{}, "")
	f([]int{2, 4}, "")
	f([]int{5}, "10")
}

func TestReverse(t *testing.T) {
	f := func(input, expected []int) {
		iter := ft.SliceIter(input).(ft.ReversibleIter[int])
		result := ft.Collect(ft.Reverse(iter))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4}, []int{4, 3, 2, 1})
	f([]int{}, []int{})
	f([]int{1}, []int{1})
}

func TestSkip(t *testing.T) {
	f := func(input, expected []int, skip int) {
		iter := ft.SliceIter(input)
		result := ft.Collect(ft.Skip(iter, skip))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4}, []int{2, 3, 4}, 1)
	f([]int{1, 2, 3, 4}, []int{3, 4}, 2)
	f([]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, 0)
	f([]int{}, []int{}, 2)
}

func TestChunk(t *testing.T) {
	f := func(input []int, expected [][]int, size int) {
		iter := ft.SliceIter(input)
		result := ft.Collect(ft.Chunk(iter, size))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3, 4, 5}, [][]int{{1, 2}, {3, 4}, {5}}, 2)
	f([]int{}, [][]int{}, 2)
	f([]int{1, 2, 3, 4}, [][]int{{1, 2}, {3, 4}}, 2)
	f([]int{1}, [][]int{{1}}, 5)
}

func TestScan(t *testing.T) {
	f := func(input []int, expected []string, initial ...string) {
		var i string
		if len(initial) > 0 {
			i = initial[0]
		}
		iter := ft.SliceIter(input)
		result := ft.Collect(ft.Scan(iter, func(o string, i int) string {
			return o + strconv.Itoa(i)
		}, i))
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3}, []string{"1", "12", "123"})
	f([]int{1, 2, 3}, []string{"__1", "__12", "__123"}, "__")
	f([]int{}, []string{}, "__")
}

func TestProduct(t *testing.T) {
	type output struct {
		i int
		s string
		f float64
	}
	f := func(input1 []int, input2 []string, input3 []float64, expected []output) {
		iter1 := ft.SliceIter(input1)
		iter2 := ft.SliceIter(input2)
		iter3 := ft.SliceIter(input3)
		result := make([]output, 0)
		for _, p := range ft.Collect(ft.Product(ft.Product(iter1, iter2), iter3)) {
			result = append(result, output{i: p.First.First, s: p.First.Second, f: p.Second})
		}
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2}, []string{"one", "two"}, []float64{1.1, 2.2}, []output{
		{1, "one", 1.1},
		{1, "one", 2.2},
		{1, "two", 1.1},
		{1, "two", 2.2},
		{2, "one", 1.1},
		{2, "one", 2.2},
		{2, "two", 1.1},
		{2, "two", 2.2},
	})
}

func TestCycle(t *testing.T) {
	slice := []int{1, 2, 3}
	iter := ft.Cycle(ft.SliceIter(slice))
	for i := 0; i < 5; i++ {
		next := iter.Next()
		assert.NotNil(t, next)
		assert.Equal(t, i%len(slice)+1, *next)
	}
}

func TestCycle_Empty(t *testing.T) {
	slice := []int{}
	iter := ft.Cycle(ft.SliceIter(slice))
	for i := 0; i < 5; i++ {
		next := iter.Next()
		assert.Nil(t, next)
	}
}

func TestZip(t *testing.T) {
	type output struct {
		i int
		s string
		f float64
	}
	f := func(input1 []int, input2 []string, input3 []float64, expected []output) {
		iter1 := ft.SliceIter(input1)
		iter2 := ft.SliceIter(input2)
		iter3 := ft.SliceIter(input3)
		result := make([]output, 0)
		for _, p := range ft.Collect(ft.Zip(ft.Zip(iter1, iter2), iter3)) {
			result = append(result, output{i: p.First.First, s: p.First.Second, f: p.Second})
		}
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2}, []string{"one", "two"}, []float64{1.1, 2.2}, []output{
		{1, "one", 1.1},
		{2, "two", 2.2},
	})
}

func TestAny(t *testing.T) {
	f := func(input []int, expected bool) {
		iter := ft.SliceIter(input)
		result := ft.Any(iter, func(t int) bool {
			return t%2 == 0
		})
		assert.Equal(t, expected, result)
	}
	f([]int{1, 1, 1}, false)
	f([]int{2, 2, 2}, true)
	f([]int{1, 1, 2}, true)
	f([]int{2, 1, 1}, true)
	f([]int{}, false)
}

func TestAll(t *testing.T) {
	f := func(input []int, expected bool) {
		iter := ft.SliceIter(input)
		result := ft.All(iter, func(t int) bool {
			return t%2 == 0
		})
		assert.Equal(t, expected, result)
	}
	f([]int{1, 1, 1}, false)
	f([]int{2, 2, 2}, true)
	f([]int{1, 1, 2}, false)
	f([]int{2, 1, 1}, false)
	f([]int{}, false)
}

func TestCount(t *testing.T) {
	f := func(input []int) {
		expected := len(input)
		iter := ft.SliceIter(input)
		result := ft.Count(iter)
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3})
	f([]int{})
}

func TestCount_Predicate(t *testing.T) {
	f := func(input []int, expected int) {
		iter := ft.SliceIter(input)
		result := ft.Count(iter, func(t int) bool {
			return t%2 == 0
		})
		assert.Equal(t, expected, result)
	}
	f([]int{1, 2, 3}, 1)
	f([]int{}, 0)
	f([]int{1, 1, 1}, 0)
	f([]int{2, 2, 2}, 3)
}

func testSum[T ft.Number](input []T, expected T, t *testing.T, initial ...T) {
	iter := ft.SliceIter(input)
	var result T
	if len(initial) > 0 {
		result = ft.Sum(iter, initial[0])
	} else {
		result = ft.Sum(iter)
	}
	assert.Equal(t, expected, result)
}

func TestSum(t *testing.T) {
	testSum([]int{1, 2, 3}, 6, t)
	testSum([]float64{1.1, 2.2, 3.3}, 6.6, t)
	testSum([]float64{}, 0, t)
	testSum([]int{}, 5, t, 5)
	testSum([]int{1, 2, 3}, 16, t, 10)
}

func TestForEach(t *testing.T) {
	f := func(input []int, fe func(int)) {
		iter := ft.SliceIter(input)
		ft.ForEach(iter, fe)
	}
	i := 0
	f([]int{1, 2, 3}, func(el int) {
		i += el
	})
	assert.Equal(t, 6, i)
}

func TestIntoChannel(t *testing.T) {
	f := func(input []int, ctx context.Context) {
		iter := ft.SliceIter(input)
		c := ft.IntoChannel(iter, ctx)
		result := make([]int, 0)
		for elem := range c {
			result = append(result, elem)
		}
		assert.Equal(t, input, result)
	}
	f([]int{1, 2, 3}, nil)
	f([]int{}, nil)
	f([]int{1}, nil)
}

func TestIntoChannel_Context(t *testing.T) {
	f := func(input []int, ctx context.Context) <-chan int {
		iter := ft.SliceIter(input)
		c := ft.IntoChannel(iter, ctx)
		return c
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := f([]int{1, 2, 3}, ctx)
	itersCount := 0
	for elem := range c {
		itersCount++
		assert.Equal(t, 1, elem)
		cancel()
	}
	assert.Equal(t, 1, itersCount)
	cancel()
}

func TestIntoChannel_WithValue(t *testing.T) {
	f := func(input []int, ctx context.Context) <-chan int {
		iter := ft.SliceIter(input)
		c := ft.IntoChannel(iter, ctx)
		return c
	}
	ctx := context.WithValue(context.Background(), "size", 3)
	c := f([]int{1, 2, 3}, ctx)
	runtime.Gosched()
	time.Sleep(10 * time.Millisecond) // give some time to spawned goroutine fill buffer
	assert.Equal(t, 3, len(c))        // all 3 values was writen into channel
	<-c
	assert.Equal(t, 2, len(c)) // remain 2 elements in channel
}

func TestEnumerate(t *testing.T) {
	f := func(input []int, start ...int) {
		iter := ft.SliceIter(input)
		var enumIter ft.Iter[ft.EnumeratePair[int]]
		var countFrom int
		if len(start) > 0 {
			countFrom = start[0]
			enumIter = ft.Enumerate(iter, start[0])
		} else {
			enumIter = ft.Enumerate(iter)
		}
		for i, p := range ft.Collect(enumIter) {
			assert.Equal(t, i+countFrom, p.Idx)
			assert.Equal(t, input[i], p.Value)
		}
	}
	f([]int{1, 2, 3})
	f([]int{})
	f([]int{1, 2, 3}, -10)
}

func TestMax(t *testing.T) {
	f := func(input []int, expected int) {
		iter := ft.SliceIter(input)
		result := ft.Max(iter, func(a, b int) bool {
			return a < b
		})
		if expected != 0 {
			assert.Equal(t, expected, *result)
		} else {
			assert.Nil(t, result)
		}
	}
	f([]int{5, 1, -20, 21, 3, 4}, 21)
	f([]int{}, 0)
	f([]int{5, 5, 5}, 5)
	f([]int{1, 5, 1}, 5)
	f([]int{2}, 2)
	f([]int{5, 1, 1}, 5)
	f([]int{1, 1, 5}, 5)
	f([]int{-20, -1, -5}, -1)
}

func TestMin(t *testing.T) {
	f := func(input []int, expected int) {
		iter := ft.SliceIter(input)
		result := ft.Min(iter, func(a, b int) bool {
			return a < b
		})
		if expected != 0 {
			assert.Equal(t, expected, *result)
		} else {
			assert.Nil(t, result)
		}
	}
	f([]int{5, 1, -20, 21, 3, 4}, -20)
	f([]int{}, 0)
	f([]int{5, 5, 5}, 5)
	f([]int{1, 5, 1}, 1)
	f([]int{2}, 2)
	f([]int{5, 1, 1}, 1)
	f([]int{1, 1, 5}, 1)
}

func TestSliceIter(t *testing.T) {
	f := func(input []int) {
		iter := ft.SliceIter(input)
		for i := 0; i < len(input); i++ {
			next := iter.Next()
			assert.NotNil(t, next)
			assert.Equal(t, input[i], *next)
		}
		next := iter.Next()
		assert.Nil(t, next)
	}
	f([]int{1, 2, 3, 4})
	f([]int{1})
	f([]int{})
}

func TestMapIter(t *testing.T) {
	f := func(input map[int]string) {
		iter := ft.MapIter(input)
		for range input {
			next := iter.Next()
			assert.NotNil(t, next)
			assert.Contains(t, input, next.Key)
			assert.Equal(t, input[next.Key], next.Value)
		}
		next := iter.Next()
		assert.Nil(t, next)
	}
	f(map[int]string{1: "1", 2: "2", 3: "3"})
	f(map[int]string{1: "1"})
	f(map[int]string{})
}

func TestMapIterOverSlice(t *testing.T) {
	f := func(input map[int]string) {
		iter := ft.MapIterOverSlice(input)
		for range input {
			next := iter.Next()
			assert.NotNil(t, next)
			assert.Contains(t, input, next.Key)
			assert.Equal(t, input[next.Key], next.Value)
		}
		next := iter.Next()
		assert.Nil(t, next)
	}
	f(map[int]string{1: "1", 2: "2", 3: "3"})
	f(map[int]string{1: "1"})
	f(map[int]string{})
}
