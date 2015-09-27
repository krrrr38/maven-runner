package command

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/krrrr38/maven-runner/repository"
	"github.com/krrrr38/maven-runner/utils"
	"io"
	"os"
	"os/exec"
	"strings"
)

// JarCommand exec jar
type JarCommand struct {
	Meta
}

// Run generates a new cli project. It returns exit code
func (c *JarCommand) Run(args []string) int {
	// options
	var (
		groupID      string
		artifactName string
		version      string
		option       string
		argument     string
		jarDirectory string
	)

	uflag := flag.NewFlagSet("jar", flag.ContinueOnError)
	uflag.Usage = func() { c.UI.Error(c.Help()) }

	uflag.StringVar(&groupID, "groupId", "", "groupId")
	uflag.StringVar(&groupID, "g", "", "groupId (short)")

	uflag.StringVar(&artifactName, "artifact", "", "artifact")
	uflag.StringVar(&artifactName, "a", "", "artifact (short)")

	uflag.StringVar(&version, "version", "", "version")
	uflag.StringVar(&version, "V", "", "version (short)")

	uflag.StringVar(&option, "option", "", "option")
	uflag.StringVar(&option, "o", "", "option (short)")

	uflag.StringVar(&argument, "argument", "", "argument")
	uflag.StringVar(&argument, "A", "", "argument (short)")

	uflag.StringVar(&jarDirectory, "jarDir", "", "jarDir")
	uflag.StringVar(&jarDirectory, "O", "", "jarDir (short)")

	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	uflag.SetOutput(errW)

	go func() {
		for errScanner.Scan() {
			c.UI.Error(errScanner.Text())
		}
	}()

	if err := uflag.Parse(args); err != nil {
		c.UI.Error("Failed to parse arguments")
		return ExitCodeFailed
	}

	parsedArgs := uflag.Args()
	if len(parsedArgs) != 0 {
		c.UI.Error(fmt.Sprintf("Invalid arguments: %s", strings.Join(parsedArgs, " ")))
		c.UI.Error(c.Help())
		return ExitCodeFailed
	}

	if groupID == "" {
		c.UI.Error("Invalid arguments: -groupId is required")
		c.UI.Error(c.Help())
		return ExitCodeFailed
	}

	if artifactName == "" {
		c.UI.Error("Invalid arguments: -artifact is required")
		c.UI.Error(c.Help())
		return ExitCodeFailed
	}

	if jarDirectory == "" {
		if workDir, err := os.Getwd(); err == nil {
			jarDirectory = workDir
		} else {
			c.UI.Error(err.Error())
			return ExitCodeFailed
		}
	}

	jarDir := &repository.JarDirectory{
		Dir: jarDirectory,
	}
	mavenRepo := repository.NewMavenRepository()
	localRepo, err := repository.NewLocalRepository()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeFailed
	}
	artifact := &repository.Artifact{
		GroupID: groupID,
		Name:    artifactName,
		Version: version,
	}

	if artifact.Version == "" {
		latestVersion, err := mavenRepo.LatestVersion(artifact)
		if err != nil {
			c.UI.Error(err.Error())
			return ExitCodeFailed
		}
		artifact.Version = latestVersion
	}

	utils.Log("info", fmt.Sprintf("jarDir = %s", jarDir))

	jarResolver := &repository.JarResolver{
		JarDir:    *jarDir,
		LocalRepo: *localRepo,
		MavenRepo: *mavenRepo,
	}
	jarPath, err := jarResolver.ResolveJar(artifact)
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeFailed
	}

	// run jar
	// cmd := exec.Command("java", option, "-jar", jarPath, argument) // Error: Could not find or load main class ???
	command := fmt.Sprintf("java %s -jar %s %s", option, jarPath, argument)
	cmd := exec.Command("sh", "-c", command)
	utils.Log("info", fmt.Sprintf("Exec jar: sh -c %s", cmd.Args))
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeFailed
	}
	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		c.UI.Error(err.Error())
		return ExitCodeFailed
	}
	defer stderrPipe.Close()

	go io.Copy(os.Stdout, stdoutPipe)
	go io.Copy(os.Stderr, stderrPipe)

	cmd.Run()
	return ExitCodeOK
}

// Synopsis is a one-line, short synopsis of the command.
func (c *JarCommand) Synopsis() string {
	return "Generate project design template"
}

// Help is a long-form help text that includes the command-line
// usage, a brief few sentences explaining the function of the command,
// and the complete list of flags the command accepts.
func (c *JarCommand) Help() string {
	helpText := `
Usage: maven-runner jar [option]
  Install and run jar. You can select jar file by option.

Example: maven-runner jar -groupId com.sample -artifact program -version 0.1.0 -options "-Dproperty=value" -argument "foo bar"
  above command run following command after downloading jar file from maven central
  ==> java -Dproperty=value -jar program-0.1.0.jar foo bar

Options:
  -groupId=name, -g (required)  Jar library groupId.
  -artifact=name, -a (required) Jar library artifact name.
  -version=value, -V            Jar library version.
                                If version is not set, find latest version library.
                                If the target jar file is not donwloaded in "~/.m2" direcotry, try to download it.
  -option=value, -o             Jar execution options
  -argument=value, -A           Jar execution argument
  -jarDir=value, -O             maven-runner set jar into this dir, default is current directory
`
	return strings.TrimSpace(helpText)
}
