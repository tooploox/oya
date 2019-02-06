package pack

type Pack interface {
	Vendor(vendorDir string) error
	Version() string
	ImportPath() string
	Update() error
}
