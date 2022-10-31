#!/bin/sh

#/ test.sh [--ed25519]

usage() {
  grep '^#/' "$0" | cut -c4- >&2
  exit 1
}

IDENTITY_FILE="id_rsa_client"
while [ "$#" -gt 0 ]; do
  case "$1" in
    --ed25519) IDENTITY_FILE="id_ed25519_client"; shift ;;
    *) usage ;;
  esac
done


go get .
go run . &

# wait for server to start
until nc -v -z localhost 2022 >/dev/null 2>&1; do sleep 1; done

chmod 0600 "$IDENTITY_FILE"

ssh \
  -o UserKnownHostsFile=/dev/null \
  -o StrictHostKeyChecking=no \
  -o IdentityAgent=/dev/null \
  -o IdentityFile="$IDENTITY_FILE" \
  -p 2022 \
  -v \
  localhost \
> output 2>&1

if grep -q 'no mutual signature algorithm' output; then
  echo "⛔ connection failed"
  echo "⛔ $(grep 'no mutual signature algorithm' output)"
  exit 1
fi

if grep -q 'Hello, world!' output; then
  echo "✅ Connection worked!"
  exit 0
fi
