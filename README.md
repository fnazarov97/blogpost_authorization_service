# Blogpost with postgres
# blockpost_grpc my grpc project
# this project will be finish soon
# In order for you to implement this project, you need to write your own .env file

Migrate DB:
migrate -path ./storage/migrations -database 'postgres://admin:yourDB_password@localhost:5432/blogpost_auth?sslmode=disable' up

