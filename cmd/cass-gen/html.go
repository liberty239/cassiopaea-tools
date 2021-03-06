package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pterm/pterm"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"

	"github.com/liberty239/cassiopaea-tools/pkg/html"
	"github.com/liberty239/cassiopaea-tools/pkg/source/adapter"
	"github.com/liberty239/cassiopaea-tools/pkg/source/local"
	"github.com/liberty239/cassiopaea-tools/pkg/term"
)

const htmlCss = `
html {
	font-size: 14px
}
body {
	color: black;
	font-family: Georgia, Times, "Times New Roman", serif;
	font-size: 1.1rem;
	line-height: 1.5;
	text-align: justify;
}
h2 {
	color: black;
	font-family: Georgia, Times, "Times New Roman", serif;
	font-size: 2.2rem;
	font-weight: bold;
	line-height: 1.5;
	text-align: justify;
}
.bbCodeBlock {
	box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
}
.bbCodeBlock-content {
	padding: 2px 16px;
}
.bbCodeBlock-expandLink {
	display: none;
}
.sidenav {
	background-color: white;
	height: 100%;
	width: 12.5%; 
	position: fixed;
	z-index: 1;
	top: 0;
	box-shadow: 8px 0px 8px rgba(0,0,0,0.2);
	overflow-y: scroll;
}
.sidenav.a {
	padding-left: 16px;
}
.main {
	margin-left: 19.5%;
	margin-right: 19.5%;
	overflow-x: hidden;
	z-index: -1;
}
.span {
	box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
}
small {
	font-size: smaller;
}
.answer {
	background-color: rgb(232,232,232);
}
`

const htmlTpl = `
<html lang="en">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
{{range .MetaTags}}
	<meta name="{{.Key}}" content="{{.Value}}">
{{end}}
<style>{{.Css}}</style>
</head>

<body>
<div class="sidenav">
{{range .Data}}
	{{if .Year}}
    <h3>{{.Year}}</h3>
	{{end}}
    <a href="#{{.Id}}">{{.DateSansYear}}</a><br/>
{{end}}
<p></p>
</div>

{{if .Signature}}
<div class="main">
    <center><small>{{.Signature}}</small></center>
	<center><small>{{getCount .}} transcripts: {{.From}} - {{.To}}</small></center>
</div>
{{end}}

{{range .Data}}
	<section id={{.Id}}>
	<div class="main">
	    <h2>Session {{.Date}}</h2>
	    <a href="{{.Url}}" target="_blank">Forum discussion</a></br></br>
	    {{getBody .}}
	    </br>
	</div>
	</section>
{{end}}

</body>
</html>`

type htmlMetaTag struct {
	Key   string
	Value string
}

type htmlContextData struct {
	Id           string
	Date         string
	DateSansYear string
	Year         string
	Url          string

	rc io.ReadCloser
}

type htmlConetxt struct {
	Css       string
	Signature string
	From, To  string
	MetaTags  []htmlMetaTag
	Data      []htmlContextData
}

func buildSignature() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	name := user.Name
	if len(name) == 0 {
		name = "uid:" + user.Uid
	}

	host, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Built by %s@%s at %s",
		name,
		host,
		time.Now().Format(time.RFC1123),
	), nil
}

func parseMetaTags(args args) (ret []htmlMetaTag) {
	for _, x := range strings.Split(args.MetaTags, ",") {
		if tags := strings.Split(x, ":"); len(tags) == 2 {
			ret = append(ret, htmlMetaTag{
				Key:   tags[0],
				Value: tags[1],
			})
		}
	}
	return
}

func doHtml(args args) error {
	then := time.Now()

	s, err := term.Spinner().Start("Rendering...")
	if err != nil {
		s.Fail(err)
		return err
	}

	sig, err := buildSignature()
	if err != nil {
		s.Fail(err)
		return err
	}

	tpl, err := template.New("main").
		Funcs(template.FuncMap{
			"getBody": func(s *htmlContextData) (ret string) {
				q, err := goquery.NewDocumentFromReader(s.rc)
				if err != nil {
					s.rc.Close()
					panic(err)
				}
				if err := s.rc.Close(); err != nil {
					panic(err)
				}

				q.Find("body").EachWithBreak(func(i int, s *goquery.Selection) bool {
					ret, err = s.Html()
					if err != nil {
						panic(err)
					}
					return false
				})
				return
			},
			"getCount": func(ctx *htmlConetxt) string {
				return strconv.Itoa(len(ctx.Data))
			},
		}).
		Parse(htmlTpl)
	if err != nil {
		s.Fail(err)
	}

	ctx := htmlConetxt{
		Css:       htmlCss,
		Signature: sig,
		MetaTags:  parseMetaTags(args),
	}

	src := adapter.Chain(
		local.NewSource(args.Directory),
		adapter.Ascending,
	)
	if err := src.FetchAll(func(ts time.Time, url string, rc io.ReadCloser) error {
		ctx.Data = append(ctx.Data, htmlContextData{
			Id:           fmt.Sprintf("%d", ts.Unix()),
			Date:         ts.Format("2 January 2006"),
			DateSansYear: ts.Format("02 January"),
			Year:         ts.Format("2006"),
			Url:          url,
			rc:           rc,
		})
		return nil
	}); err != nil {
		s.Fail(err)
		return err
	}

	if len(ctx.Data) == 0 {
		err := errors.New("no sessions")
		s.Fail(err)
		return err
	}
	ctx.From = ctx.Data[0].Date
	ctx.To = ctx.Data[len(ctx.Data)-1].Date

	var year string
	for i, x := range ctx.Data {
		if x.Year != year {
			year = x.Year
			continue
		}
		ctx.Data[i].Year = ""
	}

	f, err := os.Create(args.File)
	if err != nil {
		s.Fail(err)
		return nil
	}

	m := minify.New()
	m.AddFunc("text/html", mhtml.Minify)
	m.AddFunc("text/css", css.Minify)
	wc := m.Writer("text/html", f)

	// Hacks to make GoogleBot happy...
	mhtml.DefaultMinifier.KeepDocumentTags = true
	mhtml.DefaultMinifier.KeepQuotes = true

	if err := html.TemplateExecute(
		tpl,
		wc,
		&ctx,
	); err != nil {
		s.Fail(err)
		return err
	}

	if err := wc.Close(); err != nil {
		s.Fail(err)
		return err
	}

	if err := f.Close(); err != nil {
		s.Fail(err)
		return err
	}

	s.Success("Sessions written to HTML file: ", args.File)
	pterm.Info.Println("Rendering took", time.Since(then).Seconds(), "seconds")

	return nil
}
