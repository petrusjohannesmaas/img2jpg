
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/gif"
	_ "image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/webp"
	"golang.org/x/image/tiff"
)

func init() {
	image.RegisterFormat("webp", "RIFF????WEBP", webp.Decode, webp.DecodeConfig)
	image.RegisterFormat("tiff", "II*", tiff.Decode, tiff.DecodeConfig)
	image.RegisterFormat("tiff", "MM", tiff.Decode, tiff.DecodeConfig)
}

var supportedExts = map[string]bool{
	".png": true, ".webp": true, ".gif": true,
	".tiff": true, ".tif": true, ".ico": true, ".bmp": true,
}

func main() {
	input := ""
	quality := 80
	outDir := ""
	recursive := false

	// Manual argument parsing to allow flags in any order
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-q", "--quality":
			if i+1 < len(os.Args) {
				if q, err := strconv.Atoi(os.Args[i+1]); err == nil {
					quality = q
				}
				i++
			}
		case "-o", "--output":
			if i+1 < len(os.Args) {
				outDir = os.Args[i+1]
				i++
			}
		case "-r", "--recursive":
			recursive = true
		default:
			if input == "" {
				input = os.Args[i]
			}
		}
	}

	if input == "" {
		fmt.Println("Usage: img2jpg <input_path> [-q quality] [-o output_dir] [-r]")
		os.Exit(1)
	}

	if quality < 1 { quality = 1 }
	if quality > 100 { quality = 100 }

	// Default output to a "converted" subfolder relative to input
	if outDir == "" {
		if info, err := os.Stat(input); err == nil && info.IsDir() {
			outDir = filepath.Join(input, "converted")
		} else {
			outDir = filepath.Join(filepath.Dir(input), "converted")
		}
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Printf("❌ Failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(input)
	if err != nil {
		fmt.Printf("❌ Cannot access input path: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		processDir(input, outDir, quality, recursive)
	} else {
		convertSingle(input, outDir, quality)
	}

	fmt.Printf("\n✅ Conversion complete. Output: %s\n", outDir)
}

func processDir(srcDir, outDir string, quality int, recursive bool) {
	filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if d.IsDir() && !recursive && path != srcDir {
				return filepath.SkipDir
			}
			return err
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExts[ext] { return nil }

		rel, _ := filepath.Rel(srcDir, path)
		targetDir := filepath.Join(outDir, filepath.Dir(rel))
		os.MkdirAll(targetDir, 0755)

		outPath := filepath.Join(targetDir, strings.TrimSuffix(filepath.Base(path), ext)+".jpg")
		convertImage(path, outPath, quality)
		return nil
	})
}

func convertSingle(srcPath, outDir string, quality int) {
	ext := strings.ToLower(filepath.Ext(srcPath))
	if !supportedExts[ext] {
		fmt.Printf("⚠️  Skipping unsupported format: %s\n", ext)
		return
	}
	outPath := filepath.Join(outDir, strings.TrimSuffix(filepath.Base(srcPath), ext)+".jpg")
	convertImage(srcPath, outPath, quality)
}

func convertImage(src, dst string, quality int) {
	f, err := os.Open(src)
	if err != nil {
		fmt.Printf("❌ Open failed %s: %v\n", src, err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Printf("❌ Decode failed %s: %v\n", src, err)
		return
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Over)

	out, err := os.Create(dst)
	if err != nil {
		fmt.Printf("❌ Create failed %s: %v\n", dst, err)
		return
	}
	defer out.Close()

	if err := jpeg.Encode(out, rgba, &jpeg.Options{Quality: quality}); err != nil {
		fmt.Printf("❌ Encode failed %s: %v\n", dst, err)
		return
	}

	fmt.Printf("✅ %s → %s\n", src, dst)
}
