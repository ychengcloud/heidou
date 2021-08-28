// Copyright 2019-present Facebook Inc. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package heidou

import (
	"go/token"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var (
	Funcs = template.FuncMap{
		"ToUpper":  strings.ToUpper,
		"ToLower":  strings.ToLower,
		"receiver": receiver,
		"snake":    snake,
		"pascal":   pascal,
		"camel":    camel,
	}
)

// receiver returns the receiver name of the given type.
//
//	[]T       => t
//	[1]T      => t
//	User      => u
//	UserQuery => uq
//
func receiver(s string) (r string) {
	// Trim invalid tokens for identifier prefix.
	s = strings.Trim(s, "[]*&0123456789")
	parts := strings.Split(strcase.ToSnake(s), "_")
	min := len(parts[0])
	for _, w := range parts[1:] {
		if len(w) < min {
			min = len(w)
		}
	}

	//TODO 重复检测
	s = parts[0][:1]
	for _, w := range parts[1:] {
		s += w[:1]
	}

	name := strings.ToLower(s)
	if token.Lookup(name).IsKeyword() {
		name = "_" + name
	}
	return name
}

func snake(s string) string {
	return strcase.ToSnake(s)
}

func pascal(s string) string {
	return strcase.ToCamel(s)
}

func camel(s string) string {
	return strcase.ToLowerCamel(s)
}
