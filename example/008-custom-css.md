# Custom CSS Sample

## Command Option

```bash
$ md-html-cli example/custom-css.md -c example/css/custom-css.css
```

### wheel

#### Example for `wheel` class

```markdown
<img class="vinyl ep" src="img/go.png">
<img class="vinyl lp" src="img/go.png">
```

#### Custom CSS for `wheel`

```css
/* example/css/custom-css.css */

img.vinyl {
    margin: 12px;
}

img.vinyl.ep {
    animation: spin calc(60s / 45) linear infinite;
}

img.vinyl.lp {
    animation: spin calc(60s / (33 + 1 / 3)) linear infinite;
}

@keyframes spin {
    from { transform: rotate(0deg);   }
    to   { transform: rotate(360deg); }
}
```

<img class="vinyl ep" src="img/go.png">
<img class="vinyl lp" src="img/go.png">
<!-- dummy for avoid .markdown-body > :last-child { margin-bottom: 0 !important; } -->
<p/>
