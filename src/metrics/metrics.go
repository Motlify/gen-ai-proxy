package metrics

import (
	"context"

	"gen-ai-proxy/src/database"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCollector struct {
	db *database.Queries
	totalTokens *prometheus.Desc
	totalPrice *prometheus.Desc
	totalInputTokensByModel *prometheus.Desc
	totalOutputTokensByModel *prometheus.Desc
}

func NewMetricsCollector(db *database.Queries) *MetricsCollector {
	return &MetricsCollector{
		db: db,
		totalTokens: prometheus.NewDesc(
			"gen_ai_proxy_total_tokens",
			"Total number of tokens processed by provider, model, and connection.",
			[]string{"provider_id", "provider_name", "model_id", "model_name", "connection_id", "connection_name"},
			nil,
		),
		totalPrice: prometheus.NewDesc(
			"gen_ai_proxy_total_price",
			"Total price incurred by provider, model, and connection.",
			[]string{"provider_id", "provider_name", "model_id", "model_name", "connection_id", "connection_name"},
			nil,
		),
		totalInputTokensByModel: prometheus.NewDesc(
			"gen_ai_proxy_total_input_tokens_by_model",
			"Total number of input tokens processed per model.",
			[]string{"provider_id", "provider_name", "model_id", "model_name", "connection_id", "connection_name"},
			nil,
		),
		totalOutputTokensByModel: prometheus.NewDesc(
			"gen_ai_proxy_total_output_tokens_by_model",
			"Total number of output tokens processed per model.",
			[]string{"provider_id", "provider_name", "model_id", "model_name", "connection_id", "connection_name"},
			nil,
		),
	}
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.totalTokens
	ch <- c.totalPrice
	ch <- c.totalInputTokensByModel
	ch <- c.totalOutputTokensByModel
}



func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	ctx := context.Background()

	tokenStats, err := c.db.GetTotalTokensByProviderModelConnection(ctx)
	if err != nil {
		log.Printf("Error querying total tokens: %v", err)
		return
	}
	log.Printf("Retrieved token stats: %+v", tokenStats)

	for _, stat := range tokenStats {
		log.Printf("Collecting totalTokens: provider_id=%s, provider_name=%s, model_id=%s, model_name=%s, connection_id=%s, connection_name=%s, value=%f", stat.ProviderID.String(), stat.ProviderName, stat.ModelID.String(), stat.ModelName, stat.ConnectionID.String(), stat.ConnectionName, float64(stat.TotalTokens))
		ch <- prometheus.MustNewConstMetric(
			c.totalTokens,
			prometheus.CounterValue,
			float64(stat.TotalTokens),
			stat.ProviderID.String(),
			stat.ProviderName,
			stat.ModelID.String(),
			stat.ModelName,
			stat.ConnectionID.String(),
			stat.ConnectionName,
		)
	}

	priceStats, err := c.db.GetTotalPriceByProviderModelConnection(ctx)
	if err != nil {
		log.Printf("Error querying total price: %v", err)
		return
	}
	// log.Printf("Retrieved price stats: %+v", priceStats)

	for _, stat := range priceStats {
		var priceValue float64
		if stat.TotalPrice.Valid {
			pgFloat8Value, err := stat.TotalPrice.Float64Value()
			if err != nil {
				log.Printf("Error converting TotalPrice to float64: %v", err)
				continue
			}
			priceValue = pgFloat8Value.Float64
		}
		// log.Printf("Collecting totalPrice: provider_id=%s, provider_name=%s, model_id=%s, model_name=%s, connection_id=%s, connection_name=%s, value=%f",
		// 	stat.ProviderID.String(), stat.ProviderName, stat.ModelID.String(), stat.ModelName, stat.ConnectionID.String(), stat.ConnectionName, priceValue,
		// )
		ch <- prometheus.MustNewConstMetric(
			c.totalPrice,
			prometheus.CounterValue,
			priceValue,
			stat.ProviderID.String(),
			stat.ProviderName,
			stat.ModelID.String(),
			stat.ModelName,
			stat.ConnectionID.String(),
			stat.ConnectionName,
		)
}

	inputTokenStats, err := c.db.GetTotalInputTokensByProviderModelConnection(ctx)
	if err != nil {
		log.Printf("Error querying total input tokens by model: %v", err)
		return
	}
	// log.Printf("Retrieved input token stats by model: %+v", inputTokenStats)

	for _, stat := range inputTokenStats {
		// log.Printf("Collecting totalInputTokensByModel: provider_id=%s, provider_name=%s, model_id=%s, model_name=%s, connection_id=%s, connection_name=%s, value=%f", stat.ProviderID.String(), stat.ProviderName, stat.ModelID.String(), stat.ModelName, stat.ConnectionID.String(), stat.ConnectionName, float64(stat.TotalInputTokens))
		ch <- prometheus.MustNewConstMetric(
			c.totalInputTokensByModel,
			prometheus.CounterValue,
			float64(stat.TotalInputTokens),
			stat.ProviderID.String(),
			stat.ProviderName,
			stat.ModelID.String(),
			stat.ModelName,
			stat.ConnectionID.String(),
			stat.ConnectionName,
		)
	}

	outputTokenStats, err := c.db.GetTotalOutputTokensByProviderModelConnection(ctx)
	if err != nil {
		log.Printf("Error querying total output tokens by model: %v", err)
		return
	}
	log.Printf("Retrieved output token stats by model: %+v", outputTokenStats)

	for _, stat := range outputTokenStats {
		// log.Printf("Collecting totalOutputTokensByModel: provider_id=%s, provider_name=%s, model_id=%s, model_name=%s, connection_id=%s, connection_name=%s, value=%f", stat.ProviderID.String(), stat.ProviderName, stat.ModelID.String(), stat.ModelName, stat.ConnectionID.String(), stat.ConnectionName, float64(stat.TotalOutputTokens))
		ch <- prometheus.MustNewConstMetric(
			c.totalOutputTokensByModel,
			prometheus.CounterValue,
			float64(stat.TotalOutputTokens),
			stat.ProviderID.String(),
			stat.ProviderName,
			stat.ModelID.String(),
			stat.ModelName,
			stat.ConnectionID.String(),
			stat.ConnectionName,
		)
	}
}
