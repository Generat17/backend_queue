package handler

import (
	"github.com/alexandrevicenzi/go-sse"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "server/docs"
	"server/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes(s *sse.Server) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	//router.Use(cors.Default()) // отключяем CORS политику default config
	router.Use(cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedOrigins:   []string{"http://localhost:3000", "http://digitalqueue.ru", "http://212.113.117.124:3000", "https://digitalqueue.ru", "https://212.113.117.124:3000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"access-control-allow-headers", "access-control-allow-methods", "access-control-allow-origin", "authorization", "content-type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // Swagger
	router.GET("/events/:channel", func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/refresh", h.refresh)
		auth.POST("/logout", h.logout)
		auth.POST("/sign-in/:workstation", h.signInWorkstation)

		employee := auth.Group("/employee", h.userIdentityWorkstation)
		{
			employee.POST("/client", h.getNewClient)
			employee.POST("/confirmClient", h.confirmClient)
			employee.POST("/notComeClient", h.notComeClient)
			employee.POST("/endClient", h.endClient)
			employee.POST("/getStatus", h.getStatus)
			employee.POST("/update", h.updateEmployee)
			employee.POST("/remove", h.removeEmployee)
		}

		responsibility := auth.Group("/responsibility", h.userIdentityWorkstation)
		{
			responsibility.POST("/add", h.addResponsibility)
			responsibility.POST("/update", h.updateResponsibility)
			responsibility.POST("/remove", h.removeResponsibility)
		}

		workstation := auth.Group("/workstation", h.userIdentityWorkstation)
		{
			workstation.POST("/add", h.addWorkstation)
			workstation.POST("/update", h.updateWorkstation)
			workstation.POST("/remove", h.removeWorkstation)
		}

		workstationResponsibility := auth.Group("/workstationResponsibility", h.userIdentityWorkstation)
		{
			workstationResponsibility.POST("/update", h.updateWorkstationResponsibility)
		}

		employeeResponsibility := auth.Group("/employeeResponsibility", h.userIdentityWorkstation)
		{
			employeeResponsibility.POST("/update", h.updateEmployeeResponsibility)
		}
		log := auth.Group("/log")
		{
			log.GET("", h.getClientsLog)
			log.GET("/clear", h.clearLog)
		}

		queue := auth.Group("/queue")
		{
			queue.GET("", h.getQueueAdminList)
			queue.GET("/clear", h.clearQueue)
			queue.GET("/timing", h.getTimingList)
			queue.GET("/email", h.getEmailList)
			queue.GET("/restart", h.restartIdentity)

		}
		email := auth.Group("/email")
		{
			email.GET("", h.getEmailList)
			email.POST("/add", h.addEmail)
			email.POST("/remove", h.removeEmail)

		}

		timing := auth.Group("/timing")
		{
			timing.GET("", h.getTimingList)
			timing.POST("/add", h.addTiming)
			timing.POST("/update", h.updateTiming)
			timing.POST("/remove", h.removeTiming)
			timing.POST("/active", h.activeTiming)
		}

	}

	// доработать
	//api := router.Group("/api", h.userIdentity)
	api := router.Group("/api")
	{
		// api для операций с сотрудниками
		employee := api.Group("/employee")
		{
			employee.GET("", h.getEmployeeList)
			employee.GET("/status/:workstation", h.getEmployeeStatus)
		}

		workstation := api.Group("/workstation")
		{
			workstation.GET("", h.getWorkstationList)
		}

		// api для операций с обязанностями (услугами)
		responsibility := api.Group("/responsibility")
		{
			responsibility.GET("", h.getResponsibilityList)
		}

		// api для операций с очередью (ticket - это элемент массива очереди)
		queue := api.Group("/queue")
		{
			queue.GET("", h.getQueueLists)
			queue.GET(":service", h.addQueueItem)
			queue.GET("status/:workstation", h.getQueueItemStatus)
			queue.GET("quality/:client/:quality", h.updateQuality)
		}
	}

	return router
}
