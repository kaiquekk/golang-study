# Projeto de Estudos em Golang

Este é um projeto de estudo com o objetivo de servir como uma introdução à linguagem Go (Golang). A aplicação processa dados de cotações históricas de preços de títulos negociados na Bolsa em um ano específico, realizando operações de leitura, armazenamento e visualização gráfica.

## Funcionalidades

A aplicação realiza as seguintes etapas:

1. **Leitura e Processamento de Arquivo**:
   - Recebe um arquivo de texto contendo dados históricos de cotações.
   - Lê o conteúdo do arquivo e itera pelas linhas, mapeando as informações com base no posicionamento dos caracteres.

2. **Armazenamento em Banco de Dados**:
   - Conecta-se a uma base de dados PostgreSQL.
   - Armazena os dados mapeados mais relevantes em uma tabela, que será utilizada para a construção do gráfico.

3. **Geração de Gráfico**:
   - Lê os dados armazenados na tabela do banco de dados.
   - Gera um gráfico que exibe a variação do preço máximo dentre todas as ações negociadas em cada mês, com base em todos os dados processados e presentes na base.

## Pré-requisitos

Antes de executar a aplicação, certifique-se de que os seguintes requisitos estão atendidos:

- **Go**: A linguagem Go deve estar instalada no seu ambiente. Você pode baixá-la em [golang.org](https://golang.org/).
- **PostgreSQL**: Um banco de dados PostgreSQL deve estar configurado e acessível.
- **Variáveis Globais**: Configure as variáveis de ambiente com os dados de conexão ao banco de dados PostgreSQL. Por exemplo:
    ```go
    const (
        host     = <host>
        port     = <port>
        user     = <user>
        password = <password>
        dbname   = <database_name>
    )
    ```
## Como Executar
Siga os passos abaixo para executar a aplicação:

- Baixe as dependências do projeto:
    ```bash
    go get -d ./...
    ```
- Execute a aplicação, passando o caminho do arquivo de texto a ser processado como argumento:
    ```bash
    go run . <arquivo_de_texto>
    ```
## Observações
- Certifique-se de que o arquivo de texto fornecido segue o formato esperado pela aplicação.
- O gráfico gerado será exibido na porta 8089 no localhost.
## Objetivo
Este projeto foi desenvolvido com fins educacionais, para explorar conceitos fundamentais da linguagem Go, como:

- Manipulação de arquivos.
- Conexão com bancos de dados.
- Geração de gráficos.
- Estrutura básica da linguagem.


## Aprendizado
Conceitos e tópicos básicos da linguagem Go que foram introduzidos com esse projeto:

- Sintaxes básicas da linguagem;
- Gerenciamento de packages e dependências externas;
- Declaração e uso de constantes;
- Definição de rotas e exposição de rotas HTTP;
- Uso de logs e prints de mensagens;
- Conexão com base de dados e execução de queries;
- Integração com biblioteca de construção de gráficos e uso de suas funções;
- Uso de funções do file system para leitura de arquivos;
- Aplicação de padrões de concorrência e paralelismo.