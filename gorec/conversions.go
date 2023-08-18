package gorec

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"
	"database/sql"
	"github.com/google/uuid"
)

func ConvertArbitraryToBool(val interface{}) (bool, error) {
	switch newval := val.(type) {
		case string:
			newval = strings.ToLower(newval)
			switch newval {
				case "t", "y", "1":
					return true, nil
				default:
					return false, nil
			}
		case bool:
			return newval, nil
		default:
			return false, errors.New("Could not convert")
	}
}

func ConvertArbitraryToTime(val interface{}) (time.Time, error) {
	switch newval := val.(type) {
	case string:
		t, err := time.Parse("RFC3339", newval)
		if err != nil {
			return time.Unix(0, 0), err
		}
		return t, nil
	case time.Time:
		return newval, nil
	default:
		return time.Unix(0, 0), errors.New("Could not convert")
	}
}

func ConvertArbitraryToNullTime(val interface{}) (sql.NullTime, error) {
	result := sql.NullTime{}
	var err error

	if val == nil {
		return result, nil
	}

	switch newval := val.(type) {
	case sql.NullTime:
		return newval, nil
	case string:
		t, err := time.Parse("RFC3339", newval)
		if err == nil {
			result.Valid = true
			result.Time = t
		}
	case *string:
		if newval != nil {
			t, err := time.Parse("RFC3339", *newval)
			if err == nil {
				result.Valid = true
				result.Time = t
			}
		}
	case time.Time:
		result.Valid = true
		result.Time = newval
	case *time.Time:
		if newval != nil {
			result.Valid = true
			result.Time = *newval
		}
	}

	return result, err
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

func ConvertArbitraryToNullInt32(val interface{}) (sql.NullInt32, error) {
	result := sql.NullInt32{}
	var err error

	if val == nil {
		return result, nil
	}

	switch newval := val.(type) {
	case sql.NullInt32:
		return newval, nil
	case string:
		if newval == "" {
			return result, nil
		}

		tmpval, err := strconv.ParseInt(newval, 10, 64)
		if err == nil {
			result.Int32 = int32(tmpval)
			result.Valid = true
		}
	case int32:
		result.Int32 = newval
		result.Valid = true
	case int64:
		result.Int32 = int32(newval)
		result.Valid = true
	default:
		err = errors.New("Could not convert to int32")
	}

	return result, err
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

func ConvertArbitraryToNullString(val interface{}) (sql.NullString, error) {
	result := sql.NullString{}
	var err error

	if val == nil {
		return result, nil
	}

	switch newval := val.(type) {
	case sql.NullString:
		return newval, nil
	case string:
		result.String = newval
		result.Valid = true
	case *string:
		result.String = *newval
		result.Valid = true
	}

	return result, err
}

func ConvertArbitraryToArbitrary(val interface{}) (interface{}, error) {
	return val, nil
}

func ConvertArbitraryToNullUUID(val interface{}) (uuid.NullUUID, error) {
	result := uuid.NullUUID{}
	var err error

	if val == nil {
		return result, nil
	}

	switch newval := val.(type) {
	case uuid.NullUUID:
		return newval, nil
	case uuid.UUID:
		result.UUID = newval
		result.Valid = true
	case string:
		tmp, err := uuid.Parse(newval)
		if err == nil {
			result.UUID = tmp
			result.Valid = true
		}
	}

	return result, err
}
