package src

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/cen-ngc5139/BeePF/loader/lib/src/meta"
	"go.uber.org/zap"
)

type SendHandler struct {
	Logger   *zap.Logger
	SkipNets []*net.IPNet // 改为切片，支持多个CIDR
}

// BPF程序中的ipv4_key_t结构体对应的Go结构
type EventWrapper struct {
	Pid   uint32 `json:"pid"`
	Saddr uint32 `json:"saddr"`
	Daddr uint32 `json:"daddr"`
	Lport uint16 `json:"lport"`
	Dport uint16 `json:"dport"`
	Bytes uint64 `json:"bytes"`
}

// 将uint32格式的IP地址转换为可读的字符串
func (k *EventWrapper) FormatSaddr() string {
	return net.IPv4(byte(k.Saddr), byte(k.Saddr>>8), byte(k.Saddr>>16), byte(k.Saddr>>24)).String()
}

func (k *EventWrapper) FormatDaddr() string {
	return net.IPv4(byte(k.Daddr), byte(k.Daddr>>8), byte(k.Daddr>>16), byte(k.Daddr>>24)).String()
}

// 实现 EventHandler 接口
func (h *SendHandler) HandleEvent(ctx *meta.UserContext, data *meta.ReceivedEventData) error {
	switch data.Type {
	case meta.TypeJsonText:
		// 首先解析外层JSON
		var wrapper EventWrapper
		err := json.Unmarshal([]byte(data.JsonText), &wrapper)
		if err != nil {
			h.Logger.Error("解析外层JSON失败", zap.Error(err))
			return err
		}

		// 校验目标地址是否在跳过列表中
		if h.shouldSkip(wrapper.Daddr) {
			// h.Logger.Debug("地址在跳过列表中，忽略",
			// 	zap.String("源地址", wrapper.FormatSaddr()),
			// 	zap.String("目标地址", wrapper.FormatDaddr()))
			return nil
		}

		// 记录TCP连接信息
		h.Logger.Info("TCP连接数据",
			zap.String("源地址", fmt.Sprintf("%s:%d", wrapper.FormatSaddr(), wrapper.Lport)),
			zap.String("目标地址", fmt.Sprintf("%s:%d", wrapper.FormatDaddr(), wrapper.Dport)),
			zap.Uint32("进程ID", wrapper.Pid),
			zap.Uint64("字节数", wrapper.Bytes))

	case meta.TypePlainText:
		h.Logger.Info("接收到纯文本",
			zap.String("data", data.Text))
	}
	return nil
}

// 检查IP是否在需要跳过的网段中
func (h *SendHandler) shouldSkip(ipAddr uint32) bool {
	if len(h.SkipNets) == 0 {
		return false
	}

	// 将uint32转换为net.IP
	ip := net.IPv4(byte(ipAddr), byte(ipAddr>>8), byte(ipAddr>>16), byte(ipAddr>>24))

	// 检查IP是否在任何一个跳过的网段中
	for _, skipNet := range h.SkipNets {
		if skipNet.Contains(ip) {
			return true
		}
	}
	return false
}

type SkipHandler struct {
}

func (h *SkipHandler) HandleEvent(ctx *meta.UserContext, data *meta.ReceivedEventData) error {
	return nil
}

// 从逗号分隔的CIDR列表字符串解析网段
func ParseCIDRList(cidrList string) ([]*net.IPNet, error) {
	if cidrList == "" {
		return nil, nil
	}

	var nets []*net.IPNet
	for _, cidr := range strings.Split(cidrList, ",") {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}

		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("解析CIDR失败 %s: %v", cidr, err)
		}
		nets = append(nets, ipNet)
	}
	return nets, nil
}
