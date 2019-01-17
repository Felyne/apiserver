#!/usr/bin/env bash
cd ../ && openssl req -new -nodes -x509 -out conf/server.crt -keyout conf/server.key -days 3650 \
-subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=127.0.0.1/emailAddress=chenwd_135@qq.com"
