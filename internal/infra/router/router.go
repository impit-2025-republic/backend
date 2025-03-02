package router

import (
	"b8boost/backend/internal/adapters/api/action"
	"b8boost/backend/internal/adapters/api/middleware"
	"b8boost/backend/internal/adapters/repo"
	"b8boost/backend/internal/infra/ai"
	"b8boost/backend/internal/infra/jwt"
	"b8boost/backend/internal/infra/ldap"
	"b8boost/backend/internal/infra/tgbot"
	"b8boost/backend/internal/usecase"
	"context"
	"fmt"
	"net/http"
	"strings"

	_ "b8boost/backend/docs"

	"github.com/gin-gonic/gin"

	cors "github.com/rs/cors/wrapper/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type RouterHTTP struct {
	ai       ai.Vllm
	jwt      jwt.JWKSHandler
	router   *gin.Engine
	botToken string
	ldap     ldap.LDAP
	db       *gorm.DB
	tgbot    tgbot.TgBot
}

func NewRouterHTTP(
	jwt jwt.JWKSHandler,
	botToken string,
	ldap ldap.LDAP,
	db *gorm.DB,
	ai ai.Vllm,
	tgbot tgbot.TgBot,
) RouterHTTP {
	router := gin.Default()
	return RouterHTTP{
		jwt:      jwt,
		router:   router,
		botToken: botToken,
		ldap:     ldap,
		db:       db,
		ai:       ai,
		tgbot:    tgbot,
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
	r.router.GET("/events/upcoming", r.buildAuthMiddleware(r.jwt), r.GetUpcomingEvents())
	r.router.POST("/events/visit", r.buildAuthMiddleware(r.jwt), r.VisitEventAction())
	r.router.GET("/events/archived", r.GetArchivedEvents())
	r.router.GET("/users/me", r.buildAuthMiddleware(r.jwt), r.GetUserMe())
	r.router.POST("/admin/events/visit", r.AdminVisitEventAction())
	r.router.POST("/llm", r.LLMAction())

	r.router.GET("/users/top", r.GetTopBalance())

	r.router.GET("/products", r.GetProductsAction())
	r.router.POST("/products/buy", r.buildAuthMiddleware(r.jwt), r.BuyProductAction())
	r.router.POST("/products/open/case", r.buildAuthMiddleware(r.jwt), r.CaseOpenAction())
	r.router.GET("/users/transactions", r.buildAuthMiddleware(r.jwt), r.GetUserTransaction())

	r.router.GET("/jwts", r.buildValidateJwts())
}

func getJwtClaimsFromIstio(r *http.Request, jwtClaims jwt.JWKSHandler) *jwt.Claims {
	header := r.Header.Get("Authorization")
	fmt.Println(header)

	var authHeader string
	initData := strings.Split(r.Header.Get("Authorization"), " ")

	if len(initData) == 2 {
		authHeader = initData[1]
	}

	var claimsJwt *jwt.Claims
	if authHeader != "" {
		claims, err := jwtClaims.Verify(authHeader)
		if err != nil {
			return nil
		}
		claimsJwt = claims
		// err := json.Unmarshal([]byte(authHeader), &claims)
		// if err != nil {
		// 	return nil
		// }
	}

	return claimsJwt
}

func (g RouterHTTP) buildAuthMiddleware(jwt jwt.JWKSHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := getJwtClaimsFromIstio(c.Request, jwt)
		// audince := claims.Audience
		// if audince{
		// 	//TODO keycloak
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid audience"})
		// }

		ctx := context.WithValue(c.Request.Context(), middleware.UserIDKey, claims.Subject)
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

// @Summary		visit event
// @Tags			event
// @Security		BearerAuth
// @Produce		json
// @Param			input	body		usecase.VisitEventInput	true	"input"
// @Success		200
// @Failure		500
// @Router			/events/visit [post]
func (r *RouterHTTP) VisitEventAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewVisitEventInteractor(
				repo.NewEventRepo(r.db),
				repo.NewEventUserVisits(r.db),
				repo.NewUserRepo(r.db),
				r.tgbot,
			)
			act = action.NewVisitEventAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		admin visit event
// @Tags			event
// @Security		BearerAuth
// @Produce		json
// @Param			input	body		usecase.AdminVisitEventInput	true	"input"
// @Success		200
// @Failure		500
// @Router			/admin/events/visit [post]
func (r *RouterHTTP) AdminVisitEventAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewAdminVisitEventInteractor(
				repo.NewEventRepo(r.db),
				repo.NewEventUserVisits(r.db),
				repo.NewachievementUserRepo(r.db),
				repo.NewAchievementRepo(r.db),
				repo.NewUserWallet(r.db),
				repo.NewUserWalletHistoryRepo(r.db),
			)
			act = action.NewAdminVisitEventAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		buy product
// @Tags			product
// @Security		BearerAuth
// @Produce		json
// @Param			input	body		usecase.BuyProductInput	true	"input"
// @Success		200
// @Failure		500
// @Router			/products/buy [post]
func (r *RouterHTTP) BuyProductAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewBuyProductInteractor(
				repo.NewProductRepo(r.db),
				repo.NewUserWinningRepo(r.db),
				repo.NewUserWallet(r.db),
				repo.NewUserWalletHistoryRepo(r.db),
			)
			act = action.NewBuyProductAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		case open
// @Tags			product
// @Security		BearerAuth
// @Produce		json
// @Param			input	body		usecase.CaseOpenInput	true	"input"
// @Success		200 {object} usecase.CaseOpenOutput
// @Failure		500
// @Router			/products/open/case [post]
func (r *RouterHTTP) CaseOpenAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewCaseOpenInteractor(
				repo.NewCaseProductProbabilityRepo(r.db),
				repo.NewProductRepo(r.db),
				repo.NewUserWallet(r.db),
				repo.NewUserWalletHistoryRepo(r.db),
				repo.NewUserWinningRepo(r.db),
			)
			act = action.NewCaseOpenAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		chat with llm
// @Tags			llm
// @Security		BearerAuth
// @Produce		json
// @Param			input	body		usecase.LLMChatInput	true	"input"
// @Success		200		{object}	ai.StreamResponse
// @Failure		500
// @Router			/llm [post]
func (r *RouterHTTP) LLMAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewLLmChatInteractor(
				r.ai,
			)
			act = action.NewLLMChatAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get my transaction
// @Tags			user
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.GetMyHistoryWalletOutput
// @Failure		500
// @Router			/users/transactions [get]
func (r *RouterHTTP) GetUserTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewGetMyHistoryWalletInteractor(
				repo.NewUserWalletHistoryRepo(r.db),
			)
			act = action.NewGetMyHistoryWalletAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get upcoming events
// @Tags			event
// @Security		BearerAuth
// @Param			period	query		string	false	"period"
// @Produce		json
// @Success		200		{object}	usecase.UpcomingEventList
// @Failure		500
// @Router			/events/upcoming [get]
func (r *RouterHTTP) GetUpcomingEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewUpcomingEventsInteractor(
				repo.NewEventRepo(r.db),
				repo.NewEventUserVisits(r.db),
			)
			act = action.NewUpcomingEventsAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get top balance
// @Tags			user
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.TopBalanceOutput
// @Failure		500
// @Router			/users/top [get]
func (r *RouterHTTP) GetTopBalance() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewTopBalanceInteractor(
				repo.NewUserWallet(r.db),
			)
			act = action.NewTopBalanceAction(uc)
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
				repo.NewUserWallet(r.db),
				repo.NewEventUserVisits(r.db),
				repo.NewEventRepo(r.db),
				repo.NewUserWinningRepo(r.db),
			)
			act = action.NewUserMeAction(uc)
		)

		act.Execute(c.Writer, c.Request)
	}
}

// @Summary		get products
// @Tags			product
// @Security		BearerAuth
// @Produce		json
// @Success		200		{object}	usecase.FindProductOutput
// @Failure		500
// @Router			/products [get]
func (r *RouterHTTP) GetProductsAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			uc = usecase.NewFindProductInteractor(
				repo.NewProductRepo(r.db),
			)
			act = action.NewFindProductAction(uc)
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
