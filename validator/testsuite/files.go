//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package testsuite

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/katydid/parser-go/parser"
	"github.com/katydid/validator-go-jsonschema/json"
	"github.com/katydid/validator-go-jsonschema/validator"
	"github.com/katydid/validator-go-jsonschema/validator/ast"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

var testpath string
var benchpath string

func init() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = "../../../../../../"
	}
	testpath = filepath.Join(gopath, "src/github.com/katydid/validator-testsuite/validator/tests")
	benchpath = filepath.Join(gopath, "src/github.com/katydid/validator-testsuite/validator/benches")
}

func TestSuiteExists() (bool, error) {
	if exists(testpath) {
		return true, nil
	}
	if os.Getenv("TESTSUITE") == "MUST" {
		return false, fmt.Errorf("testsuite does not exist at %v", testpath)
	}
	return false, nil
}

func BenchSuiteExists() (bool, error) {
	if exists(testpath) {
		return true, nil
	}
	if os.Getenv("TESTSUITE") == "MUST" {
		return false, fmt.Errorf("testsuite does not exist at %v", testpath)
	}
	return false, fmt.Errorf("benchsuite does not exist at %v", testpath)
}

func getFolders(path string) (map[string][]string, error) {
	folders := make(map[string][]string)
	codecFileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, codecFileInfo := range codecFileInfos {
		if !codecFileInfo.IsDir() {
			continue
		}
		codecFolderName := codecFileInfo.Name()
		codecPath := filepath.Join(path, codecFolderName)
		caseDirInfos, err := os.ReadDir(codecPath)
		if err != nil {
			return nil, err
		}
		for _, caseDirInfo := range caseDirInfos {
			if !caseDirInfo.IsDir() {
				continue
			}
			casePath := filepath.Join(codecPath, caseDirInfo.Name())
			folders[codecFolderName] = append(folders[codecFolderName], casePath)
		}
	}
	return folders, nil
}

func ReadTestSuite() ([]Test, error) {
	tests := []Test{}
	codecs, err := getFolders(testpath)
	if err != nil {
		return nil, err
	}
	for codec, folders := range codecs {
		switch codec {
		case "json":
		default:
			// codec not supported
			continue
		}
		for _, folder := range folders {
			test, err := readTestFolder(folder)
			if err != nil {
				return nil, err
			}
			tests = append(tests, *test)
		}
	}
	return tests, nil
}

func ReadBenchmarkSuite() ([]Bench, error) {
	benches := []Bench{}
	codecs, err := getFolders(benchpath)
	if err != nil {
		return nil, err
	}
	for codec, folders := range codecs {
		switch codec {
		case "json":
		default:
			// codec not supported
			continue
		}
		for _, folder := range folders {
			bench, err := readBenchFolder(folder)
			if err != nil {
				return nil, err
			}
			benches = append(benches, *bench)
		}
	}
	return benches, nil
}

type Test struct {
	Name     string
	Grammar  *ast.Grammar
	Parser   parser.Interface
	Expected bool
	Record   bool
}

func readTestFolder(path string) (*Test, error) {
	name := filepath.Base(path)
	g, err := readGrammar(path)
	if err != nil {
		return nil, err
	}
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("err <%v> reading folder <%s>", err, path)
	}
	var p parser.Interface
	var expected bool
	var codecName string
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filebase := fileInfo.Name()
		filename := filepath.Join(path, filebase)
		names := strings.Split(filebase, ".")
		valid := names[0] == "valid"
		invalid := names[0] == "invalid"
		if !valid && !invalid {
			continue
		}
		expected = valid
		codecName = names[len(names)-1]
		switch codecName {
		case "json":
			p, err = newJsonParser(filename)
			if err != nil {
				return nil, err
			}
		default:
			// unsupported codec
			continue
		}
	}
	if p == nil {
		return nil, fmt.Errorf("couldn't find valid.* or invalid.* filename inside <%s>", path)
	}
	name = name + capFirst(codecName)
	return &Test{
		Name:     name,
		Grammar:  g,
		Parser:   p,
		Expected: expected,
		Record:   true,
	}, nil
}

type Bench struct {
	Name    string
	Grammar *ast.Grammar
	Parsers []ResetParser
	Record  bool
}

type ResetParser interface {
	parser.Interface
	Reset() error
}

func readBenchFolder(path string) (*Bench, error) {
	name := filepath.Base(path)
	g, err := readGrammar(path)
	if err != nil {
		return nil, err
	}
	fileInfos, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("err <%v> reading folder <%s>", err, path)
	}
	var parsers []ResetParser
	var codecName string
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filebase := fileInfo.Name()
		names := strings.Split(filebase, "_")
		valid := names[0] == "valid"
		invalid := names[0] == "invalid"
		if !valid && !invalid {
			continue
		}
		filename := filepath.Join(path, filebase)
		codecName = filepath.Ext(filename)[1:]
		switch codecName {
		case "json":
			p, err := newJsonParser(filename)
			if err != nil {
				return nil, err
			}
			parsers = append(parsers, p)
		default:
			// unsupported codec
			continue
		}
	}
	return &Bench{
		Name:    name + capFirst(codecName),
		Grammar: g,
		Parsers: parsers,
		Record:  true,
	}, nil
}

func capFirst(s string) string {
	b := []byte(s)
	b[0] ^= ' '
	return string(b)
}

func readGrammar(path string) (*ast.Grammar, error) {
	validatorTxt := filepath.Join(path, "validator.txt")
	validatorBytes, err := os.ReadFile(validatorTxt)
	if err != nil {
		return nil, fmt.Errorf("err <%v> reading file <%s>", err, validatorTxt)
	}
	g, err := validator.Parse(string(validatorBytes))
	if err != nil {
		return nil, fmt.Errorf("err <%v> parsing grammar from file <%s>", err, validatorTxt)
	}
	return g, nil
}

func newJsonParser(filename string) (ResetParser, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("err <%v> reading file <%s>", err, filename)
	}
	j := json.NewJsonParser()
	if err := j.Init(bytes); err != nil {
		return nil, fmt.Errorf("err <%v> parser.Init with bytes from filename <%s>", err, filename)
	}
	return j, nil
}
