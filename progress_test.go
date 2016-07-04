package progress

import (
	"bytes"
	"testing"
	"time"
)

func TestProgress(t *testing.T) {
	var buf bytes.Buffer
	pb := New(&buf)
	pb.clock = &fakeClock{now: time.Now()}
	pb.Update(1e3, 1e9)
	pb.clock.(*fakeClock).Tick(100 * time.Millisecond)
	pb.Update(2e3, 1e9)
	pb.clock.(*fakeClock).Tick(10 * time.Second)
	pb.Update(3e3, 1e9)

	want := "1,000/1,000,000,000\r2,000/1,000,000,000\r3,000/1,000,000,000@100/s\r"
	if buf.String() != want {
		t.Fatalf("got=%q", buf.Bytes())
	}
}
