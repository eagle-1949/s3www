# s3www
Serve static files from any S3 compatible object storage endpoints.

## Install
If you do not have a working Golang environment, please follow [How to install Golang](https://docs.minio.io/docs/how-to-install-golang).
```
go get github.com/harshavardhana/s3www
```

## Run
Make sure you have `index.html` under `website-bucket`
```
s3www -endpoint "https://s3.amazonaws.com" -accessKey "accessKey" \
      -secretKey "secretKey" -bucket "website-bucket"

s3www: Started listening on http://0.0.0.0:8080

```
If you use Docker
```
ENDPOINT="https://s3.amazonaws.com"
ACCESS_KEY="accessKey"
SECRET_KEY="secretKey"
BUCKET="website-bucket"
ADDRESS="ADDRESS:PORT"
TLS_CERT=""
TLS_KEY=""
```

## Test
Point your web browser to http://0.0.0.0:8080 ensure your `s3www` is serving your `index.html` successfully.

## License
This project is distributed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0), see [LICENSE](./LICENSE) for more information.

