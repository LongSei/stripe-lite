package v1

// RawSpecData returns the decoded raw bytes of the embedded OpenAPI specification.
func RawSpecData() ([]byte, error) {
	return rawSpec()
}
