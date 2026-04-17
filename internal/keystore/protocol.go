package keystore

import (
	"errors"
	"sort"
)

var validProtocols = map[string]bool{
	"rest":    true,
	"grpc":    true,
	"graphql": true,
	"soap":    true,
	"webhook": true,
	"other":   true,
}

func (s *Store) SetProtocol(service, protocol string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return errors.New("service not found: " + service)
	}
	if !validProtocols[protocol] {
		return errors.New("invalid protocol: " + protocol)
	}
	e := s.data.Entries[service]
	e.Protocol = protocol
	s.data.Entries[service] = e
	return s.save()
}

func (s *Store) GetProtocol(service string) (string, error) {
	e, ok := s.data.Entries[service]
	if !ok {
		return "", errors.New("service not found: " + service)
	}
	if e.Protocol == "" {
		return "rest", nil
	}
	return e.Protocol, nil
}

func (s *Store) ClearProtocol(service string) error {
	if _, ok := s.data.Entries[service]; !ok {
		return errors.New("service not found: " + service)
	}
	e := s.data.Entries[service]
	e.Protocol = ""
	s.data.Entries[service] = e
	return s.save()
}

func (s *Store) ServicesByProtocol(protocol string) []string {
	var results []string
	for name, e := range s.data.Entries {
		effective := e.Protocol
		if effective == "" {
			effective = "rest"
		}
		if effective == protocol {
			results = append(results, name)
		}
	}
	sort.Strings(results)
	return results
}
