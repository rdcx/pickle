package pickle

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
)

type Function struct {
	Name   string
	Type   string
	Method string
}

func (f Function) HasTest() bool {
	return f.Type == "mux"
}

type Config struct {
	Functions []Function
}

func compileTemplate(templateName, outputFile, out string, function Function) error {
	name := path.Base(templateName)

	tmpl := template.Must(template.New(name).Funcs(sprig.FuncMap()).ParseFiles(templateName))

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, name, function)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}

	contents := processed.Bytes()
	if strings.HasSuffix(outputFile, ".go") {
		contents, err = format.Source(contents)
		if err != nil {
			log.Fatalf("Could not format processed template: %v\n", err)
		}
	}

	if _, err := os.Stat(out + "/" + function.Name); os.IsNotExist(err) {
		os.MkdirAll(out+"/"+function.Name, 0700) // Create your file
	}

	fmt.Println("Writing file: ", outputFile)
	f, _ := os.Create(outputFile)
	w := bufio.NewWriter(f)
	w.WriteString(string(contents))
	w.Flush()
	return nil
}

func compileTest(function Function, out string) error {

	if function.HasTest() {
		outputFile := out + "/" + function.Name + "/main_test.go"
		templateName := "./templates/mux/main_test.go.tmpl"
		return compileTemplate(templateName, outputFile, out, function)
	}

	return nil
}

func compileMain(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/main.go"
	templateName := "./templates/mux/main.go.tmpl"
	return compileTemplate(templateName, outputFile, out, function)
}

func compileGoModFile(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/go.mod"
	templateName := "./templates/mux/go.mod.tmpl"

	return compileTemplate(templateName, outputFile, out, function)
}

func compileDockerfile(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/Dockefile"
	templateName := "./templates/mux/Dockerfile.tmpl"

	return compileTemplate(templateName, outputFile, out, function)
}

func compileDockerCompose(functions []Function, out string) error {

	outputFile := out + "/docker-compose.yaml"
	templateName := "./templates/docker-compose.yaml.tmpl"

	name := path.Base(templateName)

	tmpl := template.Must(template.New(name).Funcs(sprig.FuncMap()).ParseFiles(templateName))

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, name, functions)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}

	if _, err := os.Stat(out); os.IsNotExist(err) {
		os.MkdirAll(out, 0700) // Create your file
	}

	fmt.Println("Writing file: ", outputFile)
	f, _ := os.Create(outputFile)
	w := bufio.NewWriter(f)
	w.WriteString(processed.String())
	w.Flush()
	return nil
}

func Gen(in string, out string) error {

	content, err := os.ReadFile(in)

	if err != nil {
		return err
	}

	var conf Config
	_, err = toml.Decode(string(content), &conf)

	if err != nil {
		log.Fatal(err)
	}

	for _, f := range conf.Functions {
		compileDockerfile(f, out)
		compileGoModFile(f, out)
		compileMain(f, out)
		compileTest(f, out)
	}

	compileDockerCompose(conf.Functions, out)

	return nil

}
