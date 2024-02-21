package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/bakerag1/gocal"
	"gopkg.in/yaml.v2"
)

const outputFmt = `---
title: "%v"
date: %v
startdate: %v
enddate: %v
external_url: %v
layout: %v
location: %v
feature-img: "assets/img/big-trail.jpg"
---

%v
`

func main() {
	switch os.Args[1] {
	case "calendar":
		addCalendarItems()
	case "weekly":
		postWeeklyCalendar()
	case "newsletter":
		newsletter()
	default:
		log.Fatal("no valid argument passed")
	}
}

func postWeeklyCalendar() {
	newname := time.Now().Format("2006-01-02") + "-this-week.md"
	title := time.Now().Format("Jan 2, 2006")
	log.Printf("creating post %s", newname)
	out, err := os.Create("content/post/" + newname)
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
	err := os.RemoveAll("content/events/*")
	if err != nil {
		log.Fatal(err)
	}
	events := parseEvents()
	for _, e := range events {
		cal, err := os.Create("content/events/" + e.Uid + ".md")
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
				e.Uri,
				layout,
				e.Location,
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
		if e.Class != "PUBLIC" {
			log.Printf("non-public event skipped: %s\n", e.Summary)
			continue
		}
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
		if strings.Compare(ev.Summary, "Thursday Night Social Run/Walk") == 0 {
			ev.Summary = "CANCELLED - Thursday Night Social Run/Walk"
		}
		ev.Start = e.Start.Local().Format("2006-01-02T15:04:00Z")
		ev.End = e.End.Local().Format("2006-01-02T15:04:00Z")
		ev.Created = e.Created.Format("2006-01-02T15:04:00Z")
		events = append(events, ev)
	}
	sortByStart := func(a, b int) bool {
		return strings.Compare(events[a].Start, events[b].Start) < 0
	}
	sort.Slice(events, sortByStart)
	return events
}

func newsletter() {

	paths := []string{
		"util/newsletter.tmpl",
	}

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
		"modIsZero": func(i int, m int) bool {
			return i%m == 0 && i != 0
		},
	}

	var cfg config
	cfg.getConf()
	cfg.Month = time.Now().Local().Format("January")
	cfg.Year = time.Now().Local().Format("2006")
	cfg.Now = time.Now().UTC().Format(time.RFC3339)
	f, err := os.Create("content/post/" + time.Now().Local().Format("2006-01-02") + "-newsletter.html")
	if err != nil {
		panic(err)
	}
	tData := struct {
		Config    config
		MonthInfo monthData
	}{
		Config: cfg,
		MonthInfo: monthData{
			Events: parseEvents(),
		},
	}
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	t := template.Must(template.New("newsletter").Funcs(funcMap).ParseFiles(paths...))
	err = t.ExecuteTemplate(writer, "newsletter.tmpl", tData)
	if err != nil {
		panic(err)
	}
}

func (c *config) getConf() *config {

	yamlFile, err := ioutil.ReadFile("util/newsletter-config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

type link struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
type config struct {
	Links []link `yaml:"links"`
	Month string
	Year  string
	Now   string
}
type monthData struct {
	Events []event
	News   []news
}
type event struct {
	Name        string
	Uri         string
	PatrUri     string
	Description string
	Summary     string
	Start       string
	End         string
	Location    string
	Uid         string
	Created     string
}
type news struct {
}
