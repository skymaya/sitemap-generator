# Sitemap Generator For Google

This is a quick and dirty sitemap generator written in Go. The sitemap output should be fully compliant with Google's specifications found here: https://support.google.com/webmasters/answer/183668?hl=en

Google does not endorse this application.

### Prerequisites

```
Golang + golang environment configured
```

### Installing

Clone the repo somewhere:

```
git clone git@github.com:skymaya/sitemap-generator.git
```

Build the binary:

```
go build sitemap.go
```

Run it with:

```
./sitemap -a http://example.com -f sitemap.xml
```
