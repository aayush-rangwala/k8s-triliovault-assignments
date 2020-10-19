package main

import (
	"context"
	_ "context"
	"fmt"
	_ "k8s.io/client-go/informers/batch"
	_ "k8s.io/client-go/kubernetes/typed/batch/v1"
	"k8s.io/client-go/rest"
	_ "k8s.io/client-go/util/retry"

	_ "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	_ "k8s.io/api/core/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

func main() {
	config,err:=rest.InClusterConfig()
	if err!=nil {
		panic(err.Error())
	}

	newclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	const NamespaceName string = "vishu"
	jobClient := newclient.BatchV1().Jobs(NamespaceName)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vishu-job",
			Namespace: NamespaceName,
		},
		Spec: batchv1.JobSpec{
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "vishu-job",
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "vishu-job",
							Image: "ubuntu",
							Command: []string{
								"bin/bash",
								"-c",
								"for i in 9 8 7 6 5 4 3 2 1 ; do echo $i ; done",
							},
						},
					},
					RestartPolicy: apiv1.RestartPolicyOnFailure,
				},
			},
		},
	}

	_, err = jobClient.Create(context.TODO(),job,metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Job Created successfully!")

	fmt.Println("Listing job: ")
	jobList, err := jobClient.List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, jobs := range jobList.Items {
		fmt.Println("This is the job: ", jobs.Name)
	}

	//Deleting Replica-Set
	fmt.Println("Deleting Job ... ")
	if deleteErr := jobClient.Delete(context.TODO(),"vishu-job", metav1.DeleteOptions{}); deleteErr != nil {
		panic(err.Error())
	}
	fmt.Println("Job deleted")
}

func homeDir() string {
	return os.Getenv("ROOTPATH")
}