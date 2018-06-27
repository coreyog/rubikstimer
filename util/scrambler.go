package util

// Migrated, slightly modified from www.jaapsch.net/scramble_cube.htm

import (
	"math"
	"math/rand"
)

// Default settings
var seqlen = 20

// Scramble returns a scramble string
func Scramble() string {
	var seq = []int{} // move sequences
	//tl=number of allowed moves (twistable layers) on axis -- middle layer ignored
	tl := 2 // size - 1 but size is always 3

	//set up bookkeeping
	axsl := make([]int, tl)
	axam := []int{0, 0, 0}
	la := -1
	for len(seq) < seqlen {
		if len(seq) == 21 {
			// fmt.Println("DEBUG")
		}
		// choose a different axis than previous one
		ax := int(math.Floor(rand.Float64() * 3))
		for ax == la {
			ax = int(math.Floor(rand.Float64() * 3))
		}

		// reset slice/direction counters
		for i := 0; i < tl; i++ {
			axsl[i] = 0
		}
		axam = []int{0, 0, 0}
		moved := 0
		// generate moves on this axis
		thirdProb := 0        // force it to work so this for acts as do...while
		for thirdProb == 0 && // 2/3 prob for other axis next
			moved < tl && // must change if all layers moved
			moved+len(seq) < seqlen { // must change if done enough moves
			// choose random unmoved slice
			sl := int(math.Floor(rand.Float64() * float64(tl)))
			for axsl[sl] != 0 {
				sl = int(math.Floor(rand.Float64() * float64(tl)))
			}
			// choose random amount
			q := int(math.Floor(rand.Float64() * 3))
			if tl != 3 || // odd cube always ok since middle layer is reference
				(axam[q]+1)*2 < tl || // less than half the slices in same direction also ok
				((axam[q]+1)*2 == tl && axam[0]+axam[1]+axam[2]-axam[q] == 0) { // exactly half the slices move in same direction and no other slice moved
				axam[q]++ // adjust direction count
				moved++
				axsl[sl] = q + 1 // mark the slice has moved amount
			}
			thirdProb = int(math.Floor(rand.Float64() * 3))
		}

		// append these moves to current sequence in order
		for sl := 0; sl < tl; sl++ {
			if axsl[sl] != 0 {
				q := axsl[sl] - 1
				// get semi-axis of this move
				sa := ax
				m := sl
				if sl+sl+1 >= tl { // if on other half of this axis
					sa += 3        // get semi-axis (i.e. face)
					m = tl - 1 - m // slice number counting from that face
					q = 2 - q      // opposite direction when looking at that face
				}
				// store move
				seq = append(seq, (m*6+sa)*4+q)
			}
		}
		// avoid this axis next time
		la = ax
	}
	seq = append(seq, 0)

	return scrambleString(seq)
}

func scrambleString(seq []int) string {
	s := ""
	j := 0
	for i := 0; i < len(seq)-1; i++ {
		if i != 0 {
			s += " "
		}
		k := seq[i] >> 2
		s += string("DLBURFdlburf"[k])
		j = seq[i] & 3
		if j != 0 {
			s += string(" 2'"[j])
		}
	}
	return s
}
