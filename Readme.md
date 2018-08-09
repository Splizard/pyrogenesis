# Pyrogenesis library for Go

Pyrogenesis is an open-source 3D engine geared towards creating Realtime simulations (RTS).

This is a library for the Go programming language designed to assist developers who are creating software using the pyrogenisis engine.

Pyrogenesis works on an Entity model where everything is an entity made up of different components. Each enitity can also have a visual representation called an Actor.

This library allows trivial loading and modification of Entity template files, Actor files and more.

Example:
```
	package main
	
import "github.com/Splizard/pyrogenesis"
import "fmt"

func main() {	
	public, err := pyrogenesis.LoadMod("public")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	template, err := public.LoadTemplate("template_unit")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Identity:")
	
	//This will print out each subelement name of identity!
	for _, element := range template.List("Identity") {
		fmt.Println("\t", element)
	}
	
	fmt.Println()

	fmt.Println(template.Get("Identity/GenericName"))	
	template.Set("Identity/GenericName", "Unit Template")
	
	err = template.Save()
	if err != nil {
		fmt.Println("Error saving template: ", err)
		return
	}
}
```
