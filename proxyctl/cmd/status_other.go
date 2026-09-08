//go:build !linux

package cmd

import "io"

// printOmarchyStatus 在非 Linux 平台为空操作。
func printOmarchyStatus(io.Writer) {}
