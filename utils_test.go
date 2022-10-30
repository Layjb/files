package files

import (
	"fmt"
	"testing"
)

func TestLoadCommonArg(t *testing.T) {
	content, err := LoadCommonArg("1.dat1")
	if err != nil {
		println(err.Error())
	}
	fmt.Println(string(content))
}
