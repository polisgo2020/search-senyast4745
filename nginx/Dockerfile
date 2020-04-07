FROM nginx:stable-alpine
COPY default.conf /etc/nginx/conf.d
EXPOSE 80/tcp
EXPOSE 443/tcp
WORKDIR /usr/share/nginx/html
COPY static/. .

CMD ["/bin/sh", "-c", "exec nginx -g 'daemon off;';"]
