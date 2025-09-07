node {
	def app

	stage('Clone repository') {
		echo 'Cloning repository...'
		checkout scm
		echo 'Repository cloned'
	}

	stage('Build image') {
		echo 'Building image...'
		app = docker.build("horse_network_image/latest")
		echo 'Image built'
	}

	stage('Deploying horse_network app') {
		
	}
}
