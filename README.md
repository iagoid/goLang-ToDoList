# ToDoList com Golang

**Objetivo**: Projeto feito com o intuito de agilizar a escolha de produtos de usuários durante a pandemia, com objetivo é fazer com que estes fiquem o menor tempo possível dentro de locais onde possam contrair o vírus.

### Requisitos

**Golang**: É necessário ter a linguagem Go instalada em seu computador em uma versão compatível com a 1.16. Golang pode ser baixada no seguinte link:  https://golang.org/dl/

## Excução
Abra a pasta todolist no terminal e digite o seguinte código<br>
*go build todolist.go*<br>
*./todolist*<br>
Após isso abra seu navegador no link: 
http://localhost:8080/
Acesse a página *http://localhost:8080/create*, lá existe um formulário, ao ser preenchido ele adiciona as listas o que já foi adicionado, basta entrar em qualquer outra página do localhost para visualizar o retorno em json


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
