package main

import (
	"flag"
	"fmt"

	jira "github.com/andygrunwald/go-jira"
	"log"
	"io/ioutil"
	"github.com/jung-kurt/gofpdf"
)

var (
	paramInstance string
	paramUsername string
	paramPassword string
	paramQuery string
	paramVerbose bool
	paramOutputFilename string
	paramDocumentTitle string
	pdfGenerator *gofpdf.Fpdf
	pdfTR func(string) string
)

func main() {
	// parse params
	flag.StringVar(&paramInstance, "i", "", "URL of your Jira instance, like: https://your.jira-instance.com")
	flag.StringVar(&paramUsername, "u", "", "Username of your Jira instance account")
	flag.StringVar(&paramPassword, "p", "", "Password of your Jira instance account")
	flag.StringVar(&paramQuery, "q", "", "JQL query for issues search")
	flag.StringVar(&paramOutputFilename, "o", "", "Output filename")
	flag.StringVar(&paramDocumentTitle, "t", "", "Document title")
	flag.BoolVar(&paramVerbose, "v", false, "Verbose mode")
	flag.Parse()

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

	pdfGenerator.MultiCell(0, 10, "", "", "", false)

	// issues
	for _, issue := range issues {
		pdfGenerator.SetFont("Arial", "", 11)
		pdfGenerator.SetTextColor(0, 0, 0)
		pdfGenerator.MultiCell(0, 11, pdfTR(fmt.Sprintf("Issue: %v", issue.Key)), "0", "L", false)
		pdfGenerator.MultiCell(0, 11, pdfTR(issue.Fields.Summary), "0", "L", false)
		pdfGenerator.MultiCell(0, 6, "", "", "", false)
	}

	// save to output filename
	err = pdfGenerator.OutputFileAndClose(paramOutputFilename)

	if err != nil {
		log.Fatalf("Erro while save PDF: %v", err)
	}
}
