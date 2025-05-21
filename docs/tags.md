[Reference](https://go.dev/ref/mod#vcs-version)

> each tag name must be prefixed with the module subdirectory, followed by a slash.

Example tag: `http/v0.0.1`

```shell
git tag -s 'http/v0.0.2' -m 'http v0.0.2'
git push origin 'http/v0.0.2'
```
