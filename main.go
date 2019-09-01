package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	jira "github.com/irahardianto/go-jira"
)

func main() {
	base := ""

	tp := jira.BasicAuthTransport{
		Username: "",
		Password: "",
	}

	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	roadmapKey, err := createRoadMapEpic(jiraClient, "Just a demo roadmap", "2019-09-02", "2019-09-22")
	if err != nil {
		panic(err)
	}

	epicKey, err := createProjectEpic(jiraClient, roadmapKey, "DEMO", "Just a demo epic", "Just a demo epic")
	if err != nil {
		panic(err)
	}

	issueKey, err := createProjectTasks(jiraClient, epicKey, "DEMO", "Story", "Just a demo story")
	if err != nil {
		panic(err)
	}

	fmt.Println(issueKey)
}

func createRoadMapEpic(jc *jira.Client, epicSummary string, startDate string, dueDate string) (string, error) {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Type: jira.IssueType{
				Name: "Epic",
			},
			StartDate: startDate,
			Duedate:   dueDate,
			Project: jira.Project{
				Key: "DEMO",
			},
			Summary: epicSummary,
		},
	}

	issue, _, err := jc.Issue.Create(&i)
	if err != nil {
		return "", err
	}

	return issue.Key, nil
}

func createProjectEpic(jc *jira.Client, roadmapKey string, projectKey string, epicSummary string, epicName string) (string, error) {
	// WIP
	// a := []*jira.IssueLink{}
	// jil := jira.IssueLink{
	// 	Type: jira.IssueLinkType{
	// 		Name: "Blocks",
	// 	},
	// 	OutwardIssue: &jira.Issue{
	// 		Key: roadmapKey,
	// 	},
	// }
	//
	//a = append(a, &jil)

	i := jira.Issue{
		Fields: &jira.IssueFields{
			// Assignee: &jira.User{
			// 	Name: "",
			// },
			// Reporter: &jira.User{
			// 	Name: "",
			// },
			// Description: "Test Epic",
			Type: jira.IssueType{
				Name: "Epic",
			},
			EpicName: epicName,
			Project: jira.Project{
				Key: projectKey,
			},
			// IssueLinks: a, WIP
			Summary: epicSummary,
		},
	}
	issue, resp, err := jc.Issue.Create(&i)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

	if err != nil {
		panic(err)
	}

	fmt.Println(issue.Key)

	return issue.Key, nil
}

func createProjectTasks(jc *jira.Client, epicKey string, projectKey string, issueType string, ticketSummary string) (string, error) {
	i := jira.Issue{
		Fields: &jira.IssueFields{
			// Assignee: &jira.User{
			// 	Name: "",
			// },
			// Reporter: &jira.User{
			// 	Name: "",
			// },
			// Description: "Test Issue",
			Type: jira.IssueType{
				Name: issueType,
			},
			EpicLink: epicKey,
			Project: jira.Project{
				Key: projectKey,
			},
			Summary: ticketSummary,
		},
	}

	issue, _, err := jc.Issue.Create(&i)

	if err != nil {
		panic(err)
	}

	fmt.Println(issue.Key)

	return issue.Key, nil
}

func readCSV() []Issue {
	csvFile, _ := os.Open("people.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var issues []Issue
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		issues = append(issues, Issue{
			Summary:     line[0],
			IssueType:   line[1],
			Labels:      line[2],
			Assignee:    line[3],
			DueDate:     line[4],
			Description: line[5],
			Blocks:      line[6],
			Epic:        line[7],
		})
	}
	issuesJson, _ := json.Marshal(issues)
	fmt.Println(string(issuesJson))

	return issues
}

type Issue struct {
	Summary     string `json:"summary"`
	IssueType   string `json:"lastname"`
	Labels      string `json:"labels"`
	Assignee    string `json:"assignee"`
	DueDate     string `json:"duedate"`
	Description string `json:"description"`
	Blocks      string `json:"blocks"`
	Epic        string `json:"epic"`
}
