package mailermocks

// MockMailer fake implementation for testing.
type MockMailer struct {
	SendFunc func(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

// Send implementa la interfaz mailer.Mailer.
func (m *MockMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// Si se ha definido una función personalizada para el test, la llamamos.
	if m.SendFunc != nil {
		return m.SendFunc(templateFile, username, email, data, isSandbox)
	}
	// Comportamiento por defecto: simula un envío exitoso.
	return 200, nil
}
