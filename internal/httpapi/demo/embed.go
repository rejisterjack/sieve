package demo

import "embed"

//go:embed index.html corpus/*
var FS embed.FS

//go:embed corpus/*
var corpusFS embed.FS
