#!/usr/bin/env bash

echo "[z] Loading env vars..."

if [ -f .env ]; then
    # Load Environment Variables
    export $(egrep -v '^#' .env | xargs)
fi

# GH variables
AUTH_HEADER="Authorization: token $GH_ACCESS_TOKEN"
GH_BODY='{"tag_name":"'$V'","name": "Zen Release v'$V'","draft": true}'

echo ""
echo "[z] Creating Zen release v"$V "on Github"
echo 

# make initial release reqest
resp=$(curl -H "Content-Type:application/json" -H "$AUTH_HEADER"  https://api.github.com/repos/zaviermiller/zen/releases -d "$GH_BODY")
eval $(echo "$resp" | grep -m 1 "id.:" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
[ "$id" ] || { echo "Error: Failed to get release id for tag: $resp"; echo "$response" | awk 'length($0)<100' >&2; exit 1; }

echo ""
echo "[z] Uploading binaries..."
echo ""

# upload binaries for all platforms
for GOOS in darwin linux windows freebsd; do 
    for GOARCH in 386 amd64; do 
        GH_ASSET="https://uploads.github.com/repos/zaviermiller/zen/releases/$id/assets?name=zen$V-$GOOS-$GOARCH"
        

        $(curl --data-binary @./bin/zen$V-$GOOS-$GOARCH -H "Content-Type: application/octet-stream" -H "$AUTH_HEADER" $GH_ASSET)
    done
done

echo ""
echo "Done."
echo ""
