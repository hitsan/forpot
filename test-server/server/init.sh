#!/bin/bash

/usr/sbin/sshd -D &
deno task start &
tail -f /dev/null
