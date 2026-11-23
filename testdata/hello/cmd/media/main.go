// Command media writes hello's recording, as a game's make media does: hello.gif, a square
// that changes colour, playing over and over, and hello.png, its first frame. The release
// workflow's media input is rehearsed with it.
//
//	go run ./cmd/media -out dist/media
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
)

func main() {
	out := flag.String("out", "dist/media", "the folder to write the recording into")
	flag.Parse()
	if err := record(*out); err != nil {
		fmt.Fprintln(os.Stderr, "media:", err)
		os.Exit(1)
	}
}

// record writes hello.gif and hello.png into dir, the same bytes every time.
func record(dir string) error {
	colors := color.Palette{
		color.RGBA{R: 0x0b, G: 0x3d, B: 0x1c, A: 0xff},
		color.RGBA{R: 0x1f, G: 0x8f, B: 0x45, A: 0xff},
		color.RGBA{R: 0x5c, G: 0xf0, B: 0x8a, A: 0xff},
	}
	movie := &gif.GIF{} // a LoopCount of 0 plays it over and over
	for i := range colors {
		frame := image.NewPaletted(image.Rect(0, 0, 8, 8), colors)
		for p := range frame.Pix {
			frame.Pix[p] = uint8(i)
		}
		movie.Image = append(movie.Image, frame)
		movie.Delay = append(movie.Delay, 50)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := create(filepath.Join(dir, "hello.gif"), func(f *os.File) error { return gif.EncodeAll(f, movie) }); err != nil {
		return err
	}
	return create(filepath.Join(dir, "hello.png"), func(f *os.File) error { return png.Encode(f, movie.Image[0]) })
}

// create makes the file at path with what encode writes to it.
func create(path string, encode func(f *os.File) error) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := encode(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
