package pickle

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
)

type Function struct {
	Name   string
	Method string
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
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}

	if _, err := os.Stat(out + "/" + function.Name); os.IsNotExist(err) {
		os.MkdirAll(out+"/"+function.Name, 0700) // Create your file
	}

	fmt.Println("Writing file: ", outputFile)
	f, _ := os.Create(outputFile)
	w := bufio.NewWriter(f)
	w.WriteString(string(formatted))
	w.Flush()
	return nil
}

func compileTest(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/main_test.go"
	templateName := "./templates/mux/main_test.go.tmpl"
	return compileTemplate(templateName, outputFile, out, function)
}

func compileMain(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/main.go"
	templateName := "./templates/mux/main.go.tmpl"
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

func compileGoModFile(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/go.mod"
	templateName := "./templates/mux/go.mod.tmpl"

	name := path.Base(templateName)

	tmpl := template.Must(template.New(name).Funcs(sprig.FuncMap()).ParseFiles(templateName))

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, name, function)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}

	if _, err := os.Stat(out + "/" + function.Name); os.IsNotExist(err) {
		os.MkdirAll(out+"/"+function.Name, 0700) // Create your file
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
		compileGoModFile(f, out)
		compileMain(f, out)
		compileTest(f, out)
	}

	compileDockerCompose(conf.Functions, out)

	return nil

}
