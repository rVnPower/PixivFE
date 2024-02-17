# Known Quicks Of PixivFE

## Why don't my userstyles work?

Origin: https://codeberg.org/VnPower/PixivFE/pulls/62#issuecomment-1568191

This website uses <abbr title="Content Security Policy">CSP</abbr>, which blocks the loading of inline styles. In the case of Stylus, you need to enable **Advanced > Circumvent CSP 'style-­src' via adoptedSty­leSheets** in Stylus Options.

Reference: https://github.com/openstyles/stylus/issues/1685
