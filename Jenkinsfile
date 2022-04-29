@Library('vega-shared-library') _

/* properties of scmVars (example):
    - GIT_BRANCH:PR-40
    - GIT_COMMIT:05a1c6fbe7d1ff87cfc40a011a63db574edad7e6
    - GIT_PREVIOUS_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_PREVIOUS_SUCCESSFUL_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_URL:https://github.com/vegaprotocol/vegawallet.git
*/
def scmVars = null
def version = 'UNKNOWN'
def versionHash = 'UNKNOWN'
def commitHash = 'UNKNOWN'

pipeline {
    agent any
    options {
        skipDefaultCheckout true
        timestamps()
        timeout(time: 45, unit: 'MINUTES')
    }
    parameters {
        string( name: 'VEGA_CORE_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/vega repository.
                    e.g. "develop", "v0.44.0" or commit hash. Default empty: use latests published version.''')
        string( name: 'DATA_NODE_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/data-node repository.
                    e.g. "develop", "v0.44.0" or commit hash. Default empty: use latests published version.''')
        string( name: 'ETHEREUM_EVENT_FORWARDER_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/ethereum-event-forwarder repository.
                    e.g. "main", "v0.44.0" or commit hash. Default empty: use latest published version.''')
        string( name: 'DEVOPS_INFRA_BRANCH', defaultValue: 'master',
                description: 'Git branch, tag or hash of the vegaprotocol/devops-infra repository')
        string( name: 'VEGATOOLS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/vegatools repository')
        string( name: 'SYSTEM_TESTS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/system-tests repository')
        string( name: 'PROTOS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/protos repository')
    }
    environment {
        CGO_ENABLED = 0
        GO111MODULE = 'on'
    }

    stages {
        stage('Config') {
            steps {
                cleanWs()
                sh 'printenv'
                echo "params=${params}"
                echo "isPRBuild=${isPRBuild()}"
                script {
                    params = pr.injectPRParams()
                }
                echo "params (after injection)=${params}"
            }
        }

        stage('Git clone') {
            options { retry(3) }
            steps {
                script {
                    scmVars = checkout(scm)
                    versionHash = sh (returnStdout: true, script: "echo \"${scmVars.GIT_COMMIT}\"|cut -b1-8").trim()
                    version = sh (returnStdout: true, script: "git describe --tags 2>/dev/null || echo ${versionHash}").trim()
                    commitHash = getCommitHash()
                }
                echo "scmVars=${scmVars}"
                echo "commitHash=${commitHash}"
            }
        }

        stage('Download dependencies') {
            options { retry(3) }
            steps {
                sh 'go mod download -x'
            }
        }

        stage('Compile') {
            environment {
                LDFLAGS = "-X code.vegaprotocol.io/vegawallet/version.Version=\"${version}\" -X code.vegaprotocol.io/vegawallet/version.VersionHash=\"${versionHash}\""
            }
            failFast true
            parallel {
                stage('Linux build') {
                    environment {
                        GOOS    = 'linux'
                        GOARCH  = 'amd64'
                        OUTPUT  = './build/vegawallet-linux-amd64'
                    }
                    options { retry(3) }
                    steps {
                        sh label: 'Compile', script: '''
                            go build -v -o "${OUTPUT}" -ldflags "${LDFLAGS}"
                        '''
                        sh label: 'Sanity check', script: '''
                            file ${OUTPUT}
                            ${OUTPUT} version --output json
                        '''
                    }
                }
                stage('MacOS build') {
                    environment {
                        GOOS    = 'darwin'
                        GOARCH  = 'amd64'
                        OUTPUT  = './build/vegawallet-darwin-amd64'
                    }
                    options { retry(3) }
                    steps {
                        sh label: 'Compile', script: '''
                            go build -v -o "${OUTPUT}" -ldflags "${LDFLAGS}"
                        '''
                        sh label: 'Sanity check', script: '''
                            file ${OUTPUT}
                        '''
                    }
                }
                stage('Windows build') {
                    environment {
                        GOOS    = 'windows'
                        GOARCH  = 'amd64'
                        OUTPUT  = './build/vegawallet-windows-amd64'
                    }
                    options { retry(3) }
                    steps {
                        sh label: 'Compile', script: '''
                            go build -v -o "${OUTPUT}" -ldflags "${LDFLAGS}"
                        '''
                        sh label: 'Sanity check', script: '''
                            file ${OUTPUT}
                        '''
                    }
                }
            }
        }

        stage('Tests') {
            parallel {
                stage('unit tests') {
                    options { retry(3) }
                    steps {
                        sh 'go test -v ./... 2>&1 | tee unit-test-results.txt && cat unit-test-results.txt | go-junit-report > unit-test-report.xml'
                        junit checksName: 'Unit Tests', testResults: 'unit-test-report.xml'
                    }
                }
                stage('unit tests with race') {
                    environment {
                        CGO_ENABLED = 1
                    }
                    options { retry(3) }
                    steps {
                        sh 'go test -v -race ./... 2>&1 | tee unit-test-race-results.txt && cat unit-test-race-results.txt | go-junit-report > unit-test-race-report.xml'
                        junit checksName: 'Unit Tests with Race', testResults: 'unit-test-race-report.xml'
                    }
                }
                stage('linters') {
                    steps {
                        sh '''#!/bin/bash -e
                            golangci-lint run -v --config .golangci.toml
                        '''
                    }
                }
                stage('System Tests') {
                    steps {
                        script {
                            systemTests ignoreFailure: !isPRBuild(),
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                vegawallet: commitHash,
                                ethereumEventForwarder: params.ETHEREUM_EVENT_FORWARDER_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                protos: params.PROTOS_BRANCH
                        }
                    }
                }
                stage('LNL System Tests') {
                    steps {
                        script {
                            systemTestsLNL ignoreFailure: !isPRBuild(),
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                vegawallet: commitHash,
                                ethereumEventForwarder: params.ETHEREUM_EVENT_FORWARDER_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                protos: params.PROTOS_BRANCH
                        }
                    }
                }
                stage('Capsule System Tests') {
                    steps {
                        script {
                            systemTestsCapsule vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                vegawallet: commitHash,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                protos: params.PROTOS_BRANCH,
                                ignoreFailure: !isPRBuild()

                        }
                    }
                }
            }
        }
    }
    post {
        success {
            retry(3) {
                script {
                    slack.slackSendCISuccess name: 'Go Wallet CI', channel: '#tradingcore-notify'
                }
            }
        }
        unsuccessful {
            retry(3) {
                script {
                    slack.slackSendCIFailure name: 'Go Wallet CI', channel: '#tradingcore-notify'
                }
            }
        }
    }
}
