package application

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type db gorm.DB

const getEmailByLabelSqlStatementLimit = `
	SELECT * FROM 
		emails 
	WHERE 
		deleted_at is NULL 
			AND
		id in (SELECT email_id FROM
					email_labels 
				WHERE 
			label_id = ?
		)
	ORDER BY gmail_date DESC
	LIMIT ?
	OFFSET ?`

const getEmailByLabelSqlStatement = `
	SELECT count(id) FROM 
		emails 
	WHERE 
		deleted_at is NULL 
			AND
		id in (SELECT email_id FROM
					email_labels 
				WHERE 
			label_id = ?
		)`



type AppDB interface {
	GetEmails(u *User, page, limit uint) ([]*Email, error)
	GetEmailsByLabel(u *User, label string, page, limit uint) (emails []*Email, total int64, err error)
	GetSpecificHeaders(e *Email, headers ...string) (headersMap map[string]string, hasAttachment bool, err error)
	GetCredentials() (*oauth2.Config, error)
	GetCredentialsExists() (bool, error)
	GetDB() *gorm.DB
	GetEntireEmail(u *User, id uint) (*Email, error)
	GetTotalEmails(u *User) (int64, error)
	GetUser(email string) (u *User, ok bool, err error)
	UpdateCredentials(string) error
	UpsertMessage(*gmail.Message, *User, *gmail.Service) (*Email, error)
	UpsertUser(email string, token *oauth2.Token) (*User, error)
}

func NewSqlite3(filename string) (AppDB, error) {
	gormDB, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := updateDatabase(gormDB); err != nil {
		return nil, err
	}

	return (*db)(gormDB), nil
}

func updateDatabase(database *gorm.DB) error {
	if err := database.AutoMigrate(
		&User{},
		&Email{},
		&Message{},
		&MessageBody{},
		&MessageHeader{},
		&Label{},
		&Attachment{},
		&Credentials{},
	); err != nil {
		return err
	}
	return nil
}

func (d *db) GetDB() *gorm.DB {
	return (*gorm.DB)(d)
}

func (d *db) GetCredentialsExists() (bool, error) {
	var c []Credentials
	if err := d.GetDB().Find(&c).Error; err != nil {
		return false, err
	}
	return len(c) > 0, nil
}

func (d *db) GetCredentials() (*oauth2.Config, error) {
	var c Credentials
	if err := d.GetDB().First(&c).Error; err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON([]byte(c.Credentials), gmail.GmailModifyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret to config, %w", err)
	}
	return config, nil
}

func (d *db) UpdateCredentials(credentials string) error {
	if err := d.GetDB().Where("1 = 1").Delete(&Credentials{}).Error; err != nil {
		return err
	}
	return d.GetDB().Create(&Credentials{Credentials: credentials}).Error
}

func (d *db) GetUser(emailAddress string) (*User, bool, error) {
	var u []User
	if err := d.GetDB().Where("email_address = ?", emailAddress).Find(&u).Error; err != nil {
		return nil, false, err
	} else if len(u) > 0 {
		return &u[0], true, nil
	}
	return nil, false, nil
}

func (d *db) UpsertUser(email string, token *oauth2.Token) (*User, error) {
	b, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}

	u, ok, err := d.GetUser(email)
	if err != nil {
		return nil, err
	}

	if !ok {
		u = &User{
			EmailAddress: email,
		}
	}

	u.TokenExpiry = token.Expiry
	u.Token = string(b)

	var result *gorm.DB
	if ok {
		result = d.GetDB().Save(u)
	} else {
		result = d.GetDB().Create(u)
	}

	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (d *db) UpsertMessage(m *gmail.Message, u *User, srv *gmail.Service) (*Email, error) {
	var copyMessagePart func(*gmail.MessagePart)(*Message, error)
	copyMessagePart = func(part *gmail.MessagePart) (*Message, error) {
		e := &Message{
			GmailId:  part.PartId,
			MimeType: part.MimeType,
			Filename: part.Filename,
			Headers:  make([]*MessageHeader, 0, len(part.Headers)),
			Parts:    make([]*Message, 0, len(part.Parts)),
		}

		if part.Body != nil {
			e.Body = &MessageBody{
				GmailAttachmentId: part.Body.AttachmentId,
				Size:              part.Body.Size,
			}

			if e.Filename == "" {
				body, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					return nil, err
				}
				e.Body.Data = string(body)
			} else {
				attId := e.Body.GmailAttachmentId
				attachment, err := srv.Users.Messages.Attachments.Get("me", e.GmailId, attId).Do()
				if err != nil {
					return nil, err
				} else if attachment.HTTPStatusCode != 200 {
					return nil, errors.New("failed to attachments")
				}

				b, err := base64.URLEncoding.DecodeString(attachment.Data)
				if err != nil {
					return nil, err
				}

				e.Body.Attachment = &Attachment{
					Filename:      e.Filename,
					Data:          b,
				}
			}
		}

		for _, p := range part.Parts {
			cp, err := copyMessagePart(p)
			if err != nil {
				return nil, err
			}
			e.Parts = append(e.Parts, cp)
		}

		for _, header := range part.Headers {
			var headers []MessageHeader
			if err := d.GetDB().Where("name = ? and value = ?", header.Name, header.Value).Find(&headers).Error; err != nil {
				return nil, err
			}

			if len(headers) == 0 {
				headers = []MessageHeader{{Name: header.Name, Value: header.Value}}
			}

			e.Headers = append(e.Headers, &headers[0])
		}

		return e, nil
	}

	e := &Email{
		GmailId:       m.Id,
		GmailThreadId: m.ThreadId,
		Snippet:       m.Snippet,
		GmailLabels:   make([]*Label, 0, len(m.LabelIds)),
		GmailDate:     m.InternalDate,
		EstimatedSize: m.SizeEstimate,
		UserID: u.ID,
	}

	var err error
	e.Message, err = copyMessagePart(m.Payload)
	if err != nil {
		return nil, err
	}

	for _, l := range m.LabelIds {
		var labels []Label
		if err := d.GetDB().Where("gmail_id = ? and user_id = ?", l, u.ID).Find(&labels).Error; err != nil {
			return nil, err
		}

		if len(labels) == 0 {
			labelResp, err := srv.Users.Labels.Get("me", l).Do()
			if err != nil {
				return nil, err
			} else if labelResp.HTTPStatusCode != 200 {
				return nil, errors.New("failed to get label")
			}

			labels = []Label{{
				GmailId:         l,
				Label:           labelResp.Name,
				LabelType:       labelResp.Type,
				UserID: u.ID,
			}}

			if labelResp.Color != nil {
				labels[0].TextColor = labelResp.Color.TextColor
				labels[0].BackgroundColor = labelResp.Color.BackgroundColor
			}
		}

		e.GmailLabels = append(e.GmailLabels, &labels[0])
	}

	if err := d.GetDB().Create(&e).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (d *db) GetSpecificHeaders(e *Email, headers ...string) (headersMap map[string]string, hasAttachment bool, err error) {
	m := make(map[string]string)
	if e == nil {
		return m, false, nil
	}

	find := make(map[string]struct{})
	for _, h := range headers {
		find[h] = struct{}{}
	}

	var firstMsg Message
	if err := d.GetDB().Where("email_id = ?", e.ID).Preload("Headers").Preload("Parts").First(&firstMsg).Error; err != nil {
		return nil, false, err
	}

	attachment := false
	var findHeaders func(*Message) error
	findHeaders = func(msg *Message) error {
		attachment = attachment || msg.Filename != ""

		for _, h := range msg.Headers {
			if len(find) == 0  {
				break
			}

			if _, ok := find[h.Name]; ok {
				delete(find, h.Name)
				m[h.Name] = h.Value
			}
		}

		if len(find) == 0 && attachment  {
			return nil
		}

		for _, p := range msg.Parts {
			var part Message
			if err := d.GetDB().Where("id = ?", p.ID).First(&part).Error; err != nil {
				return err
			}

			if err := findHeaders(&part); err != nil {
				return err
			}
		}
		return nil
	}

	if err := findHeaders(&firstMsg); err != nil {
		return nil, attachment, err
	}

	return m, attachment, nil
}

func (d *db) GetEmails(u *User, page, limit uint) ([]*Email, error) {
	var emails []*Email
	if err := d.GetDB().Where("user_id = ?", u.ID).Preload("GmailLabels").Order("gmail_date desc").Limit(int(limit)).Offset(int(page)).Find(&emails).Error; err != nil {
		return nil, err
	}

	for _, e := range emails {
		for _, l := range e.GmailLabels {
			l.CorrectSystemLabel()
		}
	}

	return emails, nil
}

func (d *db) GetTotalEmails(u *User) (total int64, err error) {
	err = d.GetDB().Model(&Email{}).Where("user_id = ?", u.ID).Count(&total).Error
	return
}

func (d *db) GetEntireEmail(u *User, id uint) (*Email, error) {
	var emails []Email
	if err := d.GetDB().Where("user_id = ? and id = ?", u.ID, id).Preload("GmailLabels").Find(&emails).Error; err != nil {
		return nil, err
	}

	if len(emails) == 0 {
		return nil, errors.New("no emails found")
	}

	var getMessages func(emailId, msgId uint) ([]*Message, error)
	getMessages = func(emailId, msgId uint) ([]*Message, error) {
		var messages []*Message
		if err := d.GetDB().Where("message_id = ? and email_id = ?", msgId, emailId).Preload("Body").Preload("Headers").Find(&messages).Error; err != nil {
			return nil, err
		}

		var err error
		for _, m := range messages {
			if m.Body != nil && m.Body.GmailAttachmentId != "" {
				err = d.GetDB().Where("message_body_id = ?", m.Body.ID).First(&m.Body.Attachment).Error
				if err != nil {
					return nil, err
				}
			}

			m.Parts, err = getMessages(0, m.ID)
			if err != nil {
				return nil, err
			}
		}
		return messages, nil
	}

	if m, err := getMessages(emails[0].ID, 0); err != nil {
		return nil, err
	} else {
		emails[0].Message = m[0]
	}

	return &emails[0], nil
}

func (d *db) GetEmailsByLabel(u *User, label string, page, limit uint) ([]*Email, int64, error) {
	var labels []Label
	if err := d.GetDB().Where("gmail_id = ? and user_id = ?", label, u.ID).Find(&labels).Error; err != nil {
		return nil, 0, err
	} else if len(labels) == 0 {
		return nil, 0, errors.New("label does not exist")
	}

	var emails []*Email
	if err := d.GetDB().Raw(getEmailByLabelSqlStatementLimit, labels[0].ID, limit, page).Scan(&emails).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := d.GetDB().Raw(getEmailByLabelSqlStatement, labels[0].ID).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	return emails, total, nil
}