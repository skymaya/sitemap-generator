package main

import (
  "flag"
  "bufio"
  "os"
  "time"
  "fmt"
  "html"
  "net/url"
  "regexp"
  "github.com/gocolly/colly/v2"
)

func convert_lastmod(lmh string)string {
  lm, err := time.Parse(time.RFC1123, lmh)
  if err != nil {
      return time.Now().Format("2000-01-01")
  }
  return lm.Format("2000-01-01")
}

func convert_url(ul string)string {
  u, err := url.Parse(ul)
  if err != nil {
    panic(err)
  }
  return html.EscapeString(u.String())
}

func main() {

  address := flag.String("a", "http://www.example.com", "URL to scrape, must begin with http:// or https://")
  filename := flag.String("f", "sitemap.xml", "Full path to the sitemap file name, i.e. /home/user/sitemap.xml")
  flag.Parse()

  parsed_url, err := url.Parse(*address)
  if err != nil {
      panic(err)
  }

  file, err := os.OpenFile(*filename, os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
      panic(err)
  }

  datawriter := bufio.NewWriter(file)

  datawriter.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
  datawriter.WriteString("<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")

  c := colly.NewCollector(
    colly.UserAgent("Sitemap Generator"),
    colly.AllowedDomains(parsed_url.Host),
    colly.URLFilters(
      regexp.MustCompile("(http|https)://.+$"),
    ),
  )

  c.OnHTML("a[href]", func(e *colly.HTMLElement) {
    link := e.Attr("href")
    c.Visit(e.Request.AbsoluteURL(link))
  })

  c.OnResponseHeaders(func(r *colly.Response) {
    fmt.Println("Adding", r.Request.URL.String())

    lastmod_header := (*r.Headers)["Last-Modified"][0]
    lastmod := convert_lastmod(lastmod_header)

    esc_url := convert_url(r.Request.URL.String())

    datawriter.WriteString(" <url>\n")
    datawriter.WriteString("  <loc>"+esc_url+"</loc>\n")
    datawriter.WriteString("  <lastmod>"+lastmod+"</lastmod>\n")
    datawriter.WriteString(" </url>\n")
  })

  c.Visit(*address)

  datawriter.WriteString("</urlset>")

  datawriter.Flush()
  file.Close()
}
