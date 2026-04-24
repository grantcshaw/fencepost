package keystore

import "fmt"

// FormatLink returns a formatted string for a service's link field,
// suitable for use in tabular output.
func FormatLink(link string) string {
	if link == "" {
		return "-"
	}
	return link
}

// FormatLinkWithLabel returns a formatted string combining a label and link,
// suitable for use in tabular output. If the link is empty, only the label
// is returned. If the label is empty, only the link is returned.
func FormatLinkWithLabel(label, link string) string {
	if link == "" {
		return FormatLink(label)
	}
	if label == "" {
		return link
	}
	return fmt.Sprintf("%s (%s)", label, link)
}
