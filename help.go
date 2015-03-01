package main

var appHelpTemplate = `
  Usage: {{.Name}} [options] URL

  Options:

    {{range .Flags}}{{.}}
    {{end}}
`
