package designinmemoryfilesystem

import (
	"reflect"
	"testing"
)

func TestFileSystem_Example1(t *testing.T) {
	fs := NewFileSystem()
	if got := fs.Ls("/"); len(got) != 0 {
		t.Fatalf(`Ls("/") = %v, want []`, got)
	}
	fs.MkDir("/a/b/c")
	fs.AddContentToFile("/a/b/c/d", "hello")
	if got := fs.Ls("/"); !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf(`Ls("/") = %v, want [a]`, got)
	}
	if got := fs.ReadContentFromFile("/a/b/c/d"); got != "hello" {
		t.Fatalf(`ReadContentFromFile("/a/b/c/d") = %q, want "hello"`, got)
	}
}

func TestFileSystem_LsOnFileReturnsItsOwnName(t *testing.T) {
	fs := NewFileSystem()
	fs.AddContentToFile("/foo.txt", "content")
	if got := fs.Ls("/foo.txt"); !reflect.DeepEqual(got, []string{"foo.txt"}) {
		t.Fatalf(`Ls("/foo.txt") = %v, want [foo.txt]`, got)
	}
}

func TestFileSystem_MkDirCreatesIntermediateDirs(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/x/y/z")
	if got := fs.Ls("/"); !reflect.DeepEqual(got, []string{"x"}) {
		t.Fatalf(`Ls("/") = %v, want [x]`, got)
	}
	if got := fs.Ls("/x"); !reflect.DeepEqual(got, []string{"y"}) {
		t.Fatalf(`Ls("/x") = %v, want [y]`, got)
	}
	if got := fs.Ls("/x/y"); !reflect.DeepEqual(got, []string{"z"}) {
		t.Fatalf(`Ls("/x/y") = %v, want [z]`, got)
	}
	if got := fs.Ls("/x/y/z"); len(got) != 0 {
		t.Fatalf(`Ls("/x/y/z") = %v, want [] (empty directory)`, got)
	}
}

func TestFileSystem_LsSortedLexicographically(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/zebra")
	fs.MkDir("/apple")
	fs.MkDir("/mango")
	if got := fs.Ls("/"); !reflect.DeepEqual(got, []string{"apple", "mango", "zebra"}) {
		t.Fatalf(`Ls("/") = %v, want [apple mango zebra]`, got)
	}
}

func TestFileSystem_AddContentToFileAppends(t *testing.T) {
	fs := NewFileSystem()
	fs.AddContentToFile("/note.txt", "hello ")
	fs.AddContentToFile("/note.txt", "world")
	if got := fs.ReadContentFromFile("/note.txt"); got != "hello world" {
		t.Fatalf(`ReadContentFromFile("/note.txt") = %q, want "hello world"`, got)
	}
}

func TestFileSystem_MixedFilesAndDirsInSameFolder(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/docs")
	fs.AddContentToFile("/readme.txt", "hi")
	if got := fs.Ls("/"); !reflect.DeepEqual(got, []string{"docs", "readme.txt"}) {
		t.Fatalf(`Ls("/") = %v, want [docs readme.txt]`, got)
	}
}

func TestFileSystem_AddContentToFileCreatesFileInExistingDir(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/a/b")
	fs.AddContentToFile("/a/b/file.txt", "data")
	if got := fs.Ls("/a/b"); !reflect.DeepEqual(got, []string{"file.txt"}) {
		t.Fatalf(`Ls("/a/b") = %v, want [file.txt]`, got)
	}
	if got := fs.ReadContentFromFile("/a/b/file.txt"); got != "data" {
		t.Fatalf(`ReadContentFromFile("/a/b/file.txt") = %q, want "data"`, got)
	}
}

func TestFileSystem_MultipleFilesUnderNestedDirs(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/a/b/c")
	fs.AddContentToFile("/a/b/c/one.txt", "1")
	fs.AddContentToFile("/a/b/c/two.txt", "2")
	fs.AddContentToFile("/a/other.txt", "o")
	if got := fs.Ls("/a/b/c"); !reflect.DeepEqual(got, []string{"one.txt", "two.txt"}) {
		t.Fatalf(`Ls("/a/b/c") = %v, want [one.txt two.txt]`, got)
	}
	if got := fs.Ls("/a"); !reflect.DeepEqual(got, []string{"b", "other.txt"}) {
		t.Fatalf(`Ls("/a") = %v, want [b other.txt]`, got)
	}
}

func TestFileSystem_MkDirOnAlreadyPartiallyExistingPath(t *testing.T) {
	fs := NewFileSystem()
	fs.MkDir("/a/b")
	fs.MkDir("/a/c") // "/a" already exists; only "c" should be newly created under it
	if got := fs.Ls("/a"); !reflect.DeepEqual(got, []string{"b", "c"}) {
		t.Fatalf(`Ls("/a") = %v, want [b c]`, got)
	}
}

func TestFileSystem_ReadFileWithMultipleAppends(t *testing.T) {
	fs := NewFileSystem()
	fs.AddContentToFile("/log.txt", "a")
	fs.AddContentToFile("/log.txt", "b")
	fs.AddContentToFile("/log.txt", "c")
	fs.AddContentToFile("/log.txt", "d")
	if got := fs.ReadContentFromFile("/log.txt"); got != "abcd" {
		t.Fatalf(`ReadContentFromFile("/log.txt") = %q, want "abcd"`, got)
	}
}
