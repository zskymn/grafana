package metrics

import (
	"github.com/grafana/grafana/pkg/log"
	"github.com/grafana/grafana/pkg/setting"
)

type MetricPublisher interface {
	Publish(metrics []Metric)
}

type MetricSettings struct {
	Enabled         bool
	IntervalSeconds int64

	Publishers []MetricPublisher
}

func readSettings() *MetricSettings {
	var settings = &MetricSettings{
		Enabled:    false,
		Publishers: make([]MetricPublisher, 0),
	}

	var section, err = setting.Cfg.GetSection("metrics")
	if err != nil {
		log.Fatal(3, "Unable to find metrics config section")
		return nil
	}

	settings.Enabled = section.Key("enabled").MustBool(false)
	settings.IntervalSeconds = section.Key("interval_seconds").MustInt64(10)

	if !settings.Enabled {
		return settings
	}

	if graphitePublisher, err := CreateGraphitePublisher(); err != nil {
		log.Error(3, "Metrics: Failed to init Graphite metric publisher", err)
	} else if graphitePublisher != nil {
		log.Info("Metrics: Graphite publisher initialized")
		settings.Publishers = append(settings.Publishers, graphitePublisher)
	}

	return settings
}
