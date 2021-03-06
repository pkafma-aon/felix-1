package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mattn/go-runewidth"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func isChildOf(node *html.Node, name string) bool {
	node = node.Parent
	return node != nil && node.Type == html.ElementNode && strings.ToLower(node.Data) == name
}

func hasClass(node *html.Node, clazz string) bool {
	for _, attr := range node.Attr {
		if attr.Key == "class" {
			for _, c := range strings.Fields(attr.Val) {
				if c == clazz {
					return true
				}
			}
		}
	}
	return false
}

func attr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func br(node *html.Node, w io.Writer, option *Option) {
	node = node.PrevSibling
	if node == nil {
		return
	}
	switch node.Type {
	case html.TextNode:
		text := strings.Trim(node.Data, " \t")
		if text != "" && !strings.HasSuffix(text, "\n") {
			fmt.Fprint(w, "\n")
		}
	case html.ElementNode:
		switch strings.ToLower(node.Data) {
		case "br", "p", "ul", "ol", "div", "blockquote", "h1", "h2", "h3", "h4", "h5", "h6":
			fmt.Fprint(w, "\n")
		}
	}
}

func table(node *html.Node, w io.Writer, option *Option) {
	for tr := node.FirstChild; tr != nil; tr = tr.NextSibling {
		if tr.Type == html.ElementNode && strings.ToLower(tr.Data) == "tbody" {
			node = tr
			break
		}
	}
	var header bool
	var rows [][]string
	for tr := node.FirstChild; tr != nil; tr = tr.NextSibling {
		if tr.Type != html.ElementNode || strings.ToLower(tr.Data) != "tr" {
			continue
		}
		var cols []string
		if !header {
			for th := tr.FirstChild; th != nil; th = th.NextSibling {
				if th.Type != html.ElementNode || strings.ToLower(th.Data) != "th" {
					continue
				}
				var buf bytes.Buffer
				walk(th, &buf, 0, option)
				cols = append(cols, buf.String())
			}
			if len(cols) > 0 {
				rows = append(rows, cols)
				header = true
				continue
			}
		}
		for td := tr.FirstChild; td != nil; td = td.NextSibling {
			if td.Type != html.ElementNode || strings.ToLower(td.Data) != "td" {
				continue
			}
			var buf bytes.Buffer
			walk(td, &buf, 0, option)
			cols = append(cols, buf.String())
		}
		rows = append(rows, cols)
	}
	maxcol := 0
	for _, cols := range rows {
		if len(cols) > maxcol {
			maxcol = len(cols)
		}
	}
	widths := make([]int, maxcol)
	for _, cols := range rows {
		for i := 0; i < maxcol; i++ {
			if i < len(cols) {
				width := runewidth.StringWidth(cols[i])
				if widths[i] < width {
					widths[i] = width
				}
			}
		}
	}
	for i, cols := range rows {
		for j := 0; j < maxcol; j++ {
			fmt.Fprint(w, "|")
			if j < len(cols) {
				width := runewidth.StringWidth(cols[j])
				fmt.Fprint(w, cols[j])
				fmt.Fprint(w, strings.Repeat(" ", widths[j]-width))
			} else {
				fmt.Fprint(w, strings.Repeat(" ", widths[j]))
			}
		}
		fmt.Fprint(w, "|\n")
		if i == 0 && header {
			for j := 0; j < maxcol; j++ {
				fmt.Fprint(w, "|")
				fmt.Fprint(w, strings.Repeat("-", widths[j]))
			}
			fmt.Fprint(w, "|\n")
		}
	}
	fmt.Fprint(w, "\n")
}

var emptyElements = []string{
	"area",
	"base",
	"br",
	"col",
	"embed",
	"hr",
	"img",
	"input",
	"keygen",
	"link",
	"meta",
	"param",
	"source",
	"track",
	"wbr",
}

func raw(node *html.Node, w io.Writer, option *Option) {
	switch node.Type {
	case html.ElementNode:
		fmt.Fprintf(w, "<%s", node.Data)
		for _, attr := range node.Attr {
			fmt.Fprintf(w, " %s=%q", attr.Key, attr.Val)
		}
		found := false
		tag := strings.ToLower(node.Data)
		for _, e := range emptyElements {
			if e == tag {
				found = true
				break
			}
		}
		if found {
			fmt.Fprint(w, "/>")
		} else {
			fmt.Fprint(w, ">")
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				raw(c, w, option)
			}
			fmt.Fprintf(w, "</%s>", node.Data)
		}
	case html.TextNode:
		fmt.Fprint(w, node.Data)
	}
}

func bq(node *html.Node, w io.Writer, option *Option) {
	if node.Type == html.TextNode {
		fmt.Fprint(w, strings.Replace(node.Data, "\u00a0", " ", -1))
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			bq(c, w, option)
		}
	}
}

func pre(node *html.Node, w io.Writer, option *Option) {
	if node.Type == html.TextNode {
		fmt.Fprint(w, node.Data)
	} else {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			pre(c, w, option)
		}
	}
}

func walk(node *html.Node, w io.Writer, nest int, option *Option) {
	if node.Type == html.TextNode {
		if strings.TrimSpace(node.Data) != "" {
			text := regexp.MustCompile(`[[:space:]][[:space:]]*`).ReplaceAllString(strings.Trim(node.Data, "\t\r\n"), " ")
			fmt.Fprint(w, text)
		}
	}
	n := 0
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.CommentNode:
			fmt.Fprint(w, "<!--")
			fmt.Fprint(w, c.Data)
			fmt.Fprint(w, "-->\n")
		case html.ElementNode:
			switch strings.ToLower(c.Data) {
			case "a":
				fmt.Fprint(w, "[")
				walk(c, w, nest, option)
				fmt.Fprint(w, "]("+attr(c, "href")+")")
			case "b", "strong":
				fmt.Fprint(w, "**")
				walk(c, w, nest, option)
				fmt.Fprint(w, "**")
			case "i", "em":
				fmt.Fprint(w, "_")
				walk(c, w, nest, option)
				fmt.Fprint(w, "_")
			case "del":
				fmt.Fprint(w, "~~")
				walk(c, w, nest, option)
				fmt.Fprint(w, "~~")
			case "br":
				br(c, w, option)
				fmt.Fprint(w, "\n\n")
			case "p":
				br(c, w, option)
				walk(c, w, nest, option)
				br(c, w, option)
				fmt.Fprint(w, "\n\n")
			case "code":
				if !isChildOf(c, "pre") {
					fmt.Fprint(w, "`")
					pre(c, w, option)
					fmt.Fprint(w, "`")
				}
			case "pre":
				br(c, w, option)
				var buf bytes.Buffer
				pre(c, &buf, option)
				var lang string
				if option != nil && option.PreFunc != nil {
					if guess, err := option.PreFunc(buf.String()); err == nil {
						lang = guess
					}
				}
				fmt.Fprint(w, "\n\n")
				fmt.Fprint(w, "```"+lang+"\n")
				fmt.Fprint(w, buf.String())
				if !strings.HasSuffix(buf.String(), "\n") {
					fmt.Fprint(w, "\n")
				}
				fmt.Fprint(w, "```\n\n")
			case "div":
				br(c, w, option)
				walk(c, w, nest, option)
				fmt.Fprint(w, "\n")
			case "blockquote":
				br(c, w, option)
				var buf bytes.Buffer
				if hasClass(c, "code") {
					bq(c, &buf, option)
					var lang string
					if option != nil && option.PreFunc != nil {
						if guess, err := option.PreFunc(buf.String()); err == nil {
							lang = guess
						}
					}
					fmt.Fprint(w, "```"+lang+"\n")
					fmt.Fprint(w, strings.TrimLeft(buf.String(), "\n"))
					if !strings.HasSuffix(buf.String(), "\n") {
						fmt.Fprint(w, "\n")
					}
					fmt.Fprint(w, "```\n\n")
				} else {
					walk(c, &buf, nest+1, option)

					if lines := strings.Split(strings.TrimSpace(buf.String()), "\n"); len(lines) > 0 {
						for _, l := range lines {
							fmt.Fprint(w, "> "+strings.TrimSpace(l)+"\n")
						}
						fmt.Fprint(w, "\n")
					}
				}
			case "ul", "ol":
				br(c, w, option)
				var buf bytes.Buffer
				walk(c, &buf, 1, option)
				if lines := strings.Split(strings.TrimSpace(buf.String()), "\n"); len(lines) > 0 {
					for i, l := range lines {
						if i > 0 || nest > 0 {
							fmt.Fprint(w, "\n")
						}
						fmt.Fprint(w, strings.Repeat("    ", nest)+strings.TrimSpace(l))
					}
					fmt.Fprint(w, "\n")
				}
			case "li":
				br(c, w, option)
				if isChildOf(c, "ul") {
					fmt.Fprint(w, "* ")
				} else if isChildOf(c, "ol") {
					n++
					fmt.Fprint(w, fmt.Sprintf("%d. ", n))
				}
				walk(c, w, nest, option)
				fmt.Fprint(w, "\n")
			case "h1", "h2", "h3", "h4", "h5", "h6":
				br(c, w, option)
				fmt.Fprint(w, "\n")
				fmt.Fprint(w, strings.Repeat("#", int(rune(c.Data[1])-rune('0')))+" ")
				walk(c, w, nest, option)
				fmt.Fprint(w, "\n\n")
			case "img":
				alt := attr(c, "alt")
				if alt == "" {
					alt = attr(c, "title")
				}
				src := attr(c, "src")
				dataSrc := attr(c, "data-src")
				if dataSrc == "" {
					dataSrc = attr(c, "data-original-src")
				}
				if option != nil && option.PreFunc != nil {
					if iAlt, iSrc, err := option.ImgFunc(alt, src, dataSrc); err == nil {
						fmt.Fprint(w, "!["+iAlt+"]("+iSrc+")")
					} else {
						log.Println("option.ImgFunc failed with err: ", err)
					}
				} else {
					fmt.Fprint(w, "!["+alt+"]("+src+")")
				}
			case "hr":
				br(c, w, option)
				fmt.Fprint(w, "\n---\n\n")
			case "table":
				br(c, w, option)
				table(c, w, option)
			case "style":
				if option != nil && option.Style {
					br(c, w, option)
					raw(c, w, option)
					fmt.Fprint(w, "\n\n")
				}
			case "script":
				if option != nil && option.Script {
					br(c, w, option)
					raw(c, w, option)
					fmt.Fprint(w, "\n\n")
				}
			default:
				walk(c, w, nest, option)
			}
		default:
			walk(c, w, nest, option)
		}
	}
}

type Option struct {
	PreFunc func(string) (string, error)
	ImgFunc func(alt, src, dataSrc string) (iAlt, iSrc string, err error)
	BaseUrl string
	Script  bool
	Style   bool
}

// convert convert HTML to Markdown. Read HTML from r and write to w.
func convert(w io.Writer, r io.Reader, option *Option) error {

	doc, err := html.Parse(r)
	if err != nil {
		return err
	}
	walk(doc, w, 0, option)
	fmt.Fprint(w, "\n")
	return nil
}

func ParseUrlPage(href, div, jekyllDir string) error {
	urlObj, err := url.Parse(href)
	if err != nil {
		return err
	}
	baseUrl := urlObj.Scheme + "://" + urlObj.Host
	//log.Println(baseUrl)
	client := &http.Client{}
	mkdFileName := time.Now().Format("2006-01-02") + urlObj.Path + ".md"
	mkdFileName = strings.Replace(mkdFileName, "/", "-", -1)
	mkdFilePath := filepath.Join(jekyllDir, "_posts", "golang", mkdFileName)
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("the get request's response code is not 200")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	title := doc.Find("h1.h1").Text()
	title = strings.TrimSpace(title)
	divHtml, err := doc.Find(div).Html()
	if err != nil {
		return err
	}
	r := bytes.NewBuffer([]byte(divHtml))
	err = os.MkdirAll(filepath.Dir(mkdFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(mkdFilePath)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, headerFormat, title, title, href)
	opt := &Option{
		PreFunc: func(s2 string) (s string, e error) {
			return "go", nil
		},
		ImgFunc: func(alt, src, dataSrc string) (iAlt, iSrc string, err error) {
			iAlt = "tech.mojotv.cn_" + alt
			//???????????????markdown ?????????
			iSrc = "/assets/pic/" + base64.RawURLEncoding.EncodeToString([]byte("sgmf_"+dataSrc)) + ".jpg"
			err = downImg(baseUrl, dataSrc, jekyllDir, iSrc)
			if err != nil {
				log.Println("??????????????????,", dataSrc, iSrc)
				return
			}
			return
		},
	}
	return convert(f, r, opt)
}
func ParseUrlPageJianshu(href, div, jekyllDir string) error {
	urlObj, err := url.Parse(href)
	if err != nil {
		return err
	}
	baseUrl := urlObj.Scheme + "://" + urlObj.Host
	//log.Println(baseUrl)
	client := &http.Client{}
	mkdFileName := time.Now().Format("2006-01-02") + urlObj.Path + ".md"
	mkdFileName = strings.Replace(mkdFileName, "/", "-", -1)
	mkdFilePath := filepath.Join(jekyllDir, "_posts", "golang", mkdFileName)
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("the get request's response code is not 200")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	title := doc.Find("h1.h1").Text()
	title = strings.TrimSpace(title)
	divHtml, err := doc.Find(div).Html()
	if err != nil {
		return err
	}
	r := bytes.NewBuffer([]byte(divHtml))
	err = os.MkdirAll(filepath.Dir(mkdFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(mkdFilePath)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, headerFormat, title, title, href)
	opt := &Option{
		PreFunc: func(s2 string) (s string, e error) {
			return "go", nil
		},
		ImgFunc: func(alt, src, dataSrc string) (iAlt, iSrc string, err error) {
			iAlt = "tech.mojotv.cn_" + alt
			//???????????????markdown ?????????
			iSrc = "/assets/pic/" + base64.RawURLEncoding.EncodeToString([]byte("sgmf_"+dataSrc)) + ".jpg"
			err = downImg(baseUrl, dataSrc, jekyllDir, iSrc)
			if err != nil {
				log.Println("??????????????????,", dataSrc, iSrc)
				return
			}
			return
		},
	}
	return convert(f, r, opt)
}
func ParseUrlPageLibraGen(href, div, jekyllDir string) error {
	urlObj, err := url.Parse(href)
	if err != nil {
		return err
	}
	baseUrl := urlObj.Scheme + "://" + urlObj.Host
	//log.Println(baseUrl)
	client := &http.Client{}
	mkdFileName := time.Now().Format("2006-01-02") + urlObj.Path + ".md"
	mkdFileName = strings.Replace(mkdFileName, "/", "-", -1)
	mkdFilePath := filepath.Join(jekyllDir, "_posts", "wiki", mkdFileName)
	req, err := http.NewRequest("GET", href, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("the get request's response code is not 200")
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	title := doc.Find("h1.h1").Text()
	title = strings.TrimSpace(title)
	divHtml, err := doc.Find(div).Html()
	if err != nil {
		return err
	}
	r := bytes.NewBuffer([]byte(divHtml))
	err = os.MkdirAll(filepath.Dir(mkdFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(mkdFilePath)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, headerFormat, title, title, href)
	opt := &Option{
		PreFunc: func(s2 string) (s string, e error) {
			return "go", nil
		},
		ImgFunc: func(alt, src, dataSrc string) (iAlt, iSrc string, err error) {
			iAlt = "tech.mojotv.cn_" + alt
			//???????????????markdown ?????????
			iSrc = "/assets/pic/" + base64.RawURLEncoding.EncodeToString([]byte("sgmf_"+dataSrc)) + ".jpg"
			err = downImg(baseUrl, dataSrc, jekyllDir, iSrc)
			if err != nil {
				log.Println("??????????????????,", dataSrc, iSrc)
				return
			}
			return
		},
	}
	return convert(f, r, opt)
}

const headerFormat = `---
layout: post
title: %s
category: Golang
tags: Golang
keywords: go??????
description: %s
coverage: ginbro_coverage.jpg
ref: %s
---

`

func downImg(base, uri, imgDir, imgName string) error {
	parts := strings.Split(uri, "?")
	keyName := parts[0]
	filePath := filepath.Join(imgDir, imgName)
	fullUrl := uri
	if !strings.HasPrefix(uri, "http") {
		fullUrl = fmt.Sprintf("%s/%s", base, keyName)
	}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("the get request's response code is " + res.Status)
	}
	defer res.Body.Close()
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, res.Body)
	return err
}
