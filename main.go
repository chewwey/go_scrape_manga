package main

import (
	"go_kaka_scrape/Manga"
)

func main() {
	man := Manga.MangaPage("https://chapmanganato.com/manga-dn980422")
	man.DownloadMultiFromPage(3, 4)
}
