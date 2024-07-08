package git

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
)

func getRootDirectory(dir string) string {
	arr:= strings.Split(dir, "/")
	var indexOfBsf int
	for index, value:= range arr {
		if value == "bsf"{
			indexOfBsf = index
			break;
		}
	}
	arr = arr[:indexOfBsf]
	str:= strings.Join(arr, "/")
	return str
}

func initializeTestEnv(fileName string) (string, string, *os.File, error) {
	oldDir, err:= os.Getwd();
	newDir:=getRootDirectory(oldDir)
	if(err!=nil){
		return "","",nil, err
	}
	os.Mkdir(newDir+"/bsf-temp", 0777)
	os.Chdir(newDir+"/bsf-temp")
	file, err:=os.Create(fileName)
	os.Create("sample.txt")
	if err!=nil{
		return "","",nil,err
	}
	// Previous directory is returned so as to return to the project directory
	return oldDir, newDir, file, nil
}

func initiateGitEnv() (*git.Worktree, error) {
	git.PlainInit(".", false)
	r, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err!=nil{
		return nil, err
	}
	w, err := r.Worktree()
	if err!=nil{
		return nil, err
	}
	err = w.AddWithOptions(&git.AddOptions{
		Path: ".",
	})
	if err != nil {
		return nil, err
	}
	return w, nil
}

func cleanTestEnv(oldDir string, newDir string) error {
	 os.RemoveAll(newDir+"/bsf-temp")
	 os.Chdir(oldDir)
	 return nil	
}


func TestGitAdd(t *testing.T){
	t.Run("git.Add() should run without any error for go module", func(t *testing.T) {
		oldDir, newDir, file, err:=initializeTestEnv("go.mod")
		if err!=nil{
			t.Fatal()
		}
		defer cleanTestEnv(oldDir, newDir)
		goContent:=`module test`
		file.WriteString(goContent)
		initiateGitEnv()
		Add("./")
	})

	tests:=[]struct{
		langName string
		fileName string
	}{
		{
			langName: "Javascript",
			fileName: "package-lock.json",
		},
		{
			langName: "Poetry",
			fileName: "poetry.lock",
		},
		{
			langName: "Rust",
			fileName: "Cargo.lock",
		},
	}
	for _,tt:= range tests{
		testName:= "git.Add() should run without any error for"+ tt.langName
		t.Run(testName, func (t *testing.T)  {
			oldDir, newDir, _, err:=initializeTestEnv(tt.fileName)
			if err!=nil{
				t.Fatal()
			}
			defer cleanTestEnv(oldDir, newDir)
			initiateGitEnv()
			errors:=Add("./")
			if errors!=nil{
				t.Errorf("want nil but found error: %s", errors.Error())
			}
		})
	}

	t.Run("git.Add() should throw error of go.mod not added to version control", func(t *testing.T){
		oldDir, newDir, file, err:=initializeTestEnv("go.mod")
		if err!=nil{
			t.Fatal()
		}
		defer cleanTestEnv(oldDir, newDir)
		goContent:=`module test`
		file.WriteString(goContent)
		git.PlainInit(".", false)
		git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
		errors:= Add("sample.txt")
		if errors==nil{
			t.Errorf("want error but found nil")	
		}
		if errors != ErrFilesNotAddedToVersionControl{
			t.Errorf("want ErrFilesNotAddedToVersionControl but found %s", errors.Error())
		}
		})

	for _,tt:= range tests{
			testName:= "git.Add() should throw error of file not added to version control for "+tt.langName
			t.Run(testName, func (t *testing.T)  {
				oldDir, newDir, _, err:=initializeTestEnv(tt.fileName)
				if err!=nil{
					t.Fatal()
				}
				defer cleanTestEnv(oldDir, newDir)
				git.PlainInit(".", false)
				git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
				errors:= Add("sample.txt")
				if errors==nil{
					t.Errorf("want error but found nil")	
				}
				if errors != ErrFilesNotAddedToVersionControl{
					t.Errorf("want ErrFilesNotAddedToVersionControl but found %s", errors.Error())
				}
			})
		}
}

func TestGitIgnore(t *testing.T){
	t.Run("git.Ignore() creates .gitignore and adds the path to it", func(t *testing.T) {
		oldDir, newDir, _, err:=initializeTestEnv("go.mod")
		if err!=nil{
			t.Errorf(err.Error())
		}
		defer cleanTestEnv(oldDir, newDir)
		Ignore("sample.txt")
		file, err:=os.ReadFile(".gitignore")
		if err!=nil{
			t.Fatal()
		}
		want:= strings.Contains(string(file), "sample.txt")
		if want==false{
			t.Errorf("want true recieved true")
		}
	})
	t.Run("git.Ignore() appends the path in already created .gitignore", func(t *testing.T) {
		oldDir, newDir, _, err:=initializeTestEnv("go.mod")
		if err!=nil{
			t.Errorf(err.Error())
		}
		defer cleanTestEnv(oldDir, newDir)
		fl, err:= os.Create(".gitignore")
		if err!=nil{
			t.Fatal()
		}
		fl.WriteString(`
		/path/to/be/added/1
		/path/to/be/added/2
		/path/to/be/added/3
		`)
		Ignore("sample.txt")
		file, err:=os.ReadFile(".gitignore")
		if err!=nil{
			t.Fatal()
		}
		want:= strings.Contains(string(file), "sample.txt")
		if want==false{
			t.Errorf("want true recieved false")
		}
		for i:=1; i<=3; i++{
			want:=strings.Contains(string(file), "/path/to/be/added/"+strconv.Itoa(i))
			if want==false{
				t.Errorf("want true recieved false")
			}
		}
	})
	t.Run("git.Ignore() adds nothing for already added path", func(t *testing.T) {
		oldDir, newDir, _, err:=initializeTestEnv("go.mod")
		if err!=nil{
			t.Errorf(err.Error())
		}
		defer cleanTestEnv(oldDir, newDir)
		fl, err:= os.Create(".gitignore")
		if err!=nil{
			t.Fatal()
		}
		fl.WriteString(`
		/path/to/be/added/1
		/path/to/be/added/2
		/path/to/be/added/3
		sample.txt
		`)
		Ignore("sample.txt")
		file, err:=os.ReadFile(".gitignore")
		if err!=nil{
			t.Fatal()
		}
		want:= strings.Contains(string(file), "sample.txt")
		if want==false{
			t.Errorf("want true recieved false")
		}
		for i:=1; i<=3; i++{
			want:=strings.Contains(string(file), "/path/to/be/added/"+strconv.Itoa(i))
			if want==false{
				t.Errorf("want true recieved false")
			}
		}
		count:=strings.Count(string(file), "sample.txt")
		if count!=1{
			t.Errorf("want count 1 recieved %d", count)
		}
	})
}

