package chttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/divisions"
	"github.com/MaKcm14/one-team/internal/api/chttp/employees"
	"github.com/MaKcm14/one-team/internal/api/chttp/mw"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Controller struct {
	e   *echo.Echo
	log *slog.Logger
	cfg config.ControllerConfig

	auth           auth.Authenticator
	employeeRouter employees.EmployeeRouter
	divisionRouter divisions.DivisionRouter
}

func New(
	log *slog.Logger,
	cfg config.ControllerConfig,
	authService user.IAuthService,
	employeeService employee.IEmployeeService,
	divisionService division.IDivisionService,
) *Controller {
	session := auth.NewSessionConfig(cfg.AuthCfg)
	return &Controller{
		e:   echo.New(),
		log: log,
		cfg: cfg,
		auth: auth.NewAuthenticator(
			log,
			cfg.AuthCfg,
			session,
			authService,
		),
		employeeRouter: employees.NewEmployeeRouter(
			log,
			session,
			employeeService,
		),
		divisionRouter: divisions.NewDivisionRouter(
			log,
			divisionService,
		),
	}
}

func (c Controller) Run() error {
	c.log.Info("starting the app-controller")

	c.configEndpoints()

	if err := c.e.Start(c.cfg.Socket); err != nil {
		return fmt.Errorf("%w: %s", api.ErrStartController, err)
	}
	return nil
}

func (c Controller) configEndpoints() {
	c.e.Use(
		mw.Recovery(c.log),
		mw.LoggerMW(c.log),
		c.auth.ExtractSessionMW(),
		c.auth.DebugPrintCaches(),
	)

	c.e.RouteNotFound("/*", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusNotFound, server.ErrorResponse{
			Error: fmt.Sprintf("%s: the endpoint is not supported", server.ErrRequestInfo),
		})
	})

	adminGroup := c.e.Group("/admin", c.auth.VerifyAccessTokenMW())
	{
		adminGroup.POST("/signup", c.auth.HandlerSignUp)
	}

	clientGroup := c.e.Group("/client", c.auth.VerifyAccessTokenMW())
	{
		clientGroup.POST("/logout", c.auth.HandlerLogout)
	}

	authGroup := c.e.Group("/auth")
	{
		authGroup.POST("/login", c.auth.HandlerLogin)
		authGroup.POST("/token/refresh", c.auth.HandlerRefresh)
	}

	employeeGroup := c.e.Group("/employee")
	{
		employeeGroup.GET("/get/citizenships", c.employeeRouter.HandlerGetCitizenships)
		employeeGroup.GET("/get/titles", c.employeeRouter.HandlerGetTitles)
		employeeGroup.GET("/get/employee", c.employeeRouter.HandlerGetEmployeeWithFilter)

		employeeGroup.GET("/statistics/citizenship", c.employeeRouter.HandlerCountEmployeeWithCitizenship)
		employeeGroup.GET("/statistics/salary", c.employeeRouter.HandlerCountEmployeesWithSalaryBoundary)

		employeeGroup.POST("/create", c.employeeRouter.HandlerCreateEmployee)

		employeeGroup.PUT("/update", c.employeeRouter.HandlerUpdateEmployee)
		employeeGroup.DELETE("/delete", c.employeeRouter.HandlerDeleteEmployee)
	}

	divisionGroup := c.e.Group("/division")
	{
		divisionGroup.GET("/get/list", c.divisionRouter.HandlerGetDivisions)

		divisionGroup.POST("/create", c.divisionRouter.HandlerCreateDivision)

		divisionGroup.DELETE("/delete", c.divisionRouter.HandlerDeleteDivision)
	}
}
