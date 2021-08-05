/* properties of scmVars (example):
    - GIT_BRANCH:PR-40-head
    - GIT_COMMIT:05a1c6fbe7d1ff87cfc40a011a63db574edad7e6
    - GIT_PREVIOUS_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_PREVIOUS_SUCCESSFUL_COMMIT:5d02b46fdb653f789e799ff6ad304baccc32cbf9
    - GIT_URL:https://github.com/vegaprotocol/vega.git
*/
def scmVars = null
def version = 'UNKNOWN'
def versionHash = 'UNKNOWN'

pipeline {
    agent { label 'general' }
    environment {
        CGO_ENABLED = 1
        GO111MODULE = 'on'
        SLACK_MESSAGE = "VegaWallet CI » <${RUN_DISPLAY_URL}|Jenkins ${BRANCH_NAME} Job>${ env.CHANGE_URL ? " » <${CHANGE_URL}|GitHub PR #${CHANGE_ID}>" : '' }"
    }

    stages {
        stage('Git clone') {
            options { retry(3) }
            steps {
                sh 'printenv'
                echo "${params}"
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
        stage('Linux build') {
            environment {
                GOOS    = 'linux'
                GOARCH  = 'amd64'
                LDFLAGS = "-X code.vegaprotocol.io/go-wallet/version.Version=\"${version}\" -X code.vegaprotocol.io/go-wallet/version.VersionHash=\"${versionHash}\""
                OUTPUT  = './build/vegawallet-linux-amd64'
            }
            options { retry(3) }
            steps {
                sh label: 'Compile', script: '''
                    go build -o "${OUTPUT}" -ldflags "${LDFLAGS}"
                '''
                sh label: 'Sanity check', script: '''
                    file ${OUTPUT}
                    ${OUTPUT} version
                    ${OUTPUT} help
                '''
            }
        }
        stage('Test') {
            options { retry(3) }
            steps {
                sh label: 'test', script: '''
                    go test -v ./...
                '''
            }
        }
    }
    post {
        success {
            retry(3) {
                slackSend(channel: "#tradingcore-notify", color: "good", message: ":white_check_mark: ${SLACK_MESSAGE} (${currentBuild.durationString.minus(' and counting')})")
            }
        }
        failure {
            retry(3) {
                slackSend(channel: "#tradingcore-notify", color: "danger", message: ":red_circle: ${SLACK_MESSAGE} (${currentBuild.durationString.minus(' and counting')})")
            }
        }
    }
}
