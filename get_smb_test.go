package getter

import (
	"context"
	urlhelper "github.com/hashicorp/go-getter/helper/url"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestSmbGetter_impl(t *testing.T) {
	var _ Getter = new(SmbGetter)
}

// TODO:
// allow download directory (?)
// write higher level tests
// save tests results on circleci
// write docs of how to run tests locally (makefile?)
// mount folder for testing

func TestSmbGetter_Get(t *testing.T) {
	smbTestsPreCheck(t)

	tests := []struct{
		name       string
		rawURL     string
		file       string
		createFile string
		fail       bool
	}{
		{
			"smbclient with authentication",
			"smb://vagrant:vagrant@samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient with authentication and subdir",
			"smb://vagrant:vagrant@samba/shared/subdir/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient with only username authentication",
			"smb://vagrant@samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient without authentication",
			"smb://samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"local mounted smb shared file",
			"smb://samba/shared/mounted.txt",
			"mounted.txt",
			"/samba/shared/mounted.txt",
			false,
		},
		{
			"no hostname provided",
			"smb://",
			"",
			"",
			true,
		},
		{
			"no filepath provided",
			"smb://samba",
			"",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile != "" {
				// mock mounted folder by creating one
				err := os.MkdirAll(filepath.Dir(tt.createFile), 0755)
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}

				f, err := os.Create(tt.createFile)
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}
				defer f.Close()

				// Write content to assert later
				_, err = f.WriteString("Hello\n")
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}
				f.Sync()

				defer os.RemoveAll(tt.createFile)
			}

			dst := tempDir(t)
			defer os.RemoveAll(dst)

			url, err := urlhelper.Parse(tt.rawURL)
			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}
			req := &Request{
				Dst: dst,
				u:   url,
			}

			g := new(SmbGetter)
			ctx := context.Background()
			err = g.GetFile(ctx, req)
			fail := err != nil

			if tt.fail != fail {
				if fail {
					t.Fatalf("err: unexpected error %s", err.Error())
				}
				t.Fatalf("err: expecting to fail but it did not")
			}

			if !tt.fail {
				if tt.createFile != "" {
					// Verify the destination folder is a symlink to mounted folder
					fi, err := os.Lstat(dst)
					if err != nil {
						log.Printf("MOSS err 1")
						t.Fatalf("err: %s", err)
					}
					if fi.Mode()&os.ModeSymlink == 0 {
						t.Fatal("destination is not a symlink")
					}
					// Verify the main file exists
					assertContents(t, dst, "Hello\n")
				} else {
					// Verify the file exists at the destination folder
					mainPath := filepath.Join(dst, tt.file)
					if _, err := os.Stat(mainPath); err != nil {
						log.Printf("MOSS err 2")
						t.Fatalf("err: %s", err)
					}
				}
			}
		})
	}
}

func TestSmbGetter_GetFile(t *testing.T) {
	smbTestsPreCheck(t)

	tests := []struct{
		name       string
		rawURL     string
		file       string
		createFile string
		fail       bool
	}{
		{
			"smbclient with authentication",
			"smb://vagrant:vagrant@samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient with authentication and subdirectory",
			"smb://vagrant:vagrant@samba/shared/subdir/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient with only username authentication",
			"smb://vagrant@samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient without authentication",
			"smb://samba/shared/file.txt",
			"file.txt",
			"",
			false,
		},
		{
			"smbclient get directory",
			"smb://vagrant:vagrant@samba/shared/subdir",
			"",
			"",
			true,
		},
		{
			"smbclient get non existent file",
			"smb://vagrant:vagrant@samba/shared/nonexistent.txt",
			"",
			"",
			true,
		},
		{
			"local mounted smb shared file",
			"smb://samba/shared/mounted.txt",
			"mounted.txt",
			"/samba/shared/mounted.txt",
			false,
		},
		{
			"no hostname provided",
			"smb://",
			"",
			"",
			true,
		},
		{
			"no filepath provided",
			"smb://samba",
			"",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile != "" {
				// mock mounted folder by creating one
				err := os.MkdirAll(filepath.Dir(tt.createFile), 0755)
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}

				f, err := os.Create(tt.createFile)
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}
				defer f.Close()

				// Write content to assert later
				_, err = f.WriteString("Hello\n")
				if err != nil {
					t.Fatalf("err: %s", err.Error())
				}
				f.Sync()

				defer os.RemoveAll(tt.createFile)
			}

			dst := tempDir(t)
			defer os.RemoveAll(dst)

			url, err := urlhelper.Parse(tt.rawURL)
			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}
			req := &Request{
				Dst: dst,
				u:   url,
			}

			g := new(SmbGetter)
			err = g.GetFile(context.Background(), req)

			fail := err != nil
			if tt.fail != fail {
				if fail {
					t.Fatalf("err: unexpected error %s", err.Error())
				}
				t.Fatalf("err: expecting to fail but it did not")
			}

			if !tt.fail {
				if tt.createFile != "" {
					// Verify the destination folder is a symlink to mounted folder
					fi, err := os.Lstat(dst)
					if err != nil {
						log.Printf("MOSS err 1")
						t.Fatalf("err: %s", err)
					}
					if fi.Mode()&os.ModeSymlink == 0 {
						t.Fatal("destination is not a symlink")
					}
					// Verify the main file exists
					assertContents(t, dst, "Hello\n")
				} else {
					// Verify if the file was successfully download
					// and exists at the destination folder
					mainPath := filepath.Join(dst, tt.file)
					if _, err := os.Stat(mainPath); err != nil {
						log.Printf("MOSS err 2")
						t.Fatalf("err: %s", err)
					}
				}
			}
		})
	}
}

func TestSmbGetter_Mode(t *testing.T) {
	smbTestsPreCheck(t)

	tests := []struct {
		name string
		rawURL string
		expectedMode Mode
		fail bool
	}{
		{
			"smbclient modefile for existing file",
			"smb://vagrant:vagrant@samba/shared/file.txt",
			ModeFile,
			false,
		},
		{
			"smbclient modedir for existing directory",
			"smb://vagrant:vagrant@samba/shared/subdir",
			ModeDir,
			false,
		},
		{
			"smbclient mode fail for unexisting directory",
			"smb://vagrant:vagrant@samba/shared/invaliddir",
			0,
			true,
		},
		{
			"smbclient mode fail for unexisting file",
			"smb://vagrant:vagrant@samba/shared/invalidfile.txt",
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := urlhelper.Parse(tt.rawURL)
			if err != nil {
				t.Fatalf("err: %s", err.Error())
			}

			g := new(SmbGetter)
			mode, err := g.Mode(context.Background(), url)

			fail := err != nil
			if tt.fail != fail {
				if fail {
					t.Fatalf("err: unexpected error %s", err.Error())
				}
				t.Fatalf("err: expecting to fail but it did not")
			}

			if mode != tt.expectedMode {
				t.Fatalf("err: expeting mode %d, actual mode %d", tt.expectedMode, mode)
			}
		})
	}
}

func smbTestsPreCheck(t *testing.T) {
	r := os.Getenv("ACC_SMB_TEST")
	if r != "1" {
		t.Skip("Smb getter tests won't run. ACC_SMB_TEST not set")
	}
}
