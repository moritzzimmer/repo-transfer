repo-transfer [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
=============

[Transfers](https://help.github.com/articles/about-repository-transfers/) all Github repositories to another user or organization. See [Github api](https://developer.github.com/v3/repos/#transfer-a-repository) for details.


## dependencies

* (dev) [Go 1.11+](https://golang.org/dl/)


## install

```
go get -u github.com/spring-media/repo-transfer
```

## usage

```
repo-transfer                                        
  -source string
        github source organisation
  -target string
        github target organization
  -teams value
        optional team ids
  -token string
        oauth token
```

Example:

```
repo-transfer -teams 2767510 -source SourceOrga -target TargetOrga -token 72647qwe6qw7r67qwr6qw6rq7wr6
```


