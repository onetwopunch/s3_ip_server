package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"html/template"
	"log"
	"net/http"
	"os"
  "strings"
  "bufio"
)
const tmpl = `
<html><body>
<h1>Bad Ips</h1>
<ul>
  {{ range . }}
    <li>{{ . }}</li>
  {{ end }}
</ul></body></html>`

func getIpsFromS3(bucket, object string) ([]string, error) {
	// Original list from https://www.dshield.org/ipsascii.html?limit=100

  output := []string{}
	sess := GetSession()
	svc := s3.New(sess)
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return output, err
	} else {
    buf := bufio.NewReader(result.Body)
    for {
      line, err := buf.ReadString('\n')
      if !strings.HasPrefix(line, "#") {
        ip := strings.Split(line, "\t")[0]
        if len(ip) > 0 {
          output = append(output, ip)
        }
      }
      if err != nil {
        break
      }
    }
    return output, nil
  }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  bucket := os.Getenv("AWS_BUCKET")
  object := os.Getenv("AWS_OBJECT")
	ips, err := getIpsFromS3(bucket, object)
  if err != nil {
    log.Fatal(err)
  }
	t := template.Must(template.New("index.html").Parse(tmpl))
	fmt.Println(ips)
	t.Execute(w, ips)
}

func GetSession() *session.Session {
	var creds *credentials.Credentials
	sess := session.Must(session.NewSession())
	meta := ec2metadata.New(sess)
	if meta.Available() {
		creds = credentials.NewCredentials(&ec2rolecreds.EC2RoleProvider{
			Client: meta,
		})
	} else {
		creds = credentials.NewEnvCredentials()
	}
	if _, err := creds.Get(); err != nil {
		log.Fatal(err)
	}
  region := os.Getenv("AWS_REGION")
	sess.Config = defaults.Config().WithCredentials(creds).WithRegion(region)
	return sess
}

func main() {
	http.HandleFunc("/ip", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
