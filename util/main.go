package main

import (
	"fmt"
	"github.com/apognu/gocal"
	"log"
	"os"
	"strings"
	"time"
)

const outputFmt = `---
title: %v
date: %v
external_url: %v
description: %v
---
`

func main() {
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
		fmt.Println(uri)
		cal.Write([]byte(fmt.Sprintf(outputFmt, e.Summary, e.Start.Local().Format("2006-01-02 15:04"), uri, description)))
	}
}
