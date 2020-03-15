package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
{{ .Body }}
		<footer>
{{ .Footer }}
		</footer>
	</body>
</html>
`
)

// content type represents the HTML content tto add into the template
type content struct {
	Title  string
	Body   template.HTML
	Footer string
}

func main() {
	// parse flags
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	prfBrowser := flag.String("b", "", "Preferred browser")
	flag.Parse()
	// no input file show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename, *tFname, *prfBrowser, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename, tFname string, prfBrowser string, out io.Writer, skipPreview bool) error {
	// read data from file
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// temp file
	temp, err := ioutil.TempFile("", "mdp")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	outName := temp.Name()
	fmt.Fprintln(out, outName)

	htmlData, err := parseContent(input, tFname, outName)
	if err != nil {
		return err
	}

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}
	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName, prfBrowser)
}

func parseContent(input []byte, tFname string, outName string) ([]byte, error) {
	// parse markdown file through blackfriday and bluemonday
	// generate valid safe HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// prse contents of defaultTemplate
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}
	// user provided alternate template file
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	// Instantiate the content type, adding the title and body
	c := content{
		Title:  "Markdown Preview Tool",
		Body:   template.HTML(body),
		Footer: outName,
	}

	// create a buffer of bytes
	var buffer bytes.Buffer
	// Execute the template with the content type
	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	// write bytes to the file
	return ioutil.WriteFile(outFname, data, 0644)
}

func preview(fname string, prfBrowser string) error {
	browser := "firefox"
	if prfBrowser != "" {
		browser = prfBrowser
	}

	// Locate the browser in the PATH
	browserPath, err := exec.LookPath(browser)
	if err != nil {
		return err
	}
	// Open the file on the browser
	if err := exec.Command(browserPath, fname).Start(); err != nil {
		return err
	}

	// browser time to open file before deleting
	time.Sleep(2 * time.Second)
	return nil
}
