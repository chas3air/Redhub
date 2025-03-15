package tokenssaver

type MockTokenSaver struct {
	tokens map[string][]string
}

func New() *MockTokenSaver {
	return &MockTokenSaver{
		tokens: make(map[string][]string),
	}
}

func (m *MockTokenSaver) SaveTokens(uid string, accessToken string, refreshToken string) error {
	m.tokens[uid] = []string{accessToken, refreshToken}
	return nil
}

func (m *MockTokenSaver) GetTokens(uid string) (accessToken string, refreshToken string, err error) {
	key := m.tokens[uid]
	return key[0], key[1], nil
}
