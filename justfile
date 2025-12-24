all: resume-html resume-pdf hugo-build

[working-directory: 'resume']
resume-html:
    go run main.go -template index.html.tmpl -content index.yaml -minify >> ../home/static/resume/index.html

[working-directory: 'resume']
resume-pdf:
    go run main.go -template index.html.tmpl -content index.yaml -minify -pdf > ../home/static/resume.pdf

[working-directory: 'home']
hugo-build:
    rm -r ../docs/*
    hugo build --minify --destination ../docs

[working-directory: 'home']
hugo-dev:
    hugo server --buildDrafts --disableFastRender --destination ../docs
