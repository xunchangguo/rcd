package main

import (
	"flag"
	"fmt"
	"github.com/mohae/deepcopy"
	"github.com/rancher/cli/cliclient"
	"github.com/rancher/cli/config"
	"github.com/rancher/norman/types"
	client "github.com/rancher/types/client/project/v3"
	"os"
	"regexp"
	"strings"
)

var (
	token        = flag.String("src_token", getSourceRancherToken(), "source token")
	endpoint     = flag.String("src_endpoint", getSourceRancherAddress(), "endpoint of the Source rancher server")
	srcProject   = flag.String("src_project", getSourceProject(), `source project`)
	srcNamespace = flag.String("src_namepace", getSourceNamespace(), `source namespace`)

	destToken     = flag.String("dest_token", getDestRancherToken(), "dest token")
	destEndpoint  = flag.String("dest_endpoint", getDestRancherAddress(), "endpoint of the dest rancher server")
	destProject   = flag.String("dest_project", getDestProject(), `dest project`)
	destNamespace = flag.String("dest_namespace", getDestNamespace(), `dest namespace`)
	nodeSelect    = flag.String("node_select", "stage = prod", `worker node selector`)

	cmdType  = flag.Int("type", 0, `0: create,1: update`)
	ignoreWs = flag.String("ignore", "nacos", `ignore workload`)
	ipReg    = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)(:[0-9]{1,5})?$`
)

func getDestProject() string {
	value := os.Getenv("DEST_PROJECT")
	if value == "" {
		return "local:p-hdjsk"
	}
	return value
}

func getDestNamespace() string {
	value := os.Getenv("DEST_NAMEPACE")
	if value == "" {
		return "dfs-prod"
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

func getSourceRancherToken() string {
	value := os.Getenv("SRC_TOKEN")
	if value == "" {
		return "token-m7cnq:8wv6rntvdh7thl8xbt8qxqn2g9w4l2vg2vhg5xh6fnmzhjprbj7s98"
	}
	return value
}

func getDestRancherToken() string {
	value := os.Getenv("DEST_TOKEN")
	if value == "" {
		return "token-5c68p:79phkbc2mvzpgs69lg2bq7dmp8g2plq9px5vgl64phtxvxwkhvnvxl"
	}
	return value
}

func getSourceRancherAddress() string {
	value := os.Getenv("SRC_SERVER_URL")
	if value == "" {
		return "https://10.17.12.41:8443"
	}
	return value
}

func getDestRancherAddress() string {
	value := os.Getenv("DEST_SERVER_URL")
	if value == "" {
		return "https://paas.cydb.com.cn"
	}
	return value
}

func getDestCACerts() string {
	value := os.Getenv("DEST_CA")
	if value == "" {
		return "-----BEGIN CERTIFICATE-----\nMIIDKDCCAhCgAwIBAgIJAO5rPTaMVfs2MA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAoMQswCQYDVQQGEwJDTjEZMBcGA1UEAwwQcGFhcy5jeWRiLmNv\nbS5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN8F0e5UrzIH/D1x\n7ej6YvnBoyZof5qC1M3sVFHrbFJAWUPDcUxQgwv2A1ShqjQ1wSAJGU/Xtw46HJ8i\nPr5UBeQmLT2MGA9u8zLQPQHIk7RvKvDoU8rFD8Sbpt61Dv4sTtg3vsOW02XHD7Kq\nzaaisgyfAOXqaSezqUvihf9O/b31yjOqg2X7CIFf2esfTO59LWQ4t325Ba2TkUq4\nNKglfeVGKVz4dMDkx13YQVBaw/xFGJufeiA+vI0+l4dxSHoNGa/ngcsAz5N6Z8UC\nMMZn0Z3rTWh2XjuoImG37N+OkA042lRY4No8hLcPEuyPVfzyTz8k8Sl1h8TlLawA\n62e8Su8CAwEAAaNcMFowCQYDVR0TBAIwADALBgNVHQ8EBAMCBeAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwIGCCsGAQUFBwMBMCEGA1UdEQQaMBiCEHBhYXMuY3lkYi5jb20u\nY26HBAoRDC8wDQYJKoZIhvcNAQELBQADggEBAMNvQxgxOerpZnDNALYvruvNOi6q\nYPrdphQ/hKwdOk9sFNxa+l9K+ObWEBTW1yQ13E09yocsxq8d9WpAyHqOudyO0OZx\nKJK1fZtLFAHTGQKxOYZd5SmpAYrhxdDF3evrLlunWy2+Ixpnbis12l6DoZeGwnr2\nD2YRov6Zd5P4Bu8dLpGAzta1pB56xoXHEH0RUfevv9Ead0IqZgu0NhW2OXmMNEQD\nTptXVZjuRPypLWeYT//gi8EdiLuWDov1zUtZiakXdvF30TN0NsyGFBpVqPzCPv+d\nWNGu+eAPVNst18TGkSKOB2AQgBckjPzu6mKVSkec72MExa/6XB+SO4zmdww=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIJALAKspWdq3BhMA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAhMQswCQYDVQQGEwJDTjESMBAGA1UEAwwJY2F0dGxlLWNhMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy+1mMaRcxJnib1/zBFhePLnu\nVD+XHkDo/kNaa0vVf0ESukCUhhflssz5fI65UP7XpiMUj5hQkRTday3/IK+4Ormw\nzIAFyZtY655BEeQFPXGkMqva0FuAGEqXzlyRIth3KAGbE0yChVtt5Kxa3P8A/CLF\nSzE6Zb49vV4KiECWD8ZIfuLQ9k3YjG7W+v4p0f+yhjcLzBlNgp8QBgc0QDr1rj5E\n2zExP1IoBP+uRuJP2WBOUweHbLaqVHCx1HqE40V5Tj1MBrixLnN3IVb8AGBzdhIC\nr+r2KIa7VgBW6QOcZHHRQcq5DrZguGTK+WPEjtyvOByvsVci824PIZMx9RPiewID\nAQABo1AwTjAdBgNVHQ4EFgQUfSA+u82plxPubgYZgjsWpBhmgn0wHwYDVR0jBBgw\nFoAUfSA+u82plxPubgYZgjsWpBhmgn0wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B\nAQsFAAOCAQEAB6nyZVeAR1d+N4Cdztp860/BgnuoDvwfQm880KjLCbAMkOicGw/i\nHrzTlchtUFr4B8Z349QwGhAA+reEmGu3zzDg/+Yokzk0MO5UoAcsA/n13UP5LSi7\n5AhJQraSf6izPlsh5pjFxXVW+NsH/ObK4YQ3umbpYAMChfE+12crLZBuIB5YZPlQ\nzakBDx8vERAkYgvJ3Li3mVOjW08CLzVcBEAqPPwrdG/Wqd8BWScYP72FOGdg3irl\nb+WgEK8IEOuS3we6vVbarS3Bq0rhS4DDQOWntUp0NqeD67LZU3/Ux7fh9jt8OQ0O\nrmC5Lb8vVKILsiLPCsDbC8xACsqkkp5vlQ==\n-----END CERTIFICATE-----\n"
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
	destAuth := strings.Split(*destToken, ":")
	if len(destAuth) != 2 {
		fmt.Println("invalid  dest token")
		os.Exit(1)
	}

	var arr []string = nil
	if len(*ignoreWs) == 0 {
		arr = strings.Split(*ignoreWs, ",")
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
	srcMap := map[string]map[string]string{}
	var srcWorkLoad = map[string]client.Workload{}
	for _, workload := range workloads.Data {
		if workload.NamespaceId == *srcNamespace {
			srcWorkLoad[workload.Name] = workload
			for _, container := range workload.Containers {
				envName := fmt.Sprintf("%s_%s", workload.Name, container.Name)
				fmt.Println(envName, " = ", container.Image)
				m := map[string]string{}
				m[envName] = container.Image
				for key, val := range container.Environment {
					m[key] = val
				}
				srcMap[envName] = m
			}
		}
	}

	//dest
	destc := &config.ServerConfig{
		URL:       *destEndpoint,
		AccessKey: destAuth[0],
		SecretKey: destAuth[1],
		TokenKey:  *destToken,
	}
	data, err := os.ReadFile("destCa.txt")
	if err == nil {
		destc.CACerts = string(data)
	} else {
		destc.CACerts = "-----BEGIN CERTIFICATE-----\nMIIDKDCCAhCgAwIBAgIJAO5rPTaMVfs2MA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAoMQswCQYDVQQGEwJDTjEZMBcGA1UEAwwQcGFhcy5jeWRiLmNv\nbS5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN8F0e5UrzIH/D1x\n7ej6YvnBoyZof5qC1M3sVFHrbFJAWUPDcUxQgwv2A1ShqjQ1wSAJGU/Xtw46HJ8i\nPr5UBeQmLT2MGA9u8zLQPQHIk7RvKvDoU8rFD8Sbpt61Dv4sTtg3vsOW02XHD7Kq\nzaaisgyfAOXqaSezqUvihf9O/b31yjOqg2X7CIFf2esfTO59LWQ4t325Ba2TkUq4\nNKglfeVGKVz4dMDkx13YQVBaw/xFGJufeiA+vI0+l4dxSHoNGa/ngcsAz5N6Z8UC\nMMZn0Z3rTWh2XjuoImG37N+OkA042lRY4No8hLcPEuyPVfzyTz8k8Sl1h8TlLawA\n62e8Su8CAwEAAaNcMFowCQYDVR0TBAIwADALBgNVHQ8EBAMCBeAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwIGCCsGAQUFBwMBMCEGA1UdEQQaMBiCEHBhYXMuY3lkYi5jb20u\nY26HBAoRDC8wDQYJKoZIhvcNAQELBQADggEBAMNvQxgxOerpZnDNALYvruvNOi6q\nYPrdphQ/hKwdOk9sFNxa+l9K+ObWEBTW1yQ13E09yocsxq8d9WpAyHqOudyO0OZx\nKJK1fZtLFAHTGQKxOYZd5SmpAYrhxdDF3evrLlunWy2+Ixpnbis12l6DoZeGwnr2\nD2YRov6Zd5P4Bu8dLpGAzta1pB56xoXHEH0RUfevv9Ead0IqZgu0NhW2OXmMNEQD\nTptXVZjuRPypLWeYT//gi8EdiLuWDov1zUtZiakXdvF30TN0NsyGFBpVqPzCPv+d\nWNGu+eAPVNst18TGkSKOB2AQgBckjPzu6mKVSkec72MExa/6XB+SO4zmdww=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIJALAKspWdq3BhMA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAhMQswCQYDVQQGEwJDTjESMBAGA1UEAwwJY2F0dGxlLWNhMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy+1mMaRcxJnib1/zBFhePLnu\nVD+XHkDo/kNaa0vVf0ESukCUhhflssz5fI65UP7XpiMUj5hQkRTday3/IK+4Ormw\nzIAFyZtY655BEeQFPXGkMqva0FuAGEqXzlyRIth3KAGbE0yChVtt5Kxa3P8A/CLF\nSzE6Zb49vV4KiECWD8ZIfuLQ9k3YjG7W+v4p0f+yhjcLzBlNgp8QBgc0QDr1rj5E\n2zExP1IoBP+uRuJP2WBOUweHbLaqVHCx1HqE40V5Tj1MBrixLnN3IVb8AGBzdhIC\nr+r2KIa7VgBW6QOcZHHRQcq5DrZguGTK+WPEjtyvOByvsVci824PIZMx9RPiewID\nAQABo1AwTjAdBgNVHQ4EFgQUfSA+u82plxPubgYZgjsWpBhmgn0wHwYDVR0jBBgw\nFoAUfSA+u82plxPubgYZgjsWpBhmgn0wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B\nAQsFAAOCAQEAB6nyZVeAR1d+N4Cdztp860/BgnuoDvwfQm880KjLCbAMkOicGw/i\nHrzTlchtUFr4B8Z349QwGhAA+reEmGu3zzDg/+Yokzk0MO5UoAcsA/n13UP5LSi7\n5AhJQraSf6izPlsh5pjFxXVW+NsH/ObK4YQ3umbpYAMChfE+12crLZBuIB5YZPlQ\nzakBDx8vERAkYgvJ3Li3mVOjW08CLzVcBEAqPPwrdG/Wqd8BWScYP72FOGdg3irl\nb+WgEK8IEOuS3we6vVbarS3Bq0rhS4DDQOWntUp0NqeD67LZU3/Ux7fh9jt8OQ0O\nrmC5Lb8vVKILsiLPCsDbC8xACsqkkp5vlQ==\n-----END CERTIFICATE-----\n"
	}

	destc.Project = *destProject
	destPrjCli, err := cliclient.NewProjectClient(destc)
	if err != nil {
		fmt.Printf("%v \n", err)
		os.Exit(1)
	}
	if *cmdType == 0 {
		fmt.Println("================ create cmd type ================")
		cfgList, err := prjcli.ProjectClient.ConfigMap.List(baseListOpts())
		if err != nil {
			fmt.Printf("list project configmap error, %v \n", err)
			os.Exit(1)
		}

		for _, cfg := range cfgList.Data {
			if cfg.NamespaceId == *srcNamespace {
				cfg.NamespaceId = *destNamespace
				cfg.ProjectID = *destProject
				_, err = destPrjCli.ProjectClient.ConfigMap.Create(&cfg)
				if err != nil {
					fmt.Printf("create config %s err, %v \n", cfg.Name, err)
				} else {
					fmt.Printf("create config %s sucessed \n", cfg.Name)
				}
			}
		}
		fmt.Println(" ********** create workload *******************")
		ig := false
		s := new(int64)
		*s = 0
		for _, workload := range workloads.Data {
			ig = false
			workload.NamespaceId = *destNamespace
			workload.ProjectID = *destProject
			for _, v := range arr {
				if strings.Contains(v, workload.Name) {
					fmt.Printf("ignore workload %s \n", workload.Name)
					ig = true
					break
				}
			}
			workload.Scale = s
			workload.Scheduling.Node.RequireAll = []string{*nodeSelect}
			if ig == false {
				_, err = destPrjCli.ProjectClient.Workload.Create(&workload)
				if err != nil {
					fmt.Printf("create workload %s err, %v \n", workload.Name, err)
				} else {
					fmt.Printf("create workload %s sucessed \n", workload.Name)
				}
			}
		}
	} else {
		fmt.Println("================ update cmd type ================")
		workloads, err = destPrjCli.ProjectClient.Workload.List(baseListOpts())
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
					c2, ok := srcMap[envName]
					if ok {
						w, ok := srcWorkLoad[workload.Name]
						if ok {
							for _, c := range w.Containers {
								if container.Name == c.Name {
									if container.Image != c2[envName] || len(container.Environment) != len(c.Environment) {
										chang = true
										container.Image = c2[envName]
										for key, val := range c2 {
											if key == envName || key == "SPRING_PROFILES_ACTIVE" || key == "SPRING_CLOUD_NACOS_CONFIG_NAMESPACE" || key == "SPRING_CLOUD_NACOS_DISCOVERY_NAMESPACE" {
												continue
											}
											match, _ := regexp.MatchString(ipReg, val)
											if match == false {
												container.Environment[key] = val
											}
										}
										container.ReadinessProbe = c.ReadinessProbe
										container.LivenessProbe = c.LivenessProbe
										cimages[idx] = deepcopy.Copy(container)
									}
									break
								}
							}
						}
					}
				}
				if chang {
					_, err = destPrjCli.ProjectClient.Workload.Update(&workload, map[string][]interface{}{
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
}
