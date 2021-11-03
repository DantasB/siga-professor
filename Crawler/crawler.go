package crawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	utils "github.com/DantasB/Siga-Professor/Utils"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

var siraUrl = "https://siga.ufrj.br/sira/gradeHoraria/"

func getSiraCoursesList() ([]string, error) {
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

	htmlDocument, err := htmlquery.Parse(strings.NewReader(utils.ToUtf8(body)))
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	list := htmlquery.Find(htmlDocument, "//td//a[contains(@href, '.html')]")
	for _, node := range list {
		courses = append(courses, htmlquery.SelectAttr(node, "href"))
	}

	return courses, nil
}

func AccessSiraCourses() {
	courses, err := getSiraCoursesList()
	if err != nil {
		log.Fatalln(err)
	}

	for _, course := range courses {
		accessSiraCourse(course)
	}
}

func accessSiraCourse(courseUrl string) ([]string, error) {
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

	htmlDocument, err := htmlquery.Parse(strings.NewReader(utils.ToUtf8(body)))
	if err != nil {
		log.Fatalln(err)
		return []string{}, err
	}

	if !courseUpdated(htmlDocument) {
		return []string{}, nil
	}

	if !parseCourseInformation(htmlDocument, &disciplines) {
		return []string{}, nil
	}

	return disciplines, nil
}

func parseCourseInformation(htmlDocument *html.Node, disciplines *[]string) bool {
	parseDisciplines(htmlDocument, disciplines)
	return true
}

func parseDisciplines(htmlDocument *html.Node, disciplines *[]string) {
	disciplineNodes := htmlquery.Find(htmlDocument, "//td//td//table[@class='lineBorder']//table//tr[@class != 'tableTitle']")
	for _, node := range disciplineNodes {
		line := htmlquery.InnerText(node)
		line = utils.ReplaceMultipleSpacesByPipe(utils.RemoveSeparators(line))
		if strings.Contains(line, "Lista de Disciplinas") || strings.Contains(line, "Nome Turma") {
			continue
		}

		fmt.Printf("%s\n", line)
	}
}

func courseUpdated(htmlDocument *html.Node) bool {
	courseNode := htmlquery.FindOne(htmlDocument, "//font[@size=3]/text()")
	courseName := strings.TrimSpace(htmlquery.InnerText(courseNode))
	courseName = strings.TrimSpace(strings.Split(courseName, "-")[0])

	htmlNode := htmlquery.FindOne(htmlDocument, "//td[@align='right']/text()[2]")
	datetime := strings.TrimSpace(htmlquery.InnerText(htmlNode))
	parsedTime, err := utils.ParseDate(datetime)
	if err != nil {
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
