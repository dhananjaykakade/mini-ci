package runner

type DockerfileConfig struct {
	Language     string
	InstallCmd   string
	BuildCmd     string
	StartCmd     string
	ExposePort   int
	OutputFolder string
	Env          map[string]string
	AppType      string // Added to identify the app type
}

func GetDefaultsByAppType(appType string) DockerfileConfig {
	switch appType {
	case "react":
		return DockerfileConfig{
			Language:     "node",
			InstallCmd:   "npm install",
			BuildCmd:     "npm run build",
			StartCmd:     "npx serve dist",
			ExposePort:   3000,
			OutputFolder: "dist",
		}
	case "nextjs":
		return DockerfileConfig{
			Language:     "node",
			InstallCmd:   "npm install",
			BuildCmd:     "npm run build",
			StartCmd:     "npm run start",
			ExposePort:   3000,
			OutputFolder: ".next",
		}
	case "flask":
		return DockerfileConfig{
			Language:   "python",
			InstallCmd: "pip install -r requirements.txt",
			BuildCmd:   "",
			StartCmd:   "python app.py",
			ExposePort: 5000,
		}
	case "go":
		return DockerfileConfig{
			Language:   "go",
			InstallCmd: "",
			BuildCmd:   "go build -o app .",
			StartCmd:   "./app",
			ExposePort: 8080,
		}
	case "java":
		return DockerfileConfig{
			Language:   "java",
			InstallCmd: "mvn install",
			BuildCmd:   "mvn package",
			StartCmd:   "java -jar target/app.jar",
			ExposePort: 8080,
		}
	case "node":
		return DockerfileConfig{
			Language:   "node",
			InstallCmd: "npm install",
			StartCmd:   "npm start",
			ExposePort: 3000,
		}
	case "vite":
		return DockerfileConfig{
			Language:  "node",
			InstallCmd: "npm install",
			BuildCmd:   "npm run build",
			StartCmd:   "npx serve dist",
			ExposePort: 3000,
		}
	default:
		return DockerfileConfig{
			Language:   "node",
			InstallCmd: "npm install",
			StartCmd:   "npm start",
			ExposePort: 3000,
		}
	}
}
