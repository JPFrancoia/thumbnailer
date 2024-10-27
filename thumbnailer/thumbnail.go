package thumbnailer

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"golang.org/x/image/draw"
)

// Creates a resized image from the reader and writes it to the writer.
// The 'size' parameter is the longest dimension of the new image, i.e: if the
// original image is 100x200 and 'size' is 50, the new image will be 25x50. The aspect
// ratio is maintained.
func Thumbnail(r io.Reader, w io.Writer, size int) error {
	buf := bufio.NewReader(r)
	headerBytes, err := buf.Peek(512)
    if len(headerBytes) == 0 && err != nil {
        return fmt.Errorf("Error reading header bytes: %v", err)
    }
	mimetype := detectContentType(headerBytes)

	var src image.Image

	switch mimetype {
	case "image/jpeg":
		src, err = jpeg.Decode(buf)
	case "image/png":
		src, err = png.Decode(buf)
	}

	if err != nil {
		return fmt.Errorf("Error decoding image: %v", err)
	}

	// Calculate new dimensions while maintaining the aspect ratio
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Max.X
	srcHeight := srcBounds.Max.Y

	scale := float64(size) / float64(srcWidth)
	if srcHeight > srcWidth {
		scale = float64(size) / float64(srcHeight)
	}

	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	draw.CatmullRom.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	err = jpeg.Encode(w, dst, nil)
	if err != nil {
		return err
	}

	return nil
}

// Use the first 512 bytes of the file to determine the content type.
// Returns "application/octet-stream" if no other content type is detected.
func detectContentType(fb []byte) string {
	return http.DetectContentType(fb[:512])
}
