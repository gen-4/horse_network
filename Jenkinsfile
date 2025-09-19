node {
	def app

	script {
		System.setProperty("org.jenkinsci.plugins.durabletask.BourneShellScript.HEARTBEAT_CHECK_INTERVAL", "86400")
	}

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Compile project') {
		go build cmd/main.go
	}

	stage('Build image') {
		echo 'Building image...'
		retry(3) {
			app = docker.build("horse_network_image:latest")
		}
		echo 'Image built'
	}

	stage('Deploying horse_network app') {
		
	}
}
