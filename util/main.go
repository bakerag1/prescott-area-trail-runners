package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/apognu/gocal"
)

const outputFmt = `---
title: %v
date: %v
external_url: %v
description: %v
---
`

func main() {
	switch os.Args[1] {
	case "calendar":
		addCalendarItems()
	case "weekly":
		postWeeklyCalendar()
	default:
		log.Fatal("no valid argument passed")
	}
}

func postWeeklyCalendar() {
	in, err := os.Open("util/this-week.md")
	if err != nil {
		log.Fatal("unable to open template", err)
	}
	defer in.Close()
	newname := time.Now().Format("2006-01-02") + "-this-week.md"
	log.Printf("creating post %s", newname)
	out, err := os.Create("_posts/" + newname)
	if err != nil {
		log.Fatal("unable to create post", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Fatal("unable to create post", err)
	}
	log.Printf("created file: %s", newname)
}

func addCalendarItems() {
	time.LoadLocation("America/Phoenix")
	err := os.RemoveAll("_calendar/*")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open("events.ics")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := gocal.NewParser(f)
	start, end := time.Now().Add(-24*time.Hour), time.Now().Add(120*30*24*time.Hour)
	c.Start, c.End = &start, &end
	c.Strict.Mode = gocal.StrictModeFailAttribute
	c.Parse()
	for _, e := range c.Events {
		// if e.Class != "PUBLIC" {
		// 	log.Printf("non-public event skipped: %s\n", e.Summary)
		// 	continue
		// }
		uri := e.URL
		uid := e.Uid
		uid = uid[:strings.IndexByte(uid, '@')]
		description := strings.ReplaceAll(e.Description, "\\n", "<br>\n  ")
		description = strings.ReplaceAll(description, uri, "")
		description = strings.ReplaceAll(description, ":", "&#58;")
		description = strings.ReplaceAll(description, "\n\n", "<br>\n  ")
		cal, err := os.Create("_calendar/" + uid + ".md")
		defer cal.Close()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("creating event: %s: %s - %s\n", e.Start.Local().Format("2006-01-02"), e.Uid, e.Summary)
		cal.Write([]byte(fmt.Sprintf(outputFmt, e.Summary, e.Start.Local().Format("2006-01-02 15:04"), uri, description)))
	}
}
