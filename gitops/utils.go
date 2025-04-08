package gitops

import (
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func CopyFiles(source, dest string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dest, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}
		return copyFile(path, destPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func getSSHAuth(keyPath string) ssh.AuthMethod {
	sshKey, err := os.ReadFile(keyPath)
	if err != nil {
		panic("Failed to read SSH key: " + err.Error())
	}

	signer, err := ssh.ParsePrivateKey(sshKey)
	if err != nil {
		panic("Failed to parse SSH key: " + err.Error())
	}

	// Use ssh.PublicKeys function to create an AuthMethod
	return ssh.PublicKeys(signer)
}
