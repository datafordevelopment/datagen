package codegen


func (r SortedStringToStringMap) compare(a, b string) int {
    if a < b {
        return -1
    }
    if a > b {
        return 1
    }
    return 0
}

// SortedStringToStringMap is a sorted map built on a left leaning red black balanced
// search sorted map. It stores string values, keyed by string.
type SortedStringToStringMap struct {
	root *nodeStringToString
}

// NewSortedStringToStringMap creates a sorted map.
func NewSortedStringToStringMap() *SortedStringToStringMap { return &SortedStringToStringMap{} }

// IsEmpty tells if the sorted map contains no key/value.
func (r SortedStringToStringMap) IsEmpty() bool {
	return r.root == nil
}

// Size of the sorted map.
func (r SortedStringToStringMap) Size() int { return r.root.size() }

// Clear all the values in the sorted map.
func (r *SortedStringToStringMap) Clear() { r.root = nil }

// Put a value in the sorted map at key `k`. The old value at `k` is returned
// if the key was already present.
func (r *SortedStringToStringMap) Put(k string, v string) (old string, overwrite bool) {
	r.root, old, overwrite = r.put(r.root, k, v)
	return
}

func (r *SortedStringToStringMap) put(h *nodeStringToString, k string, v string) (_ *nodeStringToString, old string, overwrite bool) {
	if h == nil {
		n := &nodeStringToString{key: k, val: v, n: 1, colorRed: true}
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
func (r SortedStringToStringMap) Get(k string) (string, bool) {
	return r.loopGet(r.root, k)
}

func (r SortedStringToStringMap) loopGet(h *nodeStringToString, k string) (v string, ok bool) {
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
func (r SortedStringToStringMap) Has(k string) bool {
	_, ok := r.loopGet(r.root, k)
	return ok
}

// Min returns the smallest key/value in the sorted map, if it exists.
func (r SortedStringToStringMap) Min() (k string, v string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.min(r.root)
	return h.key, h.val, true
}

func (r SortedStringToStringMap) min(x *nodeStringToString) *nodeStringToString {
	if x.left == nil {
		return x
	}
	return r.min(x.left)
}

// Max returns the largest key/value in the sorted map, if it exists.
func (r SortedStringToStringMap) Max() (k string, v string, ok bool) {
	if r.root == nil {
		return
	}
	h := r.max(r.root)
	return h.key, h.val, true
}

func (r SortedStringToStringMap) max(x *nodeStringToString) *nodeStringToString {
	if x.right == nil {
		return x
	}
	return r.max(x.right)
}

// Floor returns the largest key/value in the sorted map that is smaller than
// `k`.
func (r SortedStringToStringMap) Floor(key string) (k string, v string, ok bool) {
	x := r.floor(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedStringToStringMap) floor(h *nodeStringToString, k string) *nodeStringToString {
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
func (r SortedStringToStringMap) Ceiling(key string) (k string, v string, ok bool) {
	x := r.ceiling(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedStringToStringMap) ceiling(h *nodeStringToString, k string) *nodeStringToString {
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

// Select key of rank k, meaning the k-th biggest string in the sorted map.
func (r SortedStringToStringMap) Select(key int) (k string, v string, ok bool) {
	x := r.nodeselect(r.root, key)
	if x == nil {
		return
	}
	return x.key, x.val, true
}

func (r SortedStringToStringMap) nodeselect(x *nodeStringToString, k int) *nodeStringToString {
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
func (r SortedStringToStringMap) Rank(k string) int {
	return r.keyrank(k, r.root)
}

func (r SortedStringToStringMap) keyrank(k string, h *nodeStringToString) int {
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
func (r SortedStringToStringMap) Keys(visit func(string, string) bool) {
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
func (r SortedStringToStringMap) RangedKeys(lo, hi string, visit func(string, string) bool) {
	r.keys(r.root, visit, lo, hi)
}

func (r SortedStringToStringMap) keys(h *nodeStringToString, visit func(string, string) bool, lo, hi string) bool {
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
func (r *SortedStringToStringMap) DeleteMin() (oldk string, oldv string, ok bool) {
	r.root, oldk, oldv, ok = r.deleteMin(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringToStringMap) deleteMin(h *nodeStringToString) (_ *nodeStringToString, oldk string, oldv string, ok bool) {
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
func (r *SortedStringToStringMap) DeleteMax() (oldk string, oldv string, ok bool) {
	r.root, oldk, oldv, ok = r.deleteMax(r.root)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringToStringMap) deleteMax(h *nodeStringToString) (_ *nodeStringToString, oldk string, oldv string, ok bool) {
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
func (r *SortedStringToStringMap) Delete(k string) (old string, ok bool) {
	if r.root == nil {
		return
	}
	r.root, old, ok = r.delete(r.root, k)
	if !r.IsEmpty() {
		r.root.colorRed = false
	}
	return
}

func (r *SortedStringToStringMap) delete(h *nodeStringToString, k string) (_ *nodeStringToString, old string, ok bool) {

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

		var subk string
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

func (r *SortedStringToStringMap) moveRedLeft(h *nodeStringToString) *nodeStringToString {
	r.flipColors(h)
	if h.right.left.isRed() {
		h.right = r.rotateRight(h.right)
		h = r.rotateLeft(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedStringToStringMap) moveRedRight(h *nodeStringToString) *nodeStringToString {
	r.flipColors(h)
	if h.left.left.isRed() {
		h = r.rotateRight(h)
		r.flipColors(h)
	}
	return h
}

func (r *SortedStringToStringMap) balance(h *nodeStringToString) *nodeStringToString {
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

func (r *SortedStringToStringMap) rotateLeft(h *nodeStringToString) *nodeStringToString {
	x := h.right
	h.right = x.left
	x.left = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedStringToStringMap) rotateRight(h *nodeStringToString) *nodeStringToString {
	x := h.left
	h.left = x.right
	x.right = h
	x.colorRed = h.colorRed
	h.colorRed = true
	x.n = h.n
	h.n = 1 + h.left.size() + h.right.size()
	return x
}

func (r *SortedStringToStringMap) flipColors(h *nodeStringToString) {
	h.colorRed = !h.colorRed
	h.left.colorRed = !h.left.colorRed
	h.right.colorRed = !h.right.colorRed
}

// nodes

type nodeStringToString struct {
	key         string
	val         string
	left, right *nodeStringToString
	n           int
	colorRed    bool
}

func (x *nodeStringToString) isRed() bool { return (x != nil) && (x.colorRed == true) }

func (x *nodeStringToString) size() int {
	if x == nil {
		return 0
	}
	return x.n
}

