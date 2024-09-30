package manager

type Options interface {
	apply(v *VersionManager)
}
type optionFunc func(v *VersionManager)

func (o optionFunc) apply(v *VersionManager) {
	o(v)
}

func WithVersion(version string) Options {
	return optionFunc(func(o *VersionManager) {
		o.version = version
	})
}
