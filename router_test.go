package mygin

import (
	"fmt"
	"testing"
)

/*func Test_getNodeFromTemplate(t *testing.T) {
	cases := map[string]string{
		"test1":"/hello",
		"test2":"/*",
		"test3":"/hello/:id/:name",
		// "test4":"/hello/:id/name", // should return error
		"test5":"/hello/*",
	}
	want := map[string][][]string{
		"test1":[][]string{[]string{"hello"},[]string{}},
		"test2":[][]string{[]string{},[]string{"*"}},
		"test3":[][]string{[]string{"hello"},[]string{"id","name"}},
		// "test4":[][]string{[]string{"hello"},[]string{}},
		"test5":[][]string{[]string{"hello"},[]string{"*"}},
	}
	for name,input := range cases {
		gotName,gotArgs := getNodeFromTemplate(input)
		want := want[name]
		for i,_ := range gotName{
			if gotName[i] != want[0][i] {
				t.Fatalf("Error: %v",name)
			}
		}
		for i,_ := range gotArgs {
			if gotArgs[i] != want[1][i] {
				t.Fatalf("Error: %v",name)
			}
		}
	}
}*/

func TestArr(t *testing.T) {
	arr := []int{1, 2, 3}
	fmt.Println(arr[:1], arr[1:])
}
