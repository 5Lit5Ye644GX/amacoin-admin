package main

import (
	"fmt"
	"os"
	"os/exec"
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
func IssueMoney(asset string) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res int
	var qt float64

	tableau := GetGlobalAdresses() // Get Addresses
	fmt.Printf("______________________________\nLes adresses disponibles sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}
	fmt.Printf("============================ \n Quelle adresse créditer? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return false
	}
	res1 := tableau[res]

	fmt.Printf("Quelle quantité d'argent créer ?\n")
	_, err2 := fmt.Scanf("%f\n", &qt)
	if err2 != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n Erreur:%s \n", err2)
		return false
	}

	rei, err54 := client.IssueMore(res1, asset, qt)
	if err54 != nil {
		fmt.Printf("Impossible de créer la monnaie sur l'adresse choisie.\n %s \n", rei)
		return false
	}
	return true
}

// GetGlobalAdresses is a function that returns an array of the available adresses
func GetGlobalAdresses() []string {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()

	res, err := client.ListPermissions([]string{"receive"}, []string{}, false)
	if err != nil {
		fmt.Printf("Erreur cli post %s \n", err)
	}

	tabret := make([]string, 0)

	for j := range res.Result().([]interface{}) { // Here we want to extract the addresses
		fmt.Printf(" ===================== \n %d ) ", j) // From the structure in coucou
		plop := res.Result().([]interface{})[j].(map[string]interface{})
		plip := plop["address"].(string)
		tabret = append(tabret, plip) // Adding the addresses
		fmt.Printf("%s \n ==================== \n", tabret[j])
	}
	var input string
	fmt.Scanln(&input)
	return tabret
}
