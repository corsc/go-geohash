package geohash

import (
	"testing"
	"math"
)

func TestEncodeIntBasic(t *testing.T) {
	var expected GeoHashInt = 4064984913515641

	result := EncodeInt(37.8324, 112.5584, 52)

	if expected != result {
		t.Errorf("Expected %+v but was %+v", expected, result)
	}
}

func TestDecodeIntBasic(t *testing.T) {
	var expectedLat float64 = 37.8324
	var expectedLng float64 = 112.5584

	resultLat, resultLng, _, _ := DecodeInt(4064984913515641, 52)

	if math.Abs(expectedLat-resultLat) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedLat, resultLat)
	}
	if math.Abs(expectedLng-resultLng) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedLng, resultLng)
	}
}

func TestDecodeBboxInt(t *testing.T) {
	var expectedMinLat float64 = 37.8324
	var expectedMinLng float64 = 112.5584
	var expectedMaxLat float64 = 37.8324
	var expectedMaxLng float64 = 112.5584

	minLat, minLng, maxLat, maxLng := DecodeBboxInt(4064984913515641, 52)

	if math.Abs(expectedMinLat-minLat) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedMinLat, minLat)
	}
	if math.Abs(expectedMinLng-minLng) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedMinLng, minLng)
	}
	if math.Abs(expectedMaxLat-maxLat) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedMaxLat, maxLat)
	}
	if math.Abs(expectedMaxLng-maxLng) > 0.0001 {
		t.Errorf("Expected %+v but was %+v", expectedMaxLng, maxLng)
	}
}

func TestNeighborInt(t *testing.T) {
	result := NeighborInt(1702789509, North, 32)
	var expected GeoHashInt = 1702789520
	if expected != result {
		t.Errorf("Expected %+v but was %+v", expected, result)
	}

	result = NeighborInt(27898503327470, SouthWest, 46)
	expected = 27898503327465
	if expected != result {
		t.Errorf("Expected %+v but was %+v", expected, result)
	}
}

func TestNeighborsInt(t *testing.T) {
	expected := []GeoHashInt{1702789520, 1702789522, 1702789511, 1702789510, 1702789508, 1702789422, 1702789423, 1702789434, 1702789509}
	results := NeighborsInt(1702789509, 32)

	for _, expectedValue := range expected {
		found := false
		for _, resultValue := range results {
			if expectedValue == resultValue {
				found = true
			}
		}
		if !found {
			t.Errorf("Expected value %+v not found.", expectedValue)
		}
	}
}

func TestBBoxesInt(t * testing.T) {
	results := BboxesInt(30, 120, 30.0001, 120.0001, 50);
	expected := EncodeInt(30.0001, 120.0001, 50)

	found := false
	for _, resultValue := range results {
		if expected == resultValue {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected value %+v not found.", expected)
	}
}
