# Fences Divs Sample

```markdown
<style>
.green, .yellow { border-radius: 5px; }
.green { background-color: #42b983; padding: 15px; }
.yellow { background-color: #fdbc40; }
#inside-me { color: #fc625d; }
</style>

:::{.green}
## Life Inside Fences

We are now inside a div with the css-class "green". This can be used to style this block

:::{#inside-me .yellow data="important"}
fences can be nested and given ids as well as classes
:::
:::
```

<style>
.green, .yellow { border-radius: 5px; }
.green { background-color: #42b983; padding: 15px; }
.yellow { background-color: #fdbc40; text-align: center; }
#inside-me { color: #fc625d; }
</style>

:::{.green}
## Life Inside Fences

We are now inside a div with the css-class "green". This can be used to style this block

:::{#inside-me .yellow data="important"}
fences can be nested and given ids as well as classes
:::
:::
