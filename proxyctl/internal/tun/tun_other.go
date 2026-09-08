//go:build !linux

package tun

import "errors"

// On 在其他平台不支持。
func On(host, port string) error {
	return errors.New("TUN 模式仅在 Linux 上支持")
}

// Off 在其他平台不支持。
func Off() error {
	return errors.New("TUN 模式仅在 Linux 上支持")
}

// ReadStatus 在其他平台返回空状态。
func ReadStatus() (Status, error) {
	return Status{}, nil
}
