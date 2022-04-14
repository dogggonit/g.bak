package application

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"time"
)

type Credentials struct {
	gorm.Model
	Credentials string
}

type User struct {
	gorm.Model
	EmailAddress string `gorm:"uniqueIndex"`
	Token string
	TokenExpiry time.Time
	Emails []*Email
}

type Email struct {
	gorm.Model
	GmailId, GmailThreadId, Snippet string
	GmailLabels []*Label `gorm:"many2many:email_labels;"`
	GmailDate, EstimatedSize int64
	Message *Message
	UserID uint
}

type Label struct {
	gorm.Model
	GmailId   string
	Label     string
	LabelType string
	TextColor string
	BackgroundColor string
	UserID uint
}

type Message struct {
	gorm.Model
	GmailId, MimeType, Filename string
	Body *MessageBody
	Headers []*MessageHeader `gorm:"many2many:message_headers_join"`
	Parts []*Message
	MessageID uint
	EmailID uint
}

type MessageBody struct {
	gorm.Model
	MessageID uint
	GmailAttachmentId, Data string
	Size int64
	Attachment *Attachment
}

type MessageHeader struct {
	gorm.Model
	Name, Value string
}

type Attachment struct {
	gorm.Model
	Filename string
	Data []byte
	MessageBodyID uint
}

func (u *User) GetOAuth2Token() (*oauth2.Token, error) {
	var o oauth2.Token
	if err := json.Unmarshal([]byte(u.Token), &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func (l *Label) CorrectSystemLabel() {
	systemGmailLabels := map[string]string{
		"INBOX": "Inbox",
		"SPAM": "Spam",
		"TRASH": "Trash",
		"UNREAD": "Unread",
		"STARRED": "Starred",
		"IMPORTANT": "Important",
		"SENT": "Sent",
		"DRAFT": "Draft",
		"CATEGORY_PERSONAL": "Category/Personal",
		"CATEGORY_SOCIAL": "Category/Social",
		"CATEGORY_PROMOTIONS": "Category/Promotions",
		"CATEGORY_UPDATES": "Category/Updates",
		"CATEGORY_FORUMS": "Forums",
	}

	if n, ok := systemGmailLabels[l.GmailId]; ok {
		l.Label = n
	}
}