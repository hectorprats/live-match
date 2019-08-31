docker run --name livematch_rabbit -p 15672:15672 -p 5672:5672 -it livematch_rabbit
docker run --name livematch -p 9090:9090 -it livematch
docker build -t livematch .
docker build -t livematch_rabbit .


docker network connect --alias livematch_postgres livematch_net livematch_postgres
docker network connect --alias livematch livematch_net livematch
docker network connect --alias livematch_rabbit livematch_net livematch_rabbit
