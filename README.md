## deploy command

### Create
```
gcloud functions deploy bw-fetch \
	--runtime=go111 \
	--trigger-http \
	--entry-point=Fetch \
	--timeout=30 \
```

### Update
```
gcloud functions deploy bw-fetch
```

## Call function
sn default value=927
```
curl --url {endpoint} \
	-d '{"sn": 5530}' \
	-H 'Content-Type: application/json'
```
