# 🚀  CI/CD System (Built with Golang)

A  CI/CD system built with Golang that supports **any programming language**. It automatically detects the repository's language and executes the pipeline 

## 📋 Features
- 🌐 Supports repositories in **any language**
- ⚙️ Customizable pipeline logic in `/handlers/postcommit.go`
- 🔄 Auto-triggered pipelines on new commits
- ✅ MUST USE FORWARDING SERVICE LIKE SMEE 

 ## 🧪 How It Works
- Connect your GitHub App to the repository.
- Smee captures GitHub webhook events and forwards them to your local server.
- The system detects the repository language and executes the steps defined in ` /handlers/postcommit.go `
- Pipeline triggers automatically on each commit or pull request.
## 📁 Pipeline Configuration
- All pipeline logic is handled in ` /handlers/postcommit.go `
- Customize this file to define build, test, and deployment steps for any supported language.
##  💡 Example Smee Setup
- Visit smee.io and create a new channel.
- and run ` npx smee -u https://smee.io/your-custom-channel `


## ⚙️ Environment Setup
- Rename the environment file:
`mv .env.k .env`

- Fill in the .env file with necessary configurations:

##  ⚡ Triggering the Pipeline
- Push changes
- Smee captures the GitHub webhook and forwards it to the CI/CD system.
- The system detects the language and runs the pipeline steps.
##  💾 Clean Up
- To stop Smee and the CI/CD system:

### Stop the local Golang server
```
  Ctrl + C
```

### Stop Smee
```
 Ctrl + C
```


## 🌟 Future Enhancements
### 🖥️ Web Interface — 
Build a user-friendly dashboard to visualize pipelines, logs, and deployment status.
### 🔐 GitHub & Git Login — 
Integrate OAuth for secure GitHub/Git login to manage repositories and CI/CD configurations directly from the web interface.
