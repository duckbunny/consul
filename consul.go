package consul

// This package meets the requirements of the herald Pool and Declare interfaces for Consul.
import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/duckbunny/herald"
	"github.com/duckbunny/service"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/consul"
)

var (
	// TTL time to life for service in consul
	TTL int = 15

	// Where the ServiceKVPath resides
	KVpath string = "services"

	// Title for specifying herald in flags
	Title string = "consul"

	// Config falls back to client default config
	ConsulConfig api.Config = api.DefaultConfig()
)

func init() {
	ttl := os.Getenv("CONSUL_TTL")
	if ttl != "" {
		newttl, err := strconv.Atoi(ttl)
		if err != nil {
			log.Fatal(err)
		}
		TTL = newttl
	}
	flag.IntVar(&TTL, "consul-ttl", TTL, "TTL for consul microservice heartbeats.")
}

//  Consul structure
type Consul struct {
	// Agent to register service
	Agent *api.Agent
	// KV to save service definition
	KV            *api.KV
	heartBeatKill chan bool
}

// New Consul
func New() *Consul {
	c := new(Consul)
	c.heartBeatKill = make(chan bool)
	return c
}

// Start Register the service in the consul pool of services
func (c *Consul) Start(s *service.Service) error {
	p, err := strconv.ParseInt(s.Port, 10, 0)
	if err != nil {
		return err
	}
	AgentService := api.AgentServiceRegistration{
		ID:   FormattedID(s),
		Name: FormattedName(s),
		Port: int(p),
		Check: &api.AgentServiceCheck{
			TTL: fmt.Sprintf("%vs", TTL),
		},
	}
	// Register the service
	err = c.Agent.ServiceRegister(&AgentService)
	if err != nil {
		return err
	}
	// Initial run for TTL
	c.Agent.PassTTL(fmt.Sprintf("service:%v", FormattedID(s)), "TTL heartbeat")

	// Begin TTL refresh
	go c.Heartbeat(s)
	return nil
}

// Kill the Hearteat and remove the service
func (c *Consul) Stop(s *service.Service) error {
	c.heartBeatKill <- true
	return c.Agent.ServiceDeregister(FormattedID(s))
}

// Init Consul herald with Default Settings
func (c *Consul) Init() error {

	client, err := api.NewClient(ConsulConfig)
	if err != nil {
		return err
	}
	if c.Agent == nil {
		c.Agent = client.Agent()
	}

	if c.KV == nil {
		c.KV = client.KV()
	}
	return nil

}

// Send service definition to consul
func (c *Consul) Declare(s *service.Service) error {
	js, err := json.Marshal(s)
	if err != nil {
		return err
	}
	key := FormattedKey(s)
	pair := api.KVPair{
		Key:   key,
		Flags: 0,
		Value: js,
	}
	_, err = c.KV.Put(&pair, nil)
	return err
}

// Retrieve the consul service definition.  Requires Domain, Title and Version be set.  Returns err if not found.
func (c *Consul) Get(s *service.Service) error {
	key := FormattedKey(s)
	qo := api.QueryOptions{}
	v, _, err := c.KV.Get(key, &qo)
	if err != nil {
		return err
	}
	return json.Unmarshal(v.Value, s)
}

// FormattedName returns correctly formatted name of the service
func FormattedName(s *service.Service) string {
	name := fmt.Sprintf("%v-%v-%v", s.Domain, s.Title, s.Version)
	return strings.Replace(name, ".", "-", -1)
}

// FormattedID returns correctly formatted id of the service
func FormattedID(s *service.Service) string {
	return fmt.Sprintf("%v-%v-%v", FormattedName(s), s.Host, s.Port)
}

// FormattedKey returns correctly formatted key of the service
func FormattedKey(s *service.Service) string {
	return fmt.Sprintf("%v/%v/%v/%v/definition", KVpath, s.Domain, s.Title, s.Version)
}

// Heartbeat begins heart beat of health check.
func (c *Consul) Heartbeat(s *service.Service) {
	for _ = range time.Tick(time.Duration(TTL-1) * time.Second) {
		select {
		case <-c.heartBeatKill:
			return
		default:
		}
		c.Agent.PassTTL(fmt.Sprintf("service:%v", FormattedID(s)), "TTL heartbeat")
	}
}

// Register this herald with consul
func Register() {
	c := consul.New()
	herald.AddPool(Title, c)
	herald.AddDeclaration(Title, c)
}
