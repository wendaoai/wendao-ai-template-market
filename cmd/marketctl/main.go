package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
		return fmt.Errorf("用法：marketctl <validate|build|sign|verify|keygen>")
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
	case "sign":
		flags := flag.NewFlagSet("sign", flag.ContinueOnError)
		root := flags.String("root", ".", "仓库根目录")
		keyID := flags.String("key-id", "", "签名密钥 ID")
		privateKey := flags.String("private-key", "", "Base64 Ed25519 私钥文件")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *keyID == "" || *privateKey == "" {
			return fmt.Errorf("sign 需要 --key-id 和 --private-key")
		}
		signature, err := catalog.Sign(*root, *keyID, *privateKey)
		if err != nil {
			return err
		}
		fmt.Printf("signed stable.json with %s (%s)\n", signature.KeyID, signature.SHA256)
		return nil
	case "verify":
		flags := flag.NewFlagSet("verify", flag.ContinueOnError)
		root := flags.String("root", ".", "仓库根目录")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		result, signature, err := catalog.Verify(*root)
		if err != nil {
			return err
		}
		fmt.Printf("verified catalog %s with %s: %d templates\n", result.CatalogVersion, signature.KeyID, len(result.Templates))
		return nil
	case "keygen":
		flags := flag.NewFlagSet("keygen", flag.ContinueOnError)
		root := flags.String("root", ".", "仓库根目录")
		keyID := flags.String("key-id", "", "密钥 ID")
		privateKey := flags.String("private-key", "", "私钥输出路径")
		if err := flags.Parse(args[1:]); err != nil {
			return err
		}
		if *keyID == "" {
			return fmt.Errorf("keygen 需要 --key-id")
		}
		privatePath := *privateKey
		if privatePath == "" {
			privatePath = filepath.Join(*root, ".local-keys", *keyID+".key")
		}
		publicPath := filepath.Join(*root, "keys", *keyID+".pub")
		if err := catalog.GenerateKeyPair(privatePath, publicPath); err != nil {
			return err
		}
		fmt.Printf("generated public key %s; private key is ignored at %s\n", publicPath, privatePath)
		return nil
	default:
		return fmt.Errorf("未知命令：%s", args[0])
	}
}
