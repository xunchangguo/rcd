package main

import (
	"flag"
	"fmt"
	"github.com/mohae/deepcopy"
	"github.com/rancher/cli/cliclient"
	"github.com/rancher/cli/config"
	"github.com/rancher/norman/types"
	"os"
	"strings"
)

var (
	token         = flag.String("token", getRancherToken(), "address:port to listen on")
	endpoint      = flag.String("endpoint", getRancherAddress(), "endpoint of the rancher server")
	srcProject    = flag.String("src_project", getSourceProject(), `source project`)
	srcNamespace  = flag.String("src_namepace", getSourceNamespace(), `source namespace`)
	destProject   = flag.String("dest_project", getDestProject(), `dest project`)
	destNamespace = flag.String("dest_namespace", getDestNamespace(), `dest namespace`)
)

func getDestProject() string {
	value := os.Getenv("DEST_PROJECT")
	if value == "" {
		return "c-tsqvs:p-qpmwg"
	}
	return value
}

func getDestNamespace() string {
	value := os.Getenv("DEST_NAMEPACE")
	if value == "" {
		return "dfs-uat"
	}
	return value
}

func getSourceNamespace() string {
	value := os.Getenv("SRC_NAMEPACE")
	if value == "" {
		return "test"
	}
	return value
}

func getSourceProject() string {
	value := os.Getenv("SRC_PROJECT")
	if value == "" {
		return "c-tsqvs:p-sv76l"
	}
	return value
}

func getRancherToken() string {
	value := os.Getenv("RANCHER_TOKEN")
	if value == "" {
		return "token-m7cnq:8wv6rntvdh7thl8xbt8qxqn2g9w4l2vg2vhg5xh6fnmzhjprbj7s98"
	}
	return value
}

func getRancherAddress() string {
	value := os.Getenv("RANCHER_SERVER_URL")
	if value == "" {
		return "https://10.17.12.41:8443"
	}
	return value
}

func baseListOpts() *types.ListOpts {
	return &types.ListOpts{
		Filters: map[string]interface{}{
			"limit": -1,
			"all":   true,
		},
	}
}

func main() {
	flag.Parse()
	auth := strings.Split(*token, ":")
	if len(auth) != 2 {
		fmt.Println("invalid token")
		os.Exit(1)
	}

	c := &config.ServerConfig{
		URL:       *endpoint,
		AccessKey: auth[0],
		SecretKey: auth[1],
		TokenKey:  *token,
		//Insecure: true,
	}
	c.CACerts = "-----BEGIN CERTIFICATE-----\nMIIBqDCCAU2gAwIBAgIBADAKBggqhkjOPQQDAjA7MRwwGgYDVQQKExNkeW5hbWlj\nbGlzdGVuZXItb3JnMRswGQYDVQQDExJkeW5hbWljbGlzdGVuZXItY2EwHhcNMjIx\nMTAzMDE1OTM2WhcNMzIxMDMxMDE1OTM2WjA7MRwwGgYDVQQKExNkeW5hbWljbGlz\ndGVuZXItb3JnMRswGQYDVQQDExJkeW5hbWljbGlzdGVuZXItY2EwWTATBgcqhkjO\nPQIBBggqhkjOPQMBBwNCAAQ1h2wZGIFEjZtWr+zqRb22OYPV+kzi1Aew+bm6KWIu\nKzSdHxDSQxwzbHTl+XnPZySoZkdAHmPtImgprfrmNfFDo0IwQDAOBgNVHQ8BAf8E\nBAMCAqQwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUXTpI4t9fHPHAO/hvx6PO\n8wHv4HkwCgYIKoZIzj0EAwIDSQAwRgIhALLzRMFtE8FvLkys2YVKqJX7B6ePEaPK\nlgq5fH2OD1sHAiEA+FntIqCGOK2tGHVP6FyMTWe59C0p+T0C8IBh9P/AQf0=\n-----END CERTIFICATE-----"

	c.Project = *srcProject
	//project client
	prjcli, err := cliclient.NewProjectClient(c) //, "c-skmrj:p-s42lb")
	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	workloads, err := prjcli.ProjectClient.Workload.List(baseListOpts())
	if err != nil {
		fmt.Printf("list project workloads error, %v \n", err)
		os.Exit(1)
	}
	fmt.Println("-----------------project workloads ---------------------")
	srcMap := map[string]string{}
	for _, workload := range workloads.Data {
		if workload.NamespaceId == *srcNamespace {
			for _, container := range workload.Containers {
				envName := fmt.Sprintf("%s_%s", workload.Name, container.Name)
				fmt.Println(envName, " = ", container.Image)
				srcMap[envName] = container.Image
			}
		}
	}
	//dest
	c.Project = *destProject
	prjcli, err = cliclient.NewProjectClient(c)
	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}

	workloads, err = prjcli.ProjectClient.Workload.List(baseListOpts())
	if err != nil {
		fmt.Printf("list project workloads error, %v \n", err)
		os.Exit(1)
	}
	fmt.Println("-----------------project workloads ---------------------")

	for _, workload := range workloads.Data {
		if workload.NamespaceId == *destNamespace {
			cimages := make([]interface{}, len(workload.Containers), len(workload.Containers))
			chang := false
			for idx, container := range workload.Containers {
				envName := fmt.Sprintf("%s_%s", workload.Name, container.Name)
				fmt.Println(envName, " = ", container.Image)
				if len(srcMap[envName]) > 0 && container.Image != srcMap[envName] {
					chang = true
					container.Image = srcMap[envName]
				}
				cimages[idx] = deepcopy.Copy(container)
			}
			if chang {
				_, err = prjcli.ProjectClient.Workload.Update(&workload, map[string][]interface{}{
					"containers": cimages,
				})
				if err != nil {
					fmt.Printf("Update project %s workload app error, %v \n", workload.Name, err)
				} else {
					fmt.Printf("Update project %s workload sucess\n", workload.Name)
				}
			}
		}
	}
}
