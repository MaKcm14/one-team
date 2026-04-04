package chttp

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/MaKcm14/one-team/internal/api"
	"github.com/MaKcm14/one-team/internal/api/chttp/admin"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth"
	"github.com/MaKcm14/one-team/internal/api/chttp/auth/token"
	"github.com/MaKcm14/one-team/internal/api/chttp/divisions"
	"github.com/MaKcm14/one-team/internal/api/chttp/employees"
	"github.com/MaKcm14/one-team/internal/api/chttp/mw"
	"github.com/MaKcm14/one-team/internal/api/chttp/server"
	"github.com/MaKcm14/one-team/internal/config"
	"github.com/MaKcm14/one-team/internal/services/usecase/division"
	"github.com/MaKcm14/one-team/internal/services/usecase/employee"
	"github.com/MaKcm14/one-team/internal/services/usecase/root"
	"github.com/MaKcm14/one-team/internal/services/usecase/user"
)

type Controller struct {
	e   *echo.Echo
	log *slog.Logger
	cfg config.ControllerConfig

	auth           auth.Authenticator
	employeeRouter employees.EmployeeRouter
	divisionRouter divisions.DivisionRouter
	adminRouter    admin.AdminRouter
}

func New(
	log *slog.Logger,
	cfg config.ControllerConfig,
	authService user.IAuthService,
	rootService root.IRootService,
	employeeService employee.IEmployeeService,
	divisionService division.IDivisionService,
) *Controller {
	session := auth.NewSessionConfig(cfg.AuthCfg)
	tokenStorage := token.NewTokenStorage()
	return &Controller{
		e:   echo.New(),
		log: log,
		cfg: cfg,
		auth: auth.NewAuthenticator(
			log,
			cfg.AuthCfg,
			session,
			tokenStorage,
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
		adminRouter: admin.NewAdminRouter(
			log,
			tokenStorage,
			session,
			rootService,
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
		c.auth.DebugPrintCaches(),
	)

	c.e.RouteNotFound("/*", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusNotFound, server.ErrorResponse{
			Error: fmt.Sprintf("%s: the endpoint is not supported", server.ErrRequestInfo),
		})
	})

	c.e.POST("/init", c.auth.HandlerInit)

	c.configWebStaticPoints()

	adminGroup := c.e.Group("/admin", c.auth.VerifyAccessTokenMW(), c.auth.GrantAdminAccessMW())
	{
		adminGroup.GET("/get/users", c.adminRouter.HandlerAdminGetUsers)
		adminGroup.GET("/get/roles", c.adminRouter.HandlerAdminGetRoles)

		adminGroup.DELETE("/session/flush", c.adminRouter.HandlerAdminSessionFlush)
		adminGroup.DELETE("/user/delete", c.adminRouter.HandlerAdminDeleteUser)

		adminGroup.PATCH("/user/assign/role", c.adminRouter.HandlerAdminUpdateUserRole)
	}

	authGroup := c.e.Group("/auth")
	{
		authGroup.POST("/signup", c.auth.HandlerSignUp, c.auth.VerifyAccessTokenMW(), c.auth.GrantAdminAccessMW())

		authGroup.POST("/login", c.auth.HandlerLogin)
		authGroup.POST("/logout", c.auth.HandlerLogout, c.auth.VerifyAccessTokenMW())

		authGroup.POST("/token/refresh", c.auth.HandlerRefresh)

		authGroup.PATCH("/password/change", c.auth.HandlerPasswordChange)
	}

	employeeGroup := c.e.Group("/employee", c.auth.VerifyAccessTokenMW())
	{
		employeeGroup.GET("/get/citizenships", c.employeeRouter.HandlerGetCitizenships, c.auth.GrantAllAccessMW())
		employeeGroup.GET("/get/titles", c.employeeRouter.HandlerGetTitles, c.auth.GrantAllAccessMW())
		employeeGroup.GET("/get/list", c.employeeRouter.HandlerGetEmployeeWithFilter, c.auth.GrantAllAccessMW())

		employeeGroup.GET("/statistics/citizenship", c.employeeRouter.HandlerCountEmployeeWithCitizenship, c.auth.GrantAllAccessMW())
		employeeGroup.GET("/statistics/salary", c.employeeRouter.HandlerCountEmployeesWithSalaryBoundary, c.auth.GrantAllAccessMW())

		employeeGroup.POST("/create", c.employeeRouter.HandlerCreateEmployee, c.auth.GrantAdminOrHRManagerAccessMW())

		employeeGroup.PUT("/update", c.employeeRouter.HandlerUpdateEmployee, c.auth.GrantAdminOrHRManagerAccessMW())
		employeeGroup.DELETE("/delete", c.employeeRouter.HandlerDeleteEmployee, c.auth.GrantAdminOrHRManagerAccessMW())
	}

	divisionGroup := c.e.Group("/division", c.auth.VerifyAccessTokenMW())
	{
		divisionGroup.GET("/get/list", c.divisionRouter.HandlerGetDivisions, c.auth.GrantAllAccessMW())

		divisionGroup.GET("/statistics/salary", c.divisionRouter.HandlerGetSalaryStatisticsOfDivision, c.auth.GrantAllAccessMW())
		divisionGroup.GET("/statistics/statesize", c.divisionRouter.HandlerGetStateSizeStatisticsOfDivisions, c.auth.GrantAllAccessMW())

		divisionGroup.POST("/create", c.divisionRouter.HandlerCreateDivision, c.auth.GrantAdminAccessMW())

		divisionGroup.PUT("/update", c.divisionRouter.HandlerUpdateDivision, c.auth.GrantAdminAccessMW())
		divisionGroup.DELETE("/delete", c.divisionRouter.HandlerDeleteDivision, c.auth.GrantAdminAccessMW())
	}
}

func (c Controller) configWebStaticPoints() {
	c.e.Static("/index.html", "./frontend/index.html")
	// c.e.Static("/styles/main.css", "./frontend/styles/main.css")
	// c.e.Static("/js/config.js", "./frontend/js/config.js")
	// c.e.Static("/js/admin.js", "./frontend/js/admin.js")
	// c.e.Static("/js/api.js", "./frontend/js/api.js")
	// c.e.Static("/js/app.js", "./frontend/js/app.js")
	// c.e.Static("/js/auth.js", "./frontend/js/auth.js")
	// c.e.Static("/js/divisions.js", "./frontend/js/divisions.js")
	// c.e.Static("/js/employees.js", "./frontend/js/employees.js")
	// c.e.Static("/js/login.js", "./frontend/js/login.js")
	// c.e.Static("/js/router.js", "./frontend/js/router.js")
	// c.e.Static("/js/utils.js", "./frontend/js/utils.js")
}
