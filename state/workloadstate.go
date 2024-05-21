package main

import (
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/rancher/cli/cliclient"
	"github.com/rancher/cli/config"
	"github.com/rancher/norman/types"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	token     = flag.String("token", getRncherToken(), "token")
	endpoint  = flag.String("endpoint", getRancherAddress(), "endpoint of the Source rancher server")
	project   = flag.String("project", getProject(), `project`)
	namespace = flag.String("namepace", getNamespace(), `namespace`)

	defaultAuth = "gHjIoLpor9oiyugfvcsiolrd4434fde"
)

func getProject() string {
	value := os.Getenv("PROJECT")
	if value == "" {
		return "local:p-hdjsk"
	}
	return value
}

func getRncherToken() string {
	value := os.Getenv("TOKEN")
	if value == "" {
		return "cffstql8mh5vzq5lt8dt9mw24qqqm2k2t2vbnrtfqgm5qsmw6jd7gb"
	}
	return value
}

func getRancherAddress() string {
	value := os.Getenv("SERVER_URL")
	if value == "" {
		return "https://paas.cydb.com.cn"
	}
	return value
}

func getNamespace() string {
	value := os.Getenv("NAMEPACE")
	if value == "" {
		return "dfs-prod"
	}
	return value
}

func main() {
	flag.Parse()
	var authKey string
	fmt.Print("please input you auth info: ")
	psw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("get auth info error")
		os.Exit(1)
	}
	authKey = string(psw)
	if len(authKey) <= 0 {
		fmt.Println("auth info is empty")
		os.Exit(1)
	}

	if defaultAuth != authKey {
		fmt.Println("invalid auth info")
		os.Exit(1)
	}
	*token = "token-pl8wm:" + *token
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
	c.CACerts = "-----BEGIN CERTIFICATE-----\nMIIDKDCCAhCgAwIBAgIJAO5rPTaMVfs2MA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAoMQswCQYDVQQGEwJDTjEZMBcGA1UEAwwQcGFhcy5jeWRiLmNv\nbS5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN8F0e5UrzIH/D1x\n7ej6YvnBoyZof5qC1M3sVFHrbFJAWUPDcUxQgwv2A1ShqjQ1wSAJGU/Xtw46HJ8i\nPr5UBeQmLT2MGA9u8zLQPQHIk7RvKvDoU8rFD8Sbpt61Dv4sTtg3vsOW02XHD7Kq\nzaaisgyfAOXqaSezqUvihf9O/b31yjOqg2X7CIFf2esfTO59LWQ4t325Ba2TkUq4\nNKglfeVGKVz4dMDkx13YQVBaw/xFGJufeiA+vI0+l4dxSHoNGa/ngcsAz5N6Z8UC\nMMZn0Z3rTWh2XjuoImG37N+OkA042lRY4No8hLcPEuyPVfzyTz8k8Sl1h8TlLawA\n62e8Su8CAwEAAaNcMFowCQYDVR0TBAIwADALBgNVHQ8EBAMCBeAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwIGCCsGAQUFBwMBMCEGA1UdEQQaMBiCEHBhYXMuY3lkYi5jb20u\nY26HBAoRDC8wDQYJKoZIhvcNAQELBQADggEBAMNvQxgxOerpZnDNALYvruvNOi6q\nYPrdphQ/hKwdOk9sFNxa+l9K+ObWEBTW1yQ13E09yocsxq8d9WpAyHqOudyO0OZx\nKJK1fZtLFAHTGQKxOYZd5SmpAYrhxdDF3evrLlunWy2+Ixpnbis12l6DoZeGwnr2\nD2YRov6Zd5P4Bu8dLpGAzta1pB56xoXHEH0RUfevv9Ead0IqZgu0NhW2OXmMNEQD\nTptXVZjuRPypLWeYT//gi8EdiLuWDov1zUtZiakXdvF30TN0NsyGFBpVqPzCPv+d\nWNGu+eAPVNst18TGkSKOB2AQgBckjPzu6mKVSkec72MExa/6XB+SO4zmdww=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIDFTCCAf2gAwIBAgIJALAKspWdq3BhMA0GCSqGSIb3DQEBCwUAMCExCzAJBgNV\nBAYTAkNOMRIwEAYDVQQDDAljYXR0bGUtY2EwHhcNMjMwMTEyMDEzMjMwWhcNMzMw\nMTA5MDEzMjMwWjAhMQswCQYDVQQGEwJDTjESMBAGA1UEAwwJY2F0dGxlLWNhMIIB\nIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy+1mMaRcxJnib1/zBFhePLnu\nVD+XHkDo/kNaa0vVf0ESukCUhhflssz5fI65UP7XpiMUj5hQkRTday3/IK+4Ormw\nzIAFyZtY655BEeQFPXGkMqva0FuAGEqXzlyRIth3KAGbE0yChVtt5Kxa3P8A/CLF\nSzE6Zb49vV4KiECWD8ZIfuLQ9k3YjG7W+v4p0f+yhjcLzBlNgp8QBgc0QDr1rj5E\n2zExP1IoBP+uRuJP2WBOUweHbLaqVHCx1HqE40V5Tj1MBrixLnN3IVb8AGBzdhIC\nr+r2KIa7VgBW6QOcZHHRQcq5DrZguGTK+WPEjtyvOByvsVci824PIZMx9RPiewID\nAQABo1AwTjAdBgNVHQ4EFgQUfSA+u82plxPubgYZgjsWpBhmgn0wHwYDVR0jBBgw\nFoAUfSA+u82plxPubgYZgjsWpBhmgn0wDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B\nAQsFAAOCAQEAB6nyZVeAR1d+N4Cdztp860/BgnuoDvwfQm880KjLCbAMkOicGw/i\nHrzTlchtUFr4B8Z349QwGhAA+reEmGu3zzDg/+Yokzk0MO5UoAcsA/n13UP5LSi7\n5AhJQraSf6izPlsh5pjFxXVW+NsH/ObK4YQ3umbpYAMChfE+12crLZBuIB5YZPlQ\nzakBDx8vERAkYgvJ3Li3mVOjW08CLzVcBEAqPPwrdG/Wqd8BWScYP72FOGdg3irl\nb+WgEK8IEOuS3we6vVbarS3Bq0rhS4DDQOWntUp0NqeD67LZU3/Ux7fh9jt8OQ0O\nrmC5Lb8vVKILsiLPCsDbC8xACsqkkp5vlQ==\n-----END CERTIFICATE-----\n"

	c.Project = *project
	//project client
	proclitic, err := cliclient.NewProjectClient(c)
	if err != nil {
		fmt.Printf("connect server error。")
		os.Exit(1)
	}

	plist, err := proclitic.ProjectClient.Pod.List(&types.ListOpts{
		Filters: map[string]interface{}{
			"namespaceId": *namespace,
			"limit":       -1,
			"all":         true,
		},
	})

	if err != nil {
		fmt.Printf("find project pod error, %v \n", err)
		os.Exit(1)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "State", "restart", "Start Time", "Image"})
	location, err := time.LoadLocation("Asia/Shanghai")
	for _, pod := range plist.Data {
		startTime, err := time.ParseInLocation(time.RFC3339, pod.Status.StartTime, time.UTC)
		index := strings.Index(pod.Containers[0].Image, "/")
		var s string
		if index > 0 {
			s = pod.Containers[0].Image[index:]
		} else {
			s = pod.Containers[0].Image
		}
		if err != nil {
			table.Append([]string{pod.Name, pod.State, strconv.FormatInt(pod.Containers[0].RestartCount, 10), pod.Status.StartTime, s})
		} else {
			table.Append([]string{pod.Name, pod.State, strconv.FormatInt(pod.Containers[0].RestartCount, 10), startTime.In(location).Format("2006-01-02 15:04:05"), s})
		}
	}
	table.Render()
}
