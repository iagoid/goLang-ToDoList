module todolist.com/todolist

go 1.16

require (
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gorilla/mux v1.8.0
	todolist.com/utils v1.0.1
)

replace todolist.com/utils => ../utils
