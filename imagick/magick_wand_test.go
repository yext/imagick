// Copyright 2013 Herbert G. Fischer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imagick

import (
	"testing"
)

var (
	mw *MagickWand
)

// SUPPORTED_FORMATS are verified to be present in the available formats.
var SUPPORTED_FORMATS = []string{
	"BMP",
	"GIF",
	"JPEG",
	"JPG",
	"PNG",
	"SVG",
	"TGA",
	"TIFF",
	"WEBP",
}

func Init() {
	Initialize()
}

func TestNewMagickWand(t *testing.T) {
	mw := NewMagickWand()
	defer mw.Destroy()

	if !mw.IsVerified() {
		t.Fatal("MagickWand not verified")
	}
}

func TestCloningAndDestroying(t *testing.T) {
	t.Skip()
	clone := mw.Clone()
	if !clone.IsVerified() {
		t.Fatal("Unsuccessful clone")
	}
	clone.Destroy()
	if clone.IsVerified() || !mw.IsVerified() {
		t.Fatal("MagickWand not properly destroyed")
	}
}

func TestQueryConfigureOptions(t *testing.T) {
	opts := mw.QueryConfigureOptions("*")
	if len(opts) == 0 {
		t.Fatal("QueryConfigureOptions returned an empty array")
	}
	for _, opt := range opts {
		mw.QueryConfigureOption(opt)
	}
}

func TestNonExistingConfigureOption(t *testing.T) {
	_, err := mw.QueryConfigureOption("4321foobaramps1234")
	if err == nil {
		t.Fatal("Missing error when trying to get non-existing configure option")
	}
}

func TestQueryFonts(t *testing.T) {
	t.Log("font support is disabled")
	t.Skip()
	fonts := mw.QueryFonts("*")
	if len(fonts) == 0 {
		t.Fatal("ImageMagick have not identified a single font in this system")
	}
}

func TestQueryFormats(t *testing.T) {
	formats := mw.QueryFormats("*")
	if len(formats) == 0 {
		t.Fatal("ImageMagick have not identified a single image format in this system")
	}
	for _, requiredFormat := range SUPPORTED_FORMATS {
		if !contains(formats, requiredFormat) {
			t.Errorf("Format %s is missing", requiredFormat)
		}
	}
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}

func TestDeleteImageArtifact(t *testing.T) {
	t.Skip()
	err := mw.DeleteImageArtifact("*")
	t.Log(err.Error())
}

func TestReadImageBlob(t *testing.T) {
	mw = NewMagickWand()
	defer mw.Destroy()

	// Read an invalid blob
	blob := []byte{}
	if err := mw.ReadImageBlob(blob); err == nil {
		t.Fatal("Expected a failure when passing a zero length blob")
	}

	mw.ReadImage(`logo:`)
	blob = mw.GetImageBlob()

	// Read a valid blob
	if err := mw.ReadImageBlob(blob); err != nil {
		t.Fatal(err.Error())
	}
}

func TestGetImageFloats(t *testing.T) {
	Initialize()
	mw := NewMagickWand()
	defer mw.Destroy()

	var err error
	if err = mw.ReadImage(`logo:`); err != nil {
		t.Fatal("Failed to read internal logo: image")
	}

	width, height := mw.GetImageWidth(), mw.GetImageHeight()

	val, err := mw.ExportImagePixels(0, 0, width, height, "RGB", PIXEL_FLOAT)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels := val.([]float32)
	actual := len(pixels)
	expected := (width * height * 3)
	if actual != int(expected) {
		t.Fatalf("Expected RGB image to have %d float vals; Got %d", expected, actual)
	}

	val, err = mw.ExportImagePixels(0, 0, width, height, "RGBA", PIXEL_DOUBLE)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels64 := val.([]float64)
	actual = len(pixels64)
	expected = (width * height * 4)
	if actual != int(expected) {
		t.Fatalf("Expected RGBA image to have %d float vals; Got %d", expected, actual)
	}

	val, err = mw.ExportImagePixels(0, 0, width, height, "R", PIXEL_FLOAT)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels = val.([]float32)
	actual = len(pixels)
	expected = (width * height * 1)
	if actual != int(expected) {
		t.Fatalf("Expected RNN image to have %d float vals; Got %d", expected, actual)
	}

	val, err = mw.ExportImagePixels(0, 0, width, height, "GB", PIXEL_FLOAT)
	if err != nil {
		t.Fatal(err.Error())
	}
	pixels = val.([]float32)
	actual = len(pixels)
	expected = (width * height * 2)
	if actual != int(expected) {
		t.Fatalf("Expected NGB image to have %d float vals; Got %d", expected, actual)
	}
}

func BenchmarkExportImagePixels(b *testing.B) {
	wand := NewMagickWand()
	defer wand.Destroy()

	wand.ReadImage("logo:")
	wand.ScaleImage(1024, 1024)

	var val interface{}
	var pixels []float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		val, _ = wand.ExportImagePixels(0, 0, 1024, 1024, "RGB", PIXEL_FLOAT)
		pixels = val.([]float32)
	}

	b.StopTimer()

	if len(pixels) == 0 {
		b.Fatal("Pixel slice is 0")
	}
}

func BenchmarkImportImagePixels(b *testing.B) {
	wand := NewMagickWand()
	defer wand.Destroy()

	wand.ReadImage("logo:")
	wand.ScaleImage(1024, 1024)

	val, _ := wand.ExportImagePixels(0, 0, 1024, 1024, "RGB", PIXEL_FLOAT)
	pixels := val.([]float32)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wand.ImportImagePixels(0, 0, 1024, 1024, "RGB", PIXEL_UNDEFINED, pixels)
	}

	b.StopTimer()
}
