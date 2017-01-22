package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
)

func main() {
	if err := _main(); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

type iterspec struct {
	Name      string
	Interface bool
}

func _main() error {
	var iters = []iterspec{
		{Name: "Argument"},
		{Name: "Directive"},
		{Name: "Definition", Interface: true},
		{Name: "NamedType", Interface: true},
		{Name: "Selection", Interface: true},
		{Name: "Type", Interface: true},
		{Name: "VariableDefinition"},
		{Name: "ObjectDefinition", Interface: true},
		{Name: "ObjectField"},
		{Name: "ObjectFieldDefinition"},
		{Name: "EnumElementDefinition"},
		{Name: "InterfaceFieldDefinition"},
		{Name: "InputFieldDefinition"},
		{Name: "ObjectFieldArgumentDefinition"},
	}

	if err := genIterators(iters, "model/iterators.go"); err != nil {
		return err
	}
	return nil
}

func genIterators(iters []iterspec, dstfn string) error {
	var buf bytes.Buffer

	buf.WriteString("package model")
	buf.WriteString("\n\n// Auto-generated by internal/cmd/geniters/geniters.go. DO NOT EDIT")

	for _, iter := range iters {
		buf.WriteString("\n\ntype ")
		buf.WriteString(iter.Name)
		buf.WriteString("List []")
		if !iter.Interface {
			buf.WriteByte('*')
		}
		buf.WriteString(iter.Name)

		buf.WriteString("\n\nfunc (l *")
		buf.WriteString(iter.Name)
		buf.WriteString("List) Add(list ...")
		if !iter.Interface {
			buf.WriteByte('*')
		}
		buf.WriteString(iter.Name)
		buf.WriteString(") {")
		buf.WriteString("\n*l = append(*l, list...)")
		buf.WriteString("\n}")

		buf.WriteString("\n\nfunc (v ")
		buf.WriteString(iter.Name)
		buf.WriteString("List) Iterator() chan ")
		if !iter.Interface {
			buf.WriteByte('*')
		}
		buf.WriteString(iter.Name)
		buf.WriteString(" {")
		buf.WriteString("\nch := make(chan ")
		if !iter.Interface {
			buf.WriteByte('*')
		}
		buf.WriteString(iter.Name)
		buf.WriteString(", len(v))")
		buf.WriteString("\nfor _, e := range v {")
		buf.WriteString("\nch<-e")
		buf.WriteString("\n}")
		buf.WriteString("\nclose(ch)")
		buf.WriteString("\nreturn ch")
		buf.WriteString("\n}")
	}

	b, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("%s\n", buf.Bytes())
		return err
	}

	f, err := os.Create(dstfn)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write(b)
	return nil
}
