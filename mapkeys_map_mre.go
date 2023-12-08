package main

import (
	"fmt"
	"reflect"
)

type IF interface {
	f1(map[IF]string, string)
}

type T struct {
	Name string
}

func (t *T) f1(m map[IF]string, val string) {
	m[t] = val
}

func (t *T) f2(m map[IF]string, val string) {
	//m[*t] = val
	//DOES NOT COMPILE
	//Cannot use '*t' (type T) as the type IF. Type does not implement 'IF' as the 'f1' method has a pointer receiver
	//Q: How to use pointed value of t as key in map, i.e. the object of 't' hashed?
}

func main() {
	m := make(map[IF]string)
	t1 := T{"name"}
	t2 := T{"name"}

	t1.f1(m, "s1")
	t2.f1(m, "s2")

	for k, v := range m {
		fmt.Printf("%s -> %s\n", k, v)
		fmt.Printf("key type: %s\n", reflect.TypeOf(k).String())
		fmt.Printf("key type: %s\n", reflect.TypeOf(k).Kind())
	}
}
