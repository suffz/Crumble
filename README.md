
<h1 align="center" class="icon">
  <a>
    <img src="https://avatars.githubusercontent.com/u/84757238?v=4"></img>
  </a>
</h1>

[(API UPD NOTES) v1.1.0_BETA](https://namemc.info/) [Minecraft Auction + Namemc Info Bot](https://discord.com/api/oauth2/authorize?client_id=1157922376495943822&permissions=0&scope=applications.commands%20bot)
- Added a cloudflare clearance gen for namemc, uses YOUR proxys and YOUR capmonster api key to make the requests needed to return results.
- Added name checking to the /data/info endpoint.. and other endpoints relying on name > data conversion.. situation being names greater than 16 or invalid names (a-Z 0-9 _ names that dont contain these ranges of characters) shouldnt be requested on my side due to network usage.
- deprecated the /data/namemc/head endpoint due to /data/info giving the same info and more.

upcoming updates:
- complete rework of websocket system for tasks
- sync/connect namemc account on-site
- add ability to use your own api via supplying your own proxys and capmonster key
- forums, global chat, profiles, general things to allow others to view logs of your previous successes and any posts about ads.
- in relation to the ability to sync/connect a namemc account, the ability to follow bot and manage said profile from namemc may become a feature.
- namemc account generator (possibly)

| This api bypasses cloudflare without the assistance of chrome drivers or external services, it is free and gives up to date data when the name isn't cached in the backend, it will remain online for as long as I'm willing to pay for a monthly vps. |

- https://namemc.info/data/info/Dream (profile endpoint removed due to /info giving the same data)
- https://namemc.info/data/3c (?3n or ?3l)
- https://namemc.info/data/namemc/frontpage?pages=1&searches=1
- https://namemc.info/data/namemc/skinart/logo/samouraisniper
- https://namemc.info/data/namemc/skins?pages=10
- https://namemc.info/data/namemc/skin/data/dream
