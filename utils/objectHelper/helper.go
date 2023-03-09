package objecthelper

import (
	"math"
	mRand "math/rand"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

// GenerateUUID is used to generate uuid
func GenerateUUID() (uuid.UUID, error) {
	return uuid.NewUUID()
}

// MarshalJSON - Marshal returns the JSON encoding of data
func MarshalJSON(data interface{}) ([]byte, error) {
	return jsoniter.Marshal(data)
}

// UnmarshalJSON -- returns the parsed JSON-encoded data in `val`
func UnmarshalJSON(data []byte, val interface{}) error {
	return jsoniter.Unmarshal(data, val)
}

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

func IsCurrentDay(t time.Time) bool {
	y1, m1, d1 := time.Now().Date()
	y2, m2, d2 := t.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetQuarterFromMonth(month int) int {
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	case month >= 10 && month <= 12:
		return 4
	}

	return -1
}

func GenerateOTP(max int) int {

	low := 100000
	high := 999999

	return low + mRand.Intn(high-low)
}
