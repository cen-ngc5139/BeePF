package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/cen-ngc5139/BeePF/server/models"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PromClient 封装 Prometheus 客户端
type PromClient struct {
	api     v1.API
	address string
}

// NewPromClient 创建新的 Prometheus 客户端
func NewPromClient(address string) (*PromClient, error) {
	client, err := api.NewClient(api.Config{
		Address: address,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 Prometheus 客户端失败: %v", err)
	}

	return &PromClient{
		api:     v1.NewAPI(client),
		address: address,
	}, nil
}

// RangeQuery 执行范围查询
func (c *PromClient) RangeQuery(query string, start, end time.Time, step time.Duration) ([]models.MetricPoint, error) {
	ctx := context.Background()
	result, warnings, err := c.api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  step,
	})

	if err != nil {
		return nil, fmt.Errorf("范围查询失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("查询警告: %v\n", warnings)
	}

	if matrix, ok := result.(model.Matrix); ok && len(matrix) > 0 {
		points := make([]models.MetricPoint, len(matrix[0].Values))
		for i, sample := range matrix[0].Values {
			points[i] = models.MetricPoint{
				Timestamp: sample.Timestamp.Time(),
				Value:     float64(sample.Value),
			}
		}
		return points, nil
	}

	return nil, fmt.Errorf("查询结果类型错误，期望 Matrix，得到 %T", result)
}

// InstantQuery 执行即时查询
func (c *PromClient) InstantQuery(query string) (model.Value, error) {
	ctx := context.Background()
	result, warnings, err := c.api.Query(ctx, query, time.Now())

	if err != nil {
		return nil, fmt.Errorf("即时查询失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("查询警告: %v\n", warnings)
	}

	return result, nil
}

// QueryLabelValues 查询标签值
func (c *PromClient) QueryLabelValues(label string) (model.LabelValues, error) {
	ctx := context.Background()
	values, warnings, err := c.api.LabelValues(ctx, label, nil, time.Now().Add(-1*time.Hour), time.Now())

	if err != nil {
		return nil, fmt.Errorf("查询标签值失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("查询警告: %v\n", warnings)
	}

	return values, nil
}

// QuerySeries 查询时间序列
func (c *PromClient) QuerySeries(query string, start, end time.Time) ([]model.LabelSet, error) {
	ctx := context.Background()
	matches, warnings, err := c.api.Series(ctx, []string{query}, start, end)

	if err != nil {
		return nil, fmt.Errorf("查询时间序列失败: %v", err)
	}

	if len(warnings) > 0 {
		fmt.Printf("查询警告: %v\n", warnings)
	}

	return matches, nil
}

// GetTargets 获取监控目标状态
func (c *PromClient) GetTargets() (*v1.TargetsResult, error) {
	ctx := context.Background()
	targets, err := c.api.Targets(ctx)

	if err != nil {
		return nil, fmt.Errorf("获取监控目标失败: %v", err)
	}

	return &targets, nil
}
