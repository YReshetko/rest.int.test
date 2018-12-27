package suite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SuiteIterator interface {
	HasNext()bool
	Next() (suite *Suite, fileName string)
}


type testFile struct {
	Name string
	Path string
}

func (f *testFile)String() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

type loadedSuites struct {
	files []*testFile
	index int
}

func (s *loadedSuites)HasNext() bool {
	return s.index < len(s.files)
}
func (s *loadedSuites)Next() (*Suite, string) {
	file := s.files[s.index]
	s.index++
	nextSuite, err := load(file.String())
	if err != nil {
		panic(err)
	}
	return nextSuite, file.String()
}


func LoadSuites(basePath string) SuiteIterator {
	return &loadedSuites{
		extractAllFilesRecursively(keyPath(basePath)),
		0,
	}
}

func extractAllFilesRecursively(basePath string) []*testFile{
	files, err := ioutil.ReadDir(basePath)
	testFiles := []*testFile{}
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if f.IsDir() {
			testFiles = append(testFiles, extractAllFilesRecursively(basePath + "/" + f.Name())...)
		} else {
			testFiles = append(testFiles, &testFile{f.Name(), basePath})
		}
	}
	//log.Println("Return files:", testFiles, "for base path:", basePath)
	return testFiles
}

func keyPath(path string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	return fmt.Sprintf("%s%s%s", exePath, "/", path)
}

func load(fileName string) (*Suite, error) {
	file, ok := ioutil.ReadFile(fileName)
	if ok != nil {
		err := errors.New("Can't load " + fileName + "; Origin error: " + ok.Error())
		return nil, err
	}
	s := new(Suite)
	err := json.Unmarshal(file, s)
	return s, err
}