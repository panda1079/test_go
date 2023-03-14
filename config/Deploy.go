package config

type Deploy struct {
}

func (r *Deploy) Run() map[string]string {
	return map[string]string{
		"LISTEN_ADDRESS": "0.0.0.0",
		"PORT":           "9091",
		"LOG_DIR":        "./cache", //日志写入位置
	}
}
