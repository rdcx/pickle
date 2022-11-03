package pickle

import (
	"bufio"
	"bytes"
	_ "embed"
	"errors"
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
		return fmt.Errorf("unable to parse data into template: %v", err)
	}

	contents := processed.Bytes()
	if strings.HasSuffix(outputFile, ".go") {
		contents, err = format.Source(contents)
		if err != nil {
			return fmt.Errorf("could not format processed template: %v", err)
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

func getMainTemplate(function Function) string {
	if function.Type == "gateway" {
		return templateGatewayMain
	}

	if function.Type == "mux" && function.HasRedis() {

		if function.Action == "show" {
			return templateRedisMainShow
		}

		if function.Action == "store" {
			return templateRedisMainStore
		}

		if function.Action == "update" {
			return templateRedisMainUpdate
		}

		if function.Action == "destroy" {
			return templateRedisMainDestroy
		}

		return templateRedisMainShow
	}

	return templateMuxMain

}

func compileMain(function Function, out string) error {

	outputFile := out + "/" + function.Name + "/main.go"
	tem := getMainTemplate(function)

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

func compileFiles(conf Config, out string) error {

	err := compileDockerCompose(conf, out)
	if err != nil {
		return err
	}
	for _, function := range conf.Functions {
		err = compileMain(function, out)
		if err != nil {
			return err
		}
		err = compileTest(function, out)
		if err != nil {
			return err
		}
		err = compileGoModFile(function, out)
		if err != nil {
			return err
		}
		err = compileDockerfile(function, out)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateConfig(config Config) *Validator {

	validator := NewValidator()

	for _, function := range config.Functions {

		if function.Type == "gateway" && function.Action != "" {
			validator.AddError("Gateway functions cannot have an action")
		}

		// Redis functions must have a model
		if function.HasRedis() && !function.HasModel() {
			validator.AddError("Redis functions must have a model")
		}

		// Model must have a name
		if function.HasModel() && function.Model.Name == "" {
			validator.AddError("Model must have a name")
		}

	}

	return validator
}

func decodeConfig(configFile string) (Config, error) {

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return config, err
	}

	return config, nil
}

func Gen(confFile string, outPath string) error {

	conf, err := decodeConfig(confFile)

	if err != nil {
		return err
	}

	validator := validateConfig(conf)

	if validator.HasErrors() {
		validator.PrintErrors()
		return errors.New("invalid config, see above errors")
	}

	return compileFiles(conf, outPath)
}
