package pack

type Pack interface {
	Vendor(vendorDir string) error
	Version() string
	ImportUrl() string
	Update() error
}
