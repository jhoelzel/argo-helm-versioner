package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

// struct to hold the result for each application
type AppCheckResult struct {
	FilePath       string
	Application    string
	CurrentVersion string
	LatestVersion  string
	Status         string
}

// New function to recursively walk through the directories
func walkDir(dirPath string, processFile func(filename string) error) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			return processFile(path)
		}
		return nil
	})
}

// Function to process each YAML file and check if it's an ArgoCD Application
func processArgoCDApplicationFile(filePath string, results *[]AppCheckResult) error {
	app, err := loadArgoApplicationFromFile(filePath)
	if err != nil {
		return nil // File might not be an Argo CD Application, skip it
	}

	if app.Spec.Source.Chart != "" {
		fmt.Printf("Checking Chart: %v\n", app.Spec.Source.Chart)
		status := "Up-to-date"
		latestVersion, err := checkLatestHelmVersion(app.Spec.Source.RepoURL, app.Spec.Source.Chart)
		if err != nil {
			status = "Error"
		}

		if latestVersion != app.Spec.Source.TargetRevision {
			status = "Update available"
		}

		*results = append(*results, AppCheckResult{
			FilePath:       filePath,
			Application:    app.Metadata.Name,
			CurrentVersion: app.Spec.Source.TargetRevision,
			LatestVersion:  latestVersion,
			Status:         status,
		})
	}
	return nil
}

// Define a struct to match the Argo CD Application spec for Helm
type ArgoApplication struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Source struct {
			RepoURL        string `yaml:"repoURL"`
			TargetRevision string `yaml:"targetRevision"`
			Chart          string `yaml:"chart"`
		} `yaml:"source"`
	} `yaml:"spec"`
}

// Define a struct for chart versions that will allow us to use semantic version sorting
type ChartVersion struct {
	Version *semver.Version
}

// Function to parse a string into a ChartVersion, ignoring any invalid versions
func NewChartVersion(versionString string) *ChartVersion {
	v, err := semver.NewVersion(versionString)
	if err != nil {
		return nil
	}
	return &ChartVersion{Version: v}
}

type ChartVersions []*ChartVersion

// Implement sort.Interface for ChartVersions
func (c ChartVersions) Len() int           { return len(c) }
func (c ChartVersions) Less(i, j int) bool { return c[i].Version.LessThan(c[j].Version) }
func (c ChartVersions) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type HelmRepoIndex struct {
	Entries map[string][]struct {
		Version string `yaml:"version"`
	} `yaml:"entries"`
}

// Function to load YAML from a file and unmarshal into ArgoApplication struct
func loadArgoApplicationFromFile(filename string) (*ArgoApplication, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var app ArgoApplication
	err = yaml.Unmarshal(buf, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}

// Function to check for the latest Helm chart version from the Helm repository
func checkLatestHelmVersion(repoURL, chartName string) (string, error) {
	resp, err := http.Get(repoURL + "/index.yaml")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var index HelmRepoIndex
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(bodyBytes, &index)
	if err != nil {
		return "", err
	}

	// Parse versions and sort them using semantic versioning
	var versions ChartVersions
	for _, entry := range index.Entries[chartName] {
		cv := NewChartVersion(entry.Version)
		if cv != nil {
			versions = append(versions, cv)
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no valid versions found for chart %s", chartName)
	}

	sort.Sort(versions)
	latestVersion := versions[len(versions)-1].Version.String()

	return latestVersion, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: <program> <directory path>")
		os.Exit(1)
	}
	dirPath := os.Args[1]

	var results []AppCheckResult
	err := walkDir(dirPath, func(filePath string) error {
		return processArgoCDApplicationFile(filePath, &results)
	})
	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
		return
	}

	// Sort results by Application and then by Status
	sort.Slice(results, func(i, j int) bool {
		if results[i].Application == results[j].Application {
			return results[i].Status < results[j].Status
		}
		return results[i].Application < results[j].Application
	})

	// Print results in a table with Application as the first column
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Application\tFilePath\tCurrent Version\tLatest Version\tStatus\t")
	for _, result := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", result.Application, result.FilePath, result.CurrentVersion, result.LatestVersion, result.Status)
	}
	w.Flush()
}
