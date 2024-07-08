```shell
apk add iproute2-tc

# inject a 500ms delay - will result in 2 connections because of 1s recurrence
tc qdisc add dev eth0 root netem delay 500ms

# delete the rule
tc qdisc del dev eth0 root

# inject 75% packet loss
tc qdisc add dev eth0 root netem loss 75%
```
