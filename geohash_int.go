package geohash

import (
	"fmt"
	"math"
)

const (
	// MaxBitDepth defines both the maximum and default geohash accuracy.
	MaxBitDepth int64 = 52
)

// bearing defines the compass bearing/direction in matrix form relative to a center point of 0,0
//  |----------------------|
// 	|   NW  |   N   |  NE  |
// 	|  1,-1 |  1,0  |  1,1 |
//  |----------------------|
// 	|   W   |   X   |   E  |
// 	|  0,-1 |  0,0  |  0,1 |
//  |----------------------|
// 	|   SW  |   S   |  NE  |
// 	| -1,-1 |  -1,0 | -1,1 |
//  |----------------------|
type bearing struct {
	x, y int
}

// North bearing from reference point X
var North = bearing{1, 0}

// NorthEast bearing from reference point X
var NorthEast = bearing{1, 1}

// East bearing from reference point X
var East = bearing{0, 1}

// SouthEast bearing from reference point X
var SouthEast = bearing{-1, 1}

// South bearing from reference point X
var South = bearing{-1, 0}

// SouthWest bearing from reference point X
var SouthWest = bearing{-1, -1}

// West bearing from reference point X
var West = bearing{0, -1}

// NorthWest bearing from reference point X
var NorthWest = bearing{1, -1}

// bitsToDistanceInMeters provides a mapping between bitDepth values and distances
var bitsToDistanceInMeters map[int64]float64

func init() {
	// Reference: https://github.com/yinqiwen/ardb/blob/master/doc/spatial-index.md
	bitsToDistanceInMeters = map[int64]float64{
		52: 0.5971,
		50: 1.1943,
		48: 2.3889,
		46: 4.7774,
		44: 9.5547,
		42: 19.1095,
		40: 38.2189,
		38: 76.4378,
		36: 152.8757,
		34: 305.751,
		32: 611.5028,
		30: 1223.0056,
		28: 2446.0112,
		26: 4892.0224,
		24: 9784.0449,
		22: 19568.0898,
		20: 39136.1797,
		18: 78272.35938,
		16: 156544.7188,
		14: 313089.4375,
		12: 626178.875,
		10: 1252357.75,
		8:  2504715.5,
		6:  5009431,
		4:  10018863,
	}
}

// EncodeInt will encode a pair of latitude and longitude values into a geohash integer.
//
// The third argument is the bitDepth of this number, which affects the precision of the geohash
// but also must be used consistently when decoding. Bit depth must be even.
func EncodeInt(latitude float64, longitude float64, bitDepth int64) int64 {
	// input validation
	validateBitDepth(bitDepth)

	// initialize the calculation
	var bitsTotal int64
	var mid float64
	var maxLat float64 = 90.0
	var minLat float64 = -90.0
	var maxLng float64 = 180.0
	var minLng float64 = -180.0

	var geohash int64
	for bitsTotal < bitDepth {
		geohash *= 2

		if bitsTotal%2 == 0 {
			mid = (maxLng + minLng) / 2

			if longitude > mid {
				geohash += 1
				minLng = mid
			} else {
				maxLng = mid
			}
		} else {
			mid = (maxLat + minLat) / 2
			if latitude > mid {
				geohash += 1
				minLat = mid
			} else {
				maxLat = mid
			}
		}
		bitsTotal++
	}
	return geohash
}

// DecodeInt with decode a integer geohashed number into pair of latitude and longitude value approximations.
//
// Returned values include a latitude and longitude along with the maximum error of the calculation.
// This effectively means that a geohash integer will not return a location but an "area".
// The size of the area returned will be vary with different bitDepth settings.
//
// Note: You should provide the same bitDepth to decode the number as was used to produce the geohash originally.
func DecodeInt(geohash int64, bitDepth int64) (lat float64, lng float64, latErr float64, lngErr float64) {
	// input validation
	validateBitDepth(bitDepth)

	minLat, minLng, maxLat, maxLng := DecodeBboxInt(geohash, bitDepth)
	lat = (minLat + maxLat) / 2
	lng = (minLng + maxLng) / 2
	latErr = maxLat - lat
	lngErr = maxLng - lng
	return
}

// DecodeBboxInt will decode a geohash integer into the bounding box that matches it.
//
// Returned as a four corners of a square region.
func DecodeBboxInt(geohash int64, bitDepth int64) (minLat float64, minLng float64, maxLat float64, maxLng float64) {
	// input validation
	validateBitDepth(bitDepth)

	// initialize the calculation
	maxLat = 90
	minLat = -90
	maxLng = 180
	minLng = -180

	var latBit int64
	var lonBit int64
	var steps int64 = bitDepth / 2

	for thisStep := int64(0); thisStep < steps; thisStep++ {
		lonBit = getBit(geohash, ((steps-thisStep)*2)-1)
		latBit = getBit(geohash, ((steps-thisStep)*2)-2)

		if latBit == 0 {
			maxLat = (maxLat + minLat) / 2
		} else {
			minLat = (maxLat + minLat) / 2
		}

		if lonBit == 0 {
			maxLng = (maxLng + minLng) / 2
		} else {
			minLng = (maxLng + minLng) / 2
		}
	}

	return
}

// NeighborInt will find the neighbor of a integer geohash in certain bearing/direction.
//
// The bitDepth should be specified and the same as the value used to generate the hash.
func NeighborInt(geohash int64, bearing bearing, bitDepth int64) int64 {
	// input validation
	validateBitDepth(bitDepth)

	lat, lng, latErr, lngErr := DecodeInt(geohash, bitDepth)
	neighborLat := lat + float64(bearing.x)*latErr*2
	neighborLng := lng + float64(bearing.y)*lngErr*2
	return EncodeInt(neighborLat, neighborLng, bitDepth)
}

// NeighborsInt is the same as calling NeighborInt for each direction and will return all 8 neighbors and the center location.
func NeighborsInt(geohash int64, bitDepth int64) []int64 {
	// input validation
	validateBitDepth(bitDepth)

	var output []int64
	output = append(output, NeighborInt(geohash, North, bitDepth))
	output = append(output, NeighborInt(geohash, NorthEast, bitDepth))
	output = append(output, NeighborInt(geohash, East, bitDepth))
	output = append(output, NeighborInt(geohash, SouthEast, bitDepth))
	output = append(output, NeighborInt(geohash, South, bitDepth))
	output = append(output, NeighborInt(geohash, SouthWest, bitDepth))
	output = append(output, NeighborInt(geohash, West, bitDepth))
	output = append(output, NeighborInt(geohash, NorthWest, bitDepth))
	output = append(output, geohash)
	return output
}

// BboxesInt will return all the hash integers between minLat, minLon, maxLat, maxLon at the requested bitDepth
func BboxesInt(minLat float64, minLon float64, maxLat float64, maxLon float64, bitDepth int64) []int64 {
	// input validation
	validateBitDepth(bitDepth)

	// find the corners
	hashSouthWest := EncodeInt(minLat, minLon, bitDepth)
	hashNorthEast := EncodeInt(maxLat, maxLon, bitDepth)

	_, _, latErr, lngErr := DecodeInt(hashSouthWest, bitDepth)
	perLat := latErr * 2
	perLng := lngErr * 2

	swMinLat, _, _, swMaxLng := DecodeBboxInt(hashSouthWest, bitDepth)
	neMinLat, _, _, neMaxLng := DecodeBboxInt(hashNorthEast, bitDepth)

	latStep := round((neMinLat-swMinLat)/perLat, 0.5, 0)
	lngStep := round((neMaxLng-swMaxLng)/perLng, 0.5, 0)

	var output []int64
	for lat := 0; lat <= int(latStep); lat++ {
		for lng := 0; lng <= int(lngStep); lng++ {
			output = append(output, NeighborInt(hashSouthWest, bearing{lat, lng}, bitDepth))
		}
	}
	return output
}

// getBit returns the bit at the requested location
func getBit(geohash int64, position int64) int64 {
	return int64(int((float64(geohash) / math.Pow(float64(2), float64(position)))) & 0x01)
}

// FindBitDepth will attempt to find the maximum bitdepth which contains the supplied distance
func FindBitDepth(distanceMeters float64) int64 {
	for key, value := range bitsToDistanceInMeters {
		if value > distanceMeters {
			return MaxBitDepth - key
		}
	}
	return 0
}

// Shift provides a convenient way to convert from MaxBitDepth to another
func Shift(value int64, bitDepth int64) int64 {
	// input validation
	validateBitDepth(bitDepth)

	return value << uint64(MaxBitDepth-bitDepth)
}

// validateBitDepth will ensure the supplied bitDepth is valid or cause panic() otherwise.
func validateBitDepth(bitDepth int64) {
	if bitDepth > MaxBitDepth || bitDepth <= 0 {
		panic(fmt.Sprintf("bitDepth must be greater than 0 and less than or equal to %d, was %d", MaxBitDepth, bitDepth))
	}
	if math.Mod(float64(bitDepth), float64(2)) != 0 {
		panic(fmt.Sprintf("bitDepth must be even, was %d", bitDepth))
	}
}

// round is the "missing" round function from the math lib
func round(val float64, roundOn float64, places int) float64 {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
