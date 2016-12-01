package manga

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const baseURL = "http://www.mangareader.net"

// Manga is a struct that holds the logic about how download HQ from MangaReader
type Manga struct {
	URL            string
	ChapterNode    string
	PageNode       string
	ImageNode      string
	ParallelPages  int
	ParallelImages int
	OutputFunc     func(<-chan map[string]string, string) error
}

// FetchManga gets all manga pages basead on an URL
func (m *Manga) FetchManga(out string) error {
	done := make(chan struct{})
	defer close(done)

	chapters, errc := m.getChapters(done)
	pages := m.getPages(done, chapters)
	images := m.getImages(done, pages)
	err := m.outputFunc(images, out)

	if err := <-errc; err != nil {
		return err
	}

	return err
}

// buildImageName gets an image node and extract useful information about it,
// like title, chapter number, page number and return a name with it
func (m *Manga) buildImageName(imageNode *goquery.Selection) (name string, src string) {
	alt, altExist := imageNode.Attr("alt")
	src, srcExist := imageNode.Attr("src")

	if !altExist || !srcExist {
		return name, src
	}

	ext := path.Ext(src)
	fields := strings.Fields(alt)
	title := fields[0]
	chapter := fmt.Sprintf("%05s", fields[1])
	page := fmt.Sprintf("%05s", fields[4])

	return fmt.Sprintf("%s-Chap-%s-Pg-%s%s", title, chapter, page, ext), src
}

// getChapters get page links for a specific chapter
func (m *Manga) getChapters(done chan struct{}) (<-chan string, <-chan error) {
	chapters := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(chapters)

		doc, err := goquery.NewDocument(m.URL)

		if err != nil {
			errc <- err
			return
		}

		doc.Find(m.ChapterNode).Each(func(chapterID int, chapterNode *goquery.Selection) {
			chapterURL, _ := chapterNode.Attr("href")

			select {
			case chapters <- baseURL + chapterURL:
			case <-done:
				errc <- errors.New("getChapters canceled")
			}
		})

		errc <- nil
	}()

	return chapters, errc
}

// getImages spawn goroutines for fetch images URLs
func (m *Manga) getImages(done <-chan struct{}, in <-chan string) <-chan map[string]string {
	var wg sync.WaitGroup
	images := make(chan map[string]string)

	wg.Add(m.ParallelImages)
	for i := 0; i < m.ParallelImages; i++ {
		go func() {
			m.fetchImages(done, in, images)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(images)
	}()

	return images
}

// getPages spawn goroutines for fetch page URLs
func (m *Manga) getPages(done <-chan struct{}, in <-chan string) <-chan string {
	var wg sync.WaitGroup
	pages := make(chan string)

	wg.Add(m.ParallelPages)
	for i := 0; i < m.ParallelPages; i++ {
		go func() {
			m.fetchPages(done, in, pages)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(pages)
	}()

	return pages
}

// fetchImages given a page URL extract an image link from it
func (m *Manga) fetchImages(done <-chan struct{}, pages <-chan string, images chan<- map[string]string) {
	for page := range pages {
		doc, _ := goquery.NewDocument(page)
		image := doc.Find(m.ImageNode).First()
		name, src := m.buildImageName(image)

		select {
		case images <- map[string]string{name: src}:
		case <-done:
			return
		}
	}
}

// fetchPages given a chapter URL extract a page link from it
func (m *Manga) fetchPages(done <-chan struct{}, chapters <-chan string, pages chan<- string) {
	for chapter := range chapters {
		page, err := goquery.NewDocument(chapter)

		if err != nil {
			return
		}

		page.Find(m.PageNode).Each(func(pageID int, pageNode *goquery.Selection) {
			pageURL, exist := pageNode.Attr("value")

			if !exist {
				return
			}

			select {
			case pages <- baseURL + pageURL:
			case <-done:
				return
			}
		})
	}
}

// outputFunc call the OutputFunc on Manga struct or retuning nothing
func (m *Manga) outputFunc(images <-chan map[string]string, out string) error {
	if m.OutputFunc != nil {
		return m.OutputFunc(images, out)
	}

	return nil
}
