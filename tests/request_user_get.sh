curl -sX GET http://localhost:8080/user/get -H 'accept: application/json'\
	-H 'x-token: string'\
| jq .
