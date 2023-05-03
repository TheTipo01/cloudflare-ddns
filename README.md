# cloudflare-ddns
This is a tool to help you manage dynamic IPs under cloudflare's dns service. The tool periodically fetches the current IPv4 and/or IPv6
and updates it over on cloudflare if it's new.

There are multiple interfaces (plugins) for this tool, each one provides a different way to fetch the public IPs.

I welcome any and all PRs trying to expand on the features of this app so feel free to contribute with your own plugin/additions.

# Installation
Installing this tool is easy, just clone this repository and run make.
You can use this one liner:
```
git clone https://github.com/nylone/cloudflare-ddns /tmp/cloudflare-ddns && cd /tmp/cloudflare-ddns && make && mv out $HOME/cloudflare-ddns && cd $HOME/cloudflare-ddns && rm -rf /tmp/cloudflare-ddns
```
which is just:
```
git clone https://github.com/nylone/cloudflare-ddns /tmp/cloudflare-ddns
cd /tmp/cloudflare-ddns
make
mv out $HOME/cloudflare-ddns
cd $HOME/cloudflare-ddns
rm -rf /tmp/cloudflare-ddns
```
