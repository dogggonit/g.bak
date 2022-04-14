package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
	"strings"
)

type GinHelp gin.Context

func GetIdParam(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		ctx.JSON(400, "no id specified")
		ctx.Abort()
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		ctx.JSON(400, "id is invalid")
		ctx.Abort()
		return
	}

	ctx.Set("id", uint(id))
}

func Paginate(c *gin.Context) {
	pageStr := c.Request.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(400, fmt.Errorf("bad page value '%s', %w", pageStr, err).Error())
		c.Abort()
		return
	}

	limitStr := c.Request.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "25"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(400, fmt.Errorf("bad limit value '%s', %w", limitStr, err).Error())
		c.Abort()
		return
	}

	c.Set("limit", uint(math.Max(math.Min(float64(limit), 500), 0)))
	c.Set("page", uint(math.Max(float64(page), 1)))
	fmt.Printf("page %d, limit %d", c.GetUint("page"), c.GetUint("limit"))
}

func GetEmailParam(db AppDB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Param("email")
		if email == "" {
			ctx.JSON(500, "email not specified")
			ctx.Abort()
			return
		}

		user, ok, err := db.GetUser(email)
		if err != nil {
			ctx.JSON(500, "failed to get user from database")
			ctx.Abort()
			return
		} else if !ok {
			ctx.JSON(404, "user not found")
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
	}
}

func (h *GinHelp) GetContext() *gin.Context {
	return (*gin.Context)(h)
}

func (h *GinHelp) GetEmailParam() (*User, bool) {
	u, ok := h.GetContext().Get("user")
	if !ok {
		return nil, false
	}
	user, ok := u.(*User)
	if !ok {
		return nil, false
	}
	return user, true
}

func (h *GinHelp) GetPagination() (page, limit uint) {
	return h.GetContext().GetUint("page"), h.GetContext().GetUint("limit")
}

type Endpoint struct {
	Path string
	PrefixMiddleware []gin.HandlerFunc
	SuffixMiddleware []gin.HandlerFunc
	Get gin.HandlerFunc
	Post gin.HandlerFunc
	Patch gin.HandlerFunc
	Delete gin.HandlerFunc
	Put gin.HandlerFunc
}

func (e *Endpoint) Add(g *gin.RouterGroup) {
	handlers := map[string]gin.HandlerFunc{
		"GET": e.Get,
		"POST": e.Post,
		"PATCH": e.Patch,
		"DELETE": e.Delete,
		"PUT": e.Put,
	}

	var methods []string
	for k, v := range handlers {
		if v == nil {
			delete(handlers, k)
			continue
		}
		methods = append(methods, k)
	}

	cors := func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		for _, t := range []string{"Allow", "Request"} {
			s := fmt.Sprintf("Access-Control-%s-Headers", t)
			c.Writer.Header().Set(s, "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token")
			s = fmt.Sprintf("Access-Control-%s-Methods", t)
			c.Writer.Header().Set(s, strings.Join(methods, ", "))
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		//c.Next()
	}

	g.OPTIONS(e.Path, cors)

	for k, v := range handlers {
		chain := append([]gin.HandlerFunc{cors}, e.PrefixMiddleware...)
		chain = append(chain, v)
		chain = append(chain, e.SuffixMiddleware...)

		switch k {
		case "GET":
			g.GET(e.Path, chain...)
		case "POST":
			g.POST(e.Path, chain...)
		case "PATCH":
			g.PATCH(e.Path, chain...)
		case "DELETE":
			g.DELETE(e.Path, chain...)
		case "PUT":
			g.PUT(e.Path, chain...)
		}
	}
}