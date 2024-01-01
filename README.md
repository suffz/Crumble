
<h1 align="center" class="icon">
  <a>
    <img src="https://avatars.githubusercontent.com/u/84757238?v=4"></img>
  </a>
</h1>

[Minecraft Auction + Namemc Info Bot](https://discord.com/api/oauth2/authorize?client_id=1157922376495943822&permissions=0&scope=applications.commands%20bot)

Note: If you are coming into this sniper expecting to claim anything with a single account and no proxies you will be highly disappointed.

[(API UPD NOTES) v1.1.0_BETA](https://namemc.info/) 
- Added a cloudflare clearance gen for namemc, uses YOUR proxys and YOUR capmonster api key to make the requests needed to return results.
- Added name checking to the /data/info endpoint.. and other endpoints relying on name > data conversion.. situation being names greater than 16 or invalid names (a-Z 0-9 _ names that dont contain these ranges of characters) shouldnt be requested on my side due to network usage.
- deprecated the /data/namemc/head endpoint due to /data/info giving the same info and more.

upcoming updates:
- complete rework of websocket system for tasks
- sync/connect namemc account on-site
- add ability to use your own api via supplying your own proxys and capmonster key
- forums, global chat, profiles, general things to allow others to view logs of your previous successes and any posts about ads.

| This api bypasses cloudflare without the assistance of chrome drivers or external services, it is free and gives up to date data when the name isn't cached in the backend, it will remain online for as long as I'm willing to pay for a monthly vps. |

- https://namemc.info/data/info/Dream (profile endpoint removed due to /info giving the same data)
- https://namemc.info/data/3c (?3n or ?3l)
- https://namemc.info/data/namemc/frontpage?pages=1&searches=1
- https://namemc.info/data/namemc/skinart/logo/samouraisniper
- https://namemc.info/data/namemc/skins?pages=10
- https://namemc.info/data/namemc/skin/data/dream

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
    "balance": 5.9792,
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
    "balance": 5.9792,
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
