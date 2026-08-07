sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl

curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \
  | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg

curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \
  | sudo tee /etc/apt/sources.list.d/caddy-stable.list

sudo apt update
sudo apt install caddy


################################

Config

sudo vi /etc/caddy/Caddyfile

suppose your hostname = budget.alephapps.com

put this in your Caddyfile

budget.alephapps.com {
    reverse_proxy localhost:8080
}

sudo systemctl reload caddy


