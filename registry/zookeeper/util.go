package zookeeper

import (
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/registry"
)

func encode(s *registry.Service) ([]byte, error) {
	return json.Marshal(s)
}

func decode(ds []byte) (*registry.Service, error) {
	var s *registry.Service
	err := json.Unmarshal(ds, &s)
	return s, err
}

func nodePath(domain, name, id string) string {
	service := strings.Replace(name, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	p := path.Join(prefixWithDomain(domain), service, node)
	log.Infof("register path: %s", p)
	return p
}

func prefixWithDomain(domain string) string {
	return path.Join(prefix, domain)
}

func childPath(parent, child string) string {
	return path.Join(parent, strings.Replace(child, "/", "-", -1))
}

func servicePath(domain, s string) string {
	return path.Join(prefixWithDomain(domain), serializeServiceName(s))
}

func serializeServiceName(s string) string {
	return strings.ReplaceAll(s, "/", "-")
}

func createPath(path string, data []byte, client *zk.Conn, ttl time.Duration) error {
	exists, _, err := client.Exists(path)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	name := "/"
	p := strings.Split(path, "/")

	for _, v := range p[1 : len(p)-1] {
		name += v
		e, _, _ := client.Exists(name)
		if !e {
			_, err = client.Create(name, []byte{}, int32(0), zk.WorldACL(zk.PermAll))
			if err != nil {
				return err
			}
		}
		name += "/"
	}

	if ttl > -1 {
		_, err = client.CreateTTL(path, data, zk.FlagTTL, zk.WorldACL(zk.PermAll), ttl)
	} else {
		_, err = client.Create(path, data, int32(0), zk.WorldACL(zk.PermAll))
	}
	return err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
