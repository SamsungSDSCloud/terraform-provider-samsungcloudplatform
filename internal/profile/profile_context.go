package profile

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

const ConfigFilename = ".configurations"
const CredFilename = ".credentials"
const LockFilename = ".cmplock"

// ProfileContext information
type ProfileContext struct {
	// Configuration file name
	ConfigFileName string
	// Credentials file name
	CredFileName string
	// Lock file name
	LockFileName string
	// Configuration direction
	ConfigDirectory string
}

// NewProfileContext Create default ProfileContext
func NewProfileContext() ProfileContext {
	ctx := ProfileContext{}
	ctx.ConfigFileName = ConfigFilename
	ctx.CredFileName = CredFilename
	ctx.LockFileName = LockFilename

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Failed to get user home directory")
		homeDir = ""
	}
	ctx.ConfigDirectory = homeDir + string(os.PathSeparator) + ".cmp"
	return ctx
}

// GetLockFilePath Get lock file path
func (ctx *ProfileContext) GetLockFilePath() string {
	return ctx.ConfigDirectory + string(os.PathSeparator) + ctx.LockFileName
}

// GetCredFilePath Get credentials file path
func (ctx *ProfileContext) GetCredFilePath() string {
	return ctx.ConfigDirectory + string(os.PathSeparator) + ctx.CredFileName
}

// GetConfigFilePath Get configuration file path
func (ctx *ProfileContext) GetConfigFilePath() string {
	return ctx.ConfigDirectory + string(os.PathSeparator) + ctx.ConfigFileName
}

// EnsureConfigDirectory Ensure configuration directory is present
func (ctx *ProfileContext) EnsureConfigDirectory() error {
	// Create directory if not present
	return os.MkdirAll(ctx.ConfigDirectory, os.ModeDir|os.ModePerm)
}

// EnsureLockFile Ensure file lock state
func (ctx *ProfileContext) EnsureLockFile() error {
	lockFilePath := ctx.GetLockFilePath()
	_, err := os.Stat(lockFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(lockFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		currentTime := time.Now().Local()
		err = os.Chtimes(lockFilePath, currentTime, currentTime)
		if err != nil {
			fmt.Println(err)
		}
	}
	return err
}

func (ctx *ProfileContext) LoadProfiles() (Profiles, error) {
	profiles := NewProfiles()

	// Load options
	credProfile, err := loadProfileMap("credentials", ctx.GetCredFilePath())
	if err == nil {
		profiles.ProfileMap["credentials"] = credProfile
	}
	configProfile, err := loadProfileMap("configurations", ctx.GetConfigFilePath())
	if err == nil {
		profiles.ProfileMap["configurations"] = configProfile
	}

	return profiles, nil
}

func loadProfileMap(name string, filePath string) (Profile, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalln("Failed to open profile file")
		return NewProfileWithName(name), err
	}
	defer f.Close()

	re, err := regexp.Compile("^\\[(.+)\\]$")
	if err != nil {
		log.Fatalln("Failed to generate category regex parser")
		return NewProfileWithName(name), err
	}

	profile := NewProfileWithName(name)
	category := ""

	reader := bufio.NewReader(f)
	for {
		line, prefix, err := reader.ReadLine()
		if prefix || err != nil {
			break
		}
		// Empty string
		if len(line) == 0 {
			continue
		}
		// Check
		foundIndices := re.FindIndex(line)
		if foundIndices != nil {
			// Assume first
			category = string(line)
			continue
		}
		profile.AddProperty(category, string(line))
	}

	return profile, nil
}
