// pkg/helper/log2_hist.go

package helper

import (
	"fmt"
	"strings"
)

// PrintLog2Hist 打印 log2 直方图
func PrintLog2Hist(values []uint32, valType string) string {
	// 配置参数
	const starsMax = 40

	// 找到最大索引和最大值
	idxMax := -1
	var valMax uint32 = 0

	for i, v := range values {
		if v > 0 {
			idxMax = i
		}
		if v > valMax {
			valMax = v
		}
	}

	// 如果没有数据，直接返回
	if idxMax < 0 {
		return ""
	}

	var sb strings.Builder

	// 打印标题行
	width1 := 5
	width2 := 19
	if idxMax > 32 {
		width1 = 15
		width2 = 29
	}

	fmt.Fprintf(&sb, "%*s%-*s : count    distribution\n",
		width1, "",
		width2, valType,
	)

	// 计算星号数量
	stars := starsMax
	if idxMax > 32 {
		stars = starsMax / 2
	}

	// 打印每一行
	for i := 0; i <= idxMax; i++ {
		val := values[i]

		// 计算区间
		low := uint64(1) << uint64(i)
		high := (uint64(1) << uint64(i+1)) - 1

		if low == high {
			low--
		}

		// 设置宽度
		width := 10
		if idxMax > 32 {
			width = 20
		}

		// 打印区间和计数
		fmt.Fprintf(&sb, "%*d -> %-*d : %-8d |",
			width, low,
			width, high,
			val,
		)

		// 打印星号
		printStars(&sb, val, valMax, stars)

		sb.WriteString("|\n")
	}

	return sb.String()
}

// printStars 打印星号
func printStars(sb *strings.Builder, val, valMax uint32, width int) {
	// 计算星号数量
	var numStars int
	if val <= valMax {
		numStars = int(float64(val) * float64(width) / float64(valMax))
	} else {
		numStars = width
	}

	// 打印星号
	sb.WriteString(strings.Repeat("*", numStars))

	// 打印空格
	sb.WriteString(strings.Repeat(" ", width-numStars))

	// 如果值超过最大值，添加 + 号
	if val > valMax {
		sb.WriteString("+")
	}
}
