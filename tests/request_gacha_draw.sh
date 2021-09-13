curl -sX POST http://localhost:8080/gacha/draw -H 'accept: application/json' -H 'Content-Type: application/json' \
	-H 'x-token: string' \
	-d '{
  "times": 10
}' | jq .
