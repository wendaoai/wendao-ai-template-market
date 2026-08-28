package main

import (
	"flag"
	"fmt"
	"os"

	"gitcode.com/wendao-ai/wendao-ai-template-market/internal/catalog"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "marketctl:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("用法：marketctl <validate|build>")
	}
	switch args[0] {
	case "validate":
		flags := flag.NewFlagSet("validate", flag.ContinueOnError)
		root := flags.String("root", ".", "仓库根目录")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		_, templates, err := catalog.LoadSource(*root)
		if err != nil {
			return err
		}
		fmt.Printf("validated %d templates\n", len(templates))
		return nil
	case "build":
		flags := flag.NewFlagSet("build", flag.ContinueOnError)
		root := flags.String("root", ".", "仓库根目录")
		check := flags.Bool("check", false, "只检查生成文件是否最新")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		result, err := catalog.WriteBuild(*root, *check)
		if err != nil {
			return err
		}
		fmt.Printf("catalog %s: %d templates\n", result.CatalogVersion, len(result.Templates))
		return nil
	default:
		return fmt.Errorf("未知命令：%s", args[0])
	}
}
