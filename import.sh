#/bin/bash

export GITLAB_TOKEN=1qVsgb99XFst2GRwBXxn
export TAIGA_PASSWORD=botsunit8075

for i in $(cat repo.list); do
  go run run.go import --taiga-url https://taiga.botsunit.io --taiga-user admin --taiga-project Ufancyme --gitlab-url https://gitlab.botsunit.com --gitlab-project $i
  [[ $? -gt 0 ]] && exit 1
done
