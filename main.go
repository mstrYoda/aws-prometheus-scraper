package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	ecsSvc  *ecs.ECS
	cluster *string
)

func main() {
	region := flag.String("region", "eu-central-1", "")
	cluster = flag.String("cluster", "base-cluster", "")
	services := flag.String("services", "api-go-service,go-consumer-service,service-dedup-api-service", "")
	flag.Parse()

	sess := session.Must(session.NewSession(&aws.Config{Region: aws.String(*region)}))
	ecsSvc = ecs.New(sess)

	jobs := make([]Job, 0)
	for _, svc := range strings.Split(*services, ",") {
		jobs = append(jobs,
			Job{
				Name:  svc,
				Tasks: getTaskDefinitions(svc),
			})
	}

	parsed := parseTemplate(jobs)

	if err := os.WriteFile("prometheus.yml", []byte(parsed), 777); err != nil {
		fmt.Println(err)
	}
}
