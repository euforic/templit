# templit
[![Go Report Card](https://goreportcard.com/badge/github.com/euforic/templit)](https://goreportcard.com/report/github.com/euforic/templit)
[![GoDoc](https://godoc.org/github.com/euforic/templit?status.svg)](https://godoc.org/github.com/euforic/templit)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/euforic/templit)
![Build Status](https://github.com/euforic/templit/workflows/Go/badge.svg)
![GitHub](https://img.shields.io/github/license/euforic/templit)

```go
import "github.com/euforic/templit"
```

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
    "default":      defaultVal,
}
```

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

ToKebabCase converts a string to kebab-case.

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

ToSnakeCase converts a string to snake_case.

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
func NewDefaultGitClient(defaultBranch, token string) *DefaultGitClient
```

NewDefaultGitClient creates a new DefaultGitClient with the optional defaultBranch and token.

<a name="DefaultGitClient.Checkout"></a>
### func (*DefaultGitClient) Checkout

```go
func (d *DefaultGitClient) Checkout(path, branch string) error
```

Checkout checks out a branch in a Git repository.

<a name="DefaultGitClient.Clone"></a>
### func (*DefaultGitClient) Clone

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
    git GitClient
}
```

<a name="NewExecutor"></a>
### func NewExecutor

```go
func NewExecutor(client GitClient) *Executor
```

New returns a new Executor

<a name="Executor.EmbedFunc"></a>
### func (*Executor) EmbedFunc

```go
func (e *Executor) EmbedFunc() func(remotePath string, data interface{}) (string, error)
```

EmbedFunc returns a template function that can be used to process and embed a template from a remote git repository. EmbedFunc allows embedding content from a remote repository directly into a Go template.

Steps to use:

1. Add the function to the FuncMap.
2. Use the following syntax within your template:
    ```
    {{ embed "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" . }}
    {{ embed "<host>/<owner>/<repo>#<block>@<tag_or_hash_or_branch>" . }}
    ```
Placeholders:

- `<host>`: Repository hosting service (e.g., "github.com").
- `<owner>`: Repository owner or organization.
- `<repo>`: Repository name.
- `<path>`: Path to the desired file or directory within the repository.
- `<block>`: Specific template block name.
- `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).

<a name="Executor.ImportFunc"></a>
### func (*Executor) ImportFunc

```go
func (e *Executor) ImportFunc(destPath string) func(repoAndTag, path string, data interface{}) (string, error)
```

ImportFunc returns a function that can be used as a template function to import and process a template from a remote git repository. ImportFunc allows embedding content from a remote repository into a Go template.

Steps to use:

1. Add the function to the FuncMap.
2. Use the following syntax within your template:
    ```
    `{{ import "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" "<path_to_genrate_files>" . }}
    `{{ import "<host>/<owner>/<repo>/<path>@<tag_or_hash_or_branch>" "<path_to_genrate_files>" . }}
    ```
Placeholders:

- `<host>`: Repository hosting service (e.g., "github.com").
- `<owner>`: Repository owner or organization.
- `<repo>`: Repository name.
- `<path>`: Path to the desired file or directory within the repository.
- `<tag_or_hash_or_branch>`: Specific Git reference (tag, commit hash, or branch name).

<a name="Executor.ParsePath"></a>
### func (*Executor) ParsePath

```go
func (e *Executor) ParsePath(inputPath string) error
```

ParsePath parses the given path

<a name="Executor.Render"></a>
### func (Executor) Render

```go
func (e Executor) Render(name string, data interface{}) (string, error)
```

Render executes the template with the given data

<a name="Executor.StringRender"></a>
### func (Executor) StringRender

```go
func (e Executor) StringRender(templateString string, data interface{}) (string, error)
```

StringRender renders the given template string with the given data

<a name="Executor.WalkAndProcessDir"></a>
### func (*Executor) WalkAndProcessDir

```go
func (e *Executor) WalkAndProcessDir(inputDir, outputDir string, data interface{}) error
```

WalkAndProcessDir processes all files in a directory with the given data. If walkFunc is provided, it's called for each file and directory without writing the file to disk.

<a name="GitClient"></a>
## type GitClient

GitClient is an interface that abstracts Git operations.

```go
type GitClient interface {
    Clone(host, owner, repo, dest string) error
    Checkout(path, branch string) error
    DefaultBranch() string
}
```
