package tokens

import (
	"flag"
	"fmt"
	"g.bak/internal/pkg/application"
	"os"
)

func Tokens() {
	fs := flag.NewFlagSet("tokens", flag.ExitOnError)
	dbName := fs.String("db", "email.db", "filename for database")
	_ = fs.Parse(os.Args)

	db, err := application.NewSqlite3(*dbName)
	if err != nil {
		fmt.Println("failed to open database")
		os.Exit(1)
	}

	if ok, err := db.GetCredentialsExists(); err != nil {
		fmt.Println("failed to verify credentials in the database")
		os.Exit(1)
	} else if !ok {
		fmt.Println("no credentials in database")
		os.Exit(0)
	}

	creds, err := db.GetCredentials()
	if err != nil {
		fmt.Println("failed to get credentials from database")
		os.Exit(1)
	}

	tok, err := application.GetTokenFromWeb(creds, func(s string) {
		fmt.Printf("Go to '%s' and paste the code\n", s)
	}, func() (s string, err error) {
		if _, err = fmt.Scan(&s); err != nil {
			return
		}
		return
	})
	if err != nil {
		fmt.Println("could not get token")
		os.Exit(1)
	}

	email, err := application.AddToken(db, creds, tok)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("added '%s'\n", email)
	fmt.Println("done")
}
