package debug

import (
	"encoding/json"
	"fmt"
)

// PrettyPrint prints structs in a readable way in the terminal.
func PrettyPrint(i interface{}) {
	pp, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(pp))
}
