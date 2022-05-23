package cmd

func init() {
	err := zap.RegisterSink("winfile", newWinFileSink)
	if err != nil {
		return nil, "", fmt.Errorf("couldn't register the windows file sink: %w", err)
	}
}

func toZapLogPath(p string) string {
	return "winfile:///" + appLogPath
}
