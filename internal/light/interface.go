package light

type LightClient interface {
	GetAll() (map[string]Light, error)
	Toggle(lightID string, state bool) error
}
