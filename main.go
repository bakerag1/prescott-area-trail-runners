package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/bakerag1/gocal"
	"prescottareatrailrunners.com/patr/cmd"
)

const outputFmt = `---
title: "%v"
date: %v
startdate: %v
enddate: %v
patr: %v
external_url: %v
layout: %v
location: %v
feature-img: "assets/img/big-trail.jpg"
lastmod: %v
outputs:
  - html
  - calendar
ICSDescription: |+2
  %s
---

%v
`

func main() {
	switch os.Args[1] {
	case "calendar":
		addCalendarItems()
	case "weekly":
		postWeeklyCalendar()
	default:
		cmd.Execute()
	}
}

func postWeeklyCalendar() {
	newname := time.Now().Format("2006-01-02") + "-this-week.md"
	title := time.Now().Format("Jan 2, 2006")
	log.Printf("creating post %s", newname)
	out, err := os.Create("site/content/post/" + newname)
	if err != nil {
		log.Fatal("unable to create post", err)
	}
	defer out.Close()
	t := template.Must(template.New("this-week.md").ParseFiles("util/this-week.md"))
	fw := bufio.NewWriter(out)
	err = t.Execute(fw, struct {
		Date string
		Now  string
	}{
		Date: title,
		Now:  time.Now().Format("2006-01-02"),
	})
	if err != nil {
		log.Fatal("unable to apply post template", err)
	}
	err = fw.Flush()
	if err != nil {
		log.Fatal("unable to create post", err)
	}
	log.Printf("created file: %s", newname)
}

func addCalendarItems() {
	time.LoadLocation("America/Phoenix")
	err := os.RemoveAll("site/content/events/*")
	if err != nil {
		log.Fatal(err)
	}
	events := parseEvents()
	for _, e := range events {
		cal, err := os.Create("site/content/events/" + e.Uid + ".md")
		defer cal.Close()
		if err != nil {
			log.Fatal(err)
		}

		layout := "post"

		log.Printf("creating event: %s: %s - %s\n", e.Start, e.Uid, e.Summary)

		cal.Write([]byte(
			fmt.Sprintf(outputFmt,
				e.Summary,
				e.Created,
				e.Start,
				e.End,
				e.PATR,
				e.Uri,
				layout,
				e.Location,
				time.Now().UTC().Format(time.RFC3339),
				e.ICSDescription,
				e.Description)))
	}
}

func parseEvents() []event {
	f, err := os.Open("events.ics")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	c := gocal.NewParser(f)
	start, end := time.Now().Add(-48*time.Hour), time.Now().Add(120*30*24*time.Hour)
	c.Start, c.End = &start, &end
	c.Strict.Mode = gocal.StrictModeFailAttribute
	c.Parse()
	var events []event
	for _, e := range c.Events {
		var ev event
		if e.Class != "PUBLIC" ||
			e.Organizer.Cn == "Whiskey Row Brew Club" ||
			e.Organizer.Cn == "Park Plaza Liquor and Deli" {
			log.Printf("non-public event skipped: %s\n", e.Summary)
			continue
		}
		var icsDesc string
		for i, s := range e.Description {
			icsDesc += string(s)
			if i > 0 && (i+1)%32 == 0 {
				icsDesc += "\r  "
			}
		}
		ev.ICSDescription = icsDesc
		ev.Uri = e.URL
		ev.Uid = e.Uid[:strings.IndexByte(e.Uid, '@')]
		ev.Description = strings.ReplaceAll(e.Description, "\\n", "<br>\n  ")
		ev.Description = strings.ReplaceAll(ev.Description, ev.Uri, "")
		ev.Description = strings.ReplaceAll(ev.Description, ":", "&#58;")
		ev.Description = strings.ReplaceAll(ev.Description, "\n\n", "<br>\n  ")
		expr := regexp.MustCompile(`(http[s]?)&#58;(//[^ )]*)`)
		ev.Description = expr.ReplaceAllString(ev.Description, "[$1:$2]($1:$2)")
		expr2 := regexp.MustCompile(`^([^A-z0-9]*)(.*)`)
		ev.Summary = expr2.ReplaceAllString(e.Summary, "$2")
		ev.Start = e.Start.Local().Format("2006-01-02T15:04:00Z")
		ev.DayOfMonth = e.Start.Local().Format("02")
		ev.Month = e.Start.Local().Format("Jan")
		ev.DayOfWeek = e.Start.Local().Format("Mon")
		ev.End = e.End.Local().Format("2006-01-02T15:04:00Z")
		ev.Created = e.Created.Format("2006-01-02T15:04:00Z")
		ev.PATR = e.PartStat == "ACCEPTED"
		events = append(events, ev)
	}
	sortByStart := func(a, b int) bool {
		return strings.Compare(events[a].Start, events[b].Start) < 0
	}
	sort.Slice(events, sortByStart)
	return events
}

type event struct {
	Name           string
	Uri            string
	PatrUri        string
	ICSDescription string
	Description    string
	Summary        string
	Start          string
	Month          string
	DayOfWeek      string
	DayOfMonth     string
	End            string
	Location       string
	Uid            string
	Created        string
	PATR           bool
}
