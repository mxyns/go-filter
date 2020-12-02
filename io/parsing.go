package io

import (
	"image/color"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const (
	BYTE_TO_U16     = math.MaxUint16 / 0xff // x [0; 255] * BYTE_TO_U16 = x [0; 1 << 16 - 1]
	BYTE_TO_U32     = math.MaxUint32 / 0xff // x [0; 255] * BYTE_TO_U32 = x [0; 1 << 32 - 1]
	APC_BYTE_TO_U16 = BYTE_TO_U16 / 0xff    // x [0; 255] * alpha * APC_BYTE_TO_U16 = alpha-premultiplied x [0; 1 << 16 - 1]
	APC_BYTE_TO_U32 = BYTE_TO_U32 / 0xff    // x [0; 255] * alpha * APC_BYTE_TO_U32 = alpha-premultiplied x [0; 1 << 32 - 1]
	APC_U16_TO_BYTE = 1 / APC_BYTE_TO_U16   // APC_BYTE_TO_U16 inverse
	APC_U32_TO_BYTE = 1 / APC_BYTE_TO_U32   // APC_BYTE_TO_U32 inverse
)

type OptionalRGBA64 struct {
	Color *color.RGBA64
	Empty bool
}

func ParseColor(str string) *OptionalRGBA64 {

	re := regexp.MustCompile("[\\(\\)\\[\\]]")
	str = re.ReplaceAllString(str, "")

	arr := strings.Split(str, ",")

	// avoid oob access to arr when format is not respected
	if len(arr) < 4 {
		return &OptionalRGBA64{
			Color: nil,
			Empty: true,
		}
	}

	errs := make([]error, 4)

	var rr, gg, bb, aa int
	rr, errs[0] = strconv.Atoi(arr[0])
	gg, errs[1] = strconv.Atoi(arr[1])
	bb, errs[2] = strconv.Atoi(arr[2])
	aa, errs[3] = strconv.Atoi(arr[3])

	for _, val := range errs {
		if val != nil {
			return &OptionalRGBA64{
				Color: nil,
				Empty: true,
			}
		}
	}

	return &OptionalRGBA64{
		Color: &color.RGBA64{
			R: uint16(rr * aa * APC_BYTE_TO_U16),
			G: uint16(gg * aa * APC_BYTE_TO_U16),
			B: uint16(bb * aa * APC_BYTE_TO_U16),
			A: uint16(aa * BYTE_TO_U16),
		},
		Empty: false,
	}
}
