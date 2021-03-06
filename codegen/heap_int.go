package codegen

// Most of the implementation is adapted from Algorithms 4ed by Sedgewick
// and Wayne.

// Comments are adapted from `container/heap`.
// 	 Copyright 2009 The Go Authors. All rights reserved.
// 	 Use of this source code is governed by a BSD-style
// 	 license that can be found in the LICENSE file.

func (h IntHeap) compare(a, b int) int { return int(a) - int(b) }

// IntHeap is a container of int, where the elements can be efficiently
// retrieved in their decreasing order (according to their comparison
// rules).
type IntHeap struct {
	n  int
	pq []int
}

// NewIntHeap creates a heap, optionaly with keys already populating
// it. The complexity is O(n) where n = len(keys).
func NewIntHeap(keys ...int) *IntHeap {
	h := &IntHeap{
		n:  len(keys),
		pq: append(make([]int, 1), keys...),
	}
	h.Fix()
	return h
}

// Len is the number of elements stored in the heap.
func (h *IntHeap) Len() int { return h.n }

// Peek at the largest element (according to their comparison rules), without
// removing it from the heap.
func (h *IntHeap) Peek() int { return h.pq[1] }

// Fix re-establishes the heap ordering. This is useful if elements
// of the heap have had their comparison value changed. It is equivalent to,
// but less expenasive than, Pop'ing all the elements and Push'ing them
// again.
// The complexity is O(n).
func (h *IntHeap) Fix() {
	for i := (h.n) / 2; i > 0; i-- {
		h.sink(i, h.n)
	}
}

// Push pushes the element k onto the heap. The complexity is
// O(log(n)) where n == h.Len().
func (h *IntHeap) Push(k int) {
	h.n++
	h.pq = append(h.pq, k)
	h.swim(h.n)
}

// Pop removes the largest element (according to their comparison rules) from
// the heap and returns it. The complexity is O(log(n)) where n == h.Len().
func (h *IntHeap) Pop() int {
	val := h.pq[1]
	h.swap(1, h.n)
	h.pq = h.pq[:h.n]
	h.n--
	h.sink(1, h.n)

	return val
}

// Remove removes k from the heap, if it exists. Equality is defined by
// Compare == 0.
// The complexity is O(n+log(n)) where n == h.Len().
func (h *IntHeap) Remove(k int) bool {

	cmp := h.compare(h.pq[1], k)
	if cmp == 0 {
		_ = h.Pop()
		return true
	}
	if cmp < 0 {
		// larger than largest, don't try to find it
		return false
	}

	i := 0
	for _, j := range h.pq[1:] {
		i++
		if h.compare(j, k) != 0 {
			continue
		}
		h.swap(i, 1)
		h.sink(i, h.n)
		_ = h.Pop()
		return true
	}
	// not in the heap
	return false
}

func (h *IntHeap) swap(i, j int)      { h.pq[i], h.pq[j] = h.pq[j], h.pq[i] }
func (h *IntHeap) less(i, j int) bool { return h.compare(h.pq[i], h.pq[j]) < 0 }

func (h *IntHeap) swim(k int) {
	for k > 1 && h.less(k/2, k) {
		h.swap(k/2, k)
		k = k / 2
	}
}

func (h *IntHeap) sink(k, n int) {

	for k*2 <= n {
		j := 2 * k
		if j < n && h.less(j, j+1) {
			j++
		}
		if !h.less(k, j) {
			break
		}
		h.swap(k, j)
		k = j
	}
}

