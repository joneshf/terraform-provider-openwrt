package lucirpc

import "encoding/json"

// Options are the actual UCI options for each section.
// The values can be booleans, integers, lists, and strings.
type Options map[string]json.RawMessage
