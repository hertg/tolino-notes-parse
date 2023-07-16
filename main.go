package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type highlight struct {
    title string
    author string
    page int
    highlight string
    date time.Time
}

// enum for modes
const (
    MODE_TITLE = iota
    MODE_HIGHLIGHT
    MODE_DATE
)

func main() {

    file, err := os.Open("notes.txt")
    if (err != nil) {
        panic(err)
    }
    scanner := bufio.NewScanner(file)
    
    mode := MODE_TITLE

    obj := highlight{}
    lesezeichen := 0
    list := []highlight{}
    line := ""
    readLine := false

    for scanner.Scan() {
        text := scanner.Text()

        if mode == MODE_TITLE {
            // skip empty lines
            if len(strings.TrimSpace(text)) == 0{
                continue
            }
            // skip separator
            if strings.HasPrefix(text, "----") {
                continue
            }
            re := regexp.MustCompile("(.*)\\s+\\((.*),\\s(.*)\\)")
            if !re.MatchString(text) {
                panic("unexpected line: " + text)
            }
            parts := re.FindStringSubmatch(text)
            obj.title = parts[1]
            obj.author = strings.Join([]string{parts[3], parts[2]}, " ")
            /*
            fmt.Printf("Title: %s\n", obj.title)
            fmt.Printf("Author: %s\n", obj.author)
            */
            mode = MODE_HIGHLIGHT
            continue
        }

        if mode == MODE_HIGHLIGHT {
            if readLine {
                if strings.HasPrefix(text, "Hinzugefügt am") {
                    line = strings.TrimSuffix(line, "\"")
                    obj.highlight = line
                    line = ""

                    re := regexp.MustCompile("Hinzugefügt am[  ]+([\\d\\.]+) \\| (.*)")
                    parts := re.FindStringSubmatch(text)
                    loc, _ := time.LoadLocation("Europe/Zurich")
                    t, _ := time.ParseInLocation("02.01.2006 15:04", strings.Join([]string{parts[1], parts[2]}, " "), loc)
                    obj.date = t

                    /*
                    fmt.Printf("Highlight: %s\n", obj.highlight)
                    fmt.Printf("Date: %s\n", obj.date)
                    */

                    // reset
                    list = append(list, obj)
                    obj = highlight{}
                    mode = MODE_TITLE
                    readLine = false

                    continue
                } else {
                    line += text
                }
            } else {
                re := regexp.MustCompile("(Markierung|Lesezeichen)[  ]+auf Seite[  ]+(\\d+)(\\-\\d+)*: \"(.*)(\")*")
                if re.MatchString(text) {
                    parts := re.FindStringSubmatch(text)
                    obj.page, _ = strconv.Atoi(parts[2])
                    line = parts[4]
                    readLine = true

                    /*
                    fmt.Printf("Page: %d\n", obj.page)
                    */

                    continue
                } 
            }
        }

    }
    
    
    f, err := os.Create("output.csv")
    defer f.Close()
    if err != nil {
        log.Fatal(err)
    }

    w := csv.NewWriter(f)
    defer w.Flush()

    // write header
    header := []string{"Highlight", "Title", "Author", "URL", "Note", "Location", "Date"}
    if err := w.Write(header); err != nil {
        log.Fatalln("error writing record to csv:", err)
    }

    for _, record := range list {
        var line []string
        line = append(line, record.highlight)
        line = append(line, record.title)
        line = append(line, record.author)
        line = append(line, "")
        line = append(line, "")
        line = append(line, strconv.Itoa(record.page))
        line = append(line, record.date.UTC().Format("2006-01-02 15:04:05"))
        if err := w.Write(line); err != nil {
            log.Fatalln("error writing record to csv:", err)
        }
    }

/*
    for _, h := range list {
        j, _ := json.Marshal(h)
        _ = j
        fmt.Println(h)

    }
    */


    fmt.Printf("There are %d Lesezeichen\n", lesezeichen)

    println("Hello, world!")
}
