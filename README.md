
<h1 align="center" class="icon">
  <a>
    <img src="https://avatars.githubusercontent.com/u/84757238?v=4"></img>
  </a>
</h1>

⚠️ THIS DOES NOT WORK WITHOUT PROXIES

A sniper that uses proxies and mass account support, it sends periodic requests.

Note: If you are coming into this sniper expecting to claim anything with a single account and no proxies you will be highly disappointed.

# Api Documentation.
the ign "god" will be used for areas that require a value.

# /data/info/:name
https://namemc.info/data/info/god

```json
{
    "action": "message",
    "desc": "god Returned Info.",
    "code": "200",
    "data": {
        "status": "Unavailable",
        "searches": "3577",
        "headurl": "https://s.namemc.com/2d/skin/face.png?id=0ad00123b87bb341\u0026scale=4",
        "bodyurl": "https://s.namemc.com/3d/skin/body.png?id=0ad00123b87bb341\u0026model=classic\u0026width=150\u0026height=200"
    }
}
```
# /data/3c
https://namemc.info/data/3c
> https://namemc.info/data/3c?3n (returns only three number igns, example: 123)
> https://namemc.info/data/3c?3l (returns only three letter igns, example: abc)

# /data/discord/:id
https://namemc.info/data/discord/12345678 (1234567 being the discord id, replies with the users discord username etc)

# /data/namemc/head/:name
https://namemc.info/data/namemc/head/god (returns the HEAD/BODY url example)
```json
{
    "bodyurl": "https://s.namemc.com/3d/skin/body.png?id=0ad00123b87bb341\u0026model=classic\u0026width=150\u0026height=200",
    "headurl": "https://s.namemc.com/2d/skin/face.png?id=0ad00123b87bb341\u0026scale=4",
    "id": "0ad00123b87bb341"
}
```

# /data/namemc/frontpage?pages=10&searches=0
https://namemc.info/data/namemc/frontpage?pages=10&searches=0 
```json
{
    "action": "message",
    "desc": "Returned your parsed pages.",
    "code": "200",
    "data": [
        {
            "name": "zelhilfeverwen",
            "start_date": "2023-06-20T02:35:36.839Z",
            "end_date": "2023-06-20T17:47:31.497Z",
            "start_unix": 1687228536,
            "end_unix": 1687283251,
            "searches": "0"
        },
        ...
```

# /data/namemc/skin/data/:name
https://namemc.info/data/namemc/skin/data/god
```json
{
    "hearts": "4",
    "stars": "884",
    "used": "6.1y",
    "users_equiped": [
        {
            "follower_url": "https://namemc.com/profile/Jonah_.1",
            "name": "Jonah_",
            "data": {
                "download": "https://s.namemc.com/img/emoji/twitter/1f30c.svg",
                "emoji": "🌌",
                "rank": "Moderator",
                "emojiequiped": true
            }
        },
      ...
```

# /data/namemc/skins
https://namemc.info/data/namemc/skins?pages=10
```json
{
    "action": "message",
    "desc": "Returned your parsed pages.",
    "code": "200",
    "data": [
        {
            "emoji": "",
            "owner": "zavodkirpichey",
            "number": "#1",
            "stars": "135★",
            "time": "4.7d",
            "bodyurl": "https://s.namemc.com/3d/skin/body.png?id=38943e9eb9870039\u0026model=classic\u0026width=150\u0026height=200",
            "headurl": "https://s.namemc.com/2d/skin/face.png?id=38943e9eb9870039\u0026scale=4",
            "skindownload": "https://s.namemc.com/i/38943e9eb9870039.png"
        },
      ...
```

# /data/namemc/skinart/logo/:name

![Example](https://namemc.info/data/namemc/skinart/logo/SamouraiClaimer)

# /data/profile/:name
https://namemc.info/data/profile/god
```json
{
    "bio": "",
    "followers": [
        {
            "follower_url": "https://namemc.com/profile/Dqnieel.2",
            "name": "Dqnieel",
            "data": {
                "download": "https://s.namemc.com/img/emoji/twitter/2744-fe0f.svg",
                "emoji": "❄️",
                "rank": "Emerald",
                "emojiequiped": true
            }
        },
        ...
    ],
    "headurl": "https://s.namemc.com/2d/skin/face.png?id=0ad00123b87bb341\u0026scale=4",
    "namehist": [
        "1 God Original"
    ],
    "skins": [
        {
            "download": "https://s.namemc.com/i/0ad00123b87bb341.png",
            "url": "https://namemc.com/skin/0ad00123b87bb341",
            "id": "0ad00123b87bb341",
            "changedat": "2017-05-11T07:08:43.434Z",
            "headurl": "https://s.namemc.com/2d/skin/face.png?id=0ad00123b87bb341\u0026scale=4",
            "bodyurl": "https://s.namemc.com/3d/skin/body.png?id=0ad00123b87bb341\u0026model=classic\u0026width=150\u0026height=200"
        },
        ...
    ],
    "uuid": "bc27afd7-6889-4811-97c9-135ee46cdabc",
    "views": "1094"
}
```
