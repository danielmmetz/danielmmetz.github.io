package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"text/template"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/ghodss/yaml"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/yuin/goldmark"
)

var (

	//go:embed embed/css.tmpl
	cssTmplB string
	//go:embed embed/DMSansRegular.ttf
	dmSansRegular []byte
	//go:embed embed/DMSansBold.ttf
	dmSansBold []byte
	//go:embed embed/RedHatTextRegular.ttf
	redHatTextRegular []byte
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := mainE(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err.Error())
		os.Exit(1)
	}
}

func mainE(ctx context.Context) error {
	templatePath := flag.String("template", "", "path to template file")
	contentPath := flag.String("content", "", "path to content file")
	verbose := flag.Bool("verbose", false, "if true, print eggregiously")
	shrink := flag.Bool("minify", false, "if true, minify output")
	pdf := flag.Bool("pdf", false, "if true, write pdf instead of html")
	flag.Parse()

	logger := verboseLogger{enabled: *verbose}

	minifier := minify.New()
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("text/css", css.Minify)

	cssTemplate := template.New("css").Funcs(map[string]any{"markdownify": markdownify, "base64": b64})
	cssTmpl, err := cssTemplate.Parse(cssTmplB)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	cssD := cssData{
		DMSansRegular:     dmSansRegular,
		DMSansBold:        dmSansBold,
		RedHatTextRegular: redHatTextRegular,
	}
	var cssOut bytes.Buffer
	if err := cssTmpl.Execute(&cssOut, cssD); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	if *shrink {
		var tmp bytes.Buffer
		if err := minifier.Minify("text/css", &tmp, &cssOut); err != nil {
			return fmt.Errorf("minifying css: %w", err)
		}
		cssOut = tmp
	}

	contentB, err := os.ReadFile(*contentPath)
	if err != nil {
		return fmt.Errorf("reading content file: %s: %w", *contentPath, err)
	}
	var d content
	if err := yaml.Unmarshal(contentB, &d); err != nil {
		return fmt.Errorf("parsing content: %w", err)
	}
	logger.Printf("config:\n%+v\n", d)
	d.Static.CSS = cssOut.String()

	htmlTemplate := template.New("html").Funcs(map[string]any{"markdownify": markdownify, "base64": b64})
	htmlTemplateB, err := os.ReadFile(*templatePath)
	if err != nil {
		return fmt.Errorf("reading template file %s: %w", *templatePath, err)
	}
	tmpl, err := htmlTemplate.Parse(string(htmlTemplateB))
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	var out bytes.Buffer
	if err := tmpl.Execute(&out, d); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}
	if *shrink {
		var b bytes.Buffer
		if err := minifier.Minify("text/html", &b, &out); err != nil {
			return fmt.Errorf("minifying: %w", err)
		}
		out = b
	}
	if !*pdf {
		fmt.Println(out.String())
		return nil
	}

	f, err := os.CreateTemp("", "*.html")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(f.Name())

	if _, err := f.Write(out.Bytes()); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(logger.Printf))
	defer cancel()

	var pdfBuffer []byte
	grabber := pdfGrabber(fmt.Sprintf("file://%s", f.Name()), "body", &pdfBuffer, d.PDF.MarginTop, d.PDF.MarginBottom)
	if err := chromedp.Run(ctx, grabber); err != nil {
		return fmt.Errorf("generating pdf: %w", err)
	}
	fmt.Printf("%s", pdfBuffer)
	return nil
}

type cssData struct {
	DMSansRegular     []byte
	DMSansBold        []byte
	RedHatTextRegular []byte
}

type content struct {
	Header struct {
		Name  string
		Email string
		Site  string
	}
	Employment []struct {
		Title          string
		Employer       string
		Time           string
		PreviousTitles []struct {
			Title string
			Time  string
		}
		Roles []struct {
			Title   string
			Time    string
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
		CSS string
	}
	PDF struct {
		MarginTop    float64
		MarginBottom float64
	}
}

func markdownify(s string) (string, error) {
	var b bytes.Buffer
	if err := goldmark.Convert([]byte(s), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

func b64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

type verboseLogger struct {
	enabled bool
}

func (l verboseLogger) Printf(format string, a ...any) {
	if !l.enabled {
		return
	}
	fmt.Fprintf(os.Stderr, format, a...)
}

func (l verboseLogger) Println(a ...any) {
	if !l.enabled {
		return
	}
	fmt.Fprintln(os.Stderr, a...)
}

// pdfGrabber is largely taken from https://stackoverflow.com/a/68796203.
func pdfGrabber(url string, sel string, res *[]byte, topMargin, bottomMargin float64) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(sel, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithMarginTop(topMargin).WithMarginBottom(bottomMargin).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
