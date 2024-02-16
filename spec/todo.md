# To consider
- [ ] App API support  
May be painful to implement.
Required to fully replace Pixiv, if user actions won't work universally.
https://codeberg.org/VnPower/PixivFE/issues/7

- [ ] Testing  
Do we really need testing? What to test?

- [ ] Storing values  
A JSON file to store values.
This way, values set by users won't be lost after restarts.

- [ ] User discovery  
For discovery page.  
Pretty useless if user actions (following) doesn't work.

- [ ] "Popular" artworks  
Check the README of this:  
https://github.com/kokseen1/Mashiro

- [ ] i18n  
The last thing to work on, probably.

- [x] Ranking page  
A lot of options weren't implemented.

- [x] Revisit ranking calendar  
There should be a way to display R18 thumbnails now?

# To implement
- [x] Multiple tokens support  
Let the host supply multiple tokens at once to avoid overuse.

- [ ] Pixivision  
https://www.pixivision.net/en/  
Pretty good to discover new artworks n stuff.  
Implement by parsing the webpage.

- [ ] RSS support for Pixivision  
- [ ] Search page  
A page to do more extensive searching.  
Might require JavaScript for search recommendation, if wanted.

- [ ] Manga series  
Serialized web comics. Example: https://www.pixiv.net/user/13651304/series/171013

- [ ] Novel support  
Might need some ideas for the reader's UI.  
Allow options for font size and family?  
Black and white backgrounds?  
Theme support?  

- [ ] Native AI/R15/R18/R18-G... artwork filtering  
We can filter them out using values supplied by Pixiv for each artworks.

- [ ] Full landing page  
There are a lot of sections for the landing page. https://www.pixiv.net/ajax/top/illust  
The artwork parsing part has already been implemented flawlessly.  
We only have to write the frontend code for those sections.

- [x] Merge settings page with login page  
- [ ] Various interesting pages from Pixiv.net  
  - https://www.pixiv.net/idea/
  - https://www.pixiv.net/request
  - https://www.pixiv.net/contest/ (no AJAX endpoints)
