//go:build !darwin && !windows && !linux

package proxy

import "errors"

// Get 在没有系统代理后端（macOS/Windows/Linux+GNOME）的平台上无法读取。
func Get() (*Info, error) {
	return nil, errors.New("当前平台不支持系统代理管理（支持 macOS、Windows、Linux/GNOME 桌面）")
}

// Clear 在没有系统代理后端的平台上不支持。
func Clear() error {
	return errors.New("当前平台不支持系统代理管理（支持 macOS、Windows、Linux/GNOME 桌面）")
}

// SnapshotSystem 在没有系统代理后端的平台上不支持。
func SnapshotSystem() (SystemSnapshot, error) {
	return SystemSnapshot{}, errors.New("当前平台不支持系统代理管理（支持 macOS、Windows、Linux/GNOME 桌面）")
}

// RestoreSystem 在没有系统代理后端的平台上不支持。
func RestoreSystem(SystemSnapshot) error {
	return errors.New("当前平台不支持系统代理管理（支持 macOS、Windows、Linux/GNOME 桌面）")
}

// ApplyProfile 在没有系统代理后端的平台上不支持。
func ApplyProfile(*EndpointState, *EndpointState, *EndpointState, *PACState) error {
	return errors.New("当前平台不支持系统代理管理（支持 macOS、Windows、Linux/GNOME 桌面）")
}
