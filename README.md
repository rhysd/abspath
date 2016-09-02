Tiny Absolute Path Library
==========================

`abspath` is a tiny library to handle file paths a bit better.

In golang, a file path is represented with `string` ([path/filepath]() standard library is provided).

However, treating file paths with `string` has some problems.  We always go wrong between absolute paths and relative paths.
When we create an API which takes a file path, we always need to take care about the parameter (It takes relative path? absolute path? or both of them?).
Or we need to design and implement the API to take both the absolute path and relative path.

I don't want to write a program caring about that throughout a program.  So I created a type to represent an absolute path.  It's just an strong type alias of `string`.

```go
type AbsPath string
```

It's further obvious because its type name is descriptive.  It's just an alias, so runtime overhead doesn't exist, I believe.
And it can't be relative path. For example,

```go
func writeSomethingTo(path AbsPath) error
```

when we see this API, we can know the `path` parameter is an absolute path, not a relative path.  And we can't fail to pass a relative path to above API because it never accept relative paths.

```go
// Error with ''
writeSomethingTo("relative_path") error
```

## Installation

```sh
$ go get github.com/rhysd/abspath
```

## Usage

### Create `AbsPath` Instance

```go
import github.com/rhysd/abspath

var a abspath.AbsPath
var e error

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
```

### Operate The Path

`AbsPath` has some methods deriving function defined in [`path/filepath`](https://golang.org/pkg/path/filepath) standard library.  They are helper to avoid converting between `AbsPath` and `string` frequently when you want to use functions in `path/filepath` package.

For example:

```go
a, _ = abspath.New("/absolute/path")
fmt.Printf("%s\n", a.Dir())
fmt.Printf("%s\n", a.Base())
fmt.Printf("%s\n", a.HasPrefix("/absolute"))
```

All methods are non-destructive because `AbsPath` is a type alias of `string` and strings in golang are immutable.

### Convert to `string`

You can simply convert to `string` by `string()` conversion.

```go
a, _ = abspath.New("/absolute/path")
file.WriteString(string(a))  // 'a' as a string
```

Vice versa, `string` can be converted to `AbsPath`.  We can't prevent it because golang permits the conversion.  So all we can do is prohibiting it.  We MUST use factory functions described above to create an `AbsPath` instance.
