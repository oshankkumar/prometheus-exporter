package prometheus_exporter

var defaultServiceName = "service"

func RegisterService(name string) {
	defaultServiceName = name
}

func ServiceName() string {
	return defaultServiceName
}
