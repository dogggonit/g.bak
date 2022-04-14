package application

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"sync"
	"time"
)

type APIEmails struct {
	Total int64
	Emails []APIEmail
}

type APIEmail struct {
	ID uint
	Snippet string
	To string
	From string
	Subject string
	Date int64
	HasAttachment bool
	Labels []APILabel
}

type APIUser struct {
	Email string
	Labels []APILabel
}

type APILabel struct {
	GmailId string
	Label string
	TextColor string
	BackgroundColor string
}

type authReq struct {
	timeoutChan chan struct{}
	urlChan chan string
	tokChan chan string
	code string
	created time.Time
}

var (
	authReqs = map[string]*authReq{}
	authReqsLock = &sync.Mutex{}
)

func init() {
	go func() {
		d, _ := time.ParseDuration("15m")
		t := time.NewTicker(d)
		for {
			select {
			case <-t.C:
				authReqsLock.Lock()
				for _, v := range authReqs {
					if v.created.Add(d).Before(time.Now()) {
						v.close()
					}
				}
				authReqsLock.Unlock()
			}
		}
	}()
}

func getAuthReq() *authReq {
	ar := &authReq{
		timeoutChan: make(chan struct{}, 1),
		urlChan:     make(chan string, 1),
		tokChan:     make(chan string, 1),
		created:     time.Now(),
	}

	c := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()-_=[]{}+;:'\"`\\/?.>,<|~"
	length := 24
	for i := 0; i < length; i++ {
		ar.code += string(c[rand.Intn(len(c))])
	}

	authReqs[ar.code] = ar
	return ar
}

func (ar *authReq) close() {
	delete(authReqs, ar.code)
	ar.timeoutChan<- struct{}{}
	close(ar.tokChan)
	close(ar.urlChan)
	close(ar.timeoutChan)
}

func GetLabels(db AppDB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			var emails []User
			if err := db.GetDB().Find(&emails).Error; err != nil {
				ctx.JSON(500, err.Error())
				ctx.Abort()
				return
			}

			var response []APIUser
			for _, e := range emails {
				r := APIUser{
					Email:  e.EmailAddress,
				}

				var labels []Label
				if err := db.GetDB().Where("user_id = ?", e.ID).Find(&labels).Error; err != nil {
					ctx.JSON(500, err.Error())
					ctx.Abort()
					return
				}

				for _, l := range labels {
					l.CorrectSystemLabel()
					r.Labels = append(r.Labels, APILabel{
						GmailId:         l.GmailId,
						Label:           l.Label,
						TextColor:       l.TextColor,
						BackgroundColor: l.BackgroundColor,
					})
				}

				response = append(response, r)
			}

			ctx.JSON(200, response)
		},
	}
}

func GetEmail(db AppDB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		GetEmailParam(db),
		GetIdParam,
		func(ctx *gin.Context) {
			id := ctx.GetUint("id")
			u, ok := (*GinHelp)(ctx).GetEmailParam()
			if !ok {
				ctx.JSON(404, "user not found")
				ctx.Abort()
				return
			}
			e, err := db.GetEntireEmail(u, id)
			if err != nil {
				ctx.JSON(404, "email not found")
				ctx.Abort()
				return
			}
			ctx.JSON(200, e)
		},
	}
}

func GetListEmails(db AppDB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		GetEmailParam(db),
		Paginate,
		func(ctx *gin.Context) {
			label := ctx.Param("label")
			if label == "" {
				ctx.JSON(400, "label cannot be blank")
				return
			}

			var response APIEmails

			user, ok := (*GinHelp)(ctx).GetEmailParam()
			if !ok {
				ctx.JSON(500, "could not get user")
				ctx.Abort()
				return
			}
			page, limit := (*GinHelp)(ctx).GetPagination()

			emails, total, err := db.GetEmailsByLabel(user, label, page, limit)
			response.Total = total
			if err != nil {
				ctx.JSON(500, "could not get email list")
				ctx.Abort()
				return
			}

			response.Emails = make([]APIEmail, 0, len(emails))
			for _, e := range emails {
				ae := APIEmail{
					ID:      e.ID,
					Snippet: e.Snippet,
					To:      "",
					From:    "",
					Subject: "",
					Date:    e.GmailDate,
				}

				for _, l := range e.GmailLabels {
					ae.Labels = append(ae.Labels, APILabel{
						Label: l.Label,
						TextColor: l.TextColor,
						BackgroundColor: l.BackgroundColor,
						GmailId: l.GmailId,
					})
				}

				if m, att, err := db.GetSpecificHeaders(e, "To", "From", "Subject"); err != nil {
					ctx.JSON(500, "failed to find email headers")
					ctx.Abort()
					return
				} else {
					ae.To = m["To"]
					ae.From = m["From"]
					ae.Subject = m["Subject"]
					ae.HasAttachment = att
				}

				response.Emails = append(response.Emails, ae)
			}

			ctx.JSON(200, response)
		},
	}
}

func GetTokenLink(db AppDB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			if ok, err := db.GetCredentialsExists(); err != nil || !ok {
				ctx.JSON(500, "failed to get credentials")
				return
			}

			creds, err := db.GetCredentials()
			if err != nil {
				ctx.JSON(500, "failed to get credentials")
				return
			}

			t := getAuthReq()
			rsp := struct{
				URL string
				Code string
			}{Code: t.code}

			go func() {
				if tok, err := GetTokenFromWeb(creds, func(url string) {
					t.urlChan <- url
				}, func() (string, error) {
					select {
					case <-t.timeoutChan:
						return "", errors.New("")
					case r := <-t.tokChan:
						return r, nil
					}
				}); err == nil {
					_, _ = AddToken(db, creds, tok)
				}
			}()

			rsp.URL = <-t.urlChan
			ctx.JSON(200, rsp)
		},
	}
}

func PostToken(_ AppDB) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			var t struct{
				Token string
				Code string
			}
			if err := json.NewDecoder(ctx.Request.Body).Decode(&t); err != nil {
				ctx.JSON(400, "token not valid")
				return
			}

			authReqsLock.Lock()
			defer authReqsLock.Unlock()
			if ar, ok := authReqs[t.Code]; ok {
				ar.tokChan<-t.Token
				ar.close()
				return
			}
			ctx.JSON(404, "token code not valid")
		},
	}
}