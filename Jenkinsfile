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
    tools { go '1.20.x' }
    stages {
        // stage('prepare amd64') {
        //     agent {
        //         node {
        //             label 'amd64'
        //         }
        //     }            
        //     steps {
        //         // sh 'echo ${JENKINS_HOME}'
        //         sh 'ls -al'
        //         sh 'echo $(arch) $(hostname)'
        //         sh 'go build -o bin/ping-bin ping.go'
        //         sh 'tar zcvf ping-asset-amd64.tar.gz ./bin'
        //         script {
        //             def result = sh(script: 'sha512sum ping-asset-amd64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
        //             env.assethex = result
        //             echo result
        //         }
        //         sh 'echo "$assethex"'
        //     }
        // }        
        // stage('prepare arm64') {
        //     agent {
        //         node {
        //             label 'arm64'
        //         }
        //     }            
        //     steps {
        //         // sh 'echo ${JENKINS_HOME}'
        //         sh 'ls -al'
        //         sh 'echo $(arch) $(hostname)'
        //         sh 'go build -o bin/ping-bin ping.go'
        //         sh 'echo $PWD'
        //         sh 'tar zcvf ping-asset-arm64.tar.gz bin'
        //         script {
        //             def result = sh(script: 'sha512sum ping-asset-arm64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
        //             env.assethex = result
        //             echo result
        //         }
        //         sh 'echo "$assethex"'
        //     }
        // }
        stage('prepare') {
            parallel amd64 : {
                agent {
                    node {
                        label 'amd64'
                    }
                }            
                steps {
                    // sh 'echo ${JENKINS_HOME}'
                    sh 'ls -al'
                    sh 'echo $(arch) $(hostname)'
                    sh 'go build -o bin/ping-bin ping.go'
                    sh 'tar zcvf ping-asset-amd64.tar.gz ./bin'
                    script {
                        def result = sh(script: 'sha512sum ping-asset-amd64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
                        env.assethex = result
                        echo result
                    }
                    sh 'echo "$assethex"'
                }
            },
            arm64 : {
                agent {
                    node {
                        label 'arm64'
                    }
                }            
                steps {
                    // sh 'echo ${JENKINS_HOME}'
                    sh 'ls -al'
                    sh 'echo $(arch) $(hostname)'
                    sh 'go build -o bin/ping-bin ping.go'
                    sh 'echo $PWD'
                    sh 'tar zcvf ping-asset-arm64.tar.gz bin'
                    script {
                        def result = sh(script: 'sha512sum ping-asset-arm64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
                        env.assethex = result
                        echo result
                    }
                    sh 'echo "$assethex"'
                }

            }
        }
        // stage('compare asset file amd64') {
        //     agent { 
        //         node { 
        //             label 'amd64'
        //             // customWorkspace '/var/lib/jenkins/workspace/ping-build'
        //         } 
        //     }
        //     steps {                
        //         withAWS(credentials: 'ash', region: 'ap-northeast-2') {
        //             script {
        //                 echo env.assethex
        //                 def assetexists = s3DoesObjectExist(bucket:'thisiscloudfronttest', path:'test/ping-asset-amd64.tar.gz')
        //                 env.assetexists = assetexists
        //                 env.isdiffrent = true
                        
        //                 // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
        //                 if (env.assetexists == 'true') {
        //                     echo 'exists ping-asset-amd64.tar.gz'
        //                     sh 'rm -rf compare && mkdir compare'
        //                     s3Download(file:'compare/ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/ping-asset-amd64.tar.gz', force:true)
                            
        //                     def result = sh(script: '(sha512sum compare/ping-asset-amd64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
        //                     env.pastAssethex = result
        //                     sh 'echo $assethex'
        //                     sh 'echo $pastAssethex'                            
        //                     if (env.assethex == env.pastAssethex) {
        //                         env.isdiffrent = false
        //                     }
        //                 } else {
        //                     echo 'Not exists. Download ping-asset-amd64.tar.gz'
        //                 }
        //                 if (env.isdiffrent == true) {
        //                     echo 'asset file upload'
        //                     s3Upload(file:'ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
        //                 }
        //             }
        //         }
                
        //         // withAWS(credentials: 'ash', region: 'ap-northeast-2') {
        //         //     s3Upload(file:'ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
        //         // }
        //     }
        // }
        // stage('uplode asset upload') {
        //     agent { 
        //         node { 
        //             label 'amd64'
        //             customWorkspace '/var/lib/jenkins/workspace/ping-build'
        //         } 
        //     }
        //     steps{
        //         withAWS(credentials: 'ash', region: 'ap-northeast-2') {
        //             s3Upload(file:'ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
        //         }
        //     }
        // }
        // stage('Build arm64') {
        //     agent {
        //         node {
        //             label 'arm64'
        //             // 신규 사용자 대신 ssh 원격접속 되는 기본값으로 지정하다보니..
        //             customWorkspace '/home/ec2-user/workspace/ping-build'
        //         }
        //     }
            
        //     steps {
        //         sh 'echo ${JENKINS_HOME}' 
        //         sh 'echo $(arch) $(hostname)'
        //         sh 'ls -al'
        //         sh 'echo $PWD'
        //         sh 'ls -al bin/'
        //         sh 'tar zcvf ping-bin.tar.gz ./bin '
        //         sh 'echo $(sha512sum ping-bin.tar.gz)'
        //         sh 'assethex=$(sha512sum ping-bin.tar.gz)'
        //         sh 'echo $assethex'
        //     }
        // }
    }    
}