package codegen


func (r SortedStringSet) compare(a, b string) int {
    if a < b {
        return -1
    }
    if a > b {
        return 1
    }
    return 0
}

// SortedStringSet is a sorted set built on a left leaning red black balanced
// search sorted set. It stores unique string values.
type SortedStringSet struct {
	root *nodeString
}

// NewSortedStringSet creates a sorted set.
func NewSortedStringSet() *SortedStringSet { return &SortedStringSet{} }

// IsEmpty tells if the sorted set contains no key.
func (r SortedStringSet) IsEmpty() bool {
	return r.root == nil
}

// Size of the sorted set.
func (r SortedStringSet) Size() int { return r.root.size() }

// Clear all the values in the sorted set.
func (r *SortedStringSet) Clear() { r.root = nil }

// Put the key `k` in the sorted set. If the value was already there,
// true is returned.
func (r *SortedStringSet) Put(k string) (already bool) {
	r.root, already = r.put(r.root, k)
	return
}

func (r *SortedStringSet) put(h *nodeString, k string) (_ *nodeString, already bool) {
	if h == nil {
		n := &nodeString{key: k, n: 1, colorRed: true}
		return n, already
	}

	cmp := r.compare(k, h.key)
	if cmp < 0 {
		h.left, already = r.put(h.left, k)
	} else if cmp > 0 {
		h.right, already = r.put(h.right, k)
	} else {
		already = true
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
	return h, already
}

// Contains tells if `k` is a member of the set.
func (r SortedStringSet) Contains(k string) bool {
	return r.loopContains(r.root, k)
}

func (r SortedStringSet) loopContains(h *nodeString, k string) (ok bool) {
	for h != nil {
		cmp := r.compare(k, h.key)
		if cmp == 0 {
			return true
		} else if cmp < 0 {
			h = h.left
		} else if cmp > 0 {
			h = h.right
		}
	}
	return
}

// Min returns the smallest key in the sorted set, if it exists.
func (r SortedStringSet) Min() (k string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.min(r.root)
	return h.key, true
}

func (r SortedStringSet) min(x *nodeString) *nodeString {
	if x.left == nil {
		return x
	}
	return r.min(x.left)
}

// Max returns the largest key in the sorted set, if it exists.
func (r SortedStringSet) Max() (k string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.max(r.root)
	return h.key, true
}

func (r SortedStringSet) max(x *nodeString) *nodeString {
	if x.right == nil {
		return x
	}
	return r.max(x.right)
}

// Floor returns the largest key in the sorted set that is smaller than
// `k`.
func (r SortedStringSet) Floor(key string) (k string, ok bool) {
	x := r.floor(r.root, key)
	if x == nil {
		return
	}
	return x.key, true
}

func (r SortedStringSet) floor(h *nodeString, k string) *nodeString {
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

// Ceiling returns the smallest key in the sorted set that is larger than
// `k`.
func (r SortedStringSet) Ceiling(key string) (k string, ok bool) {
	x := r.ceiling(r.root, key)
	if x == nil {
		return
	}
	return x.key, true
}

func (r SortedStringSet) ceiling(h *nodeString, k string) *nodeString {
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

// Select key of rank k, meaning the k-th biggest string in the sorted set.
func (r SortedStringSet) Select(key int) (k string, ok bool) {
	x := r.nodeselect(r.root, key)
	if x == nil {
		return
	}
	return x.key, true
}

func (r SortedStringSet) nodeselect(x *nodeString, k int) *nodeString {
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
func (r SortedStringSet) Rank(k string) int {
	return r.keyrank(k, r.root)
}

func (r SortedStringSet) keyrank(k string, h *nodeString) int {
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

// Keys visit each keys in the sorted set, in order.
// It stops when visit returns false.
func (r SortedStringSet) Keys(visit func(string) bool) {
	min, ok := r.Min()
	if !ok {
		return
	}
	// if the min exists, then the max must exist
	max, _ := r.Max()
	r.RangedKeys(min, max, visit)
}

// RangedKeys visit each keys between lo and hi in the sorted set, in order.
// It stops when visit returns false.
func (r SortedStringSet) RangedKeys(lo, hi string, visit func(string) bool) {
	r.keys(r.root, visit, lo, hi)
}

func (r SortedStringSet) keys(h *nodeString, visit func(string) bool, lo, hi string) bool {
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
		if !visit(h.key) {
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

// DeleteMin removes the smallest key from the sorted set.
func (r *SortedStringSet) DeleteMin() (oldk string, ok bool) {
	r.root, oldk, ok = r.deleteMin(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringSet) deleteMin(h *nodeString) (_ *nodeString, oldk string, ok bool) {
	if h == nil {
		return nil, oldk, false
	}

	if h.left == nil {
		return nil, h.key, true
	}
	if !h.left.isRed() && !h.left.left.isRed() {
		h = r.moveRedLeft(h)
	}
	h.left, oldk, ok = r.deleteMin(h.left)
	return r.balance(h), oldk, ok
}

// DeleteMax removes the largest key from the sorted set.
func (r *SortedStringSet) DeleteMax() (oldk string, ok bool) {
	r.root, oldk, ok = r.deleteMax(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringSet) deleteMax(h *nodeString) (_ *nodeString, oldk string, ok bool) {
	if h == nil {
		return nil, oldk, ok
	}
	if h.left.isRed() {
		h = r.rotateRight(h)
	}
	if h.right == nil {
		return nil, h.key, true
	}
	if !h.right.isRed() && !h.right.left.isRed() {
		h = r.moveRedRight(h)
	}
	h.right, oldk, ok = r.deleteMax(h.right)
	return r.balance(h), oldk, ok
}

// Delete key `k` from sorted set, if it exists.
func (r *SortedStringSet) Delete(k string) (ok bool) {
	if r.root == nil {
		return
	}
	r.root, ok = r.delete(r.root, k)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringSet) delete(h *nodeString, k string) (_ *nodeString, ok bool) {

	if h == nil {
		return h, false
	}

	if r.compare(k, h.key) < 0 {
		if h.left == nil {
			return h, false
		}

		if !h.left.isRed() && !h.left.left.isRed() {
			h = r.moveRedLeft(h)
		}

		h.left, ok = r.delete(h.left, k)
		h = r.balance(h)
		return h, ok
	}

	if h.left.isRed() {
		h = r.rotateRight(h)
	}

	if r.compare(k, h.key) == 0 && h.right == nil {
		return nil, true
	}

	if h.right != nil && !h.right.isRed() && !h.right.left.isRed() {
		h = r.moveRedRight(h)
	}

	if r.compare(k, h.key) == 0 {

		var subk string
		h.right, subk, ok = r.deleteMin(h.right)
		h.key = subk
		ok = true
	} else {
		h.right, ok = r.delete(h.right, k)
	}

	h = r.balance(h)
	return h, ok
}

// deletions

func (r *SortedStringSet) moveRedLeft(h *nodeString) *nodeString {
	r.flipColors(h)
	if h.right.left.isRed() {
		h.right = r.rotateRight(h.right)
		h = r.rotateLeft(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedStringSet) moveRedRight(h *nodeString) *nodeString {
	r.flipColors(h)
	if h.left.left.isRed() {
		h = r.rotateRight(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedStringSet) balance(h *nodeString) *nodeString {
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

func (r *SortedStringSet) rotateLeft(h *nodeString) *nodeString {
	x := h.right
	h.right = x.left
	x.left = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedStringSet) rotateRight(h *nodeString) *nodeString {
	x := h.left
	h.left = x.right
	x.right = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedStringSet) flipColors(h *nodeString) {
	h.colorRed = !h.colorRed
	h.left.colorRed = !h.left.colorRed
	h.right.colorRed = !h.right.colorRed
}

// nodes

type nodeString struct {
	key         string
	left, right *nodeString
	n           int
	colorRed    bool
}

func (x *nodeString) isRed() bool { return (x != nil) && (x.colorRed == true) }

func (x *nodeString) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

