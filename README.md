# templit
[![Go Report Card](https://goreportcard.com/badge/github.com/euforic/templit)](https://goreportcard.com/report/github.com/euforic/templit)
[![GoDoc](https://godoc.org/github.com/euforic/templit?status.svg)](https://godoc.org/github.com/euforic/templit)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/euforic/templit)
![Build Status](https://github.com/euforic/templit/workflows/Go/badge.svg)
![GitHub](https://img.shields.io/github/license/euforic/templit)

```go
import "github.com/euforic/templit"
```

## Index

- [Variables](<#variables>)
- [func EmbedFunc\(client GitClient\) func\(remotePath string, data interface\{\}, funcMap template.FuncMap\) \(string, error\)](<#EmbedFunc>)
- [func ImportFunc\(client GitClient\) func\(repoAndTag, destPath string, data interface\{\}, funcMap template.FuncMap\) \(string, error\)](<#ImportFunc>)
- [func RenderTemplate\(tmpl string, data interface\{\}, funcMap template.FuncMap\) \(string, error\)](<#RenderTemplate>)
- [func ToCamelCase\(s string\) string](<#ToCamelCase>)
- [func ToKebabCase\(s string\) string](<#ToKebabCase>)
- [func ToPascalCase\(s string\) string](<#ToPascalCase>)
- [func ToSnakeCase\(s string\) string](<#ToSnakeCase>)
- [func WalkAndProcessDir\(inputDir, outputDir string, funcMap template.FuncMap, data interface\{\}\) error](<#WalkAndProcessDir>)
- [type DefaultGitClient](<#DefaultGitClient>)
  - [func NewDefaultGitClient\(token string\) \*DefaultGitClient](<#NewDefaultGitClient>)
  - [func \(d \*DefaultGitClient\) Checkout\(path, branch string\) error](<#DefaultGitClient.Checkout>)
  - [func \(d \*DefaultGitClient\) Clone\(host, owner, repo, dest string\) error](<#DefaultGitClient.Clone>)
- [type DepInfo](<#DepInfo>)
  - [func ParseDepURL\(rawURL string\) \(\*DepInfo, error\)](<#ParseDepURL>)
- [type Executor](<#Executor>)
  - [func NewExecutor\(inputPath string, funcMap template.FuncMap\) \(\*Executor, error\)](<#NewExecutor>)
  - [func \(e Executor\) Render\(name string, data interface\{\}\) \(string, error\)](<#Executor.Render>)
- [type GitClient](<#GitClient>)
- [type WalkAndProcessDirFunc](<#WalkAndProcessDirFunc>)


## Variables

<a name="DefaultFuncMap"></a>DefaultFuncMap is the default function map for templates.

```go
var DefaultFuncMap = template.FuncMap{
    "lower":        strings.ToLower,
    "upper":        strings.ToUpper,
    "trim":         strings.TrimSpace,
    "split":        strings.Split,
    "join":         strings.Join,
    "replace":      strings.ReplaceAll,
    "contains":     strings.Contains,
    "hasPrefix":    strings.HasPrefix,
    "hasSuffix":    strings.HasSuffix,
    "trimPrefix":   strings.TrimPrefix,
    "trimSuffix":   strings.TrimSuffix,
    "trimSpace":    strings.TrimSpace,
    "trimLeft":     strings.TrimLeft,
    "trimRight":    strings.TrimRight,
    "count":        strings.Count,
    "repeat":       strings.Repeat,
    "equalFold":    strings.EqualFold,
    "splitN":       strings.SplitN,
    "splitAfter":   strings.SplitAfter,
    "splitAfterN":  strings.SplitAfterN,
    "fields":       strings.Fields,
    "toTitle":      strings.ToTitle,
    "toSnakeCase":  ToSnakeCase,
    "toCamelCase":  ToCamelCase,
    "toKebabCase":  ToKebabCase,
    "toPascalCase": ToPascalCase,
    "default": DefaultVal,
}
```

<a name="EmbedFunc"></a>
## func EmbedFunc

```go
func EmbedFunc(client GitClient) func(remotePath string, data interface{}, funcMap template.FuncMap) (string, error)
```

EmbedFunc returns a template function that can be used to process and embed a template from a remote git repository. EmbedFunc allows embedding content from a remote repository directly into a Go template.

Steps to use:

1. Add the function to the FuncMap.
2. Use the following syntax within your template:
    ```go
    {{ embed "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" . }}
    {{ embed "<host>/<owner>/<repo>#<block>@<tag_or_hash_or_branch>" . }}
    ```

Placeholders:

- `<host>`: Repository hosting service (e.g., "github.com").
- `<owner>`: Repository owner or organization.
- `<repo>`: Repository name.
- `<path>`: Path to the desired template file within the repository.
- `<block>`: Specific template block name.
- `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).

<a name="ImportFunc"></a>
## func ImportFunc

```go
func ImportFunc(client GitClient) func(repoAndTag, destPath string, data interface{}, funcMap template.FuncMap) (string, error)
```

ImportFunc returns a function that can be used as a template function to import and process a template from a remote git repository. ImportFunc allows embedding content from a remote repository into a Go template.

Steps to use:

1. Add the function to the FuncMap.
2. Use the following syntax within your template:
    ```go
    {{ import "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" "<path_to_genrate_files>" . }}
    ```

Placeholders:

- `<host>`: Repository hosting service (e.g., "github.com").
- `<owner>`: Repository owner or organization.
- `<repo>`: Repository name.
- `<path>`: Path to the desired file or directory within the repository.
- `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).


<a name="RenderTemplate"></a>
## func RenderTemplate

```go
func RenderTemplate(tmpl string, data interface{}, funcMap template.FuncMap) (string, error)
```

RenderTemplate renders a template with provided data.

<a name="ToCamelCase"></a>
## func ToCamelCase

```go
func ToCamelCase(s string) string
```

ToCamelCase converts a string to CamelCase.

<a name="ToKebabCase"></a>
## func ToKebabCase

```go
func ToKebabCase(s string) string
```

ToKebabCase converts a string to kebab\-case.

<a name="ToPascalCase"></a>
## func ToPascalCase

```go
func ToPascalCase(s string) string
```

ToPascalCase converts a string to PascalCase.

<a name="ToSnakeCase"></a>
## func ToSnakeCase

```go
func ToSnakeCase(s string) string
```

ToSnakeCase converts a string to snake\_case.

<a name="WalkAndProcessDir"></a>
## func WalkAndProcessDir

```go
func WalkAndProcessDir(inputDir, outputDir string, funcMap template.FuncMap, data interface{}) error
```

WalkAndProcessDir processes all files in a directory with the given data. If walkFunc is provided, it's called for each file and directory without writing the file to disk.

<a name="DefaultGitClient"></a>
## type DefaultGitClient

DefaultGitClient provides a default implementation for the GitClient interface.

```go
type DefaultGitClient struct {
    Token   string
    BaseURL string
}
```

<a name="NewDefaultGitClient"></a>
### func NewDefaultGitClient

```go
func NewDefaultGitClient(token string) *DefaultGitClient
```

NewDefaultGitClient creates a new DefaultGitClient with the given token.

<a name="DefaultGitClient.Checkout"></a>
### func \(\*DefaultGitClient\) Checkout

```go
func (d *DefaultGitClient) Checkout(path, branch string) error
```

Checkout checks out a branch in a Git repository.

<a name="DefaultGitClient.Clone"></a>
### func \(\*DefaultGitClient\) Clone

```go
func (d *DefaultGitClient) Clone(host, owner, repo, dest string) error
```

Clone clones a Git repository to the given destination.

<a name="DepInfo"></a>
## type DepInfo

DepInfo contains information about an embed URL.

```go
type DepInfo struct {
    Host  string
    Owner string
    Repo  string
    Path  string
    Block string
    Tag   string
}
```

<a name="ParseDepURL"></a>
### func ParseDepURL

```go
func ParseDepURL(rawURL string) (*DepInfo, error)
```

ParseDepURL is a parsed embed URL.

<a name="Executor"></a>
## type Executor

Executor is a wrapper around the template.Template type

```go
type Executor struct {
    *template.Template
}
```

<a name="NewExecutor"></a>
### func NewExecutor

```go
func NewExecutor(inputPath string, funcMap template.FuncMap) (*Executor, error)
```

NewExecutor creates a new Executor with the given template and funcMap

<a name="Executor.Render"></a>
### func \(Executor\) Render

```go
func (e Executor) Render(name string, data interface{}) (string, error)
```

Render executes the template with the given data

<a name="GitClient"></a>
## type GitClient

GitClient is an interface that abstracts Git operations.

```go
type GitClient interface {
    Clone(host, owner, repo, dest string) error
    Checkout(path, branch string) error
}
```

<a name="WalkAndProcessDirFunc"></a>
## type WalkAndProcessDirFunc

WalkAndProcessDirFunc is called for each file and directory when walking a directory.

```go
type WalkAndProcessDirFunc func(path string, isDir bool, content string) error
```


