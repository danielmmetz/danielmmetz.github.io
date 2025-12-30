all: resume-html resume-pdf hugo-build

[working-directory: 'resume']
resume-html:
    cp resume-frontmatter.md ../home/content/resume/_index.html
    go run main.go -template index.html.tmpl -content index.yaml -minify >> ../home/content/resume/_index.html

[working-directory: 'resume']
resume-pdf:
    go run main.go -template index.html.tmpl -content index.yaml -minify -pdf > ../home/static/resume.pdf

[working-directory: 'home']
hugo-build:
    hugo build --destination ../docs

[working-directory: 'home']
hugo-dev:
    hugo server --buildDrafts
