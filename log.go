package main

import "fmt"

func e2(mes1 string, error2 error) {
	fmt.Println(fmt.Sprintf("Error %s value %s", mes1, error2.Error()))
}

// Use for debugging
func mesHelp(a ...interface{}) {
	fmt.Println(a)
}

func s1(p1 string) string {
	mesHelp(p1)
	return p1
}

func s2(p1 string, p2 string) string {
	s := fmt.Sprintf("%s: %s", p1, p2)
	mesHelp(s)
	return s
}

func s3(p1 string, p2 string, p3 string) string {
	s := fmt.Sprintf("%s %s %s", p1, p2, p3)
	mesHelp(s)
	return s
}

func s4(p1 string, p2 string, p3 string, p4 string) string {
	s := fmt.Sprintf("%s %s %s %s", p1, p2, p3, p4)
	mesHelp(s)
	return s
}
