package main

import (
	"fmt"
	"kubekwery/kubernetes"
	"os"
)

func main()  {
	fmt.Println("service started")
	// call program with arg 'contexts'
 	// will list all contexts in .kube/config file
 	if os.Args[1] == "contexts" {
		ListContexts := kubernetes.ListContexts()
		fmt.Println("test:",ListContexts)
		//for a := range ListContexts {
		//	CallContext := kubernetes.CallContext(a)
		//	fmt.Println("a::",a, "  CallContext::",CallContext)
		//}

		kubernetes.CallContext()


	}
	// does nothing yet
	if os.Args[1] == "notcontexts" {
		fmt.Println("not contexts arg passed in ")
	}
	if os.Args[1] == "" {
		fmt.Println("empty arg")
	}
	fmt.Println("service stopped")

}
