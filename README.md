# JIRA to PDF

This is a small tool to export JIRA issues to PDF. 

# How to use

```
jira-to-pdf -i [your-jira-instance-url] -u [your-jira-username] -p [your-jira-password] -q [any-jql-query] -t [document-title] -o [output-filename]
```

# Options

```
-i  = Your JIRA instance URL, ex: https://your-instance.atlassian.net  
-u  = Your JIRA instance username
-p  = Your JIRA instance password
-q  = Your JQL query to list the issues on JIRA
-t  = PDF document title [optional, default: JIRA Issues]
-o  = Output filename [optional, default: jira-issues.pdf]
-it = Issue template (support for basic html) [optional, default: "<b>Issue:</b> [issue.key]<br /><b>Summary:</b> [issue.fields.summary]<br /><b>Assignee:</b> [issue.fields.assignee.name]<br /><b>Status:</b> [issue.fields.status.name]<br /><b>Created:</b> [issue.fields.created]"] 
-dtf = DateTime format [optional] (default: "2006-01-02 15:04:05", check golang format options: https://golang.org/pkg/time/#pkg-examples) 
```

# Template fields that can be replaced

```
[issue.key]
[issue.id]

[issue.fields.created]
[issue.fields.description]
[issue.fields.duedate]
[issue.fields.expand]
[issue.fields.resolutiondate]
[issue.fields.summary]
[issue.fields.timeestimate]
[issue.fields.timeoriginalestimate]
[issue.fields.timespent]

[issue.fields.priority.id]
[issue.fields.priority.name]

[issue.fields.aggregateprogress.progress]
[issue.fields.aggregateprogress.total]

[issue.fields.progress.progress]
[issue.fields.progress.total]
[issue.fields.assignee.name]

[issue.fields.assignee.emailaddrress]
[issue.fields.assignee.displayname]
[issue.fields.assignee.key]

[issue.fields.creator.name]
[issue.fields.creator.emailaddrress]
[issue.fields.creator.displayname]
[issue.fields.creator.key]

[issue.fields.reporter.name]
[issue.fields.reporter.emailaddrress]
[issue.fields.reporter.displayname]
[issue.fields.reporter.key]

[issue.fields.project.name]
[issue.fields.project.description]
[issue.fields.project.id]
[issue.fields.project.key]

[issue.fields.status.name]
[issue.fields.status.description]
[issue.fields.status.id]

[issue.fields.type.name]
[issue.fields.type.description]
[issue.fields.type.id]
```
  
# Installing from source

1. go get -u github.com/prsolucoes/jira-to-pdf
2. cd $GOPATH/src/github.com/prsolucoes/jira-to-pdf
3. make deps
4. make install

# Get prebuilt executables

You can get prebuilt executables on directory "build".


# Support with donation
[![Support with donation](http://donation.pcoutinho.com/images/donate-button.png)](http://donation.pcoutinho.com/)

# Supported By Jetbrains IntelliJ IDEA

![alt text](https://github.com/prsolucoes/jira-to-pdf/raw/master/extras/jetbrains/logo.png "Supported By Jetbrains IntelliJ IDEA")

# Author WebSite

> http://www.pcoutinho.com

# License

MIT
