package funcs

func maxU32(a, b uint32) uint32 {
	if a > b {
		return a
	} else {
		return b
	}
}
func triMaxU32(a, b, c uint32) uint32 {
	return maxU32(maxU32(a, b), c)
}
func minU32(a, b uint32) uint32 {
	if a < b {
		return a
	} else {
		return b
	}
}
func triMinU32(a, b, c uint32) uint32 {
	return minU32(minU32(a, b), c)
}
func maxInt(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
func triMaxInt(a, b, c int) int {
	return maxInt(maxInt(a, b), c)
}
func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
func triMinInt(a, b, c int) int {
	return minInt(minInt(a, b), c)
}
