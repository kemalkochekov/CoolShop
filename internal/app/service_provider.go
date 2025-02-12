package app

import (
	"Coolshop/internal/auth"
	"Coolshop/internal/cache"
	"Coolshop/internal/config"
	"Coolshop/internal/connection"
	"Coolshop/internal/user"
	userHandler "Coolshop/internal/user/delivery"
	userRepository "Coolshop/internal/user/repository"
	userUseCase "Coolshop/internal/user/usecase"
	"Coolshop/logger"
)

type serviceProvider struct {
	config          config.Config
	repository      user.Repository
	cacheRepository cache.Repository
	useCase         user.UseCase
	handler         user.Handlers
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() (config.Config, error) {
	if s.config == nil {
		cfg, err := config.NewConfig()
		if err != nil {
			return nil, err
		}
		s.config = cfg
	}

	return s.config, nil
}

func (s *serviceProvider) UserRepository(psqlDB connection.DB) user.Repository {
	if s.repository == nil {
		s.repository = userRepository.NewRepository(psqlDB)
	}

	return s.repository
}

func (s *serviceProvider) CacheRepository(redisDB connection.Cache) cache.Repository {
	if s.cacheRepository == nil {
		s.cacheRepository = cache.NewRedisRepo(redisDB)
	}

	return s.cacheRepository
}

func (s *serviceProvider) UserUseCase(psqlDB connection.DB, cacheDB connection.Cache, jwtGen auth.JwtGen) user.UseCase {
	if s.useCase == nil {
		s.useCase = userUseCase.NewUseCase(s.UserRepository(psqlDB), s.CacheRepository(cacheDB), jwtGen)
	}

	return s.useCase
}

func (s *serviceProvider) UserHandler(
	psqlDB connection.DB,
	cacheDB connection.Cache,
	jwtGen auth.JwtGen,
	zapLogger logger.LoggerInterface,
) user.Handlers {
	if s.handler == nil {
		s.handler = userHandler.NewHandler(s.UserUseCase(psqlDB, cacheDB, jwtGen), zapLogger)
	}

	return s.handler
}
