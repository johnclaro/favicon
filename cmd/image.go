/*
Copyright © 2019 John Claro <jkrclaro@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

// MediaFile TODO
type MediaFile struct {
	filepath string
}

type iconDir struct {
	reserved  uint16
	imageType uint16
	numImages uint16
}

type iconDirEntry struct {
	imageWidth   uint8
	imageHeight  uint8
	numColors    uint8
	reserved     uint8
	colorPlanes  uint16
	bitsPerPixel uint16
	sizeInBytes  uint32
	offset       uint32
}

func newIconDir() iconDir {
	var id iconDir
	id.imageType = 1
	id.numImages = 1
	return id
}

func newIconDirEntry() iconDirEntry {
	var ide iconDirEntry
	ide.colorPlanes = 1
	ide.bitsPerPixel = 32
	ide.offset = 22
	return ide
}

// OpenPNG TODO
func (media MediaFile) OpenPNG() (image.Image, error) {
	file, err := os.Open(media.filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return png.Decode(file)
}

// SaveAsICO TODO
func (media MediaFile) SaveAsICO(writerICO io.Writer, pngFile image.Image) string {
	pngFileBounds := pngFile.Bounds()
	newRGBA := image.NewRGBA(pngFileBounds)
	draw.Draw(newRGBA, pngFileBounds, pngFile, pngFileBounds.Min, draw.Src)

	id := newIconDir()
	ide := newIconDirEntry()

	pngBytesBuffer := new(bytes.Buffer)
	pngWriter := bufio.NewWriter(pngBytesBuffer)
	png.Encode(pngWriter, newRGBA)
	pngWriter.Flush()
	ide.sizeInBytes = uint32(len(pngBytesBuffer.Bytes()))

	newRGBABounds := newRGBA.Bounds()
	ide.imageWidth = uint8(newRGBABounds.Dx())
	ide.imageHeight = uint8(newRGBABounds.Dy())
	bytesBuffer := new(bytes.Buffer)

	binary.Write(bytesBuffer, binary.LittleEndian, id)
	binary.Write(bytesBuffer, binary.LittleEndian, ide)

	writerICO.Write(bytesBuffer.Bytes())
	writerICO.Write(pngBytesBuffer.Bytes())

	return "Done"
}

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Convert an image to its favicon formats.",
	Long: `Convert an image that you want to use as your favicon then use this 
	tool to convert an image to its favicon formats.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Requires `source` and `target`")
		}
		if len(args) >= 3 {
			return errors.New("Too many arguments provided")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		source := args[0]
		target := args[1]

		files := map[string]int{
			"apple-touch-icon-57x57.png":   57,
			"apple-touch-icon-60x60.png":   60,
			"apple-touch-icon-72x72.png":   72,
			"apple-touch-icon-76x76.png":   76,
			"apple-touch-icon-114x114.png": 114,
			"apple-touch-icon-120x120.png": 120,
			"apple-touch-icon-144x144.png": 144,
			"apple-touch-icon-152x152.png": 152,
			"favicon-16x16.png":            16,
			"favicon-32x32.png":            32,
			"favicon-96x96.png":            96,
			"favicon-128.png":              128,
			"favicon-196x196.png":          196,
			"mstile-70x70.png":             70,
			"ms-title-144x144.png":         144,
			"mstile-150x150.png":           150,
			"mstile-310x310.png":           310,
			"favicon.ico":                  64,
			// TODO: "mstile-310x150.png":
		}

		for filename, dimension := range files {
			filepath := filepath.Join(target, filename)

			if filename == "favicon.ico" {
				mediafile := MediaFile{filepath: source}
				pngFile, err := mediafile.OpenPNG()
				if err != nil {
					fmt.Println("Error")
				}
				writerICO, _ := os.Create(filepath)
				defer writerICO.Close()
				mediafile.SaveAsICO(writerICO, pngFile)
			} else {
				imagefile, err := imaging.Open(source)
				if err != nil {
					log.Fatalf("Failed to open image: %v", err)
				}
				imagefile = imaging.Resize(imagefile, dimension, dimension, imaging.Lanczos)
				err = imaging.Save(imagefile, filepath)
				if err != nil {
					log.Fatalf("Failed to save image: %v", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
}
