# `mkvparse`: Matroska parser in Go

Fast [Matroska](https://www.matroska.org) (`.mkv`, `.mka`) parser, written in Go.

Features:

- Supports [all Matroska tags](https://www.matroska.org/technical/specs/index.html)
- Supports short-circuiting the parser, making it possible to 
read specific data (e.g. title, author) without reading the
entire file (see the `mkvtags` example)
- Also works with [WebM](https://www.webmproject.org) (`.webm`) files
- Event-based push API
- No dependencies

## API

See the [API Reference](https://godoc.org/github.com/pindrop/mkvparse).

## Examples

Besides the examples in the [API Reference](https://godoc.org/github.com/pindrop/mkvparse),
there are some larger examples in the `examples/` dir:

- `examples/mkvinfo`: Example using basic parser API to print all elements
- `examples/mkvtags`: Example using the optimized 'sections' API to print all tags without
	parsing the entire file.
