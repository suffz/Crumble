
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

May make this API key only soon, but i like to find work arounds to keep it low resource intensive, no paid services, method packed, so right now there really is no need.

| This api bypasses cloudflare without the assistance of chrome drivers or external services, it is free and gives up to date data when the name isn't cached in the backend, it will remain online for as long as I'm willing to pay for a monthly vps.

- https://namemc.info/data/info/:name (profile endpoint removed due to /info giving the same data)
- https://namemc.info/data/3c (?3n or ?3l)
- https://namemc.info/data/namemc/frontpage?pages=1&searches=1
- https://namemc.info/data/namemc/skinart/logo/samouraisniper

oh cloudflare..
