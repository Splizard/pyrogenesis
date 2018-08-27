package pyrogenesis

import "os"
import "os/exec"
import "os/user"
import "runtime"

var Command string = "pyrogenesis"

func init() {
	uid, _ := user.Current()
	home := uid.HomeDir
	
	if runtime.GOOS == "windows" {
		Command = home+"/AppData/Local/0 A.D. alpha/binaries/system/pyrogenesis.exe"
	}
	
	env := os.Getenv("PYROGENESIS")
	if env != "" {
		Command = env
	}
}

func Run(arguments ...string) error {
	var command = exec.Command(Command, arguments...)
	
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	
	if runtime.GOOS == "windows" {
		if Public != nil {
			Public.zip.Close()
		}
	}
	err := command.Run()
	
	if runtime.GOOS == "windows" {			
		if Public != nil {
			LoadPublic()
		}
	}
	
	return err
}
