
<h1 align="center" class="icon">
  <a>
    <img src="https://avatars.githubusercontent.com/u/84757238?v=4"></img>
  </a>
</h1>

⚠️ THIS DOES NOT WORK WITHOUT PROXIES | [Minecraft Auction + Namemc Info Bot](https://discord.com/api/oauth2/authorize?client_id=1157922376495943822&permissions=0&scope=applications.commands%20bot)

Note: If you are coming into this sniper expecting to claim anything with a single account and no proxies you will be highly disappointed.

[(API UPD NOTES) v1.0.0BETA](https://namemc.info/)
- Implemented real-time caching, if you request a name without a profile itll update it in the database and return the profile info
- the /data/info/:name endpoint now has a extra json entry "images" this is the base64 info for the skin images on the users profile.
- /data/namemc/skinart/logo/:name checks for a active database entry and if the entry has profile data and any "images" within it, if not it falls back on real-time information.
- The website itself has gotten a overall change of UI and ive also fixed mobile support to a degree, i plan to make it functional soon (vps management, sniping tasks etc)

| This api bypasses cloudflare without the assistance of chrome drivers or external services, it is free and gives up to date data when the name isn't cached in the backend, it will remain online for as long as I'm willing to pay for a monthly vps. |

- https://namemc.info/data/info/Dream (profile endpoint removed due to /info giving the same data)
- https://namemc.info/data/3c (?3n or ?3l)
- https://namemc.info/data/namemc/frontpage?pages=1&searches=1
- https://namemc.info/data/namemc/skinart/logo/samouraisniper

| clearance token task system PROXYS REQUIRED |

This api uses the capmonster api, and you will need to supply your own token.

- https://namemc.info/data/clearance (POST)

Request Body:
```json
{
	"useragent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
  	"capmonster_key":"67.....",
  	"proxy":{
    	"ip":"167.0.0.1",
      	"port":"100",
      	"user":"username",
      	"password":"password"
    }
}
```

Response Body:
```json
{
    "uuid": "UUIDHERE",
    "result": {
        "solution": {
            "cf_clearance": ""
        },
        "cost": 0,
        "status": "pending",
        "errorId": 0,
        "errorCode": null,
        "errorDescription": null
    },
    "done": false,
    "clearance": {
	      "useragent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
  	    "capmonster_key":"67.....",
  	    "proxy":{
    	      "ip":"167.0.0.1",
      	    "port":"100",
      	    "user":"username",
      	    "password":"password"
        }
    }
}
```

- https://namemc.info/data/clearance/UUIDHERE (GET)

Response Data:
```json
{
    "uuid": "UUIDHERE",
    "result": {
        "solution": {
            "cf_clearance": ".0H5jP9PO0Nc86xB8Ge7TJ..."
        },
        "cost": 0.0016,
        "status": "ready",
        "errorId": 0,
        "errorCode": null,
        "errorDescription": null
    },
    "done": false,
    "clearance": {
        "useragent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "capmonster_key": "679....",
        "proxy": {
            "ip": "167.0.0.1",
            "port": "100",
            "user": "user",
            "password": "password"
        }
    }
}
```

Error Response: (403)
```json
{"error":"Task doesnt exist."}
```
