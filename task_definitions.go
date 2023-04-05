package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Job struct {
	Name  string
	Tasks []TaskData
}

type TaskData struct {
	Name           string `json:"name"`
	TaskARNVersion string `json:"TaskARNVersion"`
	IPAddr         string `json:"IPAddr"`
}

func getTaskDefinitions(svcName string) []TaskData {
	tasks, err := ecsSvc.ListTasks(&ecs.ListTasksInput{
		Cluster:       cluster,
		DesiredStatus: aws.String(ecs.DesiredStatusRunning),
		ServiceName:   aws.String(svcName),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	taskDefinitions, err := ecsSvc.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: cluster,
		Tasks:   tasks.TaskArns,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	taskData := make([]TaskData, 0)
	for _, t := range taskDefinitions.Tasks {
		for _, c := range t.Containers {
			for _, ni := range c.NetworkInterfaces {
				splitted := strings.Split(*t.TaskDefinitionArn, ":")
				if len(splitted) == 0 {
					continue
				}

				taskData = append(taskData, TaskData{
					TaskARNVersion: splitted[len(splitted)-1],
					IPAddr:         *ni.PrivateIpv4Address,
				})
			}
		}
	}

	return taskData
}
