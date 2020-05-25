package gorec

import (
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}
