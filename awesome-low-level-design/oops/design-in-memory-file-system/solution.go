package designinmemoryfilesystem

/*
LeetCode 588. Design In-Memory File System

Design a data structure that simulates an in-memory file system.

Implement the FileSystem class:
  - NewFileSystem() initializes the object of the system.
  - Ls(path string) []string returns a list of strings.
      - If path is a file path, returns a list that only contains this file's name.
      - If path is a directory path, returns the list of file and directory names in
        this directory, sorted in lexicographic order.
  - MkDir(path string) makes a new directory according to the given path. The given
    directory path does not exist. If the middle directories in the path do not exist,
    you should create them as well.
  - AddContentToFile(filePath, content string) if filePath does not exist, creates that
    file containing given content. If filePath already exists, appends the given content
    to original content.
  - ReadContentFromFile(filePath string) string returns the content in the file at
    filePath.

Example 1:

	Input:
	["FileSystem", "ls", "mkdir", "addContentToFile", "ls", "readContentFromFile"]
	[[], ["/"], ["/a/b/c"], ["/a/b/c/d", "hello"], ["/"], ["/a/b/c/d"]]
	Output:
	[null, [], null, null, ["a"], "hello"]

	Explanation:
	FileSystem fs = new FileSystem();
	fs.ls("/");                                // return []
	fs.mkdir("/a/b/c");
	fs.addContentToFile("/a/b/c/d", "hello");
	fs.ls("/");                                // return ["a"]
	fs.readContentFromFile("/a/b/c/d");        // return "hello"

Constraints:
  - 1 <= path.length, filePath.length <= 100
  - path and filePath are absolute paths which begin with '/' and do not end with '/'
    except that the path is just "/".
  - You can assume that all directory names and file names only contain lowercase
    letters, and the same names will not exist in the same directory.
  - You can assume that all operations will be passed valid parameters, and users will
    not attempt to retrieve file content or list a directory or file that does not exist.
  - 1 <= content.length <= 50.
  - At most 300 calls will be made to Ls, MkDir, AddContentToFile, and
    ReadContentFromFile.

TODO:
 1. Add the fields you need on FileSystem below. Hint: model the whole tree as a single
    root "node" with a struct like:
      type node struct {
          children map[string]*node
          isFile   bool
          content  string
      }
    A directory node has isFile=false and a populated children map; a file node has
    isFile=true, content set, and an empty/nil children map.
 2. Write a helper to split a path like "/a/b/c" into its component names ["a","b","c"]
    (see the `strings` package — strings.Split on "/", filtering out the empty leading
    component and the root "/" case).
 3. Write a helper to walk from the root down through a sequence of path components,
    returning the node at that path (and whether it exists).
 4. Implement NewFileSystem: create a single root directory node.
 5. Implement MkDir: walk the path components from the root, creating any missing
    directory nodes along the way (this is the "mkdir -p" behavior).
 6. Implement AddContentToFile: walk to the parent directory (creating it if needed isn't
    required per constraints, but walking must find/create the final file node), then
    either create a new file node with content, or append to an existing file's content.
 7. Implement ReadContentFromFile: walk to the file node and return its content.
 8. Implement Ls: walk to the node at path. If it's a file, return a single-element slice
    with just its name (the last path component). If it's a directory, return the sorted
    names of its children (see `sort.Strings`).
*/

type FileSystem struct {
	// TODO: add your fields here (e.g. a root *node)
}

func NewFileSystem() *FileSystem {
	// TODO: implement
	return &FileSystem{}
}

func (fs *FileSystem) Ls(path string) []string {
	// TODO: implement
	return nil
}

func (fs *FileSystem) MkDir(path string) {
	// TODO: implement
}

func (fs *FileSystem) AddContentToFile(filePath string, content string) {
	// TODO: implement
}

func (fs *FileSystem) ReadContentFromFile(filePath string) string {
	// TODO: implement
	return ""
}
