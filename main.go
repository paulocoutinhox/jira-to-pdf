package main

import (
	"flag"
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"log"
	"io/ioutil"
	"github.com/jung-kurt/gofpdf"
	"strings"
	"strconv"
	"golang.org/x/text/encoding/charmap"
	"time"
)

var (
	paramInstance string
	paramUsername string
	paramPassword string
	paramQuery string
	paramVerbose bool
	paramOutputFilename string
	paramDocumentTitle string
	paramIssueTemplate string
	paramDateTimeFormat string
	pdfGenerator *gofpdf.Fpdf
	pdfTR func(string) string
	apiDateTimeFormat string
)

func main() {
	// parse params
	flag.StringVar(&paramInstance, "i", "", "URL of your Jira instance, like: https://your.jira-instance.com")
	flag.StringVar(&paramUsername, "u", "", "Username of your Jira instance account")
	flag.StringVar(&paramPassword, "p", "", "Password of your Jira instance account")
	flag.StringVar(&paramQuery, "q", "", "JQL query for issues search")
	flag.StringVar(&paramOutputFilename, "o", "", "Output filename")
	flag.StringVar(&paramDocumentTitle, "t", "", "Document title")
	flag.StringVar(&paramIssueTemplate, "it", "", "Issue template")
	flag.StringVar(&paramDateTimeFormat, "dft", "", "DateTime format")
	flag.BoolVar(&paramVerbose, "v", false, "Verbose mode")
	flag.Parse()

	if len(paramInstance) == 0 {
		log.Fatal("Jira instance is not defined")
	}

	if len(paramUsername) == 0 {
		log.Fatal("Jira instance username is not defined")
	}

	if len(paramPassword) == 0 {
		log.Fatal("Jira instance password is not defined")
	}

	if len(paramQuery) == 0 {
		log.Fatal("Jira issue query (JQL) is not defined")
	}

	if len(paramOutputFilename) == 0 {
		paramOutputFilename = "jira-issues.pdf"
	}

	if len(paramDocumentTitle) == 0 {
		paramDocumentTitle = "JIRA Issues"
	}

	if len(paramIssueTemplate) == 0 {
		paramIssueTemplate = "<b>Issue:</b> [issue.key]<br /><b>Summary:</b> [issue.fields.summary]<br /><b>Assignee:</b> [issue.fields.assignee.name]<br /><b>Status:</b> [issue.fields.status.name]<br /><b>Created:</b> [issue.fields.created]"
	}

	if len(paramDateTimeFormat) == 0 {
		paramDateTimeFormat = "2006-01-02 15:04:05"
	}

	apiDateTimeFormat = "2006-01-02T15:04:05.999999999-0700";

	// authenticate
	jiraClient, err := jira.NewClient(nil, paramInstance)

	if err != nil {
		log.Fatalf("Auth error! %v", err)
	}

	jiraClient.Authentication.SetBasicAuth(paramUsername, paramPassword)

	// search for issues using query
	issues, response, err := jiraClient.Issue.Search(paramQuery, nil)

	if err != nil {
		if paramVerbose {
			var body = ""

			if b, err := ioutil.ReadAll(response.Body); err == nil {
				body = string(b)
			}

			log.Fatalf("Search error! %v\n\nBody: %v", err, body)
		} else {
			log.Fatalf("Search error! %v", err)
		}
	}

	if paramVerbose {
		log.Printf("Total of issues: %v\n", len(issues))
		fmt.Printf("Issues: %s\n", issues)
	}

	// generate PDF
	pdfGenerator = gofpdf.New("P", "mm", "A4", "")
	pdfTR = pdfGenerator.UnicodeTranslatorFromDescriptor("")

	// document title
	pdfGenerator.AddPage()
	pdfGenerator.SetFont("Arial", "B", 16)
	pdfGenerator.SetFillColor(222, 222, 222)
	pdfGenerator.SetTextColor(0, 0, 0)
	pdfGenerator.MultiCell(0, 16, pdfTR(paramDocumentTitle), "1", "C", true)

	pdfGenerator.Ln(8)

	// issues
	for index, issue := range issues {
		// set issue font
		var lineHt float64 = 9
		pdfGenerator.SetFont("Arial", "", lineHt)
		pdfGenerator.SetTextColor(0, 0, 0)

		// parse issue template
		issueText := parseIssueTemplate(index, issue)
		issueText, _ = charmap.Windows1252.NewEncoder().String(issueText)

		html := pdfGenerator.HTMLBasicNew()
		html.Write(lineHt, issueText)

		// draw issue separator
		pdfGenerator.SetDrawColor(195, 195, 195)
		pdfGenerator.Ln(lineHt)
		pdfGenerator.Ln(2)

		pageWidth, _ := pdfGenerator.GetPageSize()
		x, y := pdfGenerator.GetXY()
		marginL, marginR, _, _ := pdfGenerator.GetMargins()
		pdfGenerator.Line(x, y, x + pageWidth - marginR - marginL, y)

		pdfGenerator.Ln(2)
	}

	// save to output filename
	err = pdfGenerator.OutputFileAndClose(paramOutputFilename)

	if err != nil {
		log.Fatalf("Erro while save PDF: %v", err)
	}
}

func parseIssueTemplate(index int, issue jira.Issue) string {
	issueText := paramIssueTemplate

	issueText = strings.Replace(issueText, "[issue.key]", issue.Key, -1)
	issueText = strings.Replace(issueText, "[issue.id]", issue.ID, -1)
	issueText = strings.Replace(issueText, "[issue.fields.description]", issue.Fields.Description, -1)
	issueText = strings.Replace(issueText, "[issue.fields.duedate]", issue.Fields.Duedate, -1)
	issueText = strings.Replace(issueText, "[issue.fields.expand]", issue.Fields.Expand, -1)
	issueText = strings.Replace(issueText, "[issue.fields.resolutiondate]", issue.Fields.Resolutiondate, -1)
	issueText = strings.Replace(issueText, "[issue.fields.summary]", issue.Fields.Summary, -1)
	issueText = strings.Replace(issueText, "[issue.fields.timeestimate]", strconv.Itoa(issue.Fields.TimeEstimate), -1)
	issueText = strings.Replace(issueText, "[issue.fields.timeoriginalestimate]", strconv.Itoa(issue.Fields.TimeOriginalEstimate), -1)
	issueText = strings.Replace(issueText, "[issue.fields.timespent]", strconv.Itoa(issue.Fields.TimeSpent), -1)

	issueText = strings.Replace(issueText, "[issue.fields.project.name]", issue.Fields.Project.Name, -1)
	issueText = strings.Replace(issueText, "[issue.fields.project.description]", issue.Fields.Project.Description, -1)
	issueText = strings.Replace(issueText, "[issue.fields.project.id]", issue.Fields.Project.ID, -1)
	issueText = strings.Replace(issueText, "[issue.fields.project.key]", issue.Fields.Project.Key, -1)

	issueText = strings.Replace(issueText, "[issue.fields.type.name]", issue.Fields.Type.Name, -1)
	issueText = strings.Replace(issueText, "[issue.fields.type.description]", issue.Fields.Type.Description, -1)
	issueText = strings.Replace(issueText, "[issue.fields.type.id]", issue.Fields.Type.ID, -1)

	if issue.Fields.Priority != nil {
		issueText = strings.Replace(issueText, "[issue.fields.priority.id]", issue.Fields.Priority.ID, -1)
		issueText = strings.Replace(issueText, "[issue.fields.priority.name]", issue.Fields.Priority.Name, -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.priority.id]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.priority.name]", "", -1)
	}

	if issue.Fields.AggregateProgress != nil {
		issueText = strings.Replace(issueText, "[issue.fields.aggregateprogress.progress]", strconv.Itoa(issue.Fields.AggregateProgress.Progress), -1)
		issueText = strings.Replace(issueText, "[issue.fields.aggregateprogress.total]", strconv.Itoa(issue.Fields.AggregateProgress.Total), -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.aggregateprogress.progress]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.aggregateprogress.total]", "", -1)
	}

	if issue.Fields.Progress != nil {
		issueText = strings.Replace(issueText, "[issue.fields.progress.progress]", strconv.Itoa(issue.Fields.Progress.Progress), -1)
		issueText = strings.Replace(issueText, "[issue.fields.progress.total]", strconv.Itoa(issue.Fields.Progress.Total), -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.progress.progress]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.progress.total]", "", -1)
	}

	if issue.Fields.Assignee != nil {
		issueText = strings.Replace(issueText, "[issue.fields.assignee.name]", issue.Fields.Assignee.Name, -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.emailaddrress]", issue.Fields.Assignee.EmailAddress, -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.displayname]", issue.Fields.Assignee.DisplayName, -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.key]", issue.Fields.Assignee.Key, -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.assignee.name]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.emailaddrress]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.displayname]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.assignee.key]", "", -1)
	}

	if issue.Fields.Creator != nil {
		issueText = strings.Replace(issueText, "[issue.fields.creator.name]", issue.Fields.Creator.Name, -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.emailaddrress]", issue.Fields.Creator.EmailAddress, -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.displayname]", issue.Fields.Creator.DisplayName, -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.key]", issue.Fields.Creator.Key, -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.creator.name]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.emailaddrress]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.displayname]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.creator.key]", "", -1)
	}

	if issue.Fields.Reporter != nil {
		issueText = strings.Replace(issueText, "[issue.fields.reporter.name]", issue.Fields.Reporter.Name, -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.emailaddrress]", issue.Fields.Reporter.EmailAddress, -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.displayname]", issue.Fields.Reporter.DisplayName, -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.key]", issue.Fields.Reporter.Key, -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.reporter.name]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.emailaddrress]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.displayname]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.reporter.key]", "", -1)
	}

	if issue.Fields.Status != nil {
		issueText = strings.Replace(issueText, "[issue.fields.status.name]", issue.Fields.Status.Name, -1)
		issueText = strings.Replace(issueText, "[issue.fields.status.description]", issue.Fields.Status.Description, -1)
		issueText = strings.Replace(issueText, "[issue.fields.status.id]", issue.Fields.Status.ID, -1)
	} else {
		issueText = strings.Replace(issueText, "[issue.fields.status.name]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.status.description]", "", -1)
		issueText = strings.Replace(issueText, "[issue.fields.status.id]", "", -1)
	}

	if len(issue.Fields.Created) > 0 {
		dateTime, err := time.Parse(apiDateTimeFormat, issue.Fields.Created)

		if err == nil {
			issueText = strings.Replace(issueText, "[issue.fields.created]", dateTime.Format(paramDateTimeFormat), -1)
		} else {
			if paramVerbose {
				log.Printf("Error on parse field \"created\"! %v\n", err)
			}

			issueText = strings.Replace(issueText, "[issue.fields.created]", "", -1)
		}
	}

	return issueText;
}