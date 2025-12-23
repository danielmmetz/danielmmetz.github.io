all: resume-html resume-pdf

[working-directory: 'resume']
resume-html:
    go run main.go -template index.html.tmpl -content index.yaml -minify > ../docs/index.html

[working-directory: 'resume']
resume-pdf:
    go run main.go -template index.html.tmpl -content index.yaml -minify -pdf > ../docs/resume.pdf
