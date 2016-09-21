Tiny Absolute Path Library
==========================
[![Build Status](https://travis-ci.org/rhysd/abspath.svg?branch=master)](https://travis-ci.org/rhysd/abspath)
[![Build status](https://ci.appveyor.com/api/projects/status/usfx6p4xff31sn7e/branch/master?svg=true)](https://ci.appveyor.com/project/rhysd/abspath/branch/master)
[![Coverage Status](https://coveralls.io/repos/github/rhysd/abspath/badge.svg?branch=master)](https://coveralls.io/github/rhysd/abspath?branch=master)
[![GoDoc](https://godoc.org/github.com/rhysd/abspath?status.svg)](https://godoc.org/github.com/rhysd/abspath)

`abspath` is a tiny library to handle file paths a bit better.

In golang, a file path is represented with `string` ([path/filepath]() standard library is provided).

However, treating file paths with `string` has some problems.  We always go wrong between absolute paths and relative paths.
When we create an API which takes a file path, we always need to take care about the parameter (It takes relative path? absolute path? or both of them?).
Or we need to design and implement the API to take both the absolute path and relative path.

I don't want to write a program caring about that throughout a program.  So I created a type to represent an absolute path.  It's just a wrapper of `string`.

```go
type AbsPath struct {
    underlying string
}
```

It's further obvious because its type name is descriptive.  And it can't be relative path. For example,

```go
func writeSomethingTo(path AbsPath) error
```

when we see this API, we can know the `path` parameter is an absolute path, not a relative path.  And we can't fail to pass a relative path to above API because it never accept relative paths.

```go
// Error with 'cannot use a (type AbsPath) as type string in argument'
writeSomethingTo("relative_path")
```

In addition, `AbsPath` is kept clean.

## Installation

```sh
$ go get github.com/rhysd/abspath
```

## Usage

### Create `AbsPath` Instance

```go
import github.com/rhysd/abspath

var a abspath.AbsPath
var err error

// From a string representing absolute path.
// If it doesn't stand for absolute path, an error will be returned.
a, err = abspath.New("/absolute/path")

// Expanded to $PWD/relative_path
a, err = abspath.ExpandFrom("relative_path")

// Expanded to $HOME/relative_path
a, err = abspath.ExpandFrom("~/relative_path")

// Slashes in the string will be replaced by a file separator.
// In Windows, below can be 'C:\absolute\path'
a, err = abspath.FromSlash("/absoltue/path")

// Similar to ExpandFrom().  But all slashes will be replaced with a file sperator.
a, err = abspath.ExpandFromSlash("relative/path")

// Shortcut of os.Getwd() for AbsPath.
a, err = abspath.Getwd()

// Shortcut for getting an absolute path to home directory.
a, err = abspath.HomeDir()
```

### Operate The Path

`AbsPath` has some methods deriving function defined in [`path/filepath`](https://golang.org/pkg/path/filepath) standard library.  They are helpers to avoid converting between `AbsPath` and `string` frequently when you want to use functions in `path/filepath` package. All method are non-destructive.

For example:

```go
a, _ := abspath.New("/absolute/path")
fmt.Printf("%s\n", a.Dir())
fmt.Printf("%s\n", a.Base())
fmt.Printf("%s\n", a.HasPrefix("/absolute"))
```

Please see [GoDoc](https://godoc.org/github.com/rhysd/abspath) for full documentation.

### Convert to `string`

You can simply convert to `string` by `.String()` method.

```go
a, _ = abspath.New("/absolute/path")
file.WriteString(a.String())  // 'a' as a string
```


## License

Distributed under [Public Domain 1.0](https://creativecommons.org/publicdomain/mark/1.0/)
