package general

import (
	"io"
	"os"
	"strings"
)

func Move(ori, dst string) error {
	// First, we try to rename
	err := os.Rename(ori, dst)
	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "cross-device link") {
		return err
	}
	// Rename failed because old and new are cross device files.
	err = CopyFile(ori, dst)
	if err != nil {
		return err
	}
	_ = os.Remove(ori)
	return nil
}

func CopyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer CloseWithCheck(s)

	ss, err := os.Stat(src)
	if err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer CloseWithCheck(d)

	err = os.Chmod(dst, ss.Mode())
	if err != nil {
		return err
	}

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return nil
}
