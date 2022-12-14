package pickle

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestFunctionHasModel(t *testing.T) {
	f := Function{
		Model: Model{
			Fields: []Field{
				{
					Name: "Name",
					Type: "string",
				},
			},
		},
	}

	if !f.HasModel() {
		t.Errorf("Function should have a model, got: %v, want: %v.", f.HasModel(), true)
	}
}

func TestPickle(t *testing.T) {

	os.RemoveAll("./test/output")

	// When I run function Gen and specifty a TOML input file and a output path
	// The directory should be created for the given output path
	// I should get code generated
	// And the code should compile

	in, err := filepath.Abs("./test/config/functions.toml")
	if err != nil {
		t.Error(err)
	}
	out, err := filepath.Abs("./test/output")
	if err != nil {
		t.Error(err)
	}
	err = Gen(in, out)

	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
		return
	}

	checkPaths := []string{
		"greeter",
		"gateway",
		"person.show",
		"person.store",
		"person.update",
		"person.destroy",
	}

	for _, path := range checkPaths {
		f1, err1 := os.ReadFile("./test/assert/" + path + "/main.go")

		if err1 != nil {
			log.Fatal(err1)
		}

		// per comment, better to not read an entire file into memory
		// this is simply a trivial example.
		f2, err2 := os.ReadFile("./test/output/" + path + "/main.go")

		if err1 != nil {
			log.Fatal(err2)
		}

		if !bytes.Equal(f1, f2) {
			t.Errorf("Files are not equal for %s, got: \n%v\n, want: \n%v\n.", path, string(f2), string(f1))
		}

		f1, err1 = os.ReadFile("./test/assert/" + path + "/Dockerfile")

		if err1 != nil {
			log.Fatal(err1)
		}

		// per comment, better to not read an entire file into memory
		// this is simply a trivial example.
		f2, err2 = os.ReadFile("./test/output/" + path + "/Dockerfile")

		if err1 != nil {
			log.Fatal(err2)
		}

		if !bytes.Equal(f1, f2) {
			t.Errorf("Files are not equal for %s, got: \n%v\n, want: \n%v\n.", path, string(f2), string(f1))
		}

		f1, err1 = os.ReadFile("./test/assert/" + path + "/go.mod")

		if err1 != nil {
			log.Fatal(err1)
		}

		// per comment, better to not read an entire file into memory
		// this is simply a trivial example.
		f2, err2 = os.ReadFile("./test/output/" + path + "/go.mod")

		if err1 != nil {
			log.Fatal(err2)
		}

		if !bytes.Equal(f1, f2) {
			t.Errorf("Files are not equal for %s, got: \n%v\n, want: \n%v\n.", path, string(f2), string(f1))
		}
	}

	f1, err1 := os.ReadFile("./test/assert/docker-compose.yaml")

	if err1 != nil {
		log.Fatal(err1)
	}

	// per comment, better to not read an entire file into memory
	// this is simply a trivial example.
	f2, err2 := os.ReadFile("./test/output/docker-compose.yaml")

	if err1 != nil {
		log.Fatal(err2)
	}

	if !bytes.Equal(f1, f2) {
		t.Errorf("Files are not equal, got: \n%v\n, want: \n%v\n.", string(f2), string(f1))
	}

	// os.RemoveAll("./test/output")
}
