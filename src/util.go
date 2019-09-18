package main

import "time"

// Generates a proper ISO8601 or RFC3339 formatted date string
// given a year, month, and day value. Not for accuracy, just for
// necessary formatting.
// Ex: "2019-09-22T00:00:00Z00:00"
func generateIsoDate(yyyy string, mm string, dd string) string {
	const delim string = "-"

	return yyyy + delim + mm + delim + dd + "T00:00:00Z00:00"
}

// OrderedPosts - Extended slice with sorting methods attached
type OrderedPosts []PostInfo

// Len - Builtin length method on ordered post slice
func (p OrderedPosts) Len() int {
	return len(p)
}

// Less - Comparator method for ordering posts by date including date
// conversion methods
func (p OrderedPosts) Less(i int, j int) bool {
	dateI, _ := time.Parse(time.RFC3339, generateIsoDate(p[i].Meta.Date.Year, p[i].Meta.Date.Month, p[i].Meta.Date.Day))
	dateJ, _ := time.Parse(time.RFC3339, generateIsoDate(p[j].Meta.Date.Year, p[j].Meta.Date.Month, p[j].Meta.Date.Day))

	return dateI.After(dateJ)
}

// Swap - Swap the values in two given positions in a slice/array
func (p OrderedPosts) Swap(i int, j int) {
	p[i], p[j] = p[j], p[i]
}
