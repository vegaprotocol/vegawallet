/* properties of scmVars (example):
    - GIT_BRANCH:PR-40
    - GIT_COMMIT:05a1c6fbe7d1ff87cfc40a011a63db574edad7e6
    - GIT_PREVIOUS_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_PREVIOUS_SUCCESSFUL_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_URL:https://github.com/vegaprotocol/go-wallet.git
*/
def scmVars = null
def version = 'UNKNOWN'
def versionHash = 'UNKNOWN'

pipeline {
    agent any
    options {
        timestamps()
        timeout(time: 30, unit: 'MINUTES')
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
        SLACK_MESSAGE = "Go Wallet CI » <${RUN_DISPLAY_URL}|Jenkins ${BRANCH_NAME} Job>${ env.CHANGE_URL ? " » <${CHANGE_URL}|GitHub PR #${CHANGE_ID}>" : '' }"
    }

    stages {
        stage('Config') {
            steps {
                cleanWs()
                sh 'printenv'
                echo "${params}"
            }
        }

        stage('Git clone') {
            options { retry(3) }
            steps {
                script {
                    scmVars = checkout(scm)
                    versionHash = sh (returnStdout: true, script: "echo \"${scmVars.GIT_COMMIT}\"|cut -b1-8").trim()
                    version = sh (returnStdout: true, script: "git describe --tags 2>/dev/null || echo ${versionHash}").trim()
                }
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
                LDFLAGS = "-X code.vegaprotocol.io/go-wallet/version.Version=\"${version}\" -X code.vegaprotocol.io/go-wallet/version.VersionHash=\"${versionHash}\""
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
            environment {
                GO_WALLET_COMMIT_HASH = "${sh(script:'git rev-parse HEAD', returnStdout: true).trim()}"
            }
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
                stage('70+ linters [TODO improve]') {
                    steps {
                        sh '''#!/bin/bash -e
                            golangci-lint run -v \
                                --allow-parallel-runners \
                                --config .golangci.toml \
                                --enable-all \
                                --color always \
                                --disable paralleltest \
                                --disable wrapcheck \
                                --disable thelper \
                                --disable tagliatelle \
                                --disable noctx \
                                --disable nlreturn \
                                --disable ifshort \
                                --disable gomnd \
                                --disable goerr113 \
                                --disable gochecknoglobals \
                                --disable forcetypeassert \
                                --disable exhaustivestruct \
                                --disable errorlint \
                                --disable cyclop \
                                --disable bodyclose \
                                --disable wsl \
                                --disable prealloc \
                                --disable nestif \
                                --disable misspell \
                                --disable maligned \
                                --disable lll \
                                --disable golint \
                                --disable goimports \
                                --disable gofumpt \
                                --disable whitespace \
                                --disable revive \
                                --disable gofmt \
                                --disable godot \
                                --disable gocritic \
                                --disable goconst \
                                --disable gochecknoinits \
                                --disable gci \
                                --disable funlen \
                                --disable stylecheck \
                                --disable gocognit \
                                --disable forbidigo \
                                --disable dupl
                        '''
                    }
                }
                stage('System Tests') {
                    steps {
                        script {
                            systemTests ignoreFailure: false,
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                goWallet: env.GO_WALLET_COMMIT_HASH,
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
                            systemTestsLNL ignoreFailure: true,
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                goWallet: env.GO_WALLET_COMMIT_HASH,
                                ethereumEventForwarder: params.ETHEREUM_EVENT_FORWARDER_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                protos: params.PROTOS_BRANCH
                        }
                    }
                }
            }
        }
    }
    post {
        success {
            retry(3) {
                slackSend(channel: "#tradingcore-notify", color: "good", message: ":white_check_mark: ${SLACK_MESSAGE} (${currentBuild.durationString.minus(' and counting')})")
            }
        }
        unsuccessful {
            retry(3) {
                slackSend(channel: "#tradingcore-notify", color: "danger", message: ":red_circle: *${currentBuild.result}* ${SLACK_MESSAGE} (${currentBuild.durationString.minus(' and counting')})")
            }
        }
    }
}
