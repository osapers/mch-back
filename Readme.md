### mcg-back

Contains code of backend application;

Deps: Postgresql, Yake.

### For local test: `docker-compose up`

### Extra folder contains tag cloud
#### Way to upload it to db:
`docker run -it --rm -v $(pwd)/extra/hh_skills.csv:/tmp/hh_skills.csv bitnami/postgresql:11.12.0-debian-10-r23 psql -h ${HOSTNAME} -U postgres -d mch`
#### then
`\COPY mch.tag(id, name) FROM '/tmp/hh_skills.csv' DELIMITER '|' CSV HEADER;`