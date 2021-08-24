package benchmark
import "time"

type Network struct {
	kbps int
	Latency time.Duration
	MTU int // Bytes per packet; if non-positive, infine
}

var now = time.Now()
