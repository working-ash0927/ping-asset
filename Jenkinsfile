// Sample
pipeline {    
    // Default: agent node any
    agent none
    // agent {
    //     node {
    //         label 'amd64'
    //         // customWorkspace '/home/jenkins/factory'
    //     }
    // }
    // options { timeout(time 3, unit: 'MINUTES') }
    
    environment { // Global Env 
        // GO111MODULE = 'on'
        // AWS_SECRET_ACCESS_KEY   = credentials('')
        AWS_DEFAULT_REGION      = 'ap-northeast-2'
    }
    // tools { go '1.20.x' }
    stages {
        stage('Build amd64') {
            agent {
                node {
                    label 'amd64'
                    // customWorkspace '/home/jenkins/factory'
                }
            }
            
            steps {
                // sh 'go version'
                sh 'echo ${JENKINS_HOME}' 
                // sh 'go env' 
                sh 'ls -al' // repo 최상위 경로
                sh 'echo $PWD'
                sh 'echo $(arch)'
                sh 'echo $(hostname)'
                // sh 'go build -o ./bin/ping-bin ping.go' 
                // sh 'ls -al'
                // sh 'tar zcvf ping-bin.tar.gz ./bin '
                // sh 'echo $(sha512sum ping-bin.tar.gz)'
                // sh 'assethex=$(sha512sum ping-bin.tar.gz)'
                // sh 'echo $assethex'
                // timeout(time: 3, unit: 'MINUTES') {
                //     sh 'go run ec2count.go'
                // }
            }
        }
        stage('Build arm64') {
            agent {
                node {
                    label 'arm64'
                    // customWorkspace '/home/jenkins/factory'
                }
            }
            
            steps {
                sh 'echo ${JENKINS_HOME}' 
                sh 'ls -al'
                sh 'echo $PWD'
                sh 'echo $(arch)'
                sh 'echo $(hostname)'
            }
        }
    }    
}