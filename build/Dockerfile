FROM alpine:3.15
RUN set -eux &&sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add nginx
COPY nginx.conf /etc/nginx/nginx.conf
ADD ./package/web /web
ADD ./package/service/weichatapp /services/weichatapp/weichatapp
copy entrypoint.sh /entrypoint.sh
copy MP_verify_ZyR2psMiiv8mduT7.txt /web/MP_verify_ZyR2psMiiv8mduT7.txt
ENTRYPOINT ["sh","entrypoint.sh"]