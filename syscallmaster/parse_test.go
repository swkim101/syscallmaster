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
	if res[0].Name != "nosys" {
		t.Errorf("expect %v got %v", "nosys", res[0].Name)
	}
	if len(res[0].Args) != 1 {
		t.Errorf("failed to parse args")
	}
	if res[0].Args[0] != "void" {
		t.Errorf("expect %v got %v", "void", res[0].Args[0])
	}
	if res[0].Comments != "{ indirect syscall }" {
		t.Errorf("expect %v got %v", "{ indirect syscall }", res[0].Comments)
	}
}

func TestObsol(t *testing.T) {
	path := "testfiles/obsol"
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	res := syscallmaster.Parse(dat, path)
	if len(res) != 1 {
		t.Fail()
	}
	if res[0].Number != 11 {
		t.Errorf("expect %v got %v", 11, res[0].Number)
	}
	if res[0].Audit != "AUE_NULL" {
		t.Errorf("expect %v got %v", "AUE_NULL", res[0].Audit)
	}
	if res[0].Files != "OBSOL" {
		t.Errorf("expect %v got %v", "OBSOL", res[0].Files)
	}
	if res[0].Decl != "execv" {
		t.Errorf("expect %v got %v", "execv", res[0].Decl)
	}
	if res[0].Name != "execv" {
		t.Errorf("expect %v got %v", "execv", res[0].Name)
	}
	if len(res[0].Args) != 0 {
		t.Errorf("failed to parse args")
	}
	if res[0].Comments != "" {
		t.Errorf("expect %v got %v", "", res[0].Comments)
	}
}

var mmapdecl = `{
		void *mmap(
		    _In_ void *addr,
		    size_t len,
		    int prot,
		    int flags,
		    int fd,
		    int pad,
		    off_t pos
		);
	}`

func TestPointer(t *testing.T) {
	path := "testfiles/pointer"
	dat, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	res := syscallmaster.Parse(dat, path)
	if len(res) != 1 {
		t.Fail()
	}
	if res[0].Number != 197 {
		t.Errorf("expect %v got %v", 197, res[0].Number)
	}
	if res[0].Audit != "AUE_MMAP" {
		t.Errorf("expect %v got %v", "AUE_MMAP", res[0].Audit)
	}
	if res[0].Files != "COMPAT6|CAPENABLED" {
		t.Errorf("expect %v got %v", "COMPAT6|CAPENABLED", res[0].Files)
	}
	if len(res[0].Args) != 7 {
		t.Errorf("failed to parse args")
	}
	if res[0].Args[0] != "_In_ void *addr" {
		t.Errorf("expect %v got %v", "_In_ void *addr", res[0].Args[0])
	}
	if res[0].Args[4] != "int fd" {
		t.Errorf("expect %v got %v", "int fd", res[0].Args[4])
	}
	if res[0].Decl != mmapdecl {
		t.Errorf("expect %v got %v", mmapdecl, res[0].Decl)
	}
	if res[0].Name != "mmap" {
		t.Errorf("expect %v got %v", "mmap", res[0].Name)
	}
	if res[0].Comments != "" {
		t.Errorf("expect %v got %v", "", res[0].Comments)
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
