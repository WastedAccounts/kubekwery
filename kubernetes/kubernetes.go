package kubernetes

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// kuberenetes.go will read the home dir /.kube/config file
//https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go

type Contexts struct {
	Name string
	//Other string
}
// ListContexts will list available contexts from /.kube/config
// for Rundeck we can call this to get a list and return it to a DDL for selection
func ListContexts() []Contexts{
	var contexts []Contexts
	var kubeconfigList *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigList = flag.String("kubeconfigList", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigList = flag.String("kubeconfkubeconfigListiglist", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//Dereferences the pointer and prints kubconfig as string
	strkubeconfig := DerefString(kubeconfigList)

	kubeConfig, err := clientcmd.LoadFromFile(strkubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	// Load Context names into Contexts struct
	for a  := range kubeConfig.Contexts{
		contexts = append(contexts, Contexts{Name: a})
	}
	fmt.Println("kubeConfig.Contexts",kubeConfig.Contexts)
	fmt.Println("kubeConfig.Contexts",kubeConfig.Clusters)
	fmt.Println("kubeConfig.Contexts",kubeConfig.Preferences)
	return contexts
}

// CallContext will connect to the default context for now
// Will update this to be able to pass in a context
func CallContext() { //(context string){
	var kubeconfigCall *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigCall = flag.String("kubeconfigCall", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfigCall = flag.String("kubeckubeconfigCallonfig1", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	//Dereferences the pointer and prints kubconfig as string
	strkubeconfig := DerefString(kubeconfigCall)
	fmt.Println("kubeconfigCall:",strkubeconfig)

	// uses the current context in ~/.kube/config
	// Working on how to pass in a specific context and call that.
	// this requires the URL though. Will probably have to recall ListContexts
	// and pull the masterURL and then pass that in.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfigCall)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("config:",config)

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		namespace := "default"
		pod := "mysql-d8d99fb4-k5xn4"
		_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}

// a simple function to deref string pointers
func DerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
