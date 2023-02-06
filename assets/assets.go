package assets

import _ "embed"

var (
	//go:embed mathjax/mathjax-config.min.js
	EmbedMathJaxConfig string
	//go:embed mathjax/MathJax-TeXSVG.min.js
	EmbedMathJaxTeXSVG string
)
