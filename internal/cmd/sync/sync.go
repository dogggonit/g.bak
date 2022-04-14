package sync

import (
	"context"
	"flag"
	"fmt"
	"g.bak/internal/pkg/application"
	"google.golang.org/api/gmail/v1"
	"os"
)

func Sync() {
	fs := flag.NewFlagSet("tokens", flag.ExitOnError)
	dbName := fs.String("db", "email.db", "filename for database")
	email := fs.String("e", "", "email to sync")
	_ = fs.Parse(os.Args)

	if *email == "" {
		fmt.Println("email address must be specified")
		os.Exit(1)
	}

	db, err := application.NewSqlite3(*dbName)
	if err != nil {
		fmt.Println("failed to open database")
		os.Exit(1)
	}

	u, ok, err := db.GetUser(*email)
	if err != nil {
		fmt.Println("failed to get user from database")
		os.Exit(1)
	} else if !ok {
		fmt.Printf("user '%s' does not exist", *email)
		os.Exit(1)
	}

	tok, err := u.GetOAuth2Token()
	if err != nil {
		fmt.Println("invalid user token")
		os.Exit(1)
	}

	creds, err := db.GetCredentials()
	if err != nil {
		fmt.Println("failed to get credentials from database")
		os.Exit(1)
	}

	srv, err := application.GetGmailClient(creds, tok)
	if err != nil {
		fmt.Println("failed to create gmail client")
		os.Exit(1)
	}

	page := 0
	message := 0
	err = srv.Users.Messages.List("me").MaxResults(500).Pages(context.TODO(), func(r *gmail.ListMessagesResponse) error {
		page++
		if r.HTTPStatusCode != 200 {
			fmt.Printf("got http error %d on page %d\n", r.HTTPStatusCode, page)
			return nil
		}

		fmt.Printf("syncing page %d\n", page)

		gmailIds := make([]string, 0, len(r.Messages))
		for _, e := range r.Messages {
			gmailIds = append(gmailIds, e.Id)
		}

		var emails []*application.Email
		if err := db.GetDB().Where("gmail_id IN ?", gmailIds).Select("GmailId").Find(&emails).Error; err != nil {
			fmt.Println("failed to get existing emails")
			return err
		}

		gmailIdsMap := make(map[string]struct{})
		for _, e := range emails {
			gmailIdsMap[e.GmailId] = struct{}{}
		}

		for _, id := range gmailIds {
			message++
			if _, ok := gmailIdsMap[id]; ok {
				delete(gmailIdsMap, id)
				continue
			}

			fmt.Printf("saving message %d\n", message)

			msg, err := srv.Users.Messages.Get("me", id).Do()
			if err != nil || msg.HTTPStatusCode != 200 {
				fmt.Println("failed to get message")
				continue
			}

			if _, err := db.UpsertMessage(msg, u, srv); err != nil {
				fmt.Println("failed to insert message")
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("failed to download messages")
		os.Exit(1)
	}

	fmt.Println("done.")
}
