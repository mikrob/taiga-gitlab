#/bin/bash

export GITLAB_TOKEN=GITLABTOKEN
export TAIGA_PASSWORD=TAIGAPASSWORD

for i in $(cat repo.list); do
  go run run.go import --taiga-url https://taiga.project.io --taiga-user admin --taiga-project project --gitlab-url https://gitlab.project.com --gitlab-project $i --create-task
  [[ $? -gt 0 ]] && exit 1
done
