package credentials

import (
	"encoding/json"
	"flag"
	"fmt"
	"g.bak/internal/pkg/application"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"io/ioutil"
	"os"
)

func Credentials() {
	fs := flag.NewFlagSet("credentials", flag.ExitOnError)
	dbName := fs.String("db", "email.db", "filename for database")
	importCred := fs.Bool("i", false, "import credentials")
	exportCred := fs.Bool("e", false, "export credentials")
	credFile := fs.String("f", "credentials.json", "credentials file")
	_ = fs.Parse(os.Args)

	if *importCred && *exportCred {
		fmt.Println("only i or e can be enabled")
		os.Exit(1)
	}

	db, err := application.NewSqlite3(*dbName)
	if err != nil {
		fmt.Println("failed to open database")
		os.Exit(1)
	}

	out := func(c string) {
		fmt.Println(c)
	}

	switch {
	case *importCred:
		b, err := ioutil.ReadFile(*credFile)
		if err != nil {
			fmt.Println("failed to read credentials file")
			os.Exit(1)
		}

		if _, err := google.ConfigFromJSON(b, gmail.GmailModifyScope); err != nil {
			fmt.Println("credentials are invalid")
			os.Exit(1)
		}

		if err := db.UpdateCredentials(string(b)); err != nil {
			fmt.Println("failed to import credentials")
			os.Exit(1)
		}
	case *exportCred:
		out = func(c string) {
			if err := ioutil.WriteFile(*credFile, []byte(c), os.ModePerm); err != nil {
				fmt.Println("failed to export credentials to file")
				os.Exit(1)
			}
		}
		fallthrough
	default:
		if ok, err := db.GetCredentialsExists(); err != nil {
			fmt.Println("failed to check for credentials in database")
			os.Exit(1)
		} else if !ok {
			fmt.Println("database does not contain any credentials")
			os.Exit(0)
		}

		cred, err := db.GetCredentials()
		if err != nil {
			fmt.Println("failed to get credentials from database")
			os.Exit(1)
		}

		b, err := json.Marshal(cred)
		if err != nil {
			fmt.Println("failed to output credentials")
			os.Exit(1)
		}
		out(string(b))
	}

	fmt.Println("done.")
}