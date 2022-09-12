package mode

import (
	"fmt"
	"testing"
)

func TestCheckMode(t *testing.T) {
	t.Run("extends pub_subs", func(t *testing.T) {
		mode, strings := CheckMode([]string{"asdfasfs", "/* @exclude extends pub_subs */"})
		assertTrue(t, mode == Extend, "")
		assertTrue(t, len(strings) == 1, "")
		assertEqual(t, strings[0], "pub_subs")
	})
	t.Run("default mode when patterns not match", func(t *testing.T) {
		mode, strings := CheckMode([]string{"/* @exclude aaa pubsub */"})
		assertTrue(t, mode == Independent, "")
		assertTrue(t, len(strings) == 0, "")
	})
}

func assertEqual(t *testing.T, real string, expect string) {
	t.Helper()
	if real != expect {
		t.Errorf(fmt.Sprintf("Expect %s but got %s", expect, real))
	}
}

func assertTrue(t *testing.T, b bool, msg string) {
	t.Helper()
	if !b {
		t.Errorf(msg)
	}
}
