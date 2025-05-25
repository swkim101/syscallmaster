package syscallmaster_test

import (
	"os"
	"testing"

	"github.com/swkim101/syscallmaster/syscallmaster"
)

func TestTest1(t *testing.T) {
	path := "testfiles/test1"
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	res := syscallmaster.Parse(dat, path)
	if len(res) != 1 {
		t.Fail()
	}
	if res[0].Number != 0 {
		t.Errorf("expect %v got %v", 0, res[0].Number)
	}
	if res[0].Audit != "AUE_NULL" {
		t.Errorf("expect %v got %v", "AUE_NULL", res[0].Audit)
	}
	if res[0].Files != "ALL" {
		t.Errorf("expect %v got %v", "ALL", res[0].Files)
	}
	if res[0].Decl != "{ int nosys(void); }" {
		t.Errorf("expect %v got %v", "{ int nosys(void); }", res[0].Decl)
	}
	if res[0].Comments != "{ indirect syscall }" {
		t.Errorf("expect %v got %v", "{ indirect syscall }", res[0].Comments)
	}
}

func TestFreeBSD(t *testing.T) {
	path := "testfiles/syscalls.master.freebsd"
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	res := syscallmaster.Parse(dat, path)

	present := make([]int, 592)
	for i := range 592 {
		present[i] = 0
	}
	for _, r := range res {
		present[r.Number] = 1
	}
	for i := range 592 {
		if present[i] != 1 {
			t.Errorf("syscall #%v is missing", i)
		}
	}
}

func TestXNU(t *testing.T) {
	path := "testfiles/syscalls.master.xnu"
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	res := syscallmaster.Parse(dat, path)

	present := make([]int, 558)
	for i := range 558 {
		present[i] = 0
	}
	for _, r := range res {
		present[r.Number] = 1
	}
	for i := range 558 {
		if present[i] != 1 {
			t.Errorf("syscall #%v is missing", i)
		}
	}
}
