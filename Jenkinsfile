// Sample
pipeline {    
    // 각 스텝마다 노드를 정의하기 때문에 선언 안함. 만약 any일 경우 모든 노드 중 임의 하나에서만 동작함
    agent none
    // agent {
    //     node {
    //         label 'amd64'
    //         // customWorkspace '/home/jenkins/factory'
    //     }
    // }
    // options { timeout(time 2, unit: 'MINUTES') }
    options {
        // 하나라도 실패되었을 경우 모든 병렬 실행중인 job을 중단하고 결과를 실패처리
        parallelsAlwaysFailFast()
    }
    
    environment { // Global Env 
        // GO111MODULE = 'on'
        // AWS_SECRET_ACCESS_KEY   = credentials('')
        AWS_ACCESS_KEY_ID = credentials('aws_access_key')
        AWS_SECRET_ACCESS_KEY = credentials('aws_secret_key')
        AWS_DEFAULT_REGION = 'ap-northeast-2'
    }
    tools { go '1.20.x' }
    stages {
        // checkout scm 이 코드 버전관리 명령이라는데 어케 쓰는지 확인 필요
        stage ('test') {
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
                stage('go build amd64') {
                    agent {
                        node {
                            label 'amd64'
                        }
                    }            
                    steps {
                        // sh 'echo ${JENKINS_HOME}'
                        sh 'ls -al'
                        sh 'echo $(arch) $(hostname)'
                        sh 'go build -v -o bin/ping-bin ping.go'
                        sh 'tar zcvf ping-asset-amd64.tar.gz ./bin' 
                        script {
                            def linux_amd64_hex = sh(script: 'sha512sum ping-asset-amd64.tar.gz | awk \'{print $1}\'', returnStdout: true).trim()
                            env.linux_amd64_hex = linux_amd64_hex
                            echo linux_amd64_hex
                        }
                        sh 'echo "$linux_amd64_hex"'
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
                        sh 'echo "$linux_arm64_hex"'
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
                            //def assetexists = s3DoesObjectExist(bucket:'thisiscloudfronttest', path:'test/ping-asset-amd64.tar.gz')
                            def assetexists = sh(script: 'aws s3 ls s3://thisiscloudfronttest/test/ping-asset-amd64.tar.gz >/dev/null 2>&1 && echo true', returnStdout: true).trim()
                            echo assetexists
                            env.assetexists = assetexists
                            
                            // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
                            if (env.assetexists == 'true') {
                                echo 'exists ping-asset-amd64.tar.gz'
                                sh 'rm -rf compare && mkdir compare'
                                
                                
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
                                s3Upload(file:'ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
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
                        withAWS(credentials: 'ash', region: 'ap-northeast-2') {
                            script {
                                env.isdiffrent = true
                                sh 'echo "new asset hex: $linux_arm64_hex"'
                                def assetexists = s3DoesObjectExist(bucket:'thisiscloudfronttest', path:'test/ping-asset-arm64.tar.gz')
                                env.assetexists = assetexists
                                
                                // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
                                if (env.assetexists == 'true') {
                                    echo 'exists ping-asset-arm64.tar.gz'
                                    sh 'rm -rf compare && mkdir compare'
                                    s3Download(file:'compare/ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/ping-asset-arm64.tar.gz', force:true)
                                    
                                    def result = sh(script: '(sha512sum compare/ping-asset-arm64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
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
                                    s3Upload(file:'ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
                                } else {
                                    echo 'same file'
                                }
                            }
                        }
                    }
                }
            }
        }
        // stage ('asset compare') {
        //     parallel {
        //         stage ('asset compare amd64') {
        //             agent { 
        //                 node { 
        //                     label 'amd64'
        //                 } 
        //             }
        //             steps {
        //                 withAWS(credentials: 'ash', region: 'ap-northeast-2') {
        //                     script {
        //                         env.isdiffrent = true
        //                         sh 'echo "new asset hex: $linux_amd64_hex"'
        //                         def assetexists = s3DoesObjectExist(bucket:'thisiscloudfronttest', path:'test/ping-asset-amd64.tar.gz')
        //                         env.assetexists = assetexists
                                
        //                         // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
        //                         if (env.assetexists == 'true') {
        //                             echo 'exists ping-asset-amd64.tar.gz'
        //                             sh 'rm -rf compare && mkdir compare'
        //                             s3Download(file:'compare/ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/ping-asset-amd64.tar.gz', force:true)
                                    
        //                             def result = sh(script: '(sha512sum compare/ping-asset-amd64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
        //                             env.pastAssethex = result
        //                             sh 'echo $linux_amd64_hex'
        //                             sh 'echo $pastAssethex'
        //                             if (env.linux_amd64_hex == env.pastAssethex) {
        //                                 echo 'same asset hex'
        //                                 env.isdiffrent = false
        //                             } else {
        //                                 echo 'not same asset hex'
        //                             }
        //                         } else {
        //                             echo 'Not exists. Download ping-asset-amd64.tar.gz'
        //                         }
        //                         if (env.isdiffrent == 'true') {
        //                             echo 'Asset file upload'
        //                             s3Upload(file:'ping-asset-amd64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
        //                         } else {
        //                             echo 'same file'
        //                         }
        //                     }
        //                 }
        //             }
        //         }
        //         stage ('asset compare arm64') {
        //             agent { 
        //                 node { 
        //                     label 'arm64'
        //                 } 
        //             }
        //             steps {
        //                 withAWS(credentials: 'ash', region: 'ap-northeast-2') {
        //                     script {
        //                         env.isdiffrent = true
        //                         sh 'echo "new asset hex: $linux_arm64_hex"'
        //                         def assetexists = s3DoesObjectExist(bucket:'thisiscloudfronttest', path:'test/ping-asset-arm64.tar.gz')
        //                         env.assetexists = assetexists
                                
        //                         // s3에 업로드 된 에셋 압축파일이 있다면 새로 생성된 파일이랑 내용이 달라졌는지 확인
        //                         if (env.assetexists == 'true') {
        //                             echo 'exists ping-asset-arm64.tar.gz'
        //                             sh 'rm -rf compare && mkdir compare'
        //                             s3Download(file:'compare/ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/ping-asset-arm64.tar.gz', force:true)
                                    
        //                             def result = sh(script: '(sha512sum compare/ping-asset-arm64.tar.gz | awk \'{print $1}\')', returnStdout: true).trim()
        //                             env.pastAssethex = result
        //                             sh 'echo $linux_arm64_hex'
        //                             sh 'echo $pastAssethex'
        //                             if (env.linux_arm64_hex == env.pastAssethex) {
        //                                 echo 'same asset hex'
        //                                 env.isdiffrent = false
        //                             } else {
        //                                 echo 'not same asset hex'
        //                             }
        //                         } else {
        //                             echo 'Not exists. Must be upload ping-asset-arm64.tar.gz'
        //                         }
        //                         if (env.isdiffrent == 'true') {
        //                             echo 'asset file upload'
        //                             s3Upload(file:'ping-asset-arm64.tar.gz', bucket:'thisiscloudfronttest', path:'test/')
        //                         } else {
        //                             echo 'same file'
        //                         }
        //                     }
        //                 }
        //             }
        //         }
        //     }
        // }
        
        // ping-asset.yaml 코드를 업데이트
        // stage('update ping-asset.yaml') {
        stage('update ping-asset.yaml') {
            agent any
            steps {
                sh 'echo $linux_amd64_hex'
                sh 'echo $linux_arm64_hex'
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
      url: https://thisiscloudfronttest.s3.ap-northeast-2.amazonaws.com/ping-asset-amd64.tar.gz
      filters:
      - entity.system.os == 'linux'
      - entity.system.arch == 'amd64'
    - sha512 : $linux_arm64_hex
      url: https://thisiscloudfronttest.s3.ap-northeast-2.amazonaws.com/ping-asset-arm64.tar.gz
      filters:
      - entity.system.os == 'linux'
      - entity.system.arch == 'arm64'
                    '''
                }
                sh 'cat ./ping-asset.yaml'
                sh 'aws s3 cp ./ping-asset.yaml s3://thisiscloudfronttest/test/'
                sh 'aws s3 ls s3://thisiscloudfronttest/test/ping-asset.yaml'
                // s3Upload(file:'ping-asset.yaml', bucket:'thisiscloudfronttest', path:'test/')
            }
        }
    }    
}