package main

import "embed"

//go:embed theme
var EmbedThemes embed.FS

var (
	ThemeVue  = "vue"
	ThemeSide = "side"
)
