package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	// 创建一个根命令
	rootCmd := &cobra.Command{
		Use:   "myapp",
		Short: "MyApp is a simple command-line application",
		Long:  `MyApp demonstrates how to create command-line applications using Cobra.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello, world!")
		},
	}

	templateDir := "template/"

	// 创建一个子命令
	createCmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create project",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectName := args[0]

			err := copyDir(templateDir, projectName)

			fmt.Println(err)
		},
	}

	// 将子命令添加到根命令
	rootCmd.AddCommand(createCmd)

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func copyDir(src, dest string) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, rel)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()
		_, err = io.Copy(destFile, srcFile)
		return err
	})
}
