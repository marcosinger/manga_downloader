package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/marcosinger/manga_downloader/manga"
	"github.com/marcosinger/manga_downloader/output"
)

var mangaURL, outPath string
var parallel int

func init() {
	flag.StringVar(&mangaURL, "u", "", "Manga URL at MangaReader.net")
	flag.StringVar(&outPath, "o", "mangareader", "Output folder")
	flag.IntVar(&parallel, "c", 20, "Parallel process")
}

func main() {
	flag.Parse()

	if mangaURL == "" {
		flag.Usage()
		return
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		log.Fatalf("%s directory not found. Select another one using the option -o. Call -h for help.", outPath)
	}

	start := time.Now()
	m := manga.Manga{
		URL:            mangaURL,
		ChapterNode:    "#listing a",
		PageNode:       "#selectpage select#pageMenu option",
		ImageNode:      "#img",
		ParallelPages:  parallel,
		ParallelImages: parallel,
		OutputFunc:     output.SaveImages,
	}

	if err := m.FetchManga(outPath); err != nil {
		log.Print(err)
	}

	log.Printf("Done! Please check the %s folder. Elapsed time %s\n", outPath, time.Now().Sub(start))
}
