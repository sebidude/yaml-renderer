package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/alecthomas/kingpin"
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

	app := kingpin.New("yaml-renderer", "Render templates with yaml variable input.")

	app.Flag("templates", "path to the template files or directory").
		Short('t').
		StringVar(&templateFile)
	app.Flag("yaml", "path to the yaml value file").
		Short('y').
		StringVar(&yamlFile)
	app.Flag("output", "path to the output direcory (will be created if not exists)").
		Short('o').
		Default("rendered").
		StringVar(&outputdir)

	kingpin.MustParse(app.Parse(os.Args[1:]))

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

	outputFileName := outputdir + "/" + inputfilebasename
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

	getEnvVarForMapValue(&yamlData)

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
	return nil
}

func getEnvVarForMapValue(indata interface{}) {
	data := *indata.(*map[string]interface{})
	for key, value := range *indata.(*map[string]interface{}) {
		if s, ok := value.(string); ok {
			if strings.HasPrefix(s, "$") {
				envvar := os.Getenv(strings.TrimLeft(s, "$"))
				if len(envvar) > 0 {
					data[key] = envvar
				}
			}
			continue
		}

		if v, ok := value.(interface{}); ok {
			if d, ok := v.(map[string]interface{}); ok {
				getEnvVarForMapValue(&d)
			}

		}

	}
}
