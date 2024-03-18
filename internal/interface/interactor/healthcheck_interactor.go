package interactor

type HealthcheckInteractor interface {
	IsLive() error
	IsReady() error
}
