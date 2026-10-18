// Command fixture_covers writes the placeholder blog covers that fixtures/
// reference (active_storage_blobs.yml) into the local fileblob bucket, so a
// freshly set-up development site shows images without any real uploads.
// Bytes are not committed: storage/ is git-ignored, this regenerates them.
package main

import (
	"bytes"
	"crypto/md5" //nolint:gosec // fileblob's attrs file stores the MD5, not a security hash
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

const width, height = 1200, 600

// covers maps a blob key from the fixtures to the two corners of its gradient.
var covers = map[string][2]color.RGBA{
	"blogcoverkey001": {{79, 70, 229, 255}, {14, 165, 233, 255}},
	"blogcoverkey002": {{244, 63, 94, 255}, {251, 146, 60, 255}},
	"blogcoverkey003": {{16, 185, 129, 255}, {59, 130, 246, 255}},
}

func main() {
	dir := "storage"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gosec // G301: a developer script writing public sample covers into the local storage dir
		log.Fatal(err)
	}
	for key, c := range covers {
		data := gradient(c[0], c[1])
		if err := os.WriteFile(filepath.Join(dir, key), data, 0o644); err != nil { //nolint:gosec // G306: public sample cover, served as a static file
			log.Fatal(err)
		}
		sum := md5.Sum(data) //nolint:gosec // see import
		attrs, err := json.Marshal(map[string]any{
			"user.cache_control": "", "user.content_disposition": "", "user.content_encoding": "",
			"user.content_language": "", "user.content_type": "image/png", "user.metadata": nil,
			"md5": base64.StdEncoding.EncodeToString(sum[:]),
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, key+".attrs"), attrs, 0o644); err != nil { //nolint:gosec // G306: attrs of a public sample cover
			log.Fatal(err)
		}
	}
}

func gradient(from, to color.RGBA) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := range height {
		for x := range width {
			t := (float64(x)/width + float64(y)/height) / 2
			mix := func(a, b uint8) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*t) }
			img.SetRGBA(x, y, color.RGBA{mix(from.R, to.R), mix(from.G, to.G), mix(from.B, to.B), 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}
