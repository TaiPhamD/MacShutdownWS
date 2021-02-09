package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	//"unsafe"

	"github.com/kardianos/service"
)

type shutdownHandler struct{}
type program struct{}

var logger service.Logger
var rootDir string

type AuthStruct struct {
	Password   *string
	Mode       *uint32
	Ingredient *string
}

type BootItem struct {
	OS     string
	BootID uint16
}

type ConfigData struct {
	Password string
	Port     string
	BootDict []BootItem
}

var myConfig ConfigData

func main() {

	//Get file path from where the exe is launched
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Print(dir)
	rootDir = dir
	//set up log file
	filelog, errlog := os.OpenFile(dir+"/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if errlog != nil {
		log.Fatal(errlog)
	}
	defer filelog.Close()
	log.SetOutput(filelog)

	svcConfig := &service.Config{
		Name:        "ShutdownWebService",
		DisplayName: "ShutdownWebService",
		Description: "Shutdown Webservice Listener",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}

func (p *program) Stop(s service.Service) error {
	log.Print("Stopped Shutdown service\n")
	// Stop should not block. Return with a few seconds.
	return nil
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	log.Print("Starting service from golang")
	go p.run()
	return nil
}

func (p *program) run() {
	//Read in configuration file
	log.Print("Started Shutdown service from: " + rootDir + "\n")
	jsonFile, err := os.Open(rootDir + "/config.json")
	if err != nil {
		log.Fatal(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	//log.Print(byteValue)

	json.Unmarshal(byteValue, &myConfig)

	//log.Print("My password is:", Mypassword)
	var certPath = os.Getenv("WS_PUB_FILE")
	var keyPath = os.Getenv("WS_PRIV_FILE")
	log.Print("Using cert at : " + certPath)
	log.Print("Hosting service on port: " + myConfig.Port)
	err = http.ListenAndServeTLS(":"+myConfig.Port, certPath, keyPath, shutdownHandler{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func (h shutdownHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var jsonAuth AuthStruct
	err := decoder.Decode(&jsonAuth)
	if err != nil {
		log.Print("error decoding JSON\n")
		log.Print(r.Body)
		return
	}

	if *(jsonAuth.Password) != myConfig.Password {
		log.Print("Password from JSON doesn't match\n")
		return
	}
	//log.Print("calling shutdown mode: ", jsonAuth.Mode)
	//log.Print("Ingredient : ", jsonAuth.Ingredient)
	//Call the right shutdown mode based on ingredient and mode

	var effectiveMode = *(jsonAuth.Mode)
	if effectiveMode == 0 {
		//err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
		err = syscall.Exec("/bin/systemctl", []string{"start", "shutdownWS_shutdown.service"}, os.Environ())
		if err != nil {
			log.Print("Failed to initiate shutdown:", err)
		}
		log.Print("sucessfully executed shutdown api!")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful shutdown"))
		return
	} else {

		if jsonAuth.Ingredient == nil {
			log.Print("Error: trying to restart but ingredient from JSON body is missing")
			return
		}

		if myConfig.BootDict == nil {
			log.Print("Error: config.json does not contain boot dictionary")
		}
		for _, s := range myConfig.BootDict {
			//log.Print(s.OS)
			//log.Print(s.BootID)

			if strings.ToUpper(s.OS) == strings.ToUpper(*(jsonAuth.Ingredient)) {
				var nextBoot = strconv.FormatUint(uint64(s.BootID), 10)
				log.Print("attempting to execute restart api and setting bootnext:" + nextBoot)
				//err := syscall.Exec("/usr/bin/sudo", []string{"sudo", "efibootmgr", "-n", nextBoot}, os.Environ())
				err := syscall.Exec("/bin/efibootmgr", []string{"-o", nextBoot,
					"&&", "reboot"}, os.Environ())
				log.Print(err)
				if err != nil {
					log.Print("Unable to change uefi nextboot")
				}

				//RestartFunc.Call(uintptr(unsafe.Pointer(&bootMode)))
				log.Print("sucessfully executed restart api and setting bootnext:" + nextBoot)
				//err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
				err = syscall.Exec("/bin/systemctl", []string{"start", "shutdownWS_reboot.service"}, os.Environ())
				if err != nil {
					log.Print("Unable to restart")
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("successful restart"))
				//return
			}
		}

	}

}
