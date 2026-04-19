package keystore

// FormatLink returns a formatted string for a service's link field,
// suitable for use in tabular output.
func FormatLink(link string) string {
	if link == "" {
		return "-"
	}
	return link
}
