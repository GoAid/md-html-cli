package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	d2 "github.com/FurqanSoftware/goldmark-d2"
	"github.com/GoAid/md-html-cli/assets"
	"github.com/PuerkitoBio/goquery"
	cfh "github.com/alecthomas/chroma/v2/formatters/html"
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/mermaid"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
)

type Options struct {
	InputFile  string `long:"input" short:"i" description:"input Markdown"`
	OutputFile string `long:"output" short:"o" description:"output HTML"`

	HTMLLang    string `long:"lang" short:"l" description:"html lang attribute value, default is en"`
	HTMLTitle   string `long:"title" short:"t" description:"custom html title, default is output file name"`
	HTMLFavicon string `long:"favicon" short:"f" description:"favicon image path, if embed is used, will embed by base64 encoding"`
	EmbedImage  bool   `long:"embed" short:"e" description:"embed image by base64 encoding"`

	MathJax   bool `long:"mathjax" short:"m" description:"use MathJax"`
	TableSpan bool `long:"span" short:"s" description:"enable table row/col span"`

	CustomCSS string `long:"css" short:"c" description:"custom css file path"`
	Theme     string `long:"theme" choice:"vue" choice:"side" description:"output HTML theme"`
	TOC       bool   `long:"toc" description:"generate TOC"`
	Generated bool   `long:"gen" short:"g" description:"use HTML comments to record generation time"`
}

// goldmark convert options
var (
	extensions = []goldmark.Extender{
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		extension.Typographer,
		emoji.Emoji,
		//mathjax.MathJax,
		new(fences.Extender),
		new(mermaid.Extender),
		&d2.Extender{
			Layout:  d2dagrelayout.DefaultLayout,
			ThemeID: d2themescatalog.EvergladeGreen.ID,
			Sketch:  true,
		},
		highlighting.NewHighlighting(
			// https://xyproto.github.io/splash/docs/
			// https://xyproto.github.io/splash/docs/all.html#monokailight
			highlighting.WithStyle("monokailight"),
			highlighting.WithFormatOptions(
				cfh.WithLineNumbers(true),
			),
		),
	}
	ctx = parser.NewContext(parser.WithIDs(&HTMLHeaderID{
		values: map[string]bool{},
	}))

	parserOptions = []parser.Option{
		parser.WithAutoHeadingID(),
	}
	rendererOptions = []renderer.Option{
		html.WithXHTML(),
		html.WithUnsafe(),
	}
)

type HTMLParser struct {
	Options

	Favicon bool
	CSS     bool

	FaviconHref   template.HTML
	MathJaxConfig template.HTML
	MathJaxTeXSVG template.HTML
	ConvertedCSS  template.HTML
	ConvertedHTML template.HTML
	GeneratedAt   template.HTML

	begin   time.Time
	Elapsed time.Duration
}

func (p *HTMLParser) parserMarkdown(files []string) (htmlContent string) {
	fmt.Println("⌚  Converting Markdown to HTML ...")
	p.begin = time.Now()

	mdParser := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRendererOptions(rendererOptions...),
	)

	if len(p.OutputFile) > 0 {
		ext := regexp.QuoteMeta(filepath.Ext(p.OutputFile))
		re := regexp.MustCompile(ext + "$")
		if strings.TrimSpace(p.HTMLTitle) == "" {
			p.HTMLTitle = filepath.Base(re.ReplaceAllString(p.OutputFile, ""))
		}
		htmlStr, err := p.renderHTMLConcat(files, mdParser)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := p.writeHTML(htmlStr); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	} else {
		for _, file := range files {
			htmlStr, err := p.renderHTML(file, mdParser)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if strings.TrimSpace(p.HTMLTitle) == "" {
				p.HTMLTitle = file
			}
			p.OutputFile = file + ".html"
			if err := p.writeHTML(htmlStr); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			p.begin = time.Now()
		}
	}

	bold, reset, highlight := "\033[39;1m", "\033[0m", "\033[32m"
	fmt.Printf(`✨  Convert of "%s%s%s" is completed! (%s%v%s)`, highlight, p.OutputFile, reset, bold, p.Elapsed, reset)
	return
}

func (p *HTMLParser) renderHTML(input string, markdown goldmark.Markdown) (htmlStr string, err error) {
	var fi *os.File
	if fi, err = os.Open(input); err != nil {
		return
	}
	defer func(fi *os.File) {
		_ = fi.Close()
	}(fi)

	var md []byte
	if md, err = io.ReadAll(fi); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = markdown.Convert(md, &buf, parser.WithContext(ctx)); err != nil {
		return
	}

	if htmlStr, err = parseImageOpt(buf.String()); err != nil {
		return
	}

	if p.EmbedImage {
		if htmlStr, err = embedImage(htmlStr, filepath.Dir(input)); err != nil {
			return
		}
	}
	return
}

func (p *HTMLParser) renderHTMLConcat(inputs []string, markdown goldmark.Markdown) (htmlStr string, err error) {
	for _, input := range inputs {
		var h string
		if h, err = p.renderHTML(input, markdown); err != nil {
			return
		}
		htmlStr += h
	}
	return
}

func (p *HTMLParser) writeHTML(html string) (err error) {
	fmt.Println(" ⏳  Converting", p.OutputFile, "...")
	if strings.TrimSpace(p.HTMLLang) == "" {
		p.HTMLLang = "en"
	}

	if strings.TrimSpace(p.HTMLFavicon) != "" {
		p.Favicon = true
		if p.EmbedImage {
			if f, e := os.Stat(p.HTMLFavicon); !os.IsNotExist(e) && !f.IsDir() {
				cwd, _ := os.Getwd()
				if src, err := decodeBase64(p.HTMLFavicon, cwd); err == nil {
					mimeType := getImageMime(p.HTMLFavicon)
					p.HTMLFavicon = fmt.Sprintf("data:%s;base64,%s", mimeType, src)

				}
			}
		}
		p.HTMLFavicon = fmt.Sprintf(`<link rel="shortcut icon" type="image/x-icon" href="%s">`, p.HTMLFavicon)
		p.FaviconHref = template.HTML(p.HTMLFavicon)
	}

	if p.MathJax {
		html, err = replaceMathJaxCodeBlock(html)
		if err != nil {
			return err
		}
		p.MathJaxConfig = template.HTML(fmt.Sprintf(`<script type="text/x-mathjax-config">%s</script>`, assets.EmbedMathJaxConfig))
		p.MathJaxTeXSVG = template.HTML(fmt.Sprintf(`<script type="text/javascript">%s</script>`, assets.EmbedMathJaxTeXSVG))
	}

	if html, err = replaceCheckBox(html); err != nil {
		return
	}

	if p.TableSpan {
		if html, err = replaceTableSpan(html); err != nil {
			return
		}
	}

	if len(p.CustomCSS) > 0 {
		var fi *os.File
		if fi, err = os.Open(p.CustomCSS); err != nil {
			return
		}
		defer func(fi *os.File) {
			_ = fi.Close()
		}(fi)

		var c []byte
		if c, err = io.ReadAll(fi); err != nil {
			return
		}
		p.CSS = true
		p.CustomCSS = fmt.Sprintf(`<style type="text/css">%s%s</style>`, "\n", c)
		p.ConvertedCSS = template.HTML(p.CustomCSS)
	}
	p.ConvertedHTML = template.HTML(html)

	theme := p.Theme
	if strings.TrimSpace(theme) == "" {
		theme = "vue"
	}

	var themeTmpl []byte
	themePath := fmt.Sprintf("theme/%s/%s.gohtml", theme, theme)
	if themeTmpl, err = EmbedThemes.ReadFile(themePath); err != nil {
		return
	}
	var tmpl *template.Template
	if tmpl, err = template.New(theme).Funcs(template.FuncMap{
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
		"safeCSS": func(css string) template.CSS {
			return template.CSS(css)
		},
		"safeJS": func(js string) template.JS {
			return template.JS(js)
		},
	}).Parse(string(themeTmpl)); err != nil {
		return
	}

	dir := filepath.Dir(p.OutputFile)
	if f, e := os.Stat(dir); os.IsNotExist(e) || !f.IsDir() {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
	var fileOut *os.File
	if fileOut, err = os.Create(p.OutputFile); err != nil {
		return
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(fileOut)

	elapsed := time.Since(p.begin) // 当前耗时
	if p.Generated {
		now := time.Now().Format(time.RFC3339)
		comment := fmt.Sprintf(`<!-- Generated by github.com/GoAid/md-html-cli on %s (%v) -->`, now, elapsed)
		p.GeneratedAt = template.HTML(comment)
	}
	p.Elapsed += elapsed // 总耗时

	err = tmpl.Execute(fileOut, p)
	return
}

type HTMLHeaderID struct {
	values map[string]bool
}

func (s *HTMLHeaderID) Generate(value []byte, _ ast.NodeKind) []byte {
	id := strings.ReplaceAll(strings.ToLower(string(value)), " ", "-")
	id = url.PathEscape(id)
	if s.values[id] {
		// avoid duplicate id
		id += "_"
	}
	idBytes := []byte(id)
	s.Put(idBytes)

	return idBytes
}

func (s *HTMLHeaderID) Put(value []byte) {
	s.values[util.BytesToReadOnlyString(value)] = true
}

func decodeBase64(src, parent string) (string, error) {
	path := src
	if !filepath.IsAbs(path) {
		path = filepath.Join(parent, path)
	}
	f, err := os.Open(path)
	if err != nil {
		pathErr := err.(*os.PathError)
		errno := pathErr.Err.(syscall.Errno)
		if errno != 0x7B { // suppress ERROR_INVALID_NAME
			_, _ = fmt.Fprintln(os.Stderr, err)
			return "", nil
		}
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	var d []byte
	if d, err = io.ReadAll(f); err != nil {
		return "", err
	}

	dest := base64.StdEncoding.EncodeToString(d)
	return dest, nil
}

func getImageMime(src string) string {
	ext := filepath.Ext(src)
	mimeType := mime.TypeByExtension(ext)
	if len(mimeType) <= 0 {
		mimeType = "image"
	}
	return mimeType
}

func embedImage(src, parent string) (string, error) {
	dest := src

	reFind := regexp.MustCompile(`(<img[\S\s]+?src=")([\S\s]+?)("[\S\s]*?/?>)`)
	reUrl := regexp.MustCompile(`(?i)^https?://.*`)

	imgTags := reFind.FindAllString(src, -1)
	for _, t := range imgTags {
		imgSrc := reFind.ReplaceAllString(t, "$2")

		if reUrl.MatchString(imgSrc) {
			continue
		}
		b64img, err := decodeBase64(imgSrc, parent)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			continue
		}

		reReplace, err := regexp.Compile(`(<img[\S\s]+?src=")` + regexp.QuoteMeta(imgSrc) + `("[\S\s]*?/?>)`)
		if err != nil {
			return src, err
		}

		mimeType := getImageMime(imgSrc)
		dest = reReplace.ReplaceAllString(dest, "${1}data:"+mimeType+";base64,"+b64img+"${2}")
	}
	return dest, nil
}

func parseImageOpt(src string) (string, error) {
	dest := src

	re := regexp.MustCompile(`(<img[\S\s]+?src=)"([\S\s]+?)\?(\S+?)"([\S\s]*?/?>)`)
	dest = re.ReplaceAllStringFunc(dest, func(s string) string {
		imgTag := re.FindStringSubmatch(s)
		srcQuery := strings.Join(strings.Split(imgTag[3], "&amp;"), " ")
		return fmt.Sprintf(`%s"%s" %s%s`, imgTag[1], imgTag[2], srcQuery, imgTag[4])
	})
	return dest, nil
}

func replaceMathJaxCodeBlock(src string) (string, error) {
	srcReader := strings.NewReader(src)
	doc, err := goquery.NewDocumentFromReader(srcReader)
	if err != nil {
		return src, err
	}

	code := doc.Find("pre>code.language-math")
	code.Each(func(index int, s *goquery.Selection) {
		s.Parent().ReplaceWithHtml("<p>$$" + s.Text() + "$$</p>")
	})

	return doc.Find("body").Html()
}

func replaceCheckBox(src string) (string, error) {
	sr := strings.NewReader(src)
	doc, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		return "", err
	}

	doc.Find("li").Each(func(i int, li *goquery.Selection) {
		li.Contents().Each(func(j int, c *goquery.Selection) {
			if goquery.NodeName(c) == "#text" {
				li.Find("input").Each(func(k int, input *goquery.Selection) {
					if t, exist := input.Attr("type"); exist && t == "checkbox" {
						li.AddClass("task-list-item")
					}
				})
			}
		})
	})

	return doc.Find("body").Html()
}

func replaceTableSpan(src string) (string, error) {
	sr := strings.NewReader(src)
	doc, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("\u00a6\\s*")

	doc.Find("table").Each(func(i int, tbl *goquery.Selection) {
		tbl.Find("tbody").Each(func(j int, tbody *goquery.Selection) {
			trs := tbody.Find("tr")
			// colspan
			colMax := 0
			trs.Each(func(k int, tr *goquery.Selection) {
				tds := tr.Find("td")
				colMns := tds.Length()
				if colMns > colMax {
					colMax = colMns
				}
				col := 0
				tds.Each(func(l int, td *goquery.Selection) {
					col++
					td.Contents().Each(func(m int, c *goquery.Selection) {
						cnt := len(re.FindAllIndex([]byte(c.Text()), -1))
						if cnt > 0 {
							td.SetAttr("colspan", strconv.Itoa(cnt+1))
							c.ReplaceWithHtml(re.ReplaceAllString(c.Text(), ""))
							col += cnt
						}
					})
					if col > colMns {
						td.SetAttr("hidden", "")
					}
				})
			})
			// rowspan
			for m := 0; m < colMax; m++ {
				var root *goquery.Selection
				cnt := 0
				trs.Each(func(k int, tr *goquery.Selection) {
					tr.Find("td").Each(func(l int, td *goquery.Selection) {
						if l == m {
							atd := getActualTD(tr, l)
							if k == 0 {
								root = atd
							} else {
								if atd.Text() != "" {
									cnt = 0
									root = atd
								} else {
									cnt++
									root.SetAttr("rowspan", strconv.Itoa(cnt+1))
									atd.SetAttr("hidden", "")
								}
							}
						}
					})
				})
			}
			// remove hidden <td>
			tbody.Find("tr>td").Each(func(i int, td *goquery.Selection) {
				if _, hidden := td.Attr("hidden"); hidden {
					td.Remove()
				}
			})
		})
		// remove empty header
		empty := true
		tbl.Find("thead").Each(func(i int, thead *goquery.Selection) {
			thead.Find("tr>th").EachWithBreak(func(j int, th *goquery.Selection) bool {
				if th.Text() != "" {
					empty = false
					return false
				}
				return true
			})
			if empty {
				thead.Remove()
			}
		})
	})

	return doc.Find("body").Html()
}

func getActualTD(tr *goquery.Selection, index int) *goquery.Selection {
	pos := 0
	var result *goquery.Selection
	tr.Find("td").EachWithBreak(func(i int, td *goquery.Selection) bool {
		cs, _ := strconv.Atoi(td.AttrOr("colspan", "1"))
		pos += cs
		if pos >= index+1 {
			result = td
			return false
		}
		return true
	})

	return result
}
