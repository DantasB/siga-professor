package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var siraUrl = "https://siga.ufrj.br/sira/gradeHoraria/"

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

func GetSiraCoursesList() ([]string, error) {
	var courses []string

	resp, err := http.Get(siraUrl)
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	htmlDoc, err := htmlquery.Parse(strings.NewReader(toUtf8(body)))

	list := htmlquery.Find(htmlDoc, "//td//a[contains(@href, '.html')]")
	for _, node := range list {
		courses = append(courses, htmlquery.SelectAttr(node, "href"))
	}

	return courses, nil
}

func AccessSiraCourses() {
	courses, err := GetSiraCoursesList()
	if err != nil {
		log.Fatalln(err)
	}

	for _, course := range courses {
		AccessSiraCourse(course)
	}
}

func AccessSiraCourse(courseUrl string) ([]string, error) {
	var disciplines []string
	resp, err := http.Get(siraUrl + courseUrl)
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	htmlDoc, err := htmlquery.Parse(strings.NewReader(toUtf8(body)))

	if !CourseUpdated(htmlDoc) {
		return []string{}, nil
	}

	return disciplines, nil
}

func CourseUpdated(htmlDocument *html.Node) bool {
	courseNode := htmlquery.FindOne(htmlDocument, "//font[@size=3]/text()")
	courseName := strings.TrimSpace(htmlquery.InnerText(courseNode))
	courseName = strings.TrimSpace(strings.Split(courseName, "-")[0])

	htmlNode := htmlquery.FindOne(htmlDocument, "//td[@align='right']/text()[2]")
	datetime := strings.TrimSpace(htmlquery.InnerText(htmlNode))
	splittedDate := strings.Split(datetime, " ")
	if len(splittedDate) != 4 {
		fmt.Println("Couldn't Parse the string:", datetime)
		return false
	}

	layout := "02/01/2006"
	parsedTime, err := time.Parse(layout, splittedDate[2])
	if err != nil {
		fmt.Println(err)
		return false
	}

	currentYear, currentMonth, _ := time.Now().Date()
	lastUpdateYear, lastUpdateMonth, _ := parsedTime.Date()

	if currentYear > lastUpdateYear || currentMonth > lastUpdateMonth+1 {
		fmt.Printf("[%s] This is an old course and should be ignored. %s\n", courseName, datetime)
		return false
	}

	return true
}
