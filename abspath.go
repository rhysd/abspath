package abspath

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// Type to represent absolute path.  Please do not make an instance of this struct directly.
// Instead, factory functions are available to create it.
type AbsPath struct {
	underlying string
}

// An error object type returned when specified value is not an absolute path.
type NotAbsolutePathError struct {
	specified string
}

func (err *NotAbsolutePathError) Error() string {
	return fmt.Sprintf("Not an absolute path: %s", err.specified)
}

// Create `AbsPath` struct instance from a string.  A parameter must represent an absolute path.
// If the parameter does not represent an absolute path, it returns an error as the second return value.
//
// Example:
//	a, err := abspath.New("/foo/bar")
//	if err != nil {
//		panic(err)
//	}
func New(from string) (AbsPath, error) {
	if !filepath.IsAbs(from) {
		return AbsPath{""}, &NotAbsolutePathError{from}
	}
	return AbsPath{filepath.Clean(from)}, nil
}

// Create `AbsPath` struct with expanding the parameter.  Parameter can be a full-path, relative path or a path starting with '~'
// where '~' means a home directory.  When parameter is a relative path, it will be joined with a path to current directory automatically.
//
// Example
//	a, err := abspath.ExpandFrom("/path/to/file")
//	b, err := abspath.EpandFrom("relative_path")
//	c, err := abspath.EpandFrom("~/Documents")
func ExpandFrom(maybe_relative string) (AbsPath, error) {
	if filepath.IsAbs(maybe_relative) {
		return AbsPath{filepath.Clean(maybe_relative)}, nil
	}

	if maybe_relative == "" {
		return AbsPath{""}, fmt.Errorf("Empty path cannot be expanded")
	}

	if maybe_relative[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsPath{""}, err
		}
		return AbsPath{filepath.Join(u.HomeDir, maybe_relative[1:])}, nil
	}

	p, err := filepath.Abs(maybe_relative)
	if err != nil {
		return AbsPath{""}, err
	}
	return AbsPath{p}, nil
}

// Create `AbsPath` struct instance from a string separated by slashes.  A parameter must represent an absolute path.
// It's similar to `filepath.FromSlash()` in `path/filepath` package.
//
// Ref: https://golang.org/pkg/path/filepath/#FromSlash
func FromSlash(s string) (AbsPath, error) {
	return New(filepath.FromSlash(s))
}

// Create `AbsPath` from slash separated string.  The same as `ExpandFrom()`, '~' is interpreted as a home directory
// and relative path will be joined with a path to current directory.
//
// Example:
//	// On Windows: e.g. Expanded to 'D:\path\to\cwd\relative\path'
//	a, err := ExpandFromSlash("relative/path")
func ExpandFromSlash(s string) (AbsPath, error) {
	return ExpandFrom(filepath.FromSlash(s))
}

// Create `AbsPath` for the working directory.  This is similar to `os.Getwd()` but returns `AbsPath` instead of string.
//
// Example:
//	cwd, err := os.Getwd()
//	if err != nil {
//		panic(err)
//	}
func Getwd() (AbsPath, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return AbsPath{""}, nil
	}
	return New(cwd)
}

// Equivalent to `filepath.Base()`.
//
// Ref: https://golang.org/pkg/path/filepath/#Base
func (a AbsPath) Base() AbsPath {
	return AbsPath{filepath.Base(a.underlying)}
}

// Equivalent to `filepath.Dir()`.
//
// Ref: https://golang.org/pkg/path/filepath/#Dir
func (a AbsPath) Dir() AbsPath {
	return AbsPath{filepath.Dir(a.underlying)}
}

// Equivalent to `filepath.EvalSymlinks()`.
//
// Ref: https://golang.org/pkg/path/filepath/#EvalSymlinks
func (a AbsPath) EvalSymlinks() (AbsPath, error) {
	s, err := filepath.EvalSymlinks(a.underlying)
	if err != nil {
		return AbsPath{""}, err
	}
	s, err = filepath.Abs(s)
	if err != nil {
		return AbsPath{""}, err
	}
	return AbsPath{s}, nil
}

// Equivalent to `filepath.Ext()`.
//
// Ref: https://golang.org/pkg/path/filepath/#Ext
func (a AbsPath) Ext() string {
	return filepath.Ext(a.underlying)
}

// Equivalent to `filepath.Join()`.  Parameters are joined into the absolute path.
//
// Ref: https://golang.org/pkg/path/filepath/#Join
func (a AbsPath) Join(elem ...string) AbsPath {
	switch len(elem) {
	case 0:
		return a
	case 1:
		return AbsPath{filepath.Join(a.underlying, elem[0])}
	default:
		return AbsPath{filepath.Join(a.underlying, filepath.Join(elem...))}
	}
}

// Equivalent to `filepath.Match()`.
//
// Ref: https://golang.org/pkg/path/filepath/#Match
func (a AbsPath) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, a.underlying)
}

// Equivalent to `filepath.Rel()`.  It returns a string of relative path to the absolute path.
//
// Ref: https://golang.org/pkg/path/filepath/#Rel
func (a AbsPath) Rel(targpath string) (string, error) {
	return filepath.Rel(a.underlying, targpath)
}

// Equivalent to `filepath.Split()`.  It returns an absolute path of parent directory and a string of its name.
//
// Ref: https://golang.org/pkg/path/filepath/#Split
func (a AbsPath) Split() (dir AbsPath, file string) {
	d, f := filepath.Split(a.underlying)
	return AbsPath{d}, f
}

// Equivalent to `filepath.ToSlash()`.
//
// Ref: https://golang.org/pkg/path/filepath/#ToSlash
func (a AbsPath) ToSlash() string {
	return filepath.ToSlash(a.underlying)
}

// Equivalent to `filepath.VolumeName()`.
//
// Ref: https://golang.org/pkg/path/filepath/#VolumeName
func (a AbsPath) VolumeName() string {
	return filepath.VolumeName(a.underlying)
}

// Equivalent to `filepath.Walk()`.
//
// Ref: https://golang.org/pkg/path/filepath/#Walk
func (a AbsPath) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(a.underlying, walkFn)
}

// Returns an underlying `string` value.  You can use this method to convert `AbsPath` value to string.
//
// Example
//	a, err := abspath.New("/foo/bar")
//	file.WriteString(a.String())
func (a AbsPath) String() string {
	return a.underlying
}

// Returns true if the absolute path entry exists.
func (a AbsPath) Exists() bool {
	_, err := os.Stat(a.underlying)
	return err == nil
}

// Returns true if the absolute path entry exists and is a directory.
func (a AbsPath) IsDir() bool {
	s, err := os.Stat(a.underlying)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Returns true if the absolute path entry exists and is a file.
func (a AbsPath) IsFile() bool {
	s, err := os.Stat(a.underlying)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Returns `FileInfo` of the path.  Equivalent to `os.Stat()`.
// Example:
//	a, _ := abspath.New("/path/to/file")
//	s, err := a.Stat()
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("Mode=%v\n", s.Mode())
func (a AbsPath) Stat() (os.FileInfo, error) {
	return os.Stat(a.underlying)
}
