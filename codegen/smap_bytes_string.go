package codegen

import "bytes"

// WARNING: using []byte as keys can lead to undefined behavior if the
// []byte are modified after insertion!!!
func (r SortedBytesToStringMap) compare(a, b []byte) int { return bytes.Compare(a, b) }

// SortedBytesToStringMap is a sorted map built on a left leaning red black balanced
// search sorted map. It stores string values, keyed by []byte.
type SortedBytesToStringMap struct {
	root *nodeBytesToString
}

// NewSortedBytesToStringMap creates a sorted map.
func NewSortedBytesToStringMap() *SortedBytesToStringMap { return &SortedBytesToStringMap{} }

// IsEmpty tells if the sorted map contains no key/value.
func (r SortedBytesToStringMap) IsEmpty() bool {
	return r.root == nil
}

// Size of the sorted map.
func (r SortedBytesToStringMap) Size() int { return r.root.size() }

// Clear all the values in the sorted map.
func (r *SortedBytesToStringMap) Clear() { r.root = nil }

// Put a value in the sorted map at key `k`. The old value at `k` is returned
// if the key was already present.
func (r *SortedBytesToStringMap) Put(k []byte, v string) (old string, overwrite bool) {
	r.root, old, overwrite = r.put(r.root, k, v)
	return
}

func (r *SortedBytesToStringMap) put(h *nodeBytesToString, k []byte, v string) (_ *nodeBytesToString, old string, overwrite bool) {
	if h == nil {
		n := &nodeBytesToString{key: k, val: v, n: 1, colorRed: true}
		return n, old, overwrite
	}

	cmp := r.compare(k, h.key)
	if cmp < 0 {
		h.left, old, overwrite = r.put(h.left, k, v)
	} else if cmp > 0 {
		h.right, old, overwrite = r.put(h.right, k, v)
	} else {
		overwrite = true
		old = h.val
		h.val = v
	}

	if h.right.isRed() && !h.left.isRed() {
		h = r.rotateLeft(h)
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = r.rotateRight(h)
	}
	if h.left.isRed() && h.right.isRed() {
		r.flipColors(h)
	}
	h.n = h.left.size() + h.right.size() + 1
	return h, old, overwrite
}

// Get a value from the sorted map at key `k`. Returns false
// if the key doesn't exist.
func (r SortedBytesToStringMap) Get(k []byte) (string, bool) {
	return r.loopGet(r.root, k)
}

func (r SortedBytesToStringMap) loopGet(h *nodeBytesToString, k []byte) (v string, ok bool) {
	for h != nil {
		cmp := r.compare(k, h.key)
		if cmp == 0 {
			return h.val, true
		} else if cmp < 0 {
			h = h.left
		} else if cmp > 0 {
			h = h.right
		}
	}
	return
}

// Has tells if a value exists at key `k`. This is short hand for `Get.
func (r SortedBytesToStringMap) Has(k []byte) bool {
	_, ok := r.loopGet(r.root, k)
	return ok
}

// Min returns the smallest key/value in the sorted map, if it exists.
func (r SortedBytesToStringMap) Min() (k []byte, v string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.min(r.root)
	return h.key, h.val, true
}

func (r SortedBytesToStringMap) min(x *nodeBytesToString) *nodeBytesToString {
	if x.left == nil {
		return x
	}
	return r.min(x.left)
}

// Max returns the largest key/value in the sorted map, if it exists.
func (r SortedBytesToStringMap) Max() (k []byte, v string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.max(r.root)
	return h.key, h.val, true
}

func (r SortedBytesToStringMap) max(x *nodeBytesToString) *nodeBytesToString {
	if x.right == nil {
		return x
	}
	return r.max(x.right)
}

// Floor returns the largest key/value in the sorted map that is smaller than
// `k`.
func (r SortedBytesToStringMap) Floor(key []byte) (k []byte, v string, ok bool) {
	x := r.floor(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedBytesToStringMap) floor(h *nodeBytesToString, k []byte) *nodeBytesToString {
	if h == nil {
		return nil
	}
	cmp := r.compare(k, h.key)
	if cmp == 0 {
		return h
	}
	if cmp < 0 {
		return r.floor(h.left, k)
	}
	t := r.floor(h.right, k)
	if t != nil {
		return t
	}
	return h
}

// Ceiling returns the smallest key/value in the sorted map that is larger than
// `k`.
func (r SortedBytesToStringMap) Ceiling(key []byte) (k []byte, v string, ok bool) {
	x := r.ceiling(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedBytesToStringMap) ceiling(h *nodeBytesToString, k []byte) *nodeBytesToString {
	if h == nil {
		return nil
	}
	cmp := r.compare(k, h.key)
	if cmp == 0 {
		return h
	}
	if cmp > 0 {
		return r.ceiling(h.right, k)
	}
	t := r.ceiling(h.left, k)
	if t != nil {
		return t
	}
	return h
}

// Select key of rank k, meaning the k-th biggest []byte in the sorted map.
func (r SortedBytesToStringMap) Select(key int) (k []byte, v string, ok bool) {
	x := r.nodeselect(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedBytesToStringMap) nodeselect(x *nodeBytesToString, k int) *nodeBytesToString {
	if x == nil {
		return nil
	}
	t := x.left.size()
	if t > k {
		return r.nodeselect(x.left, k)
	} else if t < k {
		return r.nodeselect(x.right, k-t-1)
	} else {
		return x
	}
}

// Rank is the number of keys less than `k`.
func (r SortedBytesToStringMap) Rank(k []byte) int {
	return r.keyrank(k, r.root)
}

func (r SortedBytesToStringMap) keyrank(k []byte, h *nodeBytesToString) int {
	if h == nil {
		return 0
	}
	cmp := r.compare(k, h.key)
	if cmp < 0 {
		return r.keyrank(k, h.left)
	} else if cmp > 0 {
		return 1 + h.left.size() + r.keyrank(k, h.right)
	} else {
		return h.left.size()
	}
}

// Keys visit each keys in the sorted map, in order.
// It stops when visit returns false.
func (r SortedBytesToStringMap) Keys(visit func([]byte, string) bool) {
	min, _, ok := r.Min()
	if !ok {
		return
	}
	// if the min exists, then the max must exist
	max, _, _ := r.Max()
	r.RangedKeys(min, max, visit)
}

// RangedKeys visit each keys between lo and hi in the sorted map, in order.
// It stops when visit returns false.
func (r SortedBytesToStringMap) RangedKeys(lo, hi []byte, visit func([]byte, string) bool) {
	r.keys(r.root, visit, lo, hi)
}

func (r SortedBytesToStringMap) keys(h *nodeBytesToString, visit func([]byte, string) bool, lo, hi []byte) bool {
	if h == nil {
		return true
	}
	cmplo := r.compare(lo, h.key)
	cmphi := r.compare(hi, h.key)
	if cmplo < 0 {
		if !r.keys(h.left, visit, lo, hi) {
			return false
		}
	}
	if cmplo <= 0 && cmphi >= 0 {
		if !visit(h.key, h.val) {
			return false
		}
	}
	if cmphi > 0 {
		if !r.keys(h.right, visit, lo, hi) {
			return false
		}
	}
	return true
}

// DeleteMin removes the smallest key and its value from the sorted map.
func (r *SortedBytesToStringMap) DeleteMin() (oldk []byte, oldv string, ok bool) {
	r.root, oldk, oldv, ok = r.deleteMin(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedBytesToStringMap) deleteMin(h *nodeBytesToString) (_ *nodeBytesToString, oldk []byte, oldv string, ok bool) {
	if h == nil {
		return nil, oldk, oldv, false
	}

	if h.left == nil {
		return nil, h.key, h.val, true
	}
	if !h.left.isRed() && !h.left.left.isRed() {
		h = r.moveRedLeft(h)
	}
	h.left, oldk, oldv, ok = r.deleteMin(h.left)
	return r.balance(h), oldk, oldv, ok
}

// DeleteMax removes the largest key and its value from the sorted map.
func (r *SortedBytesToStringMap) DeleteMax() (oldk []byte, oldv string, ok bool) {
	r.root, oldk, oldv, ok = r.deleteMax(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedBytesToStringMap) deleteMax(h *nodeBytesToString) (_ *nodeBytesToString, oldk []byte, oldv string, ok bool) {
	if h == nil {
		return nil, oldk, oldv, ok
	}
	if h.left.isRed() {
		h = r.rotateRight(h)
	}
	if h.right == nil {
		return nil, h.key, h.val, true
	}
	if !h.right.isRed() && !h.right.left.isRed() {
		h = r.moveRedRight(h)
	}
	h.right, oldk, oldv, ok = r.deleteMax(h.right)
	return r.balance(h), oldk, oldv, ok
}

// Delete key `k` from sorted map, if it exists.
func (r *SortedBytesToStringMap) Delete(k []byte) (old string, ok bool) {
	if r.root == nil {
		return
	}
	r.root, old, ok = r.delete(r.root, k)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedBytesToStringMap) delete(h *nodeBytesToString, k []byte) (_ *nodeBytesToString, old string, ok bool) {

	if h == nil {
		return h, old, false
	}

	if r.compare(k, h.key) < 0 {
		if h.left == nil {
			return h, old, false
		}

		if !h.left.isRed() && !h.left.left.isRed() {
			h = r.moveRedLeft(h)
		}

		h.left, old, ok = r.delete(h.left, k)
		h = r.balance(h)
		return h, old, ok
	}

	if h.left.isRed() {
		h = r.rotateRight(h)
	}

	if r.compare(k, h.key) == 0 && h.right == nil {
		return nil, h.val, true
	}

	if h.right != nil && !h.right.isRed() && !h.right.left.isRed() {
		h = r.moveRedRight(h)
	}

	if r.compare(k, h.key) == 0 {

		var subk []byte
		var subv string
		h.right, subk, subv, ok = r.deleteMin(h.right)

		old, h.key, h.val = h.val, subk, subv
		ok = true
	} else {
		h.right, old, ok = r.delete(h.right, k)
	}

	h = r.balance(h)
	return h, old, ok
}

// deletions

func (r *SortedBytesToStringMap) moveRedLeft(h *nodeBytesToString) *nodeBytesToString {
	r.flipColors(h)
	if h.right.left.isRed() {
		h.right = r.rotateRight(h.right)
		h = r.rotateLeft(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedBytesToStringMap) moveRedRight(h *nodeBytesToString) *nodeBytesToString {
	r.flipColors(h)
	if h.left.left.isRed() {
		h = r.rotateRight(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedBytesToStringMap) balance(h *nodeBytesToString) *nodeBytesToString {
	if h.right.isRed() {
		h = r.rotateLeft(h)
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = r.rotateRight(h)
	}
	if h.left.isRed() && h.right.isRed() {
		r.flipColors(h)
	}
	h.n = h.left.size() + h.right.size() + 1
	return h
}

func (r *SortedBytesToStringMap) rotateLeft(h *nodeBytesToString) *nodeBytesToString {
	x := h.right
	h.right = x.left
	x.left = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedBytesToStringMap) rotateRight(h *nodeBytesToString) *nodeBytesToString {
	x := h.left
	h.left = x.right
	x.right = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedBytesToStringMap) flipColors(h *nodeBytesToString) {
	h.colorRed = !h.colorRed
	h.left.colorRed = !h.left.colorRed
	h.right.colorRed = !h.right.colorRed
}

// nodes

type nodeBytesToString struct {
	key         []byte
	val         string
	left, right *nodeBytesToString
	n           int
	colorRed    bool
}

func (x *nodeBytesToString) isRed() bool { return (x != nil) && (x.colorRed == true) }

func (x *nodeBytesToString) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

