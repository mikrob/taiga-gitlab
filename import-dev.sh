#/bin/bash

export GITLAB_TOKEN=1qVsgb99XFst2GRwBXxn
export TAIGA_PASSWORD=123123

for i in $(cat repo.list); do
  go run run.go import --taiga-url http://dev:81 --taiga-user admin --taiga-project project --gitlab-url https://gitlab.project.com --gitlab-project $i --taiga-skip-user --create-task
  [[ $? -gt 0 ]] && exit 1
done
