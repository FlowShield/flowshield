package structure

import (
	"github.com/jinzhu/copier"
)

// Copy structure mapping
func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
