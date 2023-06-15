// Sample
pipeline {    
    agent { // Default: agent any
        node {
            label 'dev'
            // customWorkspace '/home/jenkins/factory'
        }
    }
    // options { timeout(time 3, unit: 'MINUTES') }
    
    environment { // Global Env
        // GO111MODULE = 'on'
        // AWS_SECRET_ACCESS_KEY   = credentials('')
        AWS_DEFAULT_REGION      = 'ap-northeast-2'
    }
    tools { go '1.20.3' }
    stages {
        stage('aws credentials test') {
            // local env this stage
            environment { testenv = 'testenv' }
            steps { 
                withAWS(credentials: 'ash', region: 'ap-northeast-2') {
                    sh 'echo "<h1> this is error page </h1>" > error.html'
                    s3Upload(file:'error.html', bucket:'thisiscloudfronttest', path:'web/')
                }
            }
        }
        stage('Build') {
            steps {
                sh 'go version'
                sh 'echo ${JENKINS_HOME}'
                sh 'go env'
                timeout(time: 3, unit: 'MINUTES') {
                    sh 'go run ec2count.go'
                }
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }    
}