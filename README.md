### [Personal resume site](https://www.danielmmetz.com) generator

Kick it off via [just](https://github.com/casey/just).

```
❯ just all  # generates both index.html and index.pdf
```

Or do it manually via the commands written under the hood:
```
❯ go run main.go -template index.html.tmpl -content index.yaml -minify > index.html
❯ go run main.go -template index.html.tmpl -content index.yaml -minify -pdf > index.pdf
```
