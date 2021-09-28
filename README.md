# cassiopaea-tools
This repository contains tools that can be used to scrape and process [Cassiopaea Session Transcripts](https://cassiopaea.org/).

## Example usage
Be sure, that you have [Go](https://golang.org/) installed, and `$GOPATH/bin` available in your `$PATH`.

Transcripts can be scraped and saved to the local directory named `sessions` by executing the following in the terminal:
```
go get github.com/liberty239/cassiopaea-tools/./...
cass-src fetch-all ./sessions
```
