package abspath

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// AbsPath is a type to represent absolute path.  Please do not make an instance of this struct directly.
// Instead, factory functions are available to create it.
type AbsPath struct {
	underlying string
}

// NotAbsolutePathError is an error object type returned when specified value is not an absolute path.
type NotAbsolutePathError struct {
	specified string
}

func (err *NotAbsolutePathError) Error() string {
	return fmt.Sprintf("Not an absolute path: '%s'", err.specified)
}

// New creates AbsPath struct instance from a string.  A parameter must represent an absolute path.
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

// ExpandFrom creates AbsPath struct with expanding the parameter.  Parameter can be a full-path, relative path or a path starting with '~'
// where '~' means a home directory.  When parameter is a relative path, it will be joined with a path to current directory automatically.
//
// Example
//	a, err := abspath.ExpandFrom("/path/to/file")
//	b, err := abspath.ExpandFrom("relative_path")
//	c, err := abspath.ExpandFrom("~/Documents")
func ExpandFrom(specified string) (AbsPath, error) {
	if filepath.IsAbs(specified) {
		return AbsPath{filepath.Clean(specified)}, nil
	}

	if specified == "" {
		return AbsPath{""}, &NotAbsolutePathError{""}
	}

	if specified[0] == '~' {
		u, err := user.Current()
		if err != nil {
			return AbsPath{""}, err
		}
		return AbsPath{filepath.Join(u.HomeDir, specified[1:])}, nil
	}

	p, err := filepath.Abs(specified)
	if err != nil {
		return AbsPath{""}, err
	}
	return AbsPath{p}, nil
}

// FromSlash creates AbsPath struct instance from a string separated by slashes.  A parameter must represent an absolute path.
// It's similar to filepath.FromSlash() in path/filepath package.
//
// Ref: https://golang.org/pkg/path/filepath/#FromSlash
func FromSlash(s string) (AbsPath, error) {
	return New(filepath.FromSlash(s))
}

// ExpandFromSlash creates AbsPath from slash separated string.  The same as ExpandFrom(), '~' is interpreted as a home directory
// and relative path will be joined with a path to current directory.
//
// Example:
//	// On Windows: e.g. Expanded to 'D:\path\to\cwd\relative\path'
//	a, err := ExpandFromSlash("relative/path")
func ExpandFromSlash(s string) (AbsPath, error) {
	return ExpandFrom(filepath.FromSlash(s))
}

// Getwd creates AbsPath for the working directory.  This is similar to os.Getwd() but returns AbsPath instead of string.
//
// Example:
//	cwd, err := abspath.Getwd()
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

// HomeDir creates AbsPath instance for the home directory.  If home directory cannot be obtained or is not an absolute path, it will return an error.
//
// Example:
//	home, err := abspath.HomeDir()
//	if err != nil {
//		panic(err)
//	}
//	println(home.String())
func HomeDir() (AbsPath, error) {
	u, err := user.Current()
	if err != nil {
		return AbsPath{""}, err
	}
	return New(u.HomeDir)
}

// Base is equivalent to filepath.Base().
//
// Ref: https://golang.org/pkg/path/filepath/#Base
func (a AbsPath) Base() AbsPath {
	return AbsPath{filepath.Base(a.underlying)}
}

// Dir is equivalent to filepath.Dir().
//
// Ref: https://golang.org/pkg/path/filepath/#Dir
func (a AbsPath) Dir() AbsPath {
	return AbsPath{filepath.Dir(a.underlying)}
}

// EvalSymlinks is equivalent to filepath.EvalSymlinks().
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

// Ext is equivalent to filepath.Ext().
//
// Ref: https://golang.org/pkg/path/filepath/#Ext
func (a AbsPath) Ext() string {
	return filepath.Ext(a.underlying)
}

// HasPrefix is equivalent to filepath.HasPrefix().
//
// Ref: https://golang.org/pkg/path/filepath/#HasPrefix
func (a AbsPath) HasPrefix(prefix string) bool {
	return filepath.HasPrefix(a.underlying, prefix)
}

// Join is equivalent to filepath.Join().  Parameters are joined into the absolute path.
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

// Match is equivalent to filepath.Match().  It returns the absolute path matches the given pattern.
//
// Ref: https://golang.org/pkg/path/filepath/#Match
func (a AbsPath) Match(pattern string) (bool, error) {
	return filepath.Match(pattern, a.underlying)
}

// Rel is equivalent to filepath.Rel().  It returns a string of relative path to the absolute path.
//
// Ref: https://golang.org/pkg/path/filepath/#Rel
func (a AbsPath) Rel(targpath string) (string, error) {
	return filepath.Rel(a.underlying, targpath)
}

// Split is equivalent to filepath.Split().  It returns an absolute path of parent directory and a string of its name.
//
// Ref: https://golang.org/pkg/path/filepath/#Split
func (a AbsPath) Split() (dir AbsPath, file string) {
	d, f := filepath.Split(a.underlying)
	return AbsPath{d}, f
}

// ToSlash is equivalent to filepath.ToSlash().
//
// Ref: https://golang.org/pkg/path/filepath/#ToSlash
func (a AbsPath) ToSlash() string {
	return filepath.ToSlash(a.underlying)
}

// VolumeName is equivalent to filepath.VolumeName().
//
// Ref: https://golang.org/pkg/path/filepath/#VolumeName
func (a AbsPath) VolumeName() string {
	return filepath.VolumeName(a.underlying)
}

// Walk is equivalent to filepath.Walk().
//
// Ref: https://golang.org/pkg/path/filepath/#Walk
func (a AbsPath) Walk(walkFn filepath.WalkFunc) error {
	return filepath.Walk(a.underlying, walkFn)
}

// String returns an underlying string value.  You can use this method to convert AbsPath value to string.
//
// Example
//	a, err := abspath.New("/foo/bar")
//	file.WriteString(a.String())
func (a AbsPath) String() string {
	return a.underlying
}
