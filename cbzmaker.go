package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sqweek/dialog"
)

func zipFolder(srcFolder, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(srcFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Bỏ chính folder root, chỉ zip nội dung bên trong
		if path == srcFolder {
			return nil
		}

		relPath, err := filepath.Rel(srcFolder, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		w, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, file)
		return err
	})

	return err
}

func main() {
	root, err := dialog.Directory().Title("Chọn folder root chứa các chapter").Browse()
	if err != nil {
		fmt.Println("Đã huỷ hoặc không chọn folder.")
		return
	}

	fmt.Println("Root:", root)
	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("Lỗi đọc folder:", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			src := filepath.Join(root, name)
			zipPath := filepath.Join(root, name+".zip")
			cbzPath := filepath.Join(root, name+".cbz")

			fmt.Println("Đang xử lý:", name)

			err := zipFolder(src, zipPath)
			if err != nil {
				fmt.Println("Lỗi khi zip:", err)
				continue
			}

			// rename zip → cbz
			err = os.Rename(zipPath, cbzPath)
			if err != nil {
				fmt.Println("Lỗi khi đổi tên sang cbz:", err)
				continue
			}

			fmt.Println("✓ Tạo:", cbzPath)
		}
	}

	fmt.Println("\nHoàn thành!")
}
