package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	pg "github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFiles             = errors.New("from and to files are sane")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var from, to *os.File
	var fromStat, toStat os.FileInfo
	var err error
	if from, err = os.Open(fromPath); err != nil {
		return fmt.Errorf("open source file error: %w", err)
	}
	defer from.Close()
	if to, err = os.Create(toPath); err != nil {
		return fmt.Errorf("open destination file error: %w", err)
	}
	defer to.Close()
	if fromStat, err = from.Stat(); err != nil {
		return fmt.Errorf("source stat file error: %w", err)
	}
	if toStat, err = to.Stat(); err != nil {
		return fmt.Errorf("destination stat file error: %w", err)
	}
	if os.SameFile(fromStat, toStat) {
		return ErrSameFiles
	}
	if limit, err = calcLimit(fromStat, offset, limit); err != nil {
		return err
	}
	if offset > 0 {
		if _, err = from.Seek(offset, io.SeekStart); err != nil {
			return fmt.Errorf("offset error: %w", err)
		}
	}
	bar := pg.DefaultBytes(limit, "copying")
	var w io.Writer = to
	w = io.MultiWriter(w, bar)
	if _, err = io.CopyN(w, from, limit); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}
	return nil
}

func calcLimit(stat os.FileInfo, offset, limit int64) (int64, error) {
	if !stat.Mode().IsRegular() {
		if limit < 1 {
			return 0, ErrUnsupportedFile
		}
		return limit, nil
	}
	size := stat.Size()
	if size < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	if limit < 1 || size < (offset+limit) {
		return size - offset, nil
	}
	return limit, nil
}
