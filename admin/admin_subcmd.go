package chadmin

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/NextChapterSoftware/chissl/share/settings"
	"github.com/NextChapterSoftware/chissl/share/utils"
	"github.com/olekukonko/tablewriter"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func fatalError(err *error) {
	if *err != nil {
		log.Fatal(*err)
	}
}

var userAddHelp = `
  Usage: chissl admin adduser [options]

  Options:
    --username, -u  Username for the new user
    --password, -p  Password for the new user
	--addresses,-a  Comma-separated list of regex expressions
  
  Flags:
	--admin, 		Flag to add admin permission to user 
`

func (c *AdminClient) AddUser(args []string) {

	flags := flag.NewFlagSet("adduser", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Print(userAddHelp)
		os.Exit(0)
	}

	var username, password string
	var regexList RegexList
	flags.StringVar(&username, "username", "", "Username for the new user")
	flags.StringVar(&username, "u", "", "Username for the new user")
	flags.StringVar(&password, "password", "", "Password for the new user")
	flags.StringVar(&password, "p", "", "Password for the new user")
	flags.Var(&regexList, "addresses", "Comma-separated list of regex expressions")
	flags.Var(&regexList, "a", "Comma-separated list of regex expressions")
	isAdmin := flags.Bool("admin", false, "")

	err := flags.Parse(args)
	if err != nil {
		flags.Usage()
		log.Fatal("Failed to create user")
	}

	user := &settings.User{
		Name:    username,
		Pass:    password,
		IsAdmin: *isAdmin,
		Addrs:   regexList.expressions,
	}

	err = user.ValidateUser()
	fatalError(&err)

	userJSON, err := user.ToJSON()
	fatalError(&err)

	url, err := url.JoinPath(c.server, "/user")
	fatalError(&err)

	_, err = utils.HttpRequestWithBodyWithBasicAuth(
		http.MethodPost,
		url,
		userJSON,
		c.config.Username,
		c.config.Password,
	)
	fatalError(&err)

	// TODO: Add user to the system (implementation dependent on your user management system)
	log.Printf("Success: User %s was added to chiSSL server", username)
}

var userDelHelp = `
  Usage: chissl admin deluser [options]

  Options:
    --username, -u  Username of the user to delete
`

func (c *AdminClient) DelUser(args []string) {
	flags := flag.NewFlagSet("deluser", flag.ExitOnError)
	username := flags.String("username", "", "Username of the user to delete")
	flags.StringVar(username, "u", "", "Username of the user to delete")
	flags.Usage = func() {
		fmt.Print(userDelHelp)
		os.Exit(0)
	}
	flags.Parse(args)

	if *username == "" {
		flags.Usage()
	}

	url, err := url.JoinPath(c.server, "/user", *username)
	fatalError(&err)

	_, err = utils.HttpRequestNoBodyWithBasicAuth(
		http.MethodDelete,
		url,
		c.config.Username,
		c.config.Password,
	)
	fatalError(&err)

	// TODO: Delete user from the system (implementation dependent on your user management system)
	log.Printf("User %s deleted\n", *username)
}

var userGetHelp = `
  Usage: chissl admin getuser [options]

  Options:
    --username, -u  Username of the user to get
  
  Flags:
	--raw,          Flag to output the raw JSON response
`

func (c *AdminClient) GerUser(args []string) {
	flags := flag.NewFlagSet("getuser", flag.ExitOnError)
	username := flags.String("username", "", "Username of the user to get")
	flags.StringVar(username, "u", "", "Username of the user to get")
	rawOutput := flags.Bool("raw", false, "")

	flags.Usage = func() {
		fmt.Print(userGetHelp)
		os.Exit(0)
	}
	flags.Parse(args)

	if *username == "" {
		flags.Usage()
	}

	url, err := url.JoinPath(c.server, "/user", *username)
	fatalError(&err)

	result, err := utils.HttpRequestNoBodyWithBasicAuth(
		http.MethodGet,
		url,
		c.config.Username,
		c.config.Password,
	)
	fatalError(&err)

	if *rawOutput {
		fmt.Println(result)
		os.Exit(0)
	}
	user := settings.User{}
	err = json.Unmarshal([]byte(result), &user)

	// Create a new table writer and set it to write to Stdout
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Is_Admin", "Addresses"})

	// Add a border around the table
	table.SetBorder(true)
	table.SetAutoWrapText(false)

	var addrStrings []string
	for _, addr := range user.Addrs {
		addrStrings = append(addrStrings, addr.String())
	}
	table.Append([]string{user.Name, fmt.Sprint(user.IsAdmin), joinStrings(addrStrings, ", ")})

	// Render the table
	table.Render()

	// TODO: Delete user from the system (implementation dependent on your user management system)
	//fmt.Printf("User %s deleted\n", *username)
}

var userListHelp = `
  Usage: chissl admin listusers

  Flags:: 
	--raw,     Flag to output the raw JSON response
`

func (c *AdminClient) ListUsers(args []string) {

	flags := flag.NewFlagSet("listusers", flag.ExitOnError)
	rawOutput := flags.Bool("raw", false, "")

	flags.Usage = func() {
		fmt.Print(userGetHelp)
		os.Exit(0)
	}
	flags.Parse(args)

	url, err := url.JoinPath(c.server, "/users")
	if err != nil {
		log.Fatal(err)
	}
	result, err := utils.HttpRequestNoBodyWithBasicAuth(http.MethodGet, url, c.config.Username, c.config.Password)
	if err != nil {
		log.Fatal(err)
	}
	users := []*settings.User{}
	err = json.Unmarshal([]byte(result), &users)
	if len(users) == 0 {
		log.Fatal("no results returned")
	}

	if *rawOutput {
		fmt.Println(result)
		os.Exit(0)
	}

	// Create a new table writer and set it to write to Stdout
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Is_Admin", "Addresses"})

	// Add a border around the table
	table.SetBorder(true)
	table.SetAutoWrapText(false)

	// Iterate over users and append rows to the table
	for _, user := range users {
		// Join all addresses into a single string
		var addrStrings []string
		for _, addr := range user.Addrs {
			addrStrings = append(addrStrings, addr.String())
		}
		table.Append([]string{user.Name, fmt.Sprint(user.IsAdmin), joinStrings(addrStrings, ", ")})
	}

	// Render the table
	table.Render()

}

func joinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
