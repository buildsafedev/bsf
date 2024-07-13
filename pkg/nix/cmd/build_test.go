package cmd

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
)



func getRandomStringForTest() []string {
		str:= "Example string"
		rand1:= fmt.Sprint(str, rand.Int())
		rand2:= fmt.Sprint(str, rand.Int())
		rand3:= fmt.Sprint(str, rand.Int())
	return []string{rand1, rand2, rand3}
}

func TestManageStdErr(t *testing.T){
	t.Run("Test for piping stderr without warning", func(t *testing.T) {
		strArr:= getRandomStringForTest()
		origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()
	r, w, err := os.Pipe()
			if err!=nil{
				t.Fatalf(err.Error())
			}
			os.Stderr = w
			inputStr:= strings.Join(strArr, "\n")
			input := strings.NewReader(inputStr)
			err = ManageStdErr(io.NopCloser(input))
			if err != nil {
				t.Fatalf("ManageStdErr returned error: %v", err)
			}
			w.Close()
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err!=nil{
				t.Fatalf(err.Error())
			}
			str:= buf.String()
			expectedStr := strings.Join(strArr, "\n") + "\n"
			if str!=expectedStr{
				t.Errorf("ManageStdErr want %v got %v", inputStr, str)
			}
	})

	t.Run("Test for piping stderr with warning", func(t *testing.T) {
		origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()
	r, w, _ := os.Pipe()
			os.Stderr = w
			dir, err:= os.Getwd()
			if err!=nil{
				t.Fatalf(err.Error())
			}
			inputStr:= fmt.Sprintf("warning: Git tree '%s' is dirty", dir)
			input := strings.NewReader(inputStr)
			err = ManageStdErr(io.NopCloser(input))
			if err!=nil{
				t.Fatalf(err.Error())
			}
			w.Close()
			var buf bytes.Buffer
			_,err = io.Copy(&buf, r)
			if err!=nil{
				t.Fatalf(err.Error())
			}
			str:= buf.String()
			expectedStr := fmt.Sprintf("warning: Git tree '%s' is dirty.\nThis implies you have not checked-in files in the git work tree (hint: git add)\n", dir)
			if str!=expectedStr{
				t.Errorf("ManageStdErr want %v got %v", inputStr, str)
			}
	})
}

func TestManageStdOutput(t *testing.T){
	t.Run("Test for piping stderr without warning", func(t *testing.T) {
		strArr:= getRandomStringForTest()
		origStdout := os.Stdout
	defer func() { os.Stdout = origStdout }()
	r, w, _ := os.Pipe()
			os.Stdout = w
			inputStr:= strings.Join(strArr, "\n")
			input := strings.NewReader(inputStr)
			err := ManageStdOutput(io.NopCloser(input))
			if err != nil {
				t.Fatalf("ManageStdOut returned error: %v", err)
			}
			w.Close()
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err!=nil{
				t.Fatalf(err.Error())
			}
			str:= buf.String()
			expectedStr := strings.Join(strArr, "\n") + "\n"
			if str!=expectedStr{
				t.Errorf("ManageStdOutput want %v got %v", inputStr, str)
			}
	})
}