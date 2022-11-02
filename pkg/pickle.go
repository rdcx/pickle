package pickle

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig"
)

type Field struct {
	Name string
	Type string
}

type Model struct {
	Name   string
	Fields []Field
}

type Redis struct {
	Host     string
	Port     string
	DB       int
	Password string
}

type Service struct {
	Name  string
	Image string
}

type Function struct {
	Name   string
	Type   string
	Action string
	Model  Model
	Redis  Redis
}

func (f Function) HasModel() bool {
	return f.Model.Fields != nil
}

func (f Function) HasRedis() bool {
	return f.Redis.Host != ""
}

func (f Function) IsGateway() bool {
	return f.Type == "gateway"
}

func (f Function) HasTest() bool {
	return f.Type == "mux"
}

type Config struct {
	Services  []Service
	Functions []Function
}

func compileTemplate(t, outputFile, out string, function Function) error {
	tmpl := template.Must(template.New("template").Funcs(sprig.FuncMap()).Parse(t))

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, "template", function)
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
		tem := ""
		if function.Type == "mux" {
			tem = templateMuxMainTest
		}
		return compileTemplate(tem, outputFile, out, function)
	}

	return nil
}

func compileMain(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/main.go"
	tem := templateMuxMain
	if function.Type == "gateway" {
		tem = templateGatewayMain
	}
	return compileTemplate(tem, outputFile, out, function)
}

func compileGoModFile(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/go.mod"
	tem := templateMuxGoMod
	if function.Type == "gateway" {
		tem = templateGatewayGoMod
	}

	return compileTemplate(tem, outputFile, out, function)
}

func compileDockerfile(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/Dockerfile"
	tem := templateMuxDockerfile
	if function.Type == "gateway" {
		tem = templateGatewayDockerfile
	}

	return compileTemplate(tem, outputFile, out, function)
}

func compileDockerCompose(config Config, out string) error {

	outputFile := out + "/docker-compose.yaml"
	tmpl := template.Must(
		template.New("template").
			Funcs(sprig.FuncMap()).
			Parse(templateDockerComposeYaml),
	)

	var processed bytes.Buffer
	err := tmpl.ExecuteTemplate(&processed, "template", config)
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

	compileDockerCompose(conf, out)

	return nil

}
