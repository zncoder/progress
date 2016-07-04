package progress

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type Progress struct {
	w     io.Writer
	last  int64
	at    time.Time
	clock Clock
}

func New(w io.Writer) *Progress {
	return &Progress{
		w:     w,
		clock: realClock{},
	}
}

func (pb *Progress) Update(cur, total int64) {
	var rate int64
	now := pb.clock.Now()
	if !pb.at.IsZero() {
		diff := cur - pb.last
		sec := now.Sub(pb.at).Nanoseconds() / 1e9
		if sec > 0 {
			rate = diff / sec
		}
	}

	pb.last = cur
	pb.at = now

	var srate string
	if rate > 0 {
		srate = "@" + fmtInt(rate) + "/s"
	}

	fmt.Fprintf(pb.w, "%s/%s%s\r", fmtInt(cur), fmtInt(total), srate)
}

func fmtInt(n int64) string {
	var ss []string
	for n > 0 {
		s := fmt.Sprintf("%03d", n%1000)
		ss = append([]string{s}, ss...)
		n /= 1000
	}
	return strings.TrimLeft(strings.Join(ss, ","), "0")
}

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (rc realClock) Now() time.Time { return time.Now() }

type fakeClock struct {
	now time.Time
}

func (fc *fakeClock) Now() time.Time { return fc.now }

func (fc *fakeClock) Tick(dur time.Duration) {
	fc.now = fc.now.Add(dur)
}
