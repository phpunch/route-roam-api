docker run \
    -e PGADMIN_DEFAULT_EMAIL=pgadmin4@pgadmin.org \
    -e PGADMIN_DEFAULT_PASSWORD=admin \
    -p 5050:80 \
    --name pgadmin \
    --network route-roam-network \
    dpage/pgadmin4