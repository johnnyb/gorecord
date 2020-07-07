package gorec

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

func ConvertArbitraryToTime(val interface{}) (time.Time, error) {
	switch newval := val.(type) {
	case string:
		t, err := time.Parse("RFC3339", newval)
		if err != nil {
			return time.Unix(0, 0), err
		}
		return t, nil
	default:
		return time.Unix(0, 0), errors.New("Could not convert")
	}
}

func ConvertArbitraryToInt32(val interface{}) (int32, error) {
	switch newval := val.(type) {
	case string:
		tmpval, err := strconv.ParseInt(newval, 10, 64)
		if err != nil {
			return 0, err
		}
		return int32(tmpval), nil
	case int32:
		return newval, nil
	case int64:
		return int32(newval), nil
	default:
		return 0, errors.New("Could not convert to int32")
	}
}

func ConvertArbitraryToString(val interface{}) (string, error) {
	switch newval := val.(type) {
	case string:
		return newval, nil
	case fmt.Stringer:
		return newval.String(), nil
	default:
		return "", errors.New("Can't convert to string")
	}
}

func ConvertArbitraryToArbitrary(val interface{}) (interface{}, error) {
	return val, nil
}
