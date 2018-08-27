package pyrogenesis

import (
	"runtime"
	"os"
	"os/user"
	"archive/zip"
	"io"
	"errors"
)

var ModsPath string
var Public *Mod

func init() {
	uid, _ := user.Current()
	home := uid.HomeDir
	
	if runtime.GOOS == "linux" {
		ModsPath = home+"/.local/share/0ad/mods/"
	}
	
	if runtime.GOOS == "windows" {
		ModsPath = home+"/Documents/My Games/0ad/mods/"
	}
}

type Mod struct {
	Name string
	zip *zip.ReadCloser
}

func LoadPublic() (mod *Mod, err error) {
	if Public != nil {
		return Public, nil
	}
	
	mod = new(Mod)
	mod.Name = "public"
	
	if runtime.GOOS == "linux" {
		mod.zip, err = zip.OpenReader("/usr/share/0ad/data/mods/public/public.zip")
		if err != nil {
		mod.zip, err = zip.OpenReader("/usr/share/games/0ad/data/mods/public/public.zip")	
		}
	}
	
	if runtime.GOOS == "windows" {
		uid, _ := user.Current()
		home := uid.HomeDir
		
		mod.zip, err = zip.OpenReader(home+"/AppData/Local/0 A.D. alpha/binaries/data/mods/public/public.zip")
	}
	
	if Public == nil {
		Public = mod
	} else {
		Public.zip = mod.zip
	}
	
	return
}

func LoadMod(name string) (*Mod, error) {
	mod := new(Mod)
	mod.Name = name
	
	if file, err := os.Open(ModsPath+name+"/mod.json"); err != nil {
		
		if name == "public" {
			return LoadPublic()
		}
		
		return nil, errors.New("No such mod: "+name)
		
	} else {
		file.Close()
		
		return mod, nil
	}
}

//Opens a file from the mod.
func (mod *Mod) Open(path string) (io.ReadCloser, error) {
	if mod.zip != nil {
		for _, f := range mod.zip.File {
			if f.Name == path {
				return f.Open()
			}
		}
	}
	
	if file, err := os.Open(ModsPath+mod.Name+"/"+path); err == nil {
		return file, nil
	} else {
		
		if Public != nil {
			return Public.Open(path)
		}
		
		return nil, err
	}
}
