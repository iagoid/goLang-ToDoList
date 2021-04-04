# ToDoList com Golang

**Objetivo**: Projeto feito com o intuito de agilizar a escolha de produtos de usuários durante a pandemia, com objetivo é fazer com que estes fiquem o menor tempo possível dentro de locais onde possam contrair o vírus.

### Requisitos

**Golang**: É necessário ter a linguagem Go instalada em seu computador em uma versão compatível com a 1.16. Golang pode ser baixada no seguinte link:  https://golang.org/dl/

**Mux**: O mux é utilizado para facilitar a listagem de rotas válidas, para mais informações acesse https://github.com/gorilla/mux

**Opcional**: Caso você queira pode também baixar alguma ferramenta que permite enviar solicitações HTTP, como o PostMan ou Insomnia.

## Excução
Abra a pasta todolist no terminal e digite o seguinte código<br>
*go run todolist.go*<br>
Após isso abra seu navegador no link: 
*http://localhost:8080/*, está é a página onde estão listadas todas as suas listas de tarefas(Observe que já existem duas pré-definidas)

Acesse a página *http://localhost:8080/create*, lá existe um formulário, ao ser preenchido ele adiciona as listas.

Para acessar uma lista acesse *http://localhost:8080/view/algum_id_existente/*. aqui será retornada a sua lista

Para editar uma lista acesse *http://localhost:8080/edit/algum_id_existente/*. aqui será mostrado um formulário que permite você editar sua lista

Para deleter uma lista acesse *http://localhost:8080/delete/algum_id_existente/*, irá aparecer uma mensagem escrito accepted, caso você volte para a página inicial sua lista já não existirá

## Testes
Para executar os testes basta rodar o código *go test* no terminal enquanto o server estiver ativo, se tudo der certo ele reornará a mensagem PASS no terminal.

### Features

- [ ] Os dados irão ser consumidos de uma API
- [ ] Servir os dados necessário para a to-do list em formato REST
- [ ] Criação de uma tarefa
- [ ] Edição de uma tarefa
- [ ] Deleção de uma tarefa
- [ ] Visualização geral de todas as tarefas
- [ ] Validação de uma tarefa em branco
- [ ] Retorno de mensagem de erro
- [ ] Retorno de mensagem de sucesso
- [ ] Considerar que já exista uma tarefa igual e solicitar ao usuário se ele quer adicionar mesmo assim
- [ ] O retorno da API deve ser em JSON
