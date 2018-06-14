package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/flibustier/multichain-client"
)

func print(msg string) {
	runes := []rune(msg)
	for _, c := range runes {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		fmt.Printf("%c", c)
	}
}

func boom() {
	print(`
		██████╗ ███████╗██╗   ██╗███████╗    ██╗   ██╗██╗   ██╗██╗  ████████╗
		██╔══██╗██╔════╝██║   ██║██╔════╝    ██║   ██║██║   ██║██║  ╚══██╔══╝
		██║  ██║█████╗  ██║   ██║███████╗    ██║   ██║██║   ██║██║     ██║   
		██║  ██║██╔══╝  ██║   ██║╚════██║    ╚██╗ ██╔╝██║   ██║██║     ██║   
		██████╔╝███████╗╚██████╔╝███████║     ╚████╔╝ ╚██████╔╝███████╗██║   
		╚═════╝ ╚══════╝ ╚═════╝ ╚══════╝      ╚═══╝   ╚═════╝ ╚══════╝╚═╝   `)
	fmt.Println()
}

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "windows" {
		log.Println("[OK] Your OS is [Windows]")

		path := user.HomeDir + "\\AppData\\Roaming\\Multichain"
		err = os.Chdir(path)
		if err != nil {
			log.Println(err)
		}
		cmd := exec.Command(".\\multichaind.exe", "Amacoin")
		err = cmd.Start()
		if err != nil {
			fmt.Printf("Est-ce que ça tourne ?: %s \n ", err)
		}

	} else if runtime.GOOS == "linux" {
		log.Println("[OK] Your OS is [Linux]")
	}

	boom()

	// little sleep (3s) before connecting
	time.Sleep(time.Duration(3) * time.Second)
	////////////////////////// Démarrage de multichaind.exe Amacoin@IP:Port
	////////////////////////// Pour se connecter au noeud Papa (pas besoin de IP:Port si on est le noeud papa)

	// Connexion to the holy blockchain hosting the noble écu
	// We need a central node, used as a DNS seed
	///////////////////////// FLAGS TO LAUNCH THE .EXE WITH OPTIONS ////////////////////////////
	chain := flag.String("chain", "Amacoin", "is the name of the chain")
	host := flag.String("host", "localhost", "is a string for the hostname")
	port := flag.Int("port", 4336, "is a number for the host port")
	username := flag.String("username", "multichainrpc", "is a string for the username")
	password := flag.String("password", "DYiL6vb71Y8qfEo9CkYr5wyZ3GqjRxrjzkYyjsA9S1k2", "is a string for the password")
	flag.Parse()

	logs := GetLogins(*chain)
	*username = logs[0]
	*password = logs[1]
	*port = GetPort(*chain)

	///////////////////////// TACTICAL CONNECTION TO THE HOLY BLOCKCHAIN /////////////////////////
	client := multichain.NewClient(
		*chain,
		*username,
		*password,
		*port,
	).ViaNode(
		*host,
		*port,
	)

	//////////////////////// Asset Definition /////////////////////////
	RewardName := *chain  // Nom de notre monnaie.
	InitialReward := 10.0 // Récompense d'entrée.
	cents := 0.01         // Unité monétaire divisionnaire de l'écu.
	///////////////////////////////////////////////////////////////////

	obj, err := client.GetAddresses(false) // Get the addresses in our wallet.
	if err != nil {                        // Impossible to reach our wallet, please ask for lost objects.
		log.Fatal("[FATAL] Could not get addresses from Multichain", err)
	}

	addresses := obj.Result().([]interface{})                                // Different addresses stored on the node
	address := addresses[0].(string)                                         // The first wallet is the principle one. End of discussion
	obj, err = client.Issue(true, address, RewardName, InitialReward, cents) // If it's the first time the node is launched, we have to create the asset for reward

	if err != nil { // Asset already existing
		log.Printf("[OK] Asset %s seems to be already existing\n", RewardName)
	} else { // Creation of the non existing asset
		obj, err = client.IssueMore(address, RewardName, 10) // Noob award ?
		if err != nil {
			log.Println("[ERREUR SUR L'ADRESSE]")
		} else {
			log.Println("[OK] ON A RAJOUTE L'ARGENT") // Award granted
		}
		log.Printf("[OK] Asset %s successfuly created\n", RewardName) // Graphical confirmation of the asset creation's success
	}
	// End of the initialization of da wallet.
	log.Printf("[OK] On a bien démarré notre noeud. La bourse est disponible à l'adresse : %s\n", address)

	//////////////////////////////////////////////////////////
	address1 := Identification(client)
	erreu := ChoiceAdmin(client, RewardName)
	fmt.Printf("%s \n %s \n", address1, erreu)
}

///////////////////////// Fonctions d'administration /////////////

//Identification is a function that asks very basically the user to inform the program his office
func Identification(client *multichain.Client) string {
	var res int
	tableau := GetLocalAddresses(client)
	print("\n ============ I D E N T I F I C A T I O N ============= \n")
	fmt.Printf("Les adresses disponibles sur le noeud sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}
	fmt.Printf("======================================================= \n Quelle adresse correspond à votre bureau? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return ""
	}
	res1 := tableau[res]
	return res1
}

//GetLocalAddresses is a function that return a list of the addresses contained in da wallet.
func GetLocalAddresses(client *multichain.Client) []string {
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

// ChoiceAdmin is a function that open the Menu for admin functions.
func ChoiceAdmin(client *multichain.Client, asset string) error {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res1 int
	print("	 /'\\_/`\\                          	\n")
	print("	/\\      \\     __    ___   __  __  	\n")
	print("	\\ \\ \\__\\ \\  /'__`\\/' _ `\\/\\ \\/\\ \\ 	\n")
	print("	 \\ \\ \\_/\\ \\/\\  __//\\ \\/\\ \\ \\ \\_\\ \\	\n")
	print("	  \\ \\_\\\\ \\_\\ \\____\\ \\_\\ \\_\\ \\____/	\n")
	print("	   \\/_/ \\/_/\\/____/\\/_/\\/_/\\/___/ \n")

	fmt.Printf("\n \n +-----------------------------------------------------+\n | Creer une nouvelle adresse dans le portefeuille 1)  | \n | Crediter une adresse 2)                             | \n | Supprimer les permissions d'une adresse 3)          | \n | Sortie 0)                                           |\n +-----------------------------------------------------+ \n")
	_, err := fmt.Scanf("%d\n", &res1)
	switch res1 {
	case 1: // Create a new address
		err := CreateAddress(client, asset)
		if err != true {
			fmt.Printf("Error in CreateAddress")
			return nil
		}
	case 2: // Issue asset
		err := IssueMoney(client, asset)
		if err != true {
			fmt.Printf("Error in IssueMoney")
			return nil
		}
	case 3: // revoke Permission
		err := RevokePermissions(client)
		if err != true {
			fmt.Printf("Error in IssueMoney")
			return nil
		}
	case 0: // Exit
		//fmt.Println("Exiting...")
		//os.Exit(0)
		return nil
	default:
		fmt.Println("Not an option")
	}
	//fmt.Printf("J'ai rentré %d, il y a erreur : %s \n", res1, err)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return err
	}

	return nil
}

//CreateAddress is a function that creates a new address within the wallet and grant them with the basic permissions
func CreateAddress(client *multichain.Client, name string) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	res, err := client.GetNewAddress()
	if err != nil {
		fmt.Printf("Impossible de créer la nouvelle adresse. \n %s \n", err)
		return false
	}
	permissions := []string{"connect", "send", "receive", "mine"}
	resTr := []string{res.Result().(string)}
	resp, erroer := client.Grant(resTr, permissions)
	if erroer != nil {
		fmt.Printf("Grant denied : \n %s \n", erroer)
	}
	fmt.Printf("Nouvelle adresse créée avec succès. \n %s \n ======================== \n", resp)
	_, err = client.IssueMore(resTr[0], name, 10)
	if err != nil {
		log.Printf("[ERREUR SUR L'ADRESSE]")
	}

	return true
}

//RevokePermissions : Revoke the permissions for an address, (quite the same as deleting)
func RevokePermissions(client *multichain.Client) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	permissions := []string{"connect", "send", "receive", "mine"}
	var res int64

	tableau := GetGlobalAdresses(client) // Get Addresses
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

//IssueMoney is a function that allows to credit some money to an user choosen address.
func IssueMoney(client *multichain.Client, asset string) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res int
	var qt float64

	tableau := GetGlobalAdresses(client) // Get Addresses
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
func GetGlobalAdresses(client *multichain.Client) []string {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	tabret := make([]string, 0)
	params := []interface{}{"receive"}
	msg := client.Command( // It will do the manual command
		"listpermissions", // listpermissions that returns the allowed to receive a transaction
		params,            // Basically all the addresses of the network
	)
	coucou, erre := client.Post(msg)
	if erre != nil {
		fmt.Printf("Erreur cli post %s \n", erre)
	}

	for j := range coucou.Result().([]interface{}) { // Here we want to extract the addresses
		fmt.Printf(" ===================== \n %d ) ", j) // From the structure in coucou
		plop := coucou.Result().([]interface{})[j].(map[string]interface{})
		plip := plop["address"].(string)
		tabret = append(tabret, plip) // Adding the addresses
		fmt.Printf("%s \n ==================== \n", tabret[j])
	}
	var input string
	fmt.Scanln(&input)
	return tabret
}

/////////////////////////// utilitaires fichiers //////////////////////

//GetLogins Is a function that will read the multichain.conf file and returns user login and password.
func GetLogins(chain string) []string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("[FATAL] Could not get user from Multichain", err)
	}

	login := "NULL"    // Case in which we cannot find any login.
	password := "NULL" // Case in which we cannot find any password.
	path4 := "multichain.conf"
	var path string
	if runtime.GOOS == "windows" { //////////////////// PATH DIRECTORY FOR WINDOWS USERS \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
		path1 := user.HomeDir + "\\"
		path2 := "AppData\\Roaming\\Multichain\\"
		path3 := chain + "\\"
		path = path1 + path2 + path3 + path4
	} else { ///////////////////////////// PATH DIRECTORY FOR LINUX MAC ... ///////////////
		path1 := user.HomeDir + "/.multichain/"
		path2 := chain + "/"
		path = path1 + path2 + path4
	}
	inFile, err1 := os.Open(path)

	if err1 != nil {
		log.Fatal("[FATAL] Could not open Multichain path", err1)
	}

	re := regexp.MustCompile("rpcpassword=([a-zA-Z0-9]+)") // Gonna search for those strings followed by alphanumerics symbols
	re1 := regexp.MustCompile("rpcuser=([a-zA-Z0-9]+)")

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile) // Scan the file
	scanner.Split(bufio.ScanLines)      // Scan by Lines
	tableau := make([]string, 0)        // Tableau will store the data
	for scanner.Scan() {                // We read the file line by line
		//%ùfmt.Println(scanner.Text())
		if re.MatchString(scanner.Text()) { // If the line matches the searched string (after the defined string)
			password = re.FindStringSubmatch(scanner.Text())[1] // Get the scanned text
			//fmt.Println(password)
		} else if re1.MatchString(scanner.Text()) {
			login = re1.FindStringSubmatch(scanner.Text())[1]
			//fmt.Println(login)
		}
	}
	tableau = append(tableau, login)
	tableau = append(tableau, password) // Keep tablea growing with the matched strings
	return tableau
}

//GetPort Is a function that will read the params.dat file and returns the default port.
func GetPort(chain string) int {
	user, err := user.Current() // Get user's name
	if err != nil {
		log.Fatal("[FATAL] Could not get user from Multichain", err)
	}

	port := "NULL" // Case in which we cannot find any port.
	path4 := "params.dat"
	var path string
	if runtime.GOOS == "windows" { //////////////////// PATH DIRECTORY FOR WINDOWS USERS \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
		path1 := user.HomeDir + "\\"
		path2 := "AppData\\Roaming\\Multichain\\"
		path3 := chain + "\\"
		path = path1 + path2 + path3 + path4
	} else { ///////////////////////////// PATH DIRECTORY FOR LINUX MAC ... ///////////////
		path1 := user.HomeDir + "/.multichain/"
		path2 := chain + "/"
		path = path1 + path2 + path4
	}
	inFile, err1 := os.Open(path) // Open path

	if err1 != nil {
		log.Fatal("[FATAL] Could not open Multichain params.dat", err)
	}

	re := regexp.MustCompile("default-rpc-port = ([0-9]+)") //We want to get the number after "default-rpc-port = "

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile) // Scanner file
	scanner.Split(bufio.ScanLines)      // Scan by line

	for scanner.Scan() { //We read the file line by line
		//fmt.Println(scanner.Text())
		if re.MatchString(scanner.Text()) { //If it matches
			port = re.FindStringSubmatch(scanner.Text())[1] //Get the matched text
			//fmt.Println(port)
		}
	}
	port1, err := strconv.Atoi(port) //convert to integers.
	return port1
}
