package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	// "os"
	"TDMD/lib"
)

const (
	module="Table Definition Markdown"
	version="0.1.0"
)


func main(){


	log.Infof("TDMD Version : %s",version)
	var Conf lib.Configure
	// Read Config
	fmt.Printf("Endpoint : ")
	nHost, _ := fmt.Scanf("%s",&Conf.Endpoint)
	if nHost == 0 {
		Conf.Endpoint = "localhost"
	}

	fmt.Printf("Port(3306): ")
	nPort, _ :=fmt.Scanf("%s",&Conf.Port)
	if nPort == 0 {
		Conf.Port = 3306
	}

	fmt.Printf("User : ")
	nUser, _ :=fmt.Scanf("%s",&Conf.User)
	if nUser == 0 {
		Conf.User = "root"
	}
	
	fmt.Printf("Pass : ")
	nPass, _ := fmt.Scanf("%s",&Conf.Pass)
	if nPass == 0 {
		Conf.Pass = ""
	}
	
	var Export lib.Export
	
	fmt.Printf("File Path (Default: ./ ): ")
	nPath, _ := fmt.Scanf("%s",&Export.FilePath)
	if nPath == 0 {
		// Export.FilePath, _ = os.Getwd()
		Export.FilePath = "./"
	}

	
	var vDB string
	fmt.Printf("Database name : ")
	nDB, _ := fmt.Scanf("%s",&vDB)
	if nDB == 0 {
		log.Errorf("invaild Database, but All Database Export.")
		os.Exit(1)
	} else {
		Export.DB = strings.Split(vDB,",")
	}
	
	Export.DBcnt = len(Export.DB)
	if Export.DBcnt == 1 {
		var vTable string
		fmt.Printf("Table name : ")
		nTable, _ := fmt.Scanf("%s",&vTable)
		if nTable == 0 {
			log.Errorf("invaild Table")
		} else {
			Export.Table = strings.Split(vTable,",")
		}
			
	}



	// Get Definition
	x := lib.GetDefinition(Conf,Export)

	// Write Markdown
	lib.MakeMarkdown(Export.FilePath,x)
}