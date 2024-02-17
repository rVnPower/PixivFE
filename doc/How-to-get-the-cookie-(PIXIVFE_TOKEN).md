# How to get the cookie (PIXIVFE_TOKEN)

This guide covers how to get your Pixiv account's cookie to authenticate.

> **Note**:
>
> You should create an entirely new account for this to avoid account theft. And also, PixivFE will get contents **from your account.** You might not want people to know what kind of illustrations you like :P. For now, the only page that may contain contents that is relevant to you is the discovery page. Be careful if you are using your main account.

## Firefox-based

1. Log in to your Pixiv account of choice. You should be greeted with the landing page with logging in. If you are already logged in, go to the landing page.

![The URL of the landing page](https://files.catbox.moe/7dbv3e.png)

2. Hit `F12` to open up the developer tools. Then, go to the `Storage` tab.

![Storage tab on Firefox](https://files.catbox.moe/mra6rs.png)

3. At the left side, open up the `Cookies` section. Then select `www.pixiv.net`, this is the place where you will get your cookie.
   The page now should look like the screenshot below. Select the cookie with the key `PHPSESSID`, the value next to it is your account's token.

![Cookie on Firefox](https://files.catbox.moe/zb16o8.png)

4. Copy it and set the environment variable! If deploying using Docker Compose, copy it into `docker/pixivfe_token.txt` instead.

## Chrome-based

1. Log in to your Pixiv account of choice. You should be greeted with the landing page with logging in. If you are already logged in, go to the landing page.

2. Hit `F12` to open up the developer tools. Then, go to the `Applications` tab.

3. At the left side, you can see the `Storage` section. Inside of that section, there is an another section called `Cookies`, open up the `Cookies` section, then select `www.pixiv.net`. This is the place where you will get your cookie.
   The page now should look like the screenshot below. Select the cookie with the key `PHPSESSID`, the value next to it is your account's token.

![PHPSESSID on Chrome-based browsers](https://files.catbox.moe/8wu9f0.png)

4. Copy it and set the environment variable! If deploying using Docker Compose, copy it into `docker/pixivfe_token.txt` instead.

## Note

- The token should look something like this: `123456_AaBbccDDeeFFggHHIiJjkkllmMnnooPP`. The part before the underline is your member ID, the part after the underline is just a random string.
- The token will reset when you logout. Please double-check that your token is still valid before reporting any issues.
- Chrome-based browsers and some content was taken from [this page by Nandaka.](https://github.com/Nandaka/PixivUtil2/wiki#pixiv-login-using-cookie)
