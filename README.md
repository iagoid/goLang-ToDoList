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

## Explicações
Todas as listas criadas são salvas no arquivo lists.txt.
Ao iniciar o programa ele inicialmente carrega todas as listas com a função LoadLists, ela percorre os resultados buscando o ultimo Id, para ao criar a nova lista ela não possuir um mesmo identificador

A parte da criação de uma Lista é realizada em 2 partes, a primeira (verifyFormCreate), valida os campos e verifica se alguma lista que se refere aquele estabelecimento já foi criada, se sim ele armazena aqueles dados em uma variavel e retorna a pergunta se o usuário deseja criar a lista, se a resposta for sim ele cria a Lista utilizando os dados já armazenados.

A página de edição busca a posição daquela Lista dentro das Listas e altera pelos valores que vieram pelo formulário.

A exclusão de uma Lista é feita simplesmente buscando a posição desta na Lists e então a remove de lá

A checagem da conclusão da lista é feita simplesmente verificando se aquela lista já foi concluida ou não, ao enviar a requisição o valor é simplesmente alterado.

Já as funções upList e downList servem para a pessoa alterar a posição da Lista no array das Lists

Após cada uma dessas fases é chamada a função Save, que salva as listas novamente no arquivo txt

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
