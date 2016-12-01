package output

import (
	"io"
	"log"
	"net/http"
	"os"
)

// SaveImages persist images in filesystem
func SaveImages(images <-chan map[string]string, out string) error {
	for image := range images {
		for imageName, imageURL := range image {
			f, err := os.Create(out + "/" + imageName)

			if err != nil {
				return err
			}

			defer f.Close()
			res, err := http.Get(imageURL)

			if err != nil {
				return err
			}

			defer res.Body.Close()

			if _, err := io.Copy(f, res.Body); err != nil {
				return err
			}

			log.Printf("%s downloaded", imageName)
		}
	}

	return nil
}
