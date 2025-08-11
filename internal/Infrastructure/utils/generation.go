package utils

import "strings"

func (ut *Utils) GenerateSlug(title string) string {
	parts := strings.Fields(title)
	return strings.Join(parts,"-")
}