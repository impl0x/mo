package main

import "strings"

func main1() {
	p1 := [...]string{"/",
		"/health/",
		"/users",
		"/users/:id",
		"/users/:id/posts",
	}
	for _,v:=range p1{
		k:=strings.Split(v, "/")
		for _,l:=range k{
			print(l)
		}
		println(len(k))
	}
}
