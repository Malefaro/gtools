package ft

type filterIter[T any] struct {
	iter Iter[T]
	f    func(T) bool
}

func (fi *filterIter[T]) Next() *T {
	for next := fi.iter.Next(); next != nil; next = fi.iter.Next() {
		if fi.f(*next) {
			return next
		}
	}
	return nil
}

func Filter[T any](iter Iter[T], f func(T) bool) Iter[T] {
	return &filterIter[T]{
		iter: iter,
		f:    f,
	}
}

type mapIter[T any, K any] struct {
	iter   Iter[T]
	mapper func(T) K
}

func (mi *mapIter[T, K]) Next() *K {
	next := mi.iter.Next()
	if next != nil {
		n := mi.mapper(*next)
		return &n
	}
	return nil
}

func Map[T any, K any](iter Iter[T], mapper func(T) K) Iter[K] {
	return &mapIter[T, K]{
		iter:   iter,
		mapper: mapper,
	}
}

type reverseIter[T any] struct {
	iter ReversibleIter[T]
}

func (fi *reverseIter[T]) Next() *T {
	return fi.iter.Prev()
}

func Reverse[T any](iter ReversibleIter[T]) Iter[T] {
	// go to last element
	for next := iter.Next(); next != nil; next = iter.Next() {
	}
	return &reverseIter[T]{
		iter: iter,
	}
}

func Skip[T any](iter Iter[T], num int) Iter[T] {
	for i := 0; i < num; i++ {
		next := iter.Next()
		if next == nil {
			break
		}
	}
	return iter
}

type chunkIter[T any, S ~[]T] struct {
	iter Iter[T]
	size int
}

func (ci *chunkIter[T, S]) Next() *S {
	s := make(S, 0, ci.size)
	next := ci.iter.Next()
	for i := 0; i < ci.size; i++ {
		if next == nil {
			break
		}
		s = append(s, *next)
		if i != ci.size-1 {
			// not last iteration
			next = ci.iter.Next()
		}
	}
	if len(s) == 0 {
		return nil
	}
	return &s
}

// Chunk split provided `iter` into several slices
// returns iterator of slices with len less or equal `size`
func Chunk[T any, S ~[]T](iter Iter[T], size int) Iter[S] {
	return &chunkIter[T, S]{
		iter: iter,
		size: size,
	}
}

type scanIter[T any, O any] struct {
	iter       Iter[T]
	f          func(O, T) O
	lastResult O
}

func (si *scanIter[T, O]) Next() *O {
	next := si.iter.Next()
	if next == nil {
		return nil
	}
	si.lastResult = si.f(si.lastResult, *next)
	return &si.lastResult
}

// Scan same as reduce but instead of returning one result it return iterator of results at every step
func Scan[T any, O any](iter Iter[T], f func(O, T) O, initial ...O) Iter[O] {
	si := &scanIter[T, O]{
		iter: iter,
		f:    f,
	}
	if len(initial) > 0 {
		si.lastResult = initial[0]
	}
	return si
}

type productIter[T ProductPair[F, S], F any, S any] struct {
	iter1         Iter[F]
	iter2         Iter[S]
	iter2elements []S
	iter2stored   bool
	prevFirst     *F
}

type ProductPair[F any, S any] struct {
	First  F
	Second S
}

func (pi *productIter[T, F, S]) Next() *T {
	if pi.prevFirst == nil {
		pi.prevFirst = pi.iter1.Next()
		if pi.prevFirst == nil {
			return nil
		}
	}
	s := pi.iter2.Next()
	if s == nil {
		pi.iter2stored = true
		pi.prevFirst = pi.iter1.Next()
		if pi.prevFirst == nil {
			return nil
		}
		pi.iter2 = SliceIter(pi.iter2elements)
		s = pi.iter2.Next()
		if s == nil {
			return nil
		}
	}
	if !pi.iter2stored {
		pi.iter2elements = append(pi.iter2elements, *s)
	}
	return &T{
		First:  *pi.prevFirst,
		Second: *s,
	}
}

// Product make cartesian product of input iters.
// if you need Product more than 2 iterator you can do following:
//	 iter1 := ft.SliceIter([]int{1, 2})
//	 iter2 := ft.SliceIter([]float64{1.1, 2.2})
//	 iter3 := ft.SliceIter([]string{"one", "two"})
//	 result := ft.Collect(ft.Product(ft.Product(iter1, iter2), iter3))
//	 for _, p1 := range result {
//	 	fmt.Printf("%v %v %v\n", p1.First.First, p1.First.Second, p1.Second)
//	 }
func Product[F any, S any](iter1 Iter[F], iter2 Iter[S]) Iter[ProductPair[F, S]] {
	return &productIter[ProductPair[F, S], F, S]{
		iter1:         iter1,
		iter2:         iter2,
		iter2stored:   false,
		iter2elements: make([]S, 0),
	}
}

type cycleIter[T any] struct {
	iter           Iter[T]
	savedElements  []T
	elementsStored bool
}

func (ci *cycleIter[T]) Next() *T {
	next := ci.iter.Next()
	if next == nil {
		ci.elementsStored = true
		ci.iter = SliceIter(ci.savedElements)
		next = ci.iter.Next()
		if next == nil {
			return nil
		}
	}
	if !ci.elementsStored {
		ci.savedElements = append(ci.savedElements, *next)
	}
	return next
}

// Cycle returns iteror that produces same elements as `iter`
// when `iter` ends, cycled iterator continue iterates from beginning
func Cycle[T any](iter Iter[T]) Iter[T] {
	return &cycleIter[T]{
		iter:           iter,
		savedElements:  make([]T, 0),
		elementsStored: false,
	}
}

type zipIter[T ZipPair[F, S], F any, S any] struct {
	iter1 Iter[F]
	iter2 Iter[S]
}

func (zi *zipIter[T, F, S]) Next() *T {
	next1 := zi.iter1.Next()
	next2 := zi.iter2.Next()
	if next1 != nil && next2 != nil {
		return &T{
			First:  *next1,
			Second: *next2,
		}
	}
	return nil
}

type ZipPair[F any, S any] struct {
	First  F
	Second S
}

// Zip create a new iterator over provided 2
// it ends when one of iter ends
// if you need zip more than 2 iterators you can do:
//	zip := Zip(Zip(iter1, iter2), iter3)
//	next := zip.Next()
// in this case:
// next.First.First -> iter1.Next()
// next.First.Second -> iter2.Next()
// next.Second -> iter3.Next()
func Zip[F any, S any](iter1 Iter[F], iter2 Iter[S]) Iter[ZipPair[F, S]] {
	return &zipIter[ZipPair[F, S], F, S]{
		iter1: iter1,
		iter2: iter2,
	}

}

type EnumeratePair[T any] struct {
	Idx   int
	Value T
}

type enumerateIter[T any, R EnumeratePair[T]] struct {
	iter Iter[T]
	idx  int
}

func (ei *enumerateIter[T, R]) Next() *R {
	next := ei.iter.Next()
	if next != nil {
		ei.idx++
		return &R{
			Idx:   ei.idx - 1,
			Value: *next,
		}
	}
	return nil
}

func Enumerate[T any, R EnumeratePair[T]](iter Iter[T], startFrom ...int) Iter[R] {
	var idx int
	if len(startFrom) > 0 {
		idx = startFrom[0]
	}
	return &enumerateIter[T, R]{
		iter: iter,
		idx:  idx,
	}
}
