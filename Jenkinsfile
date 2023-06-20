// Sample
pipeline {    
    // 각 스텝마다 노드를 정의하기 때문에 선언 안함. 만약 any일 경우 모든 노드 중 임의 하나에서만 동작함
    agent none
    options {
        // 하나라도 실패되었을 경우 모든 병렬 실행중인 job을 중단하고 결과를 실패처리
        parallelsAlwaysFailFast()
    }
    
    environment { // Global Env 
        // GO111MODULE = 'on'
        AWS_ACCESS_KEY_ID = credentials('aws_access_key')       // Jenkins Secret Text
        AWS_SECRET_ACCESS_KEY = credentials('aws_secret_key')   // Jenkins Secret Text
        AWS_DEFAULT_REGION = 'ap-northeast-2'
    }
    tools { go '1.20.x' }   // Go Plugin
    stages {
        // checkout scm 이 코드 버전관리 명령이라는데 어케 쓰는지 확인 필요
        // 기본 경로 {JENKINS_HOME}/workspace/pipeline/{pipeline이름}/ 내에 git clone 됨
        stage ('Check out') {
            agent {
                node {
                    label 'amd64'
                }
            }
            steps {
                checkout scm
                sh 'ls -ali'
            }
        }
        // 변경되지 않은 소스코드임에도 clone되면서 변경된 inode로 인해 해시값이 변경되는 걸 수정해야함
        stage ('build asset') {
            parallel {
                stage('go build amd64') { // 해당 step은 label이 amd64인 에이전트에서만 동작 선언
                    agent {
                        node {
                            label 'amd64'
                        }
                    }            
                    steps {  
                        // sh 'echo ${JENKINS_HOME}'    // 기본값 미변경시 /var/lib/jenkins/
                        sh 'go build -v -o bin/ping-bin ping.go'
                        sh 'tar zcvf ping-asset-amd64.tar.gz ./bin' 
                        script {
                            def linux_amd64_hex = sh(script: 'sha512sum ping-asset-amd64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
                            env.linux_amd64_hex = linux_amd64_hex
                            echo linux_amd64_hex
                        }
                    }
                }
                stage('go build arm64') {
                    agent {
                        node {
                            label 'arm64'
                        }
                    }            
                    steps {
                        // sh 'echo ${JENKINS_HOME}'
                        sh 'go build -v -o bin/ping-bin ping.go'
                        sh 'tar zcvf ping-asset-arm64.tar.gz ./bin' 
                        script {
                            def linux_arm64_hex = sh(script: 'sha512sum ping-asset-arm64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
                            env.linux_arm64_hex = linux_arm64_hex
                            echo linux_arm64_hex
                        }
                    }
                }
                stage('go build win_amd64') {
                    agent {
                        node {
                            label 'win_amd64'
                        }
                    }            
                    steps {
                        // sh 'echo ${JENKINS_HOME}'
                        powershell 'go build -v -o .\\bin\\ping-bin.exe ping.go'
                        powershell '''
                        tar zcvf ping-asset-win-amd64.tar.gz .\\bin
                        exit 0
                        '''
                        script {
                            def win_amd64_hex = powershell(script: '(Get-FileHash -Path ping-asset-win.tar.gz -Algorithm SHA512).Hash', returnStdout: true)
                            env.win_amd64_hex = win_amd64_hex
                            println win_amd64_hex
                        }
                    }
                }
            }
        }
        stage ('asset compare') {
            parallel {
                stage ('asset compare amd64') {
                    agent { 
                        node { 
                            label 'amd64'
                        } 
                    }
                    steps {
                        script {
                            env.isdiffrent = true
                            sh 'echo "new asset hex: $linux_amd64_hex"'
                            // 파일 유무 확인. (true/false)
                            def assetexists = sh(script: 'aws s3 ls s3://thisiscloudfronttest/test/ping-asset-amd64.tar.gz >/dev/null 2>&1 && echo true || echo false', returnStdout: true).trim()
                            echo assetexists
                            env.assetexists = assetexists
                            
                            // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
                            if (env.assetexists == 'true') {
                                echo 'exists ping-asset-amd64.tar.gz'
                                sh 'rm -rf compare && mkdir compare'
                                sh 'aws s3 cp s3://thisiscloudfronttest/test/ping-asset-amd64.tar.gz compare/ping-asset-amd64.tar.gz'
                                def result = sh(script: '(sha512sum compare/ping-asset-amd64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()                                
                                echo result
                                env.pastAssethex = result

                                sh 'echo $linux_amd64_hex'
                                sh 'echo $pastAssethex'
                                if (env.linux_amd64_hex == env.pastAssethex) {
                                    echo 'same asset hex'
                                    env.isdiffrent = false
                                } else {
                                    echo 'not same asset hex'
                                }
                            } else {
                                echo 'Not exists. Download ping-asset-amd64.tar.gz'
                            }
                            if (env.isdiffrent == 'true') {
                                echo 'Asset file upload'
                                sh 'aws s3 cp ping-asset-amd64.tar.gz s3://thisiscloudfronttest/test/ping-asset-amd64.tar.gz --acl public-read'
                            } else {
                                echo 'same file'
                            }
                        }
                    }
                }
                stage ('asset compare arm64') {
                    agent { 
                        node { 
                            label 'arm64'
                        } 
                    }
                    steps {
                        script {
                            env.isdiffrent = true
                            sh 'echo "new asset hex: $linux_arm64_hex"'
                            def assetexists = sh(script: 'aws s3 ls s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz >/dev/null 2>&1 && echo true || echo false', returnStdout: true).trim()
                            env.assetexists = assetexists
                            
                            // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
                            if (env.assetexists == 'true') {
                                echo 'exists ping-asset-arm64.tar.gz'
                                sh 'rm -rf compare && mkdir compare'
                                sh 'aws s3 cp s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz compare/ping-asset-arm64.tar.gz'
                                def result = sh(script: '(sha512sum compare/ping-asset-arm64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
                                echo result
                                env.pastAssethex = result

                                sh 'echo $linux_arm64_hex'
                                sh 'echo $pastAssethex'
                                if (env.linux_arm64_hex == env.pastAssethex) {
                                    echo 'same asset hex'
                                    env.isdiffrent = false
                                } else {
                                    echo 'not same asset hex'
                                }
                            } else {
                                echo 'Not exists. Must be upload ping-asset-arm64.tar.gz'
                            }
                            if (env.isdiffrent == 'true') {
                                echo 'asset file upload'
                                // s3Upload(file:'ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
                                sh 'aws s3 cp ping-asset-arm64.tar.gz s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz --acl public-read'
                            } else {
                                echo 'same file'
                            }
                        }
                    }
                }
                // 스크립트 윈도우 버전으로 갱신 필요
                stage ('asset compare win_amd64') {
                    agent { 
                        node { 
                            label 'win_amd64'
                        } 
                    }
                    steps {
                        script {
                            bat(script: 'aws s3 ls')
                            // def test = powershell(script: 'aws s3 cp ping-asset-win-amd64.tar.gz s3://thisiscloudfronttest/test/ping-asset-win-amd64.tar.gz --acl public-read', returnStdout: true)
                            // println test
                            // env.isdiffrent = true
                            // sh 'echo "new asset hex: $win_amd64_hex"'
                            // def assetexists = sh(script: 'aws s3 ls s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz >/dev/null 2>&1 && echo true || echo false', returnStdout: true).trim()
                            // env.assetexists = assetexists
                            
                            // // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
                            // if (env.assetexists == 'true') {
                            //     echo 'exists ping-asset-arm64.tar.gz'
                            //     sh 'rm -r compare && mkdir compare'
                            //     sh 'aws s3 cp s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz compare/ping-asset-arm64.tar.gz'
                            //     def result = sh(script: '(sha512sum compare/ping-asset-arm64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
                            //     echo result
                            //     env.pastAssethex = result

                            //     sh 'echo $win_amd64_hex'
                            //     sh 'echo $pastAssethex'
                            //     if (env.win_amd64_hex == env.pastAssethex) {
                            //         echo 'same asset hex'
                            //         env.isdiffrent = false
                            //     } else {
                            //         echo 'not same asset hex'
                            //     }
                            // } else {
                            //     echo 'Not exists. Must be upload ping-asset-arm64.tar.gz'
                            // }
                            // if (env.isdiffrent == 'true') {
                            //     echo 'asset file upload'
                            //     // s3Upload(file:'ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
                            //     sh 'aws s3 cp ping-asset-arm64.tar.gz s3://thisiscloudfronttest/test/ping-asset-arm64.tar.gz --acl public-read'
                            // } else {
                            //     echo 'same file'
                            // }
                        }
                    }
                }
            }
        }
        // ping-asset.yaml 코드를 업데이트
        // 추후 버전정보 추가 고려해보기 https://osc131.tistory.com/90
        stage('update ping-asset.yaml') {
            agent any
            steps {
                // Jenkins 내 등록된 환경변수라 아무데서나 참조 가능
                sh 'echo $linux_amd64_hex'
                sh 'echo $linux_arm64_hex'
                sh 'echo $win_amd64_hex'
                script {
                    sh '''tee ./ping-asset.yaml << EOF 
---
type: Asset
api_version: core/v2
metadata:
  name: ping-asset
spec:
  builds:
    - sha512 : $linux_amd64_hex
      url: https://thisiscloudfronttest.s3.ap-northeast-2.amazonaws.com/test/ping-asset-amd64.tar.gz
      filters:
      - entity.system.os == 'linux'
      - entity.system.arch == 'amd64'
    - sha512 : $linux_arm64_hex
      url: https://thisiscloudfronttest.s3.ap-northeast-2.amazonaws.com/test/ping-asset-arm64.tar.gz
      filters:
      - entity.system.os == 'linux'
      - entity.system.arch == 'arm64'
    - sha512 : $win_amd64_hex
      url: https://thisiscloudfronttest.s3.ap-northeast-2.amazonaws.com/test/ping-asset-win-amd64.tar.gz
      filters:
      - entity.system.os == 'windows'
      - entity.system.arch == 'amd64'
                    '''
                }
                sh 'cat ./ping-asset.yaml'
                sh 'aws s3 cp ./ping-asset.yaml s3://thisiscloudfronttest/test/ping-asset.yaml --acl public-read'
            }
        }
        // post {
        //     always {
        //         echo 'One way or another, I have finished'
        //         deleteDir() /* clean up our workspace */
        //     }
        //     success {
        //         echo 'I succeeded!'
        //     }
        //     unstable {
        //         echo 'I am unstable :/'
        //     }
        //     failure {
        //         echo 'I failed :('
        //     }
        //     changed {
        //         echo 'Things were different before...'
        //     }
        // }  
    }
    
}

