package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
)

var (
	buildtime    string
	gitcommit    string
	appversion   string
	templateFile string
	yamlFile     string
	suffix       string
	outputdir    string
)

func main() {
	log.Printf("appversion: %s", appversion)
	log.Printf("gitcommit:  %s", gitcommit)
	log.Printf("buildtime:  %s", buildtime)

	flag.StringVar(&templateFile, "t", "", "The path to the template file/dir")
	flag.StringVar(&yamlFile, "y", "", "The path to the yaml file with the data")
	flag.StringVar(&suffix, "s", "", "Suffix for the rendered file.")
	flag.StringVar(&outputdir, "o", "", "Output directory.")

	flag.Usage = func() {
		fmt.Printf("%s - Render template with yaml input.\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(suffix) < 1 {
		suffix = ".yaml"
	}

	if !strings.HasPrefix(suffix, ".") {
		suffix = "." + suffix
	}

	info, err := os.Stat(templateFile)
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		files, err := ioutil.ReadDir(templateFile)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".tpl") {
				continue
			}
			err := renderFile(templateFile + "/" + f.Name())
			if err != nil {
				panic(err)
			}
		}
	} else {
		err := renderFile(templateFile)
		if err != nil {
			panic(err)
		}
	}

}

func renderFile(inputfilename string) error {
	inputfilebasename := path.Base(inputfilename)
	if len(outputdir) < 1 {
		outputdir = "."
	}

	outputFileName := outputdir + "/" + strings.Replace(inputfilebasename, ".tpl", suffix, -1)
	if _, err := os.Stat(outputdir); os.IsNotExist(err) {
		err := os.MkdirAll(outputdir, 0755)
		if err != nil {
			return err
		}
	}

	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	yamlFileContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var yamlData map[string]interface{}
	err = yaml.Unmarshal(yamlFileContent, &yamlData)
	if err != nil {
		return err
	}
	log.Printf("using template file %s", inputfilename)
	content, err := ioutil.ReadFile(inputfilename)
	if err != nil {
		return err
	}
	t := template.Must(template.New(inputfilename).Parse(string(content)))
	var renderedBytes bytes.Buffer
	err = t.Execute(&renderedBytes, yamlData)
	if err != nil {
		return err
	}

	log.Printf("writing file %s", outputFileName)
	_, err = outputFile.Write(renderedBytes.Bytes())
	if err != nil {
		return err
	}

	//log.Printf("%s, %s, %s, %s ", inputfilename, inputfilebasename, outputFileName, outputdir)
	return nil
}
