package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

var (
	buildtime     string
	gitcommit     string
	appversion    string
	templateFile  string
	yamlFile      string
	renderedBytes bytes.Buffer
)

func main() {
	log.Printf("appversion: %s", appversion)
	log.Printf("gitcommit:  %s", gitcommit)
	log.Printf("buildtime:  %s", buildtime)

	flag.StringVar(&templateFile, "t", "", "The path to the template file")
	flag.StringVar(&yamlFile, "y", "", "The path to the yaml file with the data")

	flag.Usage = func() {
		fmt.Printf("%s - Render template with yaml input.\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	outputFileName := strings.Replace(templateFile, ".tpl", ".yaml", 1)

	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)

	yamlFileContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(yamlFileContent, &yamlData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = t.Execute(&renderedBytes, yamlData)
	if err != nil {
		fmt.Println("executing template:", err)
		os.Exit(1)
	}

	_, err = outputFile.Write(renderedBytes.Bytes())
	if err != nil {
		panic(err)
	}

}
