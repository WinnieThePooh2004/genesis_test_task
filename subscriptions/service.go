package subscriptions

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	return &Service{repository: repository}
}

func (s *Service) Add(email string) (bool, error) {
	exists, err := s.repository.Exists(email)
	if err != nil {
		return false, err
	}

	if exists {
		return true, err
	}

	return false, s.repository.Add(email)
}
