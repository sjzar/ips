package ips

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sjzar/ips/parser"
)

var (
	rootDBFormat string
	rootDBFile   string
	rootFields   string
)

func init() {
	rootCmd.Flags().StringVarP(&rootDBFormat, "format", "", "", "database format")
	rootCmd.Flags().StringVarP(&rootDBFile, "database", "d", "", "database file")
	rootCmd.Flags().StringVarP(&rootFields, "fields", "f", "", "fields")
}

var rootCmd = &cobra.Command{
	Use:   "ips",
	Short: "ips commandline tools",
	Long:  `ips is a tool for querying, scanning, and packing IP geolocation databases.`,
	Args:  cobra.MinimumNArgs(0),
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	Run: Root,
}

// Execute 用于启动命令行工具
func Execute() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Root IP查询命令, 支持pipeline查询
func Root(cmd *cobra.Command, args []string) {

	// pipeline mode
	if len(args) == 0 {
		if fi, err := os.Stdin.Stat(); err == nil {
			if fi.Mode()&os.ModeNamedPipe != 0 {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					str := scanner.Text()
					if len(str) != 0 {
						fmt.Println(ParseLine(str))
					}
				}
				return
			}
		}
		_ = cmd.Help()
		return
	}

	fmt.Println(ParseLine(strings.Join(args, " ")))
}

// ParseLine 解析文本
func ParseLine(line string) string {
	p := parser.NewTextParser(line)
	p.IPv4FillResult = func(str string) string {
		_, values, err := GetIPv4().Find(net.ParseIP(str))
		if err != nil {
			return ""
		}
		return strings.Join(values, " ")
	}
	p.IPv6FillResult = func(str string) string {
		_, values, err := GetIPv6().Find(net.ParseIP(str))
		if err != nil {
			return ""
		}
		return strings.Join(values, " ")
	}

	return p.Parse().String()

}
