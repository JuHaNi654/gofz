package system

import (
	"testing"
)

func TestWd(t *testing.T) {
	checkPath := func(t testing.TB, value, expected string) {
		t.Helper()
		if value != expected {
			t.Errorf("got '%s' :: want '%s'", value, expected)
		}
	}

	t.Run("test valid parent path", func(t *testing.T) {
		item := InitDirectoryCache("/foo/bar")
		expected := "/foo"
		checkPath(t, item.parentWd, expected)
	})

	t.Run("test valid parent path", func(t *testing.T) {
		item := InitDirectoryCache("/foo/bar/")
		expected := "/foo"
		checkPath(t, item.parentWd, expected)
	})

	t.Run("expected to get root path", func(t *testing.T) {
		item := InitDirectoryCache("/foo")
		expected := "/"
		checkPath(t, item.parentWd, expected)
	})

	t.Run("root path should not change", func(t *testing.T) {
		item := InitDirectoryCache("/")
		expected := "/"
		checkPath(t, item.parentWd, expected)
	})

	t.Run("update to previous paths on call", func(t *testing.T) {
		item := InitDirectoryCache("/foo/bar")
		item.PreviousWd()
		expectedCurrent := "/foo"
		expectedParent := "/"

		checkPath(t, item.currentWd, expectedCurrent)
		checkPath(t, item.parentWd, expectedParent)
	})

	t.Run("update to the given path on call", func(t *testing.T) {
		item := InitDirectoryCache("/foo/bar")
		item.NextWd("/baz")
		expectedCurrent := "/foo/bar/baz"
		expectedParent := "/foo/bar"

		checkPath(t, item.currentWd, expectedCurrent)
		checkPath(t, item.parentWd, expectedParent)
	})
}
