package apiutils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rakeranjan/image-service/utils"
)

func ZipFolder(folder, zipFilePath string) (string, error) {
	sourceFolder := utils.FILE_SUFFIX + folder
	// Check if zipFilePath is a directory instead of a file
	if stat, err := os.Stat(zipFilePath); err == nil && stat.IsDir() {
		return "", fmt.Errorf("zipFilePath must include the name of the .zip file, not a directory")
	}

	// Create the ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	// Create a ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the folder and add files to the ZIP archive
	err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		// Get the relative path of the file or folder
		relPath, err := filepath.Rel(sourceFolder, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// If the path is a folder, add it to the ZIP archive as a directory
		if info.IsDir() {
			if relPath != "." { // Skip adding the root folder
				_, err := zipWriter.Create(relPath + "/")
				if err != nil {
					return fmt.Errorf("failed to add folder to ZIP: %w", err)
				}
			}
			return nil
		}

		// If the path is a file, add it to the ZIP archive
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		// Create a header in the ZIP archive
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create ZIP header: %w", err)
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// Create a writer for the file in the ZIP archive
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create ZIP writer: %w", err)
		}

		// Copy the file's content into the ZIP archive
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to write file to ZIP: %w", err)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to add files to ZIP: %w", err)
	}

	return zipFilePath, nil
}

// WriteFile writes data from an io.ReadCloser to a file in the specified folder
func WriteFile(folder, fileName string, content io.ReadCloser) (string, error) {
	defer content.Close() // Ensure the reader is closed after the function completes
	imageFolder := utils.FILE_SUFFIX + folder
	// Ensure the folder exists
	err := os.MkdirAll(imageFolder, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}

	// Construct the full file path
	filePath := filepath.Join(imageFolder, fileName)

	// Open the file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy the content from the io.ReadCloser to the file
	_, err = io.Copy(file, content)
	if err != nil {
		return "", fmt.Errorf("failed to write content to file: %w", err)
	}

	return filePath, nil
}
