Manga Downloader
========

MangaReader.net crawler based on [akitaonrails](https://github.com/akitaonrails/manga-downloadr) version.
There is some features to implement like create the PDF. Feel free to submit a PR!

## Install

```bash
go get github.com/marcosinger/manga_downloader
```

## Command Use

Download One Punch Man manga
```bash
$ manga_downloader -u "http://www.mangareader.net/onepunch-man"
```

Changing the default output folder (default: mangareader)
```bash
$ manga_downloader -u [manga URL] -o [destination]
```

Increment parallel process (default: 20)
```bash
$ manga_downloader -u [manga URL] -c 40
```
