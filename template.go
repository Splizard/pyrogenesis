package pyrogenesis

import (
	"github.com/beevik/etree"
	"errors"
	"os"
	"path/filepath"
	"bufio"
	"strings"
)

const templates = "simulation/templates/"

type Template struct {
	name string
	tree *etree.Document
	parent *Template
	mod *Mod
}

func (mod *Mod) LoadTemplate(name string) (*Template, error) {
	file, err := mod.Open(templates+name+".xml")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	tree := etree.NewDocument()
	if _, err = tree.ReadFrom(file); err != nil {
		return nil, err
	}
	
	template := new(Template)
	template.name = name
	template.tree = tree
	template.mod = mod

	if tree.Root().SelectAttr("parent") != nil {
		template.parent, err = mod.LoadTemplate(tree.Root().SelectAttr("parent").Value)
	}
	
	return template, nil
}

func (template *Template) List(path string) (list []string) {
	element := template.tree.Root().FindElement(path)
	
	if element == nil {
		return
	}
	
	for _, child := range element.ChildElements() {
		list = append(list, child.Tag)
	}
	return
}

func (template *Template) Get(path string) (result string) {
	
	//TODO deal with replace and disable.
	if template.parent != nil {
		result = template.parent.Get(path)
	}
	
	element := template.tree.Root().FindElement(path)
	if element == nil {
		return result
	}
	
	return element.Text()
}

func (template *Template) State(path string) (result bool) {

	if template.parent != nil {
		result = template.parent.State(path)
	}
	
	element := template.tree.Root().FindElement(path)
	if element == nil {
		return result
	}
	
	return element.SelectAttr("disable") != nil
}



func (template *Template) Create(path string) {

	if !strings.Contains(path, "/") {
		template.tree.Root().CreateElement(path)
		return
	}
	
	var element *etree.Element
	
	reader := bufio.NewReader(strings.NewReader(path))
	last := ""
	for i:=0; i < strings.Count(path, "/"); i++ {
		dir, _ := reader.ReadString('/')
		if last != "" {
			dir = last+"/"+dir[:len(dir)-1]
		} else {
			dir = dir[:len(dir)-1]
		}

		if element = template.tree.Root().FindElement("./"+dir); element == nil {
			if last == "" {
				element = template.tree.Root().CreateElement(filepath.Base(dir))
			} else {
				element = template.tree.Root().FindElement("./"+last).CreateElement(filepath.Base(dir))
			}
		} 
		last = dir
	}
	
	element.CreateElement(filepath.Base(path))
}

func (template *Template) Set(path string, value string) {
	for {
		element := template.tree.Root().FindElement(path)
		if element == nil {
			template.Create(path)
			continue
		}
		element.SetText(value)
		break
	}
}

func (template *Template) Reset(path string) {
	element := template.tree.Root().FindElement(path)
	if element != nil {
		element.Parent().RemoveChild(element)
	}
}

func (template *Template) Disable(path string) {
	for {
		element := template.tree.Root().FindElement(path)
		if element == nil {
			template.Create(path)
			continue
		}
		element.CreateAttr("disable", "")
		break
	}
}

func (template *Template) Enable(path string) {
	element := template.tree.Root().FindElement(path)
	if element == nil {
		template.Create(path)
		return
	}
	element.RemoveAttr("disable")
}

func (template *Template) Replace(path string) {
	for {
		element := template.tree.Root().FindElement(path)
		if element == nil {
			template.Create(path)
			continue
		}
		element.CreateAttr("replace", "")
		break
	}
}

func (template *Template) Save() error {
	if template.mod.zip != nil {
		return errors.New("Cannot write to zip files!")
	}
	
	template.tree.IndentTabs()
	
	path := ModsPath+template.mod.Name+"/"+templates+template.name+".xml"
	os.MkdirAll(filepath.Dir(path), 0755)
	return template.tree.WriteToFile(path)
}
