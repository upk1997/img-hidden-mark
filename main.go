package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

func addMark(original, mark, marked string) {
	// open original image
	originalFile, err := os.Open(original)
	if err != nil {
		log.Fatalln(err)
	}
	defer originalFile.Close()
	originalImage, _, err := image.Decode(originalFile)
	if err != nil {
		log.Fatalln(err)
	}

	// open mark image
	markFile, err := os.Open(mark)
	if err != nil {
		log.Fatalln(err)
	}
	defer markFile.Close()
	markImage, _, err := image.Decode(markFile)
	if err != nil {
		log.Fatalln(err)
	}

	// prepare marked image
	bounds := originalImage.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	markedImage := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := range w {
		for y := range h {
			// mark
			oldColor := originalImage.At(x, y)
			r, g, b, a := oldColor.RGBA()

			markR, _, _, _ := markImage.At(x, y).RGBA()
			var mark uint32 = 1
			if markR < 128 {
				mark = 0
			}

			r = (r & 0b11111110) + mark
			markedImage.SetRGBA(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}

	outfile, err := os.Create(marked)
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	png.Encode(outfile, markedImage)
}

func decodeMark(marked, decoded string) {
	// open original image
	originalFile, err := os.Open(marked)
	if err != nil {
		log.Fatalln(err)
	}
	defer originalFile.Close()
	originalImage, _, err := image.Decode(originalFile)
	if err != nil {
		log.Fatalln(err)
	}

	bounds := originalImage.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	newImage := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := range w {
		for y := range h {

			// decode
			oldColor := originalImage.At(x, y)
			originalImage.At(x, y)
			r, _, _, _ := oldColor.RGBA()
			if r&1 == 1 {
				newImage.Set(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				newImage.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	outfile, err := os.Create(decoded)
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	png.Encode(outfile, newImage)
}

func main() {
	// create output directory
	err := os.MkdirAll("output", 0755)
	if err != nil {
		log.Fatalln(err)
	}

	// add a hidden mark into an image file
	addMark("img_examples/original.png", "img_examples/mark.png", "output/marked.png")
	println("Successfully added a hidden mark")

	// extract hidden mark from an image file
	decodeMark("output/marked.png", "output/decoded.png")
	println("Successfully decoded")
}
