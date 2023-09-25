
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
