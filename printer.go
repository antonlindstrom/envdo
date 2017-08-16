package main

import (
	"os"

	"github.com/a8m/tree"
)

type fs struct{}

func (f *fs) Stat(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}
func (f *fs) ReadDir(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	names, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		return nil, err
	}
	return names, nil
}

func printTree(path string) {
	opts := &tree.Options{
		Fs:         new(fs),
		FollowLink: true,
		OutFile:    os.Stdout,
		Colorize:   false,
	}

	var nd, nf int
	inf := tree.New(path)
	if d, f := inf.Visit(opts); f != 0 {
		if d > 0 {
			d -= 1
		}
		nd, nf = nd+d, nf+f
	}

	inf.Print(opts)
}
