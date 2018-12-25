package suit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SuitIterator interface {
	HasNext()bool
	Next() (suit *Suit, fileName string)
}


type testFile struct {
	Name string
	Path string
}

func (f *testFile)String() string {
	return fmt.Sprintf("%s/%s", f.Path, f.Name)
}

type loadedSuits struct {
	files []*testFile
	index int
}

func (s *loadedSuits)HasNext() bool {
	return s.index < len(s.files)
}
func (s *loadedSuits)Next() (*Suit, string) {
	file := s.files[s.index]
	s.index++
	nextSuit, err := load(file.String())
	if err != nil {
		panic(err)
	}
	return nextSuit, file.String()
}


func LoadSuits(basePath string) SuitIterator {
	return &loadedSuits{
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

func load(fileName string) (*Suit, error) {
	file, ok := ioutil.ReadFile(fileName)
	if ok != nil {
		err := errors.New("Can't load " + fileName + "; Origin error: " + ok.Error())
		return nil, err
	}
	s := new(Suit)
	err := json.Unmarshal(file, s)
	return s, err
}