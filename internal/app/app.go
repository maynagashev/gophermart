package app

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gophermart/internal/handlers"
	"gophermart/internal/repository"
	"gophermart/internal/service"
	"gophermart/internal/worker"
)

// App представляет основную структуру приложения
type App struct {
	echo          *echo.Echo
	db            *sqlx.DB
	userHandler   *handlers.UserHandler
	orderHandler  *handlers.OrderHandler
	accrualWorker *worker.AccrualWorker
	config        Config
}

// New создает новый экземпляр приложения
func New(ctx context.Context, cfg Config) (*App, error) {
	// Инициализация базы данных
	db, err := NewDB(ctx, cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}

	// Применение миграций
	if err := MigrateDB(db.DB, cfg.MigrationsDir); err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo, cfg.JWTSecret, cfg.JWTExpirationPeriod)
	orderService := service.NewOrderService(orderRepo)
	accrualService := service.NewAccrualService(cfg.AccrualSystemAddress)

	// Инициализация воркера начислений
	accrualWorker := worker.NewAccrualWorker(
		orderRepo,
		accrualService,
		5, // количество воркеров
		0, // используем значения по умолчанию для интервалов
		0,
	)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Инициализация Echo
	e := echo.New()
	e.Validator = NewValidator()

	// Промежуточное ПО (middleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	app := &App{
		echo:          e,
		db:            db,
		userHandler:   userHandler,
		orderHandler:  orderHandler,
		accrualWorker: accrualWorker,
		config:        cfg,
	}

	// Настройка маршрутов
	app.setupRoutes()

	return app, nil
}

// Start запускает приложение
func (a *App) Start(address string) error {
	// Запускаем воркер начислений в отдельной горутине
	go a.accrualWorker.Start(context.Background())

	return a.echo.Start(address)
}

// Shutdown выполняет корректное завершение работы приложения
func (a *App) Shutdown(ctx context.Context) error {
	if err := a.db.Close(); err != nil {
		return err
	}
	return a.echo.Shutdown(ctx)
}

// setupRoutes настраивает маршруты приложения
func (a *App) setupRoutes() {
	// Группа API
	api := a.echo.Group("/api")

	// Маршруты пользователя
	user := api.Group("/user")

	// Публичные маршруты
	user.POST("/register", a.userHandler.Register)
	user.POST("/login", a.userHandler.Authenticate)

	// Защищенные маршруты
	protected := user.Group("", JWTMiddleware(a.config.JWTSecret))

	// Маршруты заказов
	protected.POST("/orders", a.orderHandler.Register)
	protected.GET("/orders", a.orderHandler.GetOrders)
}
