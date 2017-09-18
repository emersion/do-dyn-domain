package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/emersion/go-upnp-igd"
	"github.com/ProtonMail/go-appdir"
	"github.com/digitalocean/godo"
	"github.com/digitalocean/godo/context"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var domainName = flag.String("domain", "", "Domain name")
var recordName = flag.String("record", "", "Record name")

type Config struct {
	AccessToken string `yaml:"access-token"`
}

func (c *Config) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: c.AccessToken}, nil
}

func main() {
	flag.Parse()
	if *domainName == "" || *recordName == "" {
		log.Fatal("Missing -domain or -record")
	}

	configFile := filepath.Join(appdir.New("doctl").UserConfig(), "config.yaml")
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config := new(Config)
	if err := yaml.Unmarshal(b, config); err != nil {
		log.Fatal(err)
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, config)
	client := godo.NewClient(oauthClient)

	devices := igd.Discover(10 * time.Second)
	if len(devices) == 0 {
		log.Fatal("No gateway found")
	}
	d := devices[0]

	log.Println("Using gateway:", d.ID())

	ip, err := d.GetExternalIPAddress()
	if err != nil {
		log.Fatal(err)
	}
	ipStr := ip.String()
	log.Println("Discovered external IP address:", ipStr)

	ctx := context.TODO()
	records, _, err := client.Domains.Records(ctx, *domainName, nil)
	if err != nil {
		log.Fatal(err)
	}

	var record *godo.DomainRecord
	for _, r := range records {
		if r.Name == *recordName {
			record = &r
			break
		}
	}
	if record == nil {
		log.Fatal("No such record")
	}

	if record.Data == ipStr {
		log.Println("IP address didn't change")
		return
	}

	ctx = context.TODO()
	_, _, err = client.Domains.EditRecord(ctx, *domainName, record.ID, &godo.DomainRecordEditRequest{
		Data: ipStr,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Record updated with new IP address")
}
