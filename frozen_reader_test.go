package fileutil

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestFrozenReader(t *testing.T) {
	forever := func() {
		wrappedStdin := FrozenReader{}
		io.Copy(os.Stdout, wrappedStdin)
	}
	go forever()
	<-time.After(10 * time.Second)
}
