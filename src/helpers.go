package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func copyFile(src, dst string) {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return // Ignore missing files
	}
	check(os.MkdirAll(filepath.Dir(dst), 0o755))
	in, err := os.Open(src)
	check(err)
	defer in.Close()
	out, err := os.Create(dst)
	check(err)
	defer out.Close()
	_, err = io.Copy(out, in)
	check(err)
}

func parseString(content, pattern string) string {
	if match := regexp.MustCompile(pattern).FindStringSubmatch(content); len(match) > 1 {
		return match[1]
	}
	return ""
}

func parseInt(content, pattern string) int {
	var num int
	fmt.Sscanf(parseString(content, pattern), "%d", &num)
	return num
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	check(os.MkdirAll(filepath.Dir(dest), 0o755))
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func createZip(src, dest string) error {
	z, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer z.Close()
	w := zip.NewWriter(z)
	defer w.Close()
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil || src == path {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		if info.IsDir() {
			_, err = w.Create(filepath.ToSlash(rel) + "/")
			return err
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		zh, _ := zip.FileInfoHeader(info)
		zh.Name = filepath.ToSlash(rel)
		zh.Method = zip.Deflate
		zw, err := w.CreateHeader(zh)
		if err == nil {
			io.Copy(zw, f)
		}
		return err
	})
}
