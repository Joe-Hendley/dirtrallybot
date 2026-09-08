# Vendored static assets

`htmx.min.js` — [htmx](https://htmx.org) v2.0.4, released under the Zero-Clause
BSD licence. It is embedded into the binary (`//go:embed`) and served by the
challenge viewer so the page makes no external requests.

To update, replace the file with a new pinned release:

```
curl -sSL -o htmx.min.js https://unpkg.com/htmx.org@<version>/dist/htmx.min.js
```

and bump the version assertion in `web_test.go`.
