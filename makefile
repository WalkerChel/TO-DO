pstg_docker_up:
	docker run --name=todo-db -e POSTGRES_PASSWORD=qwerty -p 5432:5432 -d postgres

pstg_docker_down:
	docker stop todo-db && docker rm todo-db

migrations_up:
	migrate -path ~/GOLANG/src/TO-DO//migrations -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up

migrations_down:
	migrate -path ~/GOLANG/src/TO-DO//migrations -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down 

migrations_force:
	migrate -path ~/GOLANG/src/TO-DO//migrations -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' force 000001
