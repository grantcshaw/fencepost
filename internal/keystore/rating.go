package keystore

import "fmt"

var validRatings = map[string]bool{
	"critical": true,
	"high":     true,
	"medium":   true,
	"low":      true,
}

func (s *Store) SetRating(service, rating string) error {
	if _, ok := s.Data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	if !validRatings[rating] {
		return fmt.Errorf("invalid rating %q: must be one of critical, high, medium, low", rating)
	}
	e := s.Data.Entries[service]
	e.Rating = rating
	s.Data.Entries[service] = e
	return s.save()
}

func (s *Store) GetRating(service string) (string, error) {
	e, ok := s.Data.Entries[service]
	if !ok {
		return "", fmt.Errorf("service %q not found", service)
	}
	if e.Rating == "" {
		return "medium", nil
	}
	return e.Rating, nil
}

func (s *Store) ClearRating(service string) error {
	if _, ok := s.Data.Entries[service]; !ok {
		return fmt.Errorf("service %q not found", service)
	}
	e := s.Data.Entries[service]
	e.Rating = ""
	s.Data.Entries[service] = e
	return s.save()
}

func (s *Store) ServicesByRating(rating string) ([]string, error) {
	var out []string
	for name, e := range s.Data.Entries {
		effective := e.Rating
		if effective == "" {
			effective = "medium"
		}
		if effective == rating {
			out = append(out, name)
		}
	}
	sortResults(out)
	return out, nil
}
