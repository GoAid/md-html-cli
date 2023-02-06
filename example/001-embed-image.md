# Embed Image Sample

## Command Option

```bash
$ md-html-cli example/embed-image.md -e
```

### Local file is embedded

```markdown
![logo/go](img/go.png "go")
```

![logo/go](img/go.png "go")

#### Write HTML tag immediately

```markdown
<img src="img/go.png" width="48">
<img src="img/go.png" width="32">
<img src="img/go.png" width="16">
```

<img src="img/go.png" width="48">
<img src="img/go.png" width="32">
<img src="img/go.png" width="16">

### image as URL link is not embedded

```markdown
![google/errors/robot](https://www.google.com/images/errors/robot.png "robot")
```

![google/errors/robot](https://www.google.com/images/errors/robot.png "robot")
