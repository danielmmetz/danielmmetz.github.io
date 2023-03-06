all: html pdf

html:
    go run main.go -template index.html.tmpl -content index.yaml -minify > index.html

pdf:
    go run main.go -template index.html.tmpl -content index.yaml -minify -pdf > index.pdf
