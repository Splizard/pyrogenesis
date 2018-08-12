package pyrogenesis

import "github.com/beevik/etree"

import (
	"errors"
	"os"
	"path/filepath"
)

const actors = "art/actors/"

type Actor struct {
	name string
	tree *etree.Document
	mod *Mod
}

func (mod *Mod) LoadActor(name string) (*Actor, error) {
	file, err := mod.Open(actors+name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	tree := etree.NewDocument()
	if _, err = tree.ReadFrom(file); err != nil {
		return nil, err
	}
	
	actor := new(Actor)
	actor.name = name
	actor.tree = tree
	actor.mod = mod
	
	return actor, nil
}

func (actor *Actor) Fork(name string) {
	actor.name = name+".xml"
}

func (actor *Actor) Save() error {
	if actor.mod.zip != nil {
		return errors.New("Cannot write to zip files!")
	}
	
	actor.tree.IndentTabs()
	
	path := ModsPath+actor.mod.Name+"/"+actors+actor.name
	os.MkdirAll(filepath.Dir(path), 0755)
	return actor.tree.WriteToFile(path)
}

func (actor *Actor) Tree() *etree.Document {
	return actor.tree
}


