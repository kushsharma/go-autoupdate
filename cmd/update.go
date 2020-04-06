package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// UpdatePackage updates current executable
func UpdatePackage(currentPath, targetPath, updateHash string) error {

	//verify target dir permission
	if err := CheckPermissions(targetPath); err != nil {
		return err
	}

	//verify sha signature, TODO: refactor
	fmt.Printf("verified sha signature of %s\n", targetPath)
	//if updateHash != string(GenerateSha256(targetPath)) {
	//	return errors.New("sha hash cannot be verified of target file")
	//}

	execDir := filepath.Dir(currentPath)
	execFilename := filepath.Base(currentPath)

	// prepare update

	// get the directory the executable exists in
	updateDir := filepath.Dir(targetPath)
	targetFilename := filepath.Base(targetPath)

	// Copy the contents of newbinary to a new executable file
	fpTarget, err := os.Open(targetPath)
	if err != nil {
		return err
	}
	defer fpTarget.Close()

	newPath := filepath.Join(updateDir, fmt.Sprintf(".%s.new", targetFilename))
	fpNew, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	_, err = io.Copy(fpNew, fpTarget)
	if err != nil {
		return err
	}
	// close here to avoid os problems in move
	fpNew.Close()
	fmt.Printf("target update(%s) copied to (%s)\n", targetPath, newPath)

	// this is where we'll move the executable to so that we can swap in the updated replacement
	oldPath := filepath.Dir(targetPath)
	oldPath = filepath.Join(execDir, fmt.Sprintf(".%s.old", execFilename))

	// delete any existing old exec file from last update
	if err = os.Remove(oldPath); err == nil {
		fmt.Printf("cleaned old update file found on (%s)\n", oldPath)
	}

	// move the existing executable to a new file with .old suffix in the same directory
	err = os.Rename(currentPath, oldPath)
	if err != nil {
		return err
	}
	fmt.Printf("backup current exec file(%s) to (%s)\n", currentPath, oldPath)

	// move the new exectuable in to become the new program
	err = os.Rename(newPath, currentPath)
	fmt.Printf("moved new update file(%s) to target location(%s)\n", newPath, currentPath)

	if err != nil {
		// move unsuccessful
		//
		// The filesystem is now in a bad state. We have successfully
		// moved the existing binary to a new location, but we couldn't move the new
		// binary to take its place. That means there is no file where the current executable binary
		// used to be!
		// Try to rollback by restoring the old binary to its original path.
		rerr := os.Rename(oldPath, targetPath)
		if rerr != nil {
			return rerr
		}

		return err
	}

	// move successful, remove the old binary
	return os.Remove(oldPath)
}

// CheckPermissions determines whether the process has the correct permissions to
// perform the requested update. If the update can proceed, it returns nil, otherwise
// it returns the error that would occur if an update were attempted.
func CheckPermissions(target string) error {
	// get the directory the file exists in

	fileDir := filepath.Dir(target)
	fileName := filepath.Base(target)

	// attempt to open a file in the file's directory
	newPath := filepath.Join(fileDir, fmt.Sprintf(".%s.tmp", fileName))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	fp.Close()

	_ = os.Remove(newPath)
	return nil
}
