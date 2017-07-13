// Package main is main package of dockzen agent.
// Main function calls WI_init function in webinterface package.
package main

import (
	"webinterface"	
	"log"

)

func main() {
	log.Printf("Container-Service Agent starting")

	log.Printf("WI init function !!!")
	webinterface.WI_init()

}
