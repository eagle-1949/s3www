package main

import (
	"flag"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/minio/mc/pkg/console"
	minio "github.com/minio/minio-go"
)

// S3 - A S3 implements FileSystem using the minio client
// allowing access to your S3 buckets and objects.
//
// Note that S3 will allow all access to files in your private
// buckets, If you have any sensitive information please make
// sure to not sure this project.
type S3 struct {
	*minio.Client
	bucket string
}

// Open - implements http.Filesystem implementation.
func (s3 *S3) Open(name string) (http.File, error) {
	if strings.HasSuffix(name, "/") {
		return &httpMinioObject{
			client: s3.Client,
			object: nil,
			isDir:  true,
			bucket: bucket,
			prefix: name,
		}, nil
	}

	name = strings.TrimPrefix(name, "/")
	obj, err := s3.Client.GetObject(s3.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, os.ErrNotExist
	}

	if _, err = obj.Stat(); err != nil {
		return nil, os.ErrNotExist
	}

	return &httpMinioObject{
		client: s3.Client,
		object: obj,
		isDir:  false,
		bucket: bucket,
		prefix: name,
	}, nil
}

var (
	endpoint  string
	accessKey string
	secretKey string
	address   string
	bucket    string
	tlsCert   string
	tlsKey    string
)

func init() {
	flag.StringVar(&endpoint, "endpoint", "https://s3.amazonaws.com", "S3 server endpoint.")
	flag.StringVar(&accessKey, "accessKey", "", "Access key of S3 storage.")
	flag.StringVar(&secretKey, "secretKey", "", "Secret key of S3 storage.")
	flag.StringVar(&bucket, "bucket", "", "Bucket name which hosts static files.")
	flag.StringVar(&address, "address", ":8080", "Bind to a specific ADDRESS:PORT, ADDRESS can be an IP or hostname.")
	flag.StringVar(&tlsCert, "ssl-cert", "", "TLS certificate for this server.")
	flag.StringVar(&tlsKey, "ssl-key", "", "TLS private key for this server.")

}

func main() {
	flag.Parse()
	// 获取bucket
	if strings.TrimSpace(bucket) == "" {
		bucket = os.Getenv("BUCKET")
		if strings.TrimSpace(bucket) == "" {
			console.Fatalln(`Bucket name cannot be empty, please provide 's3www -bucket "mybucket"'`)
		}

	}
	// 获取endpoint
	if strings.TrimSpace(endpoint) == "" {
		endpoint = os.Getenv("ENDPOINT")
		if strings.TrimSpace(endpoint) == "" {
			console.Fatalln(`Endpoint cannot be empty, please provide 's3www -endpoint "myendpoint"'`)
		}
	}
	// 获取accessKey
	if strings.TrimSpace(accessKey) == "" {
		accessKey = os.Getenv("ACCESS_KEY")
		if strings.TrimSpace(accessKey) == "" {
			console.Fatalln(`accessKey cannot be empty, please provide 's3www -accessKey "myaccessKey"'`)
		}
	}
	// 获取secretKey
	if strings.TrimSpace(secretKey) == "" {
		secretKey = os.Getenv("SECRET_KEY")
		if strings.TrimSpace(secretKey) == "" {
			console.Fatalln(`secretKey cannot be empty, please provide 's3www -secretKey "mysecretKey"'`)
		}
	}
	// 获取address
	if strings.TrimSpace(address) == "" {
		address = os.Getenv("ADDRESS")
		if strings.TrimSpace(address) == "" {
			address = ":8080"
		}
	}
	// 获取ssl-cert
	if strings.TrimSpace(tlsCert) == "" {
		tlsCert = os.Getenv("TLS_CERT")

	}
	// 获取ssl-key
	if strings.TrimSpace(tlsKey) == "" {
		tlsKey = os.Getenv("TLS_KEY")

	}

	u, err := url.Parse(endpoint)
	if err != nil {
		console.Fatalln(err)
	}

	client, err := minio.NewV4(u.Host, accessKey, secretKey, u.Scheme == "https")
	if err != nil {
		console.Fatalln(err)
	}

	if tlsCert != "" && tlsKey != "" {
		console.Infof("Started listening on https://%s\n", address)
		console.Fatalln(http.ListenAndServeTLS(address, tlsCert, tlsKey, http.FileServer(&S3{client, bucket})))
	} else {
		console.Infof("Started listening on http://%s\n", address)
		console.Fatalln(http.ListenAndServe(address, http.FileServer(&S3{client, bucket})))
	}
}
