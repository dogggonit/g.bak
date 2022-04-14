package application

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailClient struct {
	user string
	srv *gmail.Service
}
//
//func (gc *GmailClient) GetEmailAttachmentsRaw(emails ...*gmail.Message) error {
//	for _, e := range emails {
//		if err := gc.getAttachmentsRaw(e.Payload); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (gc *GmailClient) GetEmailAttachments(emails ...*models.Email) error {
//	for _, e := range emails {
//		if err := gc.getAttachments(e.Message); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (gc *GmailClient) getAttachmentsRaw(msg *gmail.MessagePart) error {
//	if msg.Body.AttachmentId != "" {
//		b, err := gc.srv.Users.Messages.Attachments.Get(gc.user, msg.PartId, msg.Body.AttachmentId).Do()
//		if err != nil {
//			return err
//		}
//		msg.Body = b
//	}
//
//	for _, p := range msg.Parts {
//		if err := gc.getAttachmentsRaw(p); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (gc *GmailClient) getAttachments(msg *models.Message) error {
//	if msg.Body.GmailAttachmentId != "" {
//		b, err := gc.srv.Users.Messages.Attachments.Get(gc.user, msg.GmailId, msg.Body.GmailAttachmentId).Do()
//		if err != nil {
//			return err
//		}
//		msg.Body = &models.MessageBody{
//			GmailAttachmentId: msg.Body.GmailAttachmentId,
//			Data:              b.Data,
//			Size:              b.Size,
//		}
//	}
//
//	for _, p := range msg.Parts {
//		if err := gc.getAttachments(p); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//
//func (gc *GmailClient) getEmailsFromRequest(msgFn func(func(*gmail.Message)error)error, lastEmailId string, limit *uint, pageNumber, prevKnownTotal, newKnownTotal int) (emails []*models.Email, err error) {
//	var matchedEmailId, reachedLimit bool
//
//	defer func() {
//		switch {
//		case matchedEmailId:
//			err = nil
//		case reachedLimit:
//			err = nil
//		}
//	}()
//
//	return emails, msgFn(func(msg *gmail.Message) error {
//		if gc.OnMessageDownload != nil {
//			gc.OnMessageDownload(prevKnownTotal, newKnownTotal, pageNumber)
//			prevKnownTotal++
//		}
//
//		if matchedEmailId = msg.Id == lastEmailId; matchedEmailId {
//			return errors.New("matching msg id")
//		}
//		emails = append(emails, CopyMessageToEmail(msg))
//		if *limit >= 0 {
//			*limit--
//			if reachedLimit = *limit == 0; reachedLimit {
//				return errors.New("limit reached")
//			}
//		}
//		return nil
//	})
//}
//

//
//func copyMessagePart(m *gmail.MessagePart) *models.Message {
//}
//
//func copyMessageBody(b *gmail.MessagePartBody) *models.MessageBody {
//	mb := &models.MessageBody{
//		GmailAttachmentId: b.AttachmentId,
//		Data:         b.Data,
//		Size:         b.Size,
//	}
//
//	return mb
//}
//
//func (e *models.Message) GetAttachments() (a []*models.Attachment) {
//	var err error
//	if e.Filename != "" {
//		a = []*models.Attachment{{
//			Filename: e.Filename,
//		}}
//		a[0].Data, err = base64.URLEncoding.DecodeString(e.Body.Data)
//		if err != nil {
//			a = nil
//		}
//	}
//
//	for _, p := range e.Parts {
//		a = append(a, p.GetAttachments()...)
//	}
//	return
//}
//
//func (gc *GmailClient) GetLabel(labelId string) (*models.Label, error) {
//	result, err := gc.srv.Users.Labels.Get(gc.user, labelId).Do()
//	if err != nil {
//		return nil, err
//	}
//	l := &models.Label{
//		GmailId: result.Id,
//		Label:   result.Name,
//		LabelType: result.Type,
//	}
//	if result.Color != nil {
//		l.TextColor = result.Color.TextColor
//		l.BackgroundColor = result.Color.BackgroundColor
//	}
//	return l, nil
//}

//
//// Retrieves a token from a local file.
//func tokenFromFile(file string) (*oauth2.Token, error) {
//	f, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//	tok := &oauth2.Token{}
//	err = json.NewDecoder(f).Decode(tok)
//	if err != nil {
//		return nil, err
//	}
//	return tok, nil
//}

func GetGmailClient(credentials *oauth2.Config, token *oauth2.Token) (*gmail.Service, error)  {
	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(credentials.Client(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Gmail client, %w", err)
	}

	return srv, nil
}

// GetTokenFromWeb Requests a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config, offerURL func(string), offerCode func()(string, error)) (*oauth2.Token, error) {
	offerURL(config.AuthCodeURL("state-token", oauth2.AccessTypeOffline))
	t, err := offerCode()
	if err != nil {
		return nil, err
	}
	tok, err := config.Exchange(context.TODO(), t)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web, %w", err)
	}
	return tok, nil
}

func AddToken(db AppDB, creds *oauth2.Config, tok *oauth2.Token) (string, error) {
	srv, err := GetGmailClient(creds, tok)
	if err != nil {
		return "", errors.New("could not get gmail client")
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil || profile.HTTPStatusCode != 200 {
		return "", errors.New("could not get profile")
	}

	if _, err := db.UpsertUser(profile.EmailAddress, tok); err != nil {
		return "", errors.New("failed to save token into database")
	}

	return profile.EmailAddress, nil
}