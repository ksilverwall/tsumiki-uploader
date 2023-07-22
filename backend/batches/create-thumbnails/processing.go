package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strings"
	"net/http"

	"github.com/nfnt/resize"
)

type Result struct {
	FilePaths []string
	Errs      []error
}

func ProcessZipFile(f *zip.File, outPath string) (string, error) {
	if f.FileInfo().IsDir() {
		return "", fmt.Errorf("directory detected")
	}

	rc, err := f.Open()
	buffer := make([]byte, 512)
	rc.Read(buffer)
	contentType := http.DetectContentType(buffer)

	rc, err = f.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer rc.Close()

	img, _, err := image.Decode(rc)
	if err != nil {
		return "", fmt.Errorf("failed to decode image file (type %s): %w", contentType, err)
	}

	err = CreateThumbnail(img, outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create thumbnail file: %w", err)
	}

	return outPath, nil
}

func ProcessZipFiles(content []byte, outDir string) (Result, error) {
	errs := []error{}
	paths := []string{}
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return Result{}, err
	}

	if len(reader.File) == 0 {
		return Result{}, fmt.Errorf("File not found: %w", err)
	}

	for i, f := range reader.File {
		_, fileName := path.Split(f.Name)
		outPath := path.Join(outDir, strings.Replace(fileName, ".", "_thumbnail.", 1))

		fname, err := ProcessZipFile(f, outPath)
		if err != nil{
			errs = append(errs, fmt.Errorf("failed to process file[%d]: %w", i, err))
			continue
		}
		paths = append(paths, fname)
	}

	return Result{FilePaths: paths, Errs: errs}, nil
}

func CreateThumbnail(img image.Image, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	thumbnail := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	err = jpeg.Encode(outFile, thumbnail, nil)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	return nil
}
