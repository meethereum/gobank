running the bank : make run
running the postgres database on port 5432 : docker run --name some-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres 
(you can replace 5432 with any allowed port number)
checking if you have connected to it : telnet localhost 5432
