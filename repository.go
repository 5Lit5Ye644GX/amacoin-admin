package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	cool "github.com/fatih/color"
)

// RevokePermissions : Revoke the permissions for an address, (quite the same as deleting)
func RevokePermissions() bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	permissions := []string{"connect", "send", "receive", "mine"}
	var res int64

	tableau := GetGlobalAdresses() // Get Addresses
	/*fmt.Printf("______________________________\nLes adresses disponibles sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}*/
	fmt.Printf("============================ \n Quelle adresse créditer? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return false
	}
	res1 := tableau[res]
	resTr := make([]string, 0)
	resTr = append(resTr, res1)
	resp, erroer := client.Revoke(resTr, permissions)
	if erroer != nil {
		fmt.Printf("Revoke denied : \n %s \n", erroer)
	}
	fmt.Printf("Nouvelle adresse révoquée avec succès. \n %s \n ======================== \n", resp)
	return true
}

// IssueMoney is a function that allows to credit some money to an user choosen address.
func IssueMoney(assetName string) {
	clear()
	var res int
	var qt float64

	addresses := GetGlobalAdresses() // Get Addresses

	fmt.Println("┌───┬──────────────────────────────────────┐")
	fmt.Printf("│%-3s│%-38s│\n", "No.", " Available addresses")
	fmt.Println("├───┼──────────────────────────────────────┤")
	for i, address := range addresses {
		fmt.Printf("│")
		cool.New(cool.FgHiGreen).Printf("%-3d", i)
		fmt.Printf("│")
		cool.New(cool.FgHiCyan).Printf("%-38s", address)
		fmt.Printf("│\n")
	}
	fmt.Println("└───┴──────────────────────────────────────┘")

	fmt.Printf("Which address do you want to issue? Please input the corresponding number...\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		failf("Wrong imput, please try again.\n")
	}
	res1 := addresses[res]

	fmt.Printf("How much do you want to issue?\n")
	_, err = fmt.Scanf("%f\n", &qt)
	if err != nil { // SCAN is Not OK
		failf("Wrong imput, please try again.\n Erreur:%s \n", err)
	}

	_, err = client.IssueMore(res1, assetName, qt)
	if err != nil {
		failf("Cannot issue more asset on the chosen address.\n %s \n", err)
	}

	ok("The address was successfuly credited!")
}

// GetBalance returns the asset quantity for address
func GetBalance(address string) float64 {
	res, err := client.GetAddressBalances(address)
	if err == nil {
		if len(res.Result().([]interface{})) > 0 {
			return res.Result().([]interface{})[0].(map[string]interface{})["qty"].(float64)
		}
	}
	return 0
}

// GetLocalAddresses is a function that return a list of the addresses contained in da wallet.
func GetLocalAddresses() []string {
	obj, err := client.GetAddresses(false) // Get the addresses in our wallet.
	if err != nil {                        // Impossible to reach our wallet, please ask for lost objects.
		log.Fatal("[FATAL] Could not get addresses from Multichain", err)
	}
	addresses := obj.Result().([]interface{}) // Different addresses stored on the node
	adresses := make([]string, 0)
	for i := range addresses {
		adresses = append(adresses, addresses[i].(string))
	}
	return adresses
}

// GetGlobalAdresses is a function that returns an array of the available adresses
func GetGlobalAdresses() []string {

	res, err := client.ListPermissions([]string{"receive"}, []string{}, false)
	if err != nil {
		fmt.Printf("Cannot read addresses %s \n", err)
	}

	addresses := make([]string, 0)

	for _, obj := range res.Result().([]interface{}) { // Here we want to extract the addresses
		address := obj.(map[string]interface{})["address"].(string)
		addresses = append(addresses, address)
	}

	return addresses
}
