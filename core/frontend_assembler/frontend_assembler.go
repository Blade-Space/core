package frontend_assembler

import (
	"fmt"
	"os/exec"
)

func Init(repo string) error {
	fmt.Println("\nStarting downloads front-end ⬇️")

	err := downloadGitRepository(repo)

	if err != nil {
		return err
	}

	compileFontEnd()

	fmt.Println("")
	return nil
}

func downloadGitRepository(url string) error {
	_, err := exec.Command("git", "clone", url, "out/front-end/").Output()
	if err != nil {
		fmt.Printf("🔺 Download error %s: %v\n", url, err)
		return err
	} else {
		fmt.Printf("✅ Front-end successfully download\n")
	}

	return nil
}

func compileFontEnd() error {
	cmd := exec.Command("yarn")
	cmd.Dir = "./out/front-end"
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("🔺 Install packages failed: %w\n", err)
		return err
	} else {
		fmt.Println("✅ Install packages")
	}

	cmdBuild := exec.Command("yarn", "build")
	cmdBuild.Dir = "./out/front-end"

	_, err = cmdBuild.CombinedOutput()

	if err != nil {
		fmt.Println("🔺 Compilation failed: %w\n", err)
		return err
	} else {
		fmt.Printf("✅ Compilation successfully")
	}

	return nil
}
