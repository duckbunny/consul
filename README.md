###Consul

Is the consul implementation for herald.

[![GoDoc](https://godoc.org/github.com/duckbunny/consul?status.svg)](https://godoc.org/github.com/duckbunny/consul)


# consul
--
    import "github.com/duckbunny/consul"

The consul package implements the Pool and Declare interfaces for the Herald
package. https://github.com/duckbunny/herald

The package takes one flag "consul-ttl" to set the time until the service
expires from consul using a heartbeat.

The package utilizies the Default client returned from consul api, but this can
be overriden by editing.

    ConsulConfig.Config.Address = "192.168.1.56"

This service must be registered with herald.

    consul.Register()

Or you can fall back on the herald service registry, to register all available
services.


## Usage

```go
var (
	// TTL time to life for service in consul
	TTL int = 15

	// Where the ServiceKVPath resides
	KVpath string = "services"

	// Title for specifying herald in flags
	Title string = "consul"

	// Config falls back to client default config
	ConsulConfig *api.Config = api.DefaultConfig()
)
```

#### func  FormattedID

```go
func FormattedID(s *service.Service) string
```
FormattedID returns correctly formatted id of the service

#### func  FormattedKey

```go
func FormattedKey(s *service.Service) string
```
FormattedKey returns correctly formatted key of the service

#### func  FormattedName

```go
func FormattedName(s *service.Service) string
```
FormattedName returns correctly formatted name of the service

#### func  Register

```go
func Register()
```
Register this herald with consul

#### type Consul

```go
type Consul struct {
	// Agent to register service
	Agent *api.Agent
	// KV to save service definition
	KV *api.KV
}
```

Consul structure

#### func  New

```go
func New() *Consul
```
New Consul

#### func (*Consul) Declare

```go
func (c *Consul) Declare(s *service.Service) error
```
Send service definition to consul

#### func (*Consul) Get

```go
func (c *Consul) Get(s *service.Service) error
```
Retrieve the consul service definition. Requires Domain, Title and Version be
set. Returns err if not found.

#### func (*Consul) Heartbeat

```go
func (c *Consul) Heartbeat(s *service.Service)
```
Heartbeat begins heart beat of health check.

#### func (*Consul) Init

```go
func (c *Consul) Init() error
```
Init Consul herald with Default Settings

#### func (*Consul) Start

```go
func (c *Consul) Start(s *service.Service) error
```
Start Register the service in the consul pool of services

#### func (*Consul) Stop

```go
func (c *Consul) Stop(s *service.Service) error
```
Kill the Hearteat and remove the service
