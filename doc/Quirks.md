# Known quirks

## Why aren't my userstyles working?

Origin: https://codeberg.org/VnPower/PixivFE/pulls/62#issuecomment-1568191

PixivFE implements a strong [Content-Security-Policy](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy) that prevents inline styles from being loaded. If you're using Stylus, you need to enable **Advanced > Circumvent CSP 'style-src' via adoptedStyleSheets** in Stylus Options (see [issue #1685](https://github.com/openstyles/stylus/issues/1685) on the Stylus GitHub repository).
