#!/bin/bash

# 参考 https://blog.csdn.net/xyang81/article/details/52556886
# 参考 https://blog.csdn.net/l1028386804/article/details/80098334

apt install libssl-dev openssl libpopt-dev
apt install nginx keepalived

iptables -A INPUT -p vrrp -d 224.0.0.18/32 -j ACCEPT
apt install iptables-persistent
netfilter-persistent save
netfilter-persistent reload