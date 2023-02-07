<p align="center">
<img alt="md-html-cli" src="assets/image/logo.png">
</p>
<h1 align="center">md-html-cli</h1>
<p align="center">
用于将 markdown 转换为单个 html 文件的 CLI 实用工具。
<br>
<b>🇨🇳 中文</b> | <a href="README.md">🇺🇸 English</a>
</p>

## 安装

安装 Go 语言 `1.17` 或以上版本，然后执行以下命令：

```shell
go install github.com/GoAid/md-html-cli@latest
```

## 用例

```shell
go run github.com/GoAid/md-html-cli@latest /?
```

```shell
Usage:
  go run github.com/GoAid/md-html-cli@latest [OPTIONS]

Application Options:
  /i, /input:            input Markdown
  /o, /output:           output HTML
  /l, /lang:             html lang attribute value, default is en
  /t, /title:            custom html title, default is output file name
  /f, /favicon:          favicon image path, if embed is used, will embed by base64 encoding
  /e, /embed             embed image by base64 encoding
  /m, /mathjax           use MathJax
  /s, /span              enable table row/col span
  /c, /css:              custom css file path
      /theme:[vue|side]  output HTML theme
      /toc               generate TOC

Help Options:
  /?                     Show this help message
  /h, /help              Show this help message
```

### 样品

[GitHub Pages](https://GoAid.github.io/md-html-cli/index.html)

此 html 页面由以下命令生成：

```bash
md-html-cli -i "example/*.md" -o gh-pages/index.html -l en -t "Example Page" -f example/img/go.png -ems -c example/css/custom-css.css --theme vue --toc
```

### 示例

<details>
<summary>预览</summary>

| Markdown                                                                            | HTML                                                                                    |
|-------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------|
| ![mh-highlight-md.png](assets/image/docs/mh-highlight-md.png)                       | ![mh-highlight-html.png](assets/image/docs/mh-highlight-html.png)                       |
| ![mh-image-md.png](assets/image/docs/mh-image-md.png)                               | ![mh-image-html.png](assets/image/docs/mh-image-html.png)                               |
| ![mh-image-size-md.png](assets/image/docs/mh-image-size-md.png)                     | ![mh-image-size-html.png](assets/image/docs/mh-image-size-html.png)                     |
| ![mh-link-md.png](assets/image/docs/mh-link-md.png)                                 | ![mh-link-html.png](assets/image/docs/mh-link-html.png)                                 |
| ![mh-mathjax-md.png](assets/image/docs/mh-mathjax-md.png)                           | ![mh-mathjax-html.png](assets/image/docs/mh-mathjax-html.png)                           |
| ![mh-table-span-md.png](assets/image/docs/mh-table-span-md.png)                     | ![mh-table-span-html.png](assets/image/docs/mh-table-span-html.png)                     |
| ![mh-table-without-header-md.png](assets/image/docs/mh-table-without-header-md.png) | ![mh-table-without-header-html.png](assets/image/docs/mh-table-without-header-html.png) |
| ![mh-task-list-md.png](assets/image/docs/mh-task-list-md.png)                       | ![mh-task-list-html.png](assets/image/docs/mh-task-list-html.png)                       |

</details>

## 开发

### 主题模板

在 `theme` 文件夹中创建用于存放主题模板文件的文件夹，
并在文件夹中创建同名的 `.gohtml` 类型模板文件。

### 模板变量

| 变量                           | 说明                       |
|------------------------------|--------------------------|
| `{{ .HTMLLang }}`            | HTML 语言属性值，如 `en`、`zh` 等 |
| `{{ .HTMLTitle }}`           | HTML 自定义标题，默认为输出文件名称     |
| `{{ if .Favicon }}{{ end }}` | 是否添加 `favicon.ico`       |
| `{{ .FaviconHref }}`         | `favicon.ico` 标签元素       |
| `{{ if .TOC }}{{ end }}`     | 是否生成目录                   |
| `{{ if .CSS }}{{ end }}`     | 是否添加自定义样式                |
| `{{ .ConvertedCSS }}`        | 自定义样式标签和内容               |
| `{{ if .MathJax }}{{ end }}` | 是否使用 `MathJax` 渲染数学公式    |
| `{{ .MathJaxConfig }}`       | `MathJax` 配置文件 JS 标签元素   |
| `{{ .MathJaxTeXSVG }}`       | `MathJax` 渲染工具 JS 标签元素   |
| `{{ .ConvertedHTML }}`       | 转换后的 HTML 主内容            |

## 鸣谢

- <https://github.com/nocd5/md2html>
- <https://github.com/tscanlin/tocbot>
- <https://github.com/mathjax/MathJax>
- <https://github.com/shd101wyy/markdown-preview-enhanced>
- <https://github.com/PuerkitoBio/goquery>
- <https://github.com/jessevdk/go-flags>
- <https://github.com/yuin/goldmark>
