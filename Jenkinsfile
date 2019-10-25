pipeline {
    agent { label 'upbound-gce' }

    parameters {
        booleanParam(defaultValue: true, description: 'Execute pipeline?', name: 'shouldBuild')
    }

    options {
        disableConcurrentBuilds()
        timestamps()
    }

    environment {
        RUNNING_IN_CI = 'true'
        REPOSITORY_NAME = "${env.GIT_URL.tokenize('/')[3].split('\\.')[0]}"
        REPOSITORY_OWNER = "${env.GIT_URL.tokenize('/')[2]}"
        GITHUB_UPBOUND_BOT = credentials('github-upbound-jenkins')
    }

    stages {

        stage('Prepare') {
            steps {
                script {
                    if (env.CHANGE_ID != null) {
                        def json = sh (script: "curl -s https://api.github.com/repos/crossplaneio/crossplane/pulls/${env.CHANGE_ID}", returnStdout: true).trim()
                        def body = evaluateJson(json,'${json.body}')
                        if (body.contains("[skip ci]")) {
                            echo ("'[skip ci]' spotted in PR body text.")
                            env.shouldBuild = "false"
                        }
                    }
                }
                sh 'git config --global user.name "upbound-bot"'
                sh 'echo "machine github.com login upbound-bot password $GITHUB_UPBOUND_BOT" > ~/.netrc'
            }
        }

        stage('Build'){
            when {
                expression {
                    return env.shouldBuild != "false"
                }
            }
            steps {
                sh 'make build'
            }
            post {
                always {
                    archiveArtifacts "_output/lint/**/*"
                }
            }
        }

        stage('Unit Tests') {
            when {
                expression {
                    return env.shouldBuild != "false"
                }
            }
            steps {
                sh 'make test'
            }
            post {
                always {
                    archiveArtifacts "_output/tests/**/*"
                    junit "_output/tests/**/unit-tests.xml"
                    cobertura coberturaReportFile: '_output/tests/**/cobertura-coverage.xml',
                            classCoverageTargets: '50, 0, 0',
                            conditionalCoverageTargets: '70, 0, 0',
                            lineCoverageTargets: '40, 0, 0',
                            methodCoverageTargets: '30, 0, 0',
                            packageCoverageTargets: '80, 0, 0',
                            autoUpdateHealth: false,
                            autoUpdateStability: false,
                            enableNewApi: false,
                            failUnhealthy: false,
                            failUnstable: false,
                            maxNumberOfBuilds: 0,
                            onlyStable: false,
                            sourceEncoding: 'ASCII',
                            zoomCoverageChart: false
                }
            }
        }

    }

    post {
        always {
            script {
                sh 'make -j\$(nproc) clean'
            }
        }
    }
}

@NonCPS
def evaluateJson(String json, String gpath){
    //parse json
    def ojson = new groovy.json.JsonSlurper().parseText(json)
    //evaluate gpath as a gstring template where $json is a parsed json parameter
    return new groovy.text.GStringTemplateEngine().createTemplate(gpath).make(json:ojson).toString()
}
