package main

import (
	"errors"
	"net/http"
	"strings"
)

type Context struct {
	// request  *http.Request
	// response *Response
	// Mo       *Mo
}
type HandlerFunc func(c *Context) error
type Middleware func(HandlerFunc) HandlerFunc
type methodHandlers struct {
	get, post, put, patch, delete, options, head HandlerFunc
}

func newMethodHandler(method string, handler HandlerFunc) *methodHandlers {
	mh := methodHandlers{}
	switch method {
	case http.MethodGet:
		mh.get = handler
	case http.MethodPost:
		mh.post = handler
	case http.MethodPut:
		mh.put = handler
	case http.MethodPatch:
		mh.patch = handler
	case http.MethodDelete:
		mh.delete = handler
	case http.MethodOptions:
		mh.options = handler
	case http.MethodHead:
		mh.head = handler
	}
	return &mh
}

type nodeType uint8
const (
	static nodeType = iota
	param
	wildcard
)

type node struct{
	path string
	kind nodeType
	parent *node
	children []*node
	handlers *methodHandlers
}

// pathPart is a string without any slashes.
// returns the children which matches with the pathPart
func (n *node) next(pathPart string)*node{
	for _,v:=range n.children{
		if pathPart==v.path{
			return v
		}
	}
	return nil
}

type Router struct{
	tree *node
	removeLeadingSlashes bool
}

func newRouter()*Router{
	return &Router{
		tree: &node{
			path: "",
			kind: static,
			parent: nil,
		},
	}
}

func (r *Router)add(path, method string, handler HandlerFunc,)error{
	if !strings.HasPrefix("/",path){
		return errors.New("Path should always have prefix of \"/\"")
	}
	if r.removeLeadingSlashes{
		path=strings.TrimSuffix("/",path)
	}
	pathSplits:=strings.Split(path, "/")[1:] //[1:] due to split returning empty string

	var match *node
	for _,v:=range pathSplits{
		match=r.tree.next(v)
	}
	if match==nil{
		
	}
	return nil
}