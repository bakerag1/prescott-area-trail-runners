package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// newsCmd represents the news command
var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Create a newsletter for this month",
	Long: `Create a newsletter for the current month, including events from the run date, to add
	  an announcement section, add content to the created markdown file`,
	Run: func(cmd *cobra.Command, args []string) {
		newsletter()
	},
}

func init() {
	rootCmd.AddCommand(newsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func newsletter() {

	content := `---
type: news
date: %v
title: PATR Chatter %v
feature-img: "assets/img/big-trail.jpg"
---
`

	f, err := os.Create("site/content/post/" + time.Now().Local().Format("2006-01") + "-news.md")
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	_, err = writer.WriteString(fmt.Sprintf(content, time.Now().UTC().Format(time.RFC3339), time.Now().Format("January 2006")))
	if err != nil {
		panic(err)
	}
}
