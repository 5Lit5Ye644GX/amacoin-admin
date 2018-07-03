package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"

	cool "github.com/fatih/color"
	"github.com/flibustier/multichain-client"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	// InitialReward is the amount of money that will be credited after a new address creation
	InitialReward = 10.0 // Récompense d'entrée.
	// Cents is the subdivision of a Coin, example if 0.01, then 1 Coin = 100 Cents
	cents = 0.01 // Unité monétaire divisionnaire de l'écu.
)

func loading(chain, username, password string, port int) {
	ready := make(chan bool)
	timer := 10
	go connection(chain, username, password, port, ready)
	for i, s := range banner {
		// Check if connection is ok
		select {
		// Something happened in the connection
		case c := <-ready:
			// It seems to be ok
			if c {
				timer = 3
			} else {
				// We got a problem so let's retry
				go connection(chain, username, password, port, ready)
			}
		default:
		}
		print(s, (i+1)*timer)
	}
	fmt.Println()
}

func connection(chain, username, password string, port int, ok chan bool) {
	/////////////// TACTICAL CONNECTION TO THE HOLY BLOCKCHAIN /////////////////
	client = multichain.NewClient(
		chain,
		username,
		password,
	).ViaLocal(
		port,
	)

	_, err := client.GetInfo()
	ok <- (err == nil)
}

var client *multichain.Client

func main() {

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	name := ""
	if runtime.GOOS == "windows" {
		ok("Your OS is [Windows]")

		path := user.HomeDir + "\\AppData\\Roaming\\Multichain"
		err = os.Chdir(path)
		if err != nil {
			log.Println(err)
		}

		name = Identification()
		okf("Congratulations, you chose the blockchain %s\n", name)

		cmd := exec.Command(".\\multichaind.exe", name)
		err = cmd.Start()
		if err != nil {
			fail("Unexpected failure when trying to run multichaind.exe in the AppData/Roaming/Multichain directory")
		}

	} else if runtime.GOOS == "linux" {
		ok("Your OS is [Linux]")
	}

	///////////////////////// FLAGS TO LAUNCH THE .EXE WITH OPTIONS ////////////////////////////
	chain := flag.String("chain", name, "is the name of the chain")
	flag.Parse()

	username, password := GetLogins(*chain)
	port := GetPort(*chain)

	okf("Connection to the Blockchain %s...\n", name)

	loading(*chain, username, password, port)

	//////////////////////// Asset Definition /////////////////////////
	RewardName := *chain // Name of the asset
	///////////////////////////////////////////////////////////////////

	obj, err := client.GetAddresses(false) // Get the addresses in our wallet.
	if err != nil {                        // Impossible to reach our wallet, please ask for lost objects.
		panic("[FATAL] Could not get addresses from Multichain")
	}

	addresses := obj.Result().([]interface{})                                // Different addresses stored on the node
	address := addresses[0].(string)                                         // The first wallet is the principle one. End of discussion
	obj, err = client.Issue(true, address, RewardName, InitialReward, cents) // If it's the first time the node is launched, we have to create the asset for reward

	if err != nil { // Asset already existing
		okf("Asset %s seems to be already existing\n", RewardName)
	} else { // Creation of the non existing asset
		obj, err = client.IssueMore(address, RewardName, 10) // Noob award ?
		if err != nil {
			failf("[ERROR] Could not issue %s address\n", address)
		} else {
			log.Println("[OK] ON A RAJOUTE L'ARGENT") // Award granted
		}
		log.Printf("[OK] Asset %s successfuly created\n", RewardName) // Graphical confirmation of the asset creation's success
	}
	// End of the initialization of da wallet.
	ok(fmt.Sprintf("On a bien démarré notre noeud. La bourse est disponible à l'adresse : %s\n", address))

	//////////////////////////////////////////////////////////

	for {
		menu(RewardName)
		cool.HiCyan("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

///////////////////////// Fonctions d'administration /////////////

// Identification aims to select the Blockchain to use
func Identification() string {
	// We read all the content of the Multichain Directory
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	print("┌────── I D E N T I F I C A T I O N ──────┐\n")
	fmt.Println("│ The following Blockchain are available: │")
	for i, f := range files {
		if f.IsDir() && f.Name()[0] != '.' {
			fmt.Printf("├ ")
			cool.New(cool.FgHiGreen).Printf("%-2d", i)
			fmt.Printf(": ")
			cool.New(cool.FgHiCyan).Printf("%-36s", f.Name())
			fmt.Printf("│\n")
		}
	}
	fmt.Println("└─────────────────────────────────────────┘")
	fmt.Println("Which Blockchain should we use? Please enter the corresponding number...")

	name := ""
	var res int8
	for name == "" {
		_, err = fmt.Scanf("%d\n", &res)
		if err != nil {
			fail("Wrong imput, please try again")
		} else if files[res].IsDir() && files[res].Name()[0] != '.' {
			name = files[res].Name()
		} else {
			fail("Sorry, you cannot use this number, please try again")
		}
	}
	return name
}

func explore() {
	clear()
	addresses := GetLocalAddresses()
	fmt.Println("┌──────────────────────────────────────┬────────────────────┬────────────────────────────────────────────────────────┐")
	fmt.Printf("│%-38s│%-20s│%-56s│\n", " Address", " Amount", " Private Key")
	fmt.Println("├──────────────────────────────────────┼────────────────────┼────────────────────────────────────────────────────────┤")
	for _, address := range addresses {
		priv, _ := client.DumpPrivKey(address)
		amount := GetBalance(address)
		fmt.Printf("│")
		cool.New(cool.FgHiCyan).Printf("%-38s", address)
		fmt.Printf("│")
		amountColor := cool.FgHiRed
		if amount > 0 {
			amountColor = cool.FgHiGreen
		}
		cool.New(amountColor).Printf("%-20.2f", amount)
		fmt.Printf("│")
		cool.New(cool.FgHiMagenta).Printf("%-56s", priv.Result())
		fmt.Printf("│\n")
	}
	fmt.Println("└──────────────────────────────────────┴────────────────────┴────────────────────────────────────────────────────────┘")
}

// menu is a function that open the Menu for admin functions.
func menu(asset string) error {
	var res1 int

	clear()
	print("	 /'\\_/`\\\n")
	print("	/\\      \\     __    ___   __  __\n")
	print("	\\ \\ \\__\\ \\  /'__`\\/' _ `\\/\\ \\/\\ \\\n")
	print("	 \\ \\ \\_/\\ \\/\\  __//\\ \\/\\ \\ \\ \\_\\ \\\n")
	print("	  \\ \\_\\\\ \\_\\ \\____\\ \\_\\ \\_\\ \\____/\n")
	print("	   \\/_/ \\/_/\\/____/\\/_/\\/_/\\/___/\n")

	fmt.Printf(`
┌─────────────────────────────────────────────────────┐
│ 1) Explore                                          │ 
│ 2) Generate a new address                           │
│ 3) Issue more asset to an address                   │
│ 4) Delete address' permissions                      │ 
│ F) Pay respect                                      │
│ `)
	cool.New(cool.FgHiRed).Printf("0) Ragequit")

	fmt.Printf(`                                         │
└─────────────────────────────────────────────────────┘ 
`)

	_, err := fmt.Scanf("%d\n", &res1)
	switch res1 {
	case 1:
		explore()
	case 2: // Create a new address
		err := CreateAddress(asset)
		if err != true {
			panic("Error in CreateAddress")
		}
	case 3: // Issue asset
		IssueMoney(asset)
	case 4: // revoke Permission
		err := RevokePermissions()
		if err != true {
			panic("Error in IssueMoney")
		}
	case 0: // Exit
		cool.Red("Exiting...")
		os.Exit(0)
		return nil
	default:
		fmt.Println("Not an option")
	}

	if err != nil { // SCAN is Not OK
		cool.Red("Wrong imput, please try again.\n")
		return err
	}

	return nil
}

// CreateAddress is a function that creates a new address within the wallet and grant them with the basic permissions
func CreateAddress(name string) bool {
	clear()
	res, err := client.GetNewAddress()
	if err != nil {
		failf("Impossible de créer la nouvelle adresse.\n %s \n", err)
		return false
	}
	addr := []string{res.Result().(string)}

	ok("Nouvelle adresse créée avec succès:")
	cool.Magenta(addr[0])

	permissions := []string{"receive"}
	res, err = client.Grant(addr, permissions)
	if err != nil {
		failf("Grant denied : \n %s \n", err)
	}

	ok("Accréditation accordées avec succès.")

	_, err = client.IssueMore(addr[0], name, 10)
	if err != nil {
		fail("[ERREUR SUR L'ADRESSE]")
	}

	ok("Pourboire versé avec succès.")

	// Save the Address in a QR Code
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex) + string(os.PathSeparator)
	addressPath := exPath + "address"
	if _, err := os.Stat(addressPath); os.IsNotExist(err) {
		os.MkdirAll(addressPath, 0700)
	}

	path := addressPath + string(os.PathSeparator) + addr[0] + ".png"
	err = qrcode.WriteColorFile(addr[0], qrcode.Medium, 256, color.Black, color.White, path)
	if err != nil {
		fail("Failed to save QR code")
	} else {
		okf("QR Code %s généré avec succès.\n", path)
	}

	// Save the private key in a QR Code
	privPath := exPath + "private"
	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		os.MkdirAll(privPath, 0700)
	}

	// We need to dump priv key
	res, err = client.DumpPrivKey(addr[0])
	if err != nil {
		failf("Impossible de dumper la clef privée de %s (%s)\n", addr[0], err)
	}

	url := fmt.Sprintf("http://145.239.59.99/%s$%s", addr[0], res.Result())

	ok(url)

	path = privPath + string(os.PathSeparator) + addr[0] + ".png"
	err = qrcode.WriteColorFile(url, qrcode.High, 1024, color.White, color.RGBA{239, 139, 27, 255}, path)
	if err != nil {
		fail("Failed to save QR code")
	} else {
		okf("QR Code %s généré avec succès.\n", path)
	}

	return true
}

/////////////////////////// utilitaires fichiers //////////////////////

//GetLogins Is a function that will read the multichain.conf file and returns user login and password.
func GetLogins(chain string) (string, string) {
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
	return login, password
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
