package misc

import (
	"errors"
	"fmt"
	"math"
)

// Argument: two sorted arrays and position of the sorted array we seek
// Returns: index of the pos-th position when array1 and array2 are appended
func ArrayPos(arr1 []int, arr2 []int, pos int) (int, error) {
	if arr1 == nil || arr2 == nil {
		return 0, errors.New("Empty array slices given as argument")
	}
	if len(arr1)+len(arr2) < pos {
		return 0, errors.New("Seeking position greater than array slices lengths")
	}
	r1, err := get_pos(arr1, arr2, pos)
	return r1, err
}

// Argument: two sorted integer slices
// Returns: float representing median
func Median(a []int, b []int) (float32, error) {
	if a == nil || b == nil {
		return 0, errors.New("ArrayPos of nil slices called!")
	}
	r1, err := ArrayPos(a, b, (len(a)+len(b)+1)/2)
	if err != nil {
		return 0.0, err
	}

	m, err2 := get_mean(a, b, r1)
	if err2 != nil {
		fmt.Printf("Median of two sorted arrays %v and %v returned error %s\n",
			a, b, err2)
		return 0.0, err2
	}

	fmt.Printf("Median of two sorted arrays %v and %v is %.1f\n",
		a, b, m)
	return m, err2
}

// Argument: two sorted integer slices and index of first slice
// that is candidate for median
func get_mean(arr1 []int, arr2 []int, r1 int) (float32, error) {
	len1 := len(arr1)
	len2 := len(arr2)

	m1 := (len1 + len2 + 1) / 2
	m2 := (len1 + len2 + 2) / 2
	r2 := m1 - r1

	if r1 > len1 || r2 > len2 {
		fmt.Printf("r1=%d, r2=%d, m1=%d, m2=%d, len1=%d, len2=%d\n",
			r1, r2, m1, m2, len1, len2)
		return 0, errors.New("median idx calculation error")
	}
	val1 := math.MinInt32
	val2 := math.MinInt32
	if r1 > 0 {
		val1 = arr1[r1-1]
	}
	if r2 > 0 {
		val2 = arr2[r2-1]
	}

	nval1 := math.MaxInt32
	nval2 := math.MaxInt32
	if r1 < len1 {
		nval1 = arr1[r1]
	}
	if r2 < len2 {
		nval2 = arr2[r2]
	}
	val := val1
	if val < val2 {
		val = val2
	}
	nval := nval1
	if nval > nval2 {
		nval = nval2
	}
	// total elements may be odd - we are already at median idx
	if m1 == m2 {
		return float32(val), nil
	}

	return 0.5 * float32(val+nval), nil
}

func get_pos(a []int, b []int, pos int) (int, error) {
	len1 := len(a)
	len2 := len(b)
	if len1 == 0 {
		return 0, nil
	}
	if len2 == 0 {
		return pos, nil
	}

	a1 := 0
	a2 := len1
	max_jump := pos / 2

	// pos search is done in the two arrays a[a1+1:a2), and b[b1+1:b2)
	for i := 0; i < len1; i++ {
		p1 := (a1 + a2) / 2
		if p1 > max_jump {
			p1 = max_jump
		}
		p2 := pos - p1
		// fmt.Printf("Iter %d: a1=%d, a2=%d, p1=%d, p2=%d\n", i, a1, a2, p1, p2)

		if a1 >= a2 {
			break
		}

		if p2 > len2 {
			// we are taking too short jump in arr1 - go higher in arr1
			a1 = p1 + 1
			continue
		}

		val1 := math.MinInt32
		val2 := math.MaxInt32
		nval1 := math.MaxInt32
		nval2 := math.MaxInt32
		if p1 > 0 {
			val1 = a[p1-1]
		}
		if p1 < len1 {
			nval1 = a[p1]
		}
		if p2 > 0 {
			val2 = b[p2-1]
		}
		if p2 < len2 {
			nval2 = b[p2]
		}

		v := val1
		if val2 > v {
			v = val2
		}

		nv := nval1
		if nval2 < nv {
			nv = nval2
		}

		if v <= nv {
			// fmt.Printf("Found: v(%d) <= nv(%d), p1=%d\n", v, nv, p1)
			return p1, nil
		}

		if v > val1 {
			// we have gone too far - go lower in arr1
			a2 = p1 - 1
			continue
		}
		// we have not gone far enough - go higher in arr1
		a1 = p1 + 1
	}

	fmt.Printf("Lost: Could not locate:a=%v, b=%v, pos=%d, a1=%d, a2=%d\n",
		a, b, pos, a1, a2)

	return 0, errors.New("Lost: pos not found in slices")
}
