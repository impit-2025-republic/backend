package router

import (
	"b8boost/backend/internal/adapters/api/action"
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/adapters/repo"
	"b8boost/backend/internal/infra/jwt"
	"b8boost/backend/internal/infra/ldap"
	"b8boost/backend/internal/usecase"
	"context"
	"encoding/json"
	"net/http"

	_ "b8boost/backend/docs"

	"github.com/gin-gonic/gin"

	cors "github.com/rs/cors/wrapper/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type RouterHTTP struct {
	jwt      jwt.JWKSHandler
	router   *gin.Engine
	botToken string
	ldap     ldap.LDAP
	db       *gorm.DB
}

func NewRouterHTTP(
	jwt jwt.JWKSHandler,
	botToken string,
	ldap ldap.LDAP,
	db *gorm.DB,
) RouterHTTP {
	router := gin.Default()
	return RouterHTTP{
		jwt:      jwt,
		router:   router,
		botToken: botToken,
		ldap:     ldap,
		db:       db,
	}
}

func (r *RouterHTTP) Listen() {
	r.SetupRoutes()

	c := cors.New(cors.Options{
		// AllowedOrigins:      []string{"*"}, //TODO Edit
		AllowPrivateNetwork: true,
		AllowedMethods:      []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPatch},
		AllowCredentials:    true,
	})

	r.router.Use(c)

	r.router.Static("/docs", "./docs")
	r.router.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler, ginSwagger.URL("/docs/swagger.yaml")))
	r.router.Run(":8080")
}

func (r *RouterHTTP) SetupRoutes() {
	r.router.POST("/login", r.Login())
	r.router.GET("/events/upcoming", r.GetUpcomingEvents())
	r.router.GET("/events/archived", r.GetArchivedEvents())
	r.router.GET("/users/me", r.buildAuthMiddleware(), r.GetUserMe())

	r.router.GET("/jwts", r.buildValidateJwts())
}

func getJwtClaimsFromIstio(r *http.Request) map[string]interface{} {
	// Istio передает валидированные данные в заголовке X-JWT-Claims
	jwtClaimsHeader := r.Header.Get("X-JWT-Claims")

	var claims map[string]interface{}
	if jwtClaimsHeader != "" {
		// Парсим JSON из заголовка
		err := json.Unmarshal([]byte(jwtClaimsHeader), &claims)
		if err != nil {
			return nil
		}
	}

	return claims
}

func (g RouterHTTP) buildAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := getJwtClaimsFromIstio(c.Request)
		audince := claims["aud"].(string)
		if audince != "api-audience" {
			//TODO keycloak
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid audience"})
		}

		ctx := context.WithValue(c.Request.Context(), middleware.UserIDKey, claims["sub"])
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// @Summary		login with telegram
// @Tags			auth
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.LoginOutput
// @Failure		500
// @Router			/login [post]
func (r *RouterHTTP) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewLoginInteractor(
				r.botToken,
				r.jwt,
				r.ldap,
				repo.NewUserRepo(r.db),
			)
			act = action.NewLoginAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get upcoming events
// @Tags			event
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.UpcomingEventList
// @Failure		500
// @Router			/events/upcoming [get]
func (r *RouterHTTP) GetUpcomingEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewUpcomingEventsInteractor(
				repo.NewEventRepo(r.db),
			)
			act = action.NewUpcomingEventsAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get user me
// @Tags			user
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.UserMeOutput
// @Failure		500
// @Router			/users/me [get]
func (r *RouterHTTP) GetUserMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewUserMeInteractor(
				repo.NewUserRepo(r.db),
			)
			act = action.NewUserMeAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get archived events
// @Tags			event
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.ClosedEventsOutput
// @Failure		500
// @Router			/events/archived [get]
func (r *RouterHTTP) GetArchivedEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewClosedEventsInteractor(
				repo.NewEventRepo(r.db),
			)
			act = action.NewClosedEventsAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

func (r *RouterHTTP) buildValidateJwts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc  = usecase.NewFindJwtsInteractor(r.jwt)
			act = action.NewFindJwtsAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}
