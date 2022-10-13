package parser

const PluginName = "feed_parser_service"

type Plugin struct {
}

func (p *Plugin) Init() error {
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}
