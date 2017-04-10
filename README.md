# Jira to PDF

This is a small tool to export Jira issues to PDF. 
  
# Installing from source

1. go get -u github.com/prsolucoes/jira-to-pdf
2. cd $GOPATH/src/github.com/prsolucoes/jira-to-pdf
3. make deps
4. make install

# Get prebuilt executables

You can get prebuilt executables on directory "build".

# How to use

```
jira-to-pdf -i [your-jira-instance-url] -u [your-jira-username] -p [your-jira-password] -q [any-jql-query] -t [document-title] -o [output-filename]
```

# Support with donation
[![Support with donation](http://donation.pcoutinho.com/images/donate-button.png)](http://donation.pcoutinho.com/)

# Supported By Jetbrains IntelliJ IDEA

![alt text](https://github.com/prsolucoes/jira-to-pdf/raw/master/extras/jetbrains/logo.png "Supported By Jetbrains IntelliJ IDEA")

# Author WebSite

> http://www.pcoutinho.com

# License

MIT
