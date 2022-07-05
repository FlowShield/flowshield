package caclient

import (
	"fmt"
	"testing"
	"time"
)

func TestLifespan(t *testing.T) {
	result1, resutl2 := testLifeSpan(
		time.Date(2021, 1, 22, 0, 0, 0, 0, time.Local),
		time.Date(2021, 3, 22, 0, 0, 0, 0, time.Local),
		2,
	)
	fmt.Println(result1, resutl2)
}

func testLifeSpan(notBefore, notAfter time.Time, rate int) (remain time.Duration, ava time.Duration) {
	now := time.Now()
	if now.After(notAfter) {
		return 0, 0
	}

	remain = notAfter.Sub(now)

	certLong := notAfter.Sub(notBefore)
	ava = certLong / time.Duration(rate)

	fmt.Println("Surplus hours: ", remain.Hours())
	fmt.Println("Next replacement hours: ", ava.Hours())

	return remain, ava
}
