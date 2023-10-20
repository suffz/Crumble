
<h1 align="center" class="icon">
  <a>
    <img src="https://avatars.githubusercontent.com/u/84757238?v=4"></img>
  </a>
</h1>

⚠️ THIS DOES NOT WORK WITHOUT PROXIES

A sniper that uses proxies and mass account support, it sends periodic requests.

Note: If you are coming into this sniper expecting to claim anything with a single account and no proxies you will be highly disappointed.

accounts.txt format:

email:password
OR
bearertoken

the bearer token can be just the token itself, it validates it by decrypting the string into a JWT, checks if minecraft generated it, then pulls the account info i.e is it a Giftcard, Microsoft, ETC, in of which uses its send method.

Minecraft Auction + Namemc Info Bot: https://discord.com/api/oauth2/authorize?client_id=1157922376495943822&permissions=0&scope=applications.commands%20bot

Updated discord bot to show drop info + current user profile information.

API now shows profile information and searches, status, etc.

the discord bot also has 2 new commands, /three and /skins.

- /three - send current cached 3 char list
- /skins - send a url to a image that will be parsed into skin-art format. (you apply via highest number to lowest)

the three char api also now shows the dropping 3 chars along with there old users profile information.

profile information has 2 new fields, data and hist, data shows badges like socials and location while hist gives information on the name hist of the profile.

https://namemc.info/data/info/:name
https://namemc.info/data/profile/:name
https://namemc.info/data/3c (?3n or ?3l)
https://namemc.info/data/namemc/frontpage?pages=1&searches=1
https://namemc.info/data/namemc/skinart/logo/samouraisniper

these are the more important endpoints.