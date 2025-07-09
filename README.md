# Sheet parser and downloader from some score muse site

Combines sheets into one output PDF file

Usage:

Run `go build`

Run `./sheet-parser` with `-url` argument

Works with single url and array of urls (JSON)

Examples:

```bash
./sheet-parser -url=https://example.com/some/sheeet/1231
./sheet-parser -url='["https://example.com/some/sheeet/1231","https://example.com/some/sheeet/124343"]'
```

Output files will be in related categories

