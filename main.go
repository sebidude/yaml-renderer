package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	buildtime    string
	gitcommit    string
	appversion   string
	templateFile string
	yamlFile     string
	suffix       string
	outputdir    string
	rematch      = `\$\{?([a-zA-Z0-9_]*)\}?`
	envVarMatch  *regexp.Regexp
)

func main() {
	log.Printf("appversion: %s", appversion)
	log.Printf("gitcommit:  %s", gitcommit)
	log.Printf("buildtime:  %s", buildtime)

	app := kingpin.New("yaml-renderer", "Render templates with yaml variable input.")

	app.Flag("templates", "path to the template file or directory").
		Short('t').
		StringVar(&templateFile)
	app.Flag("yaml", "path to the yaml value file").
		Short('y').
		StringVar(&yamlFile)
	app.Flag("output", "path to the output direcory (will be created if not exists)").
		Short('o').
		Default("rendered").
		StringVar(&outputdir)

	envVarMatch = regexp.MustCompile(rematch)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	info, err := os.Stat(templateFile)
	if err != nil {
		panic(err)
	}

	if outputdir == templateFile {
		panic("output directory must not be template directory")
	}

	if info.IsDir() {
		files, err := ioutil.ReadDir(templateFile)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			i, err := os.Stat(templateFile + "/" + f.Name())
			if i.IsDir() {
				continue
			}
			err = renderFile(templateFile + "/" + f.Name())
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
	t := template.Must(template.New(inputfilename).Funcs(template.FuncMap{"N": N}).Parse(string(content)))
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
		if v, ok := value.([]interface{}); ok {
			var list []interface{}
			for _, s := range v {
				if elm, ok := s.(string); ok {
					// we need to regexp match the string and iterate over the matches here
					matches := envVarMatch.FindAllStringSubmatch(elm, -1)
					for _, m := range matches {
						if len(m) < 2 {
							if elm, ok := s.(string); ok {
								if len(elm) > 0 {
									list = append(list, elm)
								}
							}
							continue
						}
						envvar := os.Getenv(m[1])
						if len(envvar) > 0 {
							for _, ie := range strings.Split(envvar, ",") {
								if len(ie) > 0 {
									list = append(list, ie)
								}
							}
						}
					}
				} else {
					list = append(list, s)
				}

			}
			if len(list) > 0 {
				data[key] = list
			} else {
				data[key] = value
			}

			continue
		}

		if s, ok := value.(string); ok {
			// we need to regexp match the string and iterate over the matches here
			matches := envVarMatch.FindAllStringSubmatch(s, -1)
			envvarval := s
			for _, m := range matches {
				if len(m) > 1 {
					envvar := os.Getenv(m[1])
					if len(envvar) > 0 {
						envvarval = strings.Replace(envvarval, m[0], envvar, -1)
					}
				}
			}
			data[key] = envvarval
			continue
		}

		if v, ok := value.(interface{}); ok {
			if d, ok := v.(map[string]interface{}); ok {
				getEnvVarForMapValue(&d)
			}

		}

	}
}

func N(end interface{}) (stream chan int) {
	e := 0
	switch end.(type) {
	case float64:
		e = int(end.(float64))
	case int:
		e = end.(int)
	}

	stream = make(chan int)
	go func() {
		for i := 0; i < e; i++ {
			stream <- i
		}
		close(stream)
	}()
	return
}
