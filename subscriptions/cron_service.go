package subscriptions

import (
	"fmt"
	"strconv"
	"test_project/rates"
	"test_project/settings"
	"time"
)

type CronEmailService struct {
	repository      IRepository
	emailSender     *EmailSender
	currencyService *rates.NbuRatesService

	isRunning   bool
	runningChan chan struct{}
}

func NewCronService(repository IRepository, settings settings.AppSettings, ratesService *rates.NbuRatesService) *CronEmailService {
	sender := NewEmailSender(settings)
	return &CronEmailService{
		repository:      repository,
		emailSender:     sender,
		runningChan:     make(chan struct{}),
		currencyService: ratesService,
	}
}

func (s *CronEmailService) Start() {
	s.isRunning = true
	go s.RunLoop()
	s.runningChan <- struct{}{}
}

func (s *CronEmailService) Stop() {
	s.isRunning = false
	<-s.runningChan // waiting for service to stop
}

func (s *CronEmailService) RunLoop() {
	for s.isRunning {
		now := time.Now()
		nextDayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Add(24 * time.Hour)
		time.Sleep(nextDayStart.Sub(now))

		if !s.isRunning {
			return
		}

		<-s.runningChan
		err := s.sendEmails(0)
		if err != nil {
			panic(err) // or log, or restart service, or wait for 5 minutes and try again, many approaches are possible
		}

		s.runningChan <- struct{}{}
	}
}

func (s *CronEmailService) sendEmails(try int32) error {
	subscriptions, err := s.repository.GetAll()

	if err != nil {
		if try > 3 {
			return err
		}

		time.Sleep(10 * time.Second)
		return s.sendEmails(try + 1)
	}

	if len(subscriptions) == 0 {
		return nil
	}

	rate, err := s.currencyService.GetRate()

	if err != nil {
		if try > 3 {
			return err
		}

		time.Sleep(10 * time.Second)
		return s.sendEmails(try + 1)
	}

	text := fmt.Sprintf("Current currency rate USD to UAH is %v", strconv.FormatFloat(rate, 'f', 2, 64))
	for _, subscription := range subscriptions {
		s.sendEmail(0, subscription, text)
	}

	return nil
}

func (s *CronEmailService) sendEmail(try int32, subscription *Subscription, text string) {
	err := s.emailSender.SendEmail(subscription.Email, text)

	if err != nil {
		if try > 3 {
			return
		}

		time.Sleep(10 * time.Second)
		s.sendEmails(try + 1)
	}
}
