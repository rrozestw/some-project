server {
	root /var/www/html;
	index index.html index.htm index.nginx-debian.html;
	server_name fh.rozestwinski.com; # managed by Certbot
	listen [::]:443 ssl ipv6only=on; # managed by Certbot
	listen 443 ssl; # managed by Certbot
	ssl_certificate /etc/letsencrypt/live/fh.rozestwinski.com/fullchain.pem; # managed by Certbot
	ssl_certificate_key /etc/letsencrypt/live/fh.rozestwinski.com/privkey.pem; # managed by Certbot
	include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
	ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot
	location / {
		proxy_pass http://127.0.0.1:8080;
		proxy_set_header X-Request-Id $msec;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	}
}
server {
    if ($host = fh.rozestwinski.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot

	listen 80 ;
	listen [::]:80 ;
	server_name fh.rozestwinski.com;
	return 404; # managed by Certbot
}
