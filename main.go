package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/yuin/goldmark"
)

//go:embed static/*
var fs embed.FS

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := mainE(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainE(_ context.Context) error {
	templatePath := flag.String("template", "", "path to template file")
	contentPath := flag.String("content", "", "path to content file")
	verbose := flag.Bool("verbose", false, "if true, print eggregiously")
	shrink := flag.Bool("minify", false, "if true, minify output")
	flag.Parse()

	logger := verboseLogger{enabled: *verbose}

	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)

	contentB, err := os.ReadFile(*contentPath)
	if err != nil {
		return fmt.Errorf("reading content file: %s: %w", *contentPath, err)
	}
	var d data
	if err := yaml.Unmarshal(contentB, &d); err != nil {
		return fmt.Errorf("parsing content: %w", err)
	}
	logger.Printf("config:\n%+v\n", d)
	cssB, err := fs.ReadFile("static/main.css")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read static/main.css: %w", err)
	}
	cssBuf := bytes.NewBuffer(cssB)
	if *shrink {
		var cssOut bytes.Buffer
		if err := minifier.Minify("text/css", &cssOut, cssBuf); err != nil {
			return fmt.Errorf("minifying css: %w", err)
		}
		cssBuf = &cssOut
	}
	d.Static.CSS = cssBuf.String()

	templateB, err := os.ReadFile(*templatePath)
	if err != nil {
		return fmt.Errorf("reading template file %s: %w", *templatePath, err)
	}
	template, err := template.New("template").Funcs(map[string]any{"markdownify": markdownify}).Parse(string(templateB))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	var out bytes.Buffer
	if err := template.Execute(&out, d); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	if *shrink {
		var b bytes.Buffer
		if err := minifier.Minify("text/html", &b, &out); err != nil {
			return fmt.Errorf("minifying: %w", err)
		}
		out = b
	}
	fmt.Println(out.String())
	return nil
}

type data struct {
	Header struct {
		Name string
		Site string
	}
	Employment []struct {
		Title    string
		Employer string
		Time     string
		Roles    []struct {
			Title   string
			Content string
		}
	}
	Education struct {
		School  string
		Time    string
		Content string
	}
	Extras []struct {
		Title   string
		Content string
	}
	Static struct {
		CSS   string
		Fonts []struct {
		}
	}
}

func markdownify(s string) (string, error) {
	var b bytes.Buffer
	if err := goldmark.Convert([]byte(s), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

type verboseLogger struct {
	enabled bool
}

func (l verboseLogger) Printf(format string, a ...any) (n int, err error) {
	if !l.enabled {
		return 0, nil
	}
	return fmt.Fprintf(os.Stderr, format, a...)
}

func (l verboseLogger) Println(a ...any) (n int, err error) {
	if !l.enabled {
		return 0, nil
	}
	return fmt.Fprintln(os.Stderr, a...)
}
