def GIT_URL = "git@github.com:AnastasiyaGapochkina01/dos-29-cicd_02.git"

pipeline {
  angent any

  parameters {
    booleanParam(name: 'RUN_DEPLOY', defaultValue: false, description: 'Run deploy job')
    gitParameter(type: 'PT_BRANCH', name: 'BRANCH', branchFilter: 'origin/(.*)', description: 'Choose a branch to checkout', sortMode: 'DESCENDING_SMART')
  }

  environment {
    REPO = "anestesia01/dos-29"
    PRJ_NAME = "ecom"
  }

  stages {
    stage('Checkout Repo'){
      steps {
        script {
          git branch: "${params.BRANCH}, url: "${GIT_URL}"
        }
      }
    }

    stage('Build Image'){
      steps {
        script {
          sh """
            docker build -t ${env.REPO}:${env.PRJ_NAME}-${BUILD_NUMBER}" .
          """
        }
    } 
    }

    stage('Push Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-token', usernameVariable: 'username', passwordVariable: 'password')]) {
          script {
            sh """
              docker login -u ${username} -p ${password}
            """
          }
        }
      }
    }

    stage('Call Deploy Job') {
      when {
        expression { return params.RUN_DEPLOY }
      }
      steps {
        script {
          build quietPeriod: 5, wait: false, job: 'deploy', parameters: [string(name: 'DOCKER_IMAGE', value: "${env.REPO}:${env.PRJ_NAME}-${BUILD_NUMBER}")]
        }
      }
    }
  }
}
