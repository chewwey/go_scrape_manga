/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"go_kaka_scrape/Manga"
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search manga in manganato(mangakakalot)",
	Run: func(cmd *cobra.Command, args []string) {
		page, err := Manga.MangaSearch(args[0])
		if err != nil {
			log.Fatal(err)
		}
		for i := range page {
			fmt.Println(strconv.Itoa(i+1) + ": " + page[i].Name)
		}
		fmt.Println("Enter Your Manga id: ")
		var intp int
		fmt.Scanln(&intp)
		man := Manga.MangaPage(page[intp-1].Url)
		man.DownloadMultiFromPage(1, 3)

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
