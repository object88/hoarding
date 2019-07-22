package hoarding

import (
	"bytes"
	"testing"
)

func Test_SendMessage(t *testing.T) {
	var buf bytes.Buffer

	h := NewHoarder(&buf, func() bool { return true })
	h.Msgf(Always, "a")

	if len(buf.Bytes()) != 0 {
		// Premature writing
		t.Errorf("The buffer was prematurely written to")
	}

	h.Close()

	if "a" != string(buf.Bytes()) {
		t.Errorf("Buffer does not have the correct contents: '%s'", string(buf.Bytes()))
	}
}
