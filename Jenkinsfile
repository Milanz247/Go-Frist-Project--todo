// Jenkinsfile for a Go Application

pipeline {
    agent any 

    stages {
        stage('Checkout') {
            steps {
                echo 'Checking out code from repository...'
                git 'https://github.com/Milanz247/Go-Frist-Project--todo.git'
            }
        }

        stage('Download Dependencies') {
            steps {
                echo 'Downloading Go modules...'
                sh 'go mod tidy'
                sh 'go mod download'
            }
        }

        // Stage 3: Code එක Test කිරීම
        stage('Test') {
            steps {
                echo 'Running Go tests...'
                sh 'go test ./...'
            }
        }

        stage('Build') {
            steps {
                echo 'Building the Go application...'
                sh 'go build -o go-todo-app'
            }
        }


        stage('Deploy') {
            steps {
                echo 'Deploying the application...'
                echo 'Stopping any old running process...'
                sh 'pkill go-todo-app || true'
                
                echo 'Starting the new application in the background...'
         
                sh 'nohup ./go-todo-app &'
            }
        }
    }
    
    post {
        always {
            echo 'Pipeline has finished.'
        }
    }
}
