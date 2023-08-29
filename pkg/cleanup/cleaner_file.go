package cleanup

import "os"

type fileCleaner struct{}

func NewFileCleaner() Cleaner {
	return &fileCleaner{}
}

func (f *fileCleaner) Remove(path string) error {
	return os.Remove(path)
}

func (f *fileCleaner) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
