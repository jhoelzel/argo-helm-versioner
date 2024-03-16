# Argo Helm Versioner

## Overview
The Argo Helm Versioner is a utility designed to help manage and maintain Argo CD Applications deployed via Helm charts. It automates the process of scanning for Argo CD Application YAML files within a specified directory (and its subdirectories), identifies those that are Helm-based, and checks their currently deployed Helm chart versions against the latest versions available in their respective Helm repositories. The findings are neatly presented in a table, sorted by application name and update status, facilitating easy identification of applications that may require updates.

## Features
- **Deep Directory Scanning**: Efficiently traverses specified directories and subdirectories to find Argo CD Application YAML files.
- **Helm Chart Version Comparison**: Determines whether deployed Helm-based applications are using the latest available chart versions, helping ensure applications are up-to-date with the latest features and security patches.
- **Sorted Tabular Output**: Organizes the results in a user-friendly table format, sorted by application name and update status, for quick and easy review.
- **Semantic Version Support**: Employs semantic versioning principles to accurately compare Helm chart versions, ensuring reliable version differentiation.

## Prerequisites
To use the Argo Helm Versioner, ensure the following prerequisites are met:
- **Go Environment**: Go version 1.15 or later must be installed on your machine.
- **Internet Access**: Required for accessing online Helm repositories to fetch the latest chart version information.
- **Helm Repository Format**: The Helm chart repositories should adhere to the standard `index.yaml` file format for version listings.

## Installation
Follow these steps to set up the Argo Helm Versioner:

1. **Obtain the Source Code**: Clone this repository or download the source code to your local machine.
2. **Change to the Project Directory**: Navigate to the directory containing the Go source file for the versioner.
3. **Install Dependencies**: Execute the following command to install necessary dependencies, such as the semantic versioning package:
   ```sh
   go get github.com/Masterminds/semver
   go get gopkg.in/yaml.v2
   ```
4. **Compile the Tool** (optional): For convenience, you may compile the versioner into an executable:
   ```sh
   go build -o argo-helm-versioner
   ```

## Usage
The Argo Helm Versioner can be run directly using `go run`, or via the compiled executable if you chose to build it.

### Running with Go Run
Execute the following command, replacing `<path_to_go_file>` with the actual path to the Go script, and `<directory_path>` with the directory containing your Argo CD Application files:
```sh
go run <path_to_go_file> <directory_path>
```

### Using the Compiled Executable
If you compiled the tool, you can run it like so:
```sh
./argo-helm-versioner <directory_path>
```
Here, `<directory_path>` should be replaced with the path to the directory you wish to scan.

### Output Format
The versioner outputs a table with the following columns for clarity and ease of use:
- **Application**: Name of the Argo CD Application.
- **FilePath**: Location of the Argo CD Application YAML file.
- **Current Version**: The version of the Helm chart currently deployed.
- **Latest Version**: The latest available version of the Helm chart in the repository.
- **Status**: Indicates if the Helm chart is "Up-to-date" or if an "Update available".

The entries in the table are sorted by the Application name and then by the Status.

## Example Output
```
Application         FilePath                            Current Version   Latest Version   Status           
my-app              /apps/my-app/argo-app.yaml          1.0.0             1.2.0            Update available  
your-app            /apps/your-app/argo-app.yaml        2.1.0             2.1.0            Up-to-date       
...
```

## Contributing
Contributions are welcome! Feel free to open issues or pull requests to propose changes, enhancements, or bug fixes to the Argo Helm Versioner.

## License
The Argo Helm Versioner is released under the MIT License. See the LICENSE file for more information.

---
