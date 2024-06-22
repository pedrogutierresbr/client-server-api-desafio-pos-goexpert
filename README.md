# Client Server API

## Sobre o projeto

Este é o repositório destinado ao desafio Client-Server-API do curso Pós Goexpert da faculdade FullCycle.

## Funcionalidades

-   O projeto possibilita ao usuário:

    -   Via client, solicitar a cotação atual do dólar ao server
    -   Consultar uma API em tempo real
    -   Salvar dados em um Banco de Dados

## Como executar o projeto

### Pré-requisitos

Antes de começar, você vai precisar ter instalado em sua máquina as seguintes ferramentas:

-   [Git](https://git-scm.com)
-   [VSCode](https://code.visualstudio.com/)

#### Acessando o repositório

```bash

# Clone este repositório
$ git clone https://github.com/pedrogutierresbr/client-server-api-desafio-pos-goexpert.git

```

#### Rodando a aplicação - Server

```bash

# Acesse a pasta do projeto no seu terminal/cmd
$ cd server

# Importe os pacotes
$ go mod tidy

# Inicie a aplicação
$ go run server.go

# A aplicação será aberta na porta:8080

```

#### Rodando a aplicação - Client

```bash

# Acesse a pasta do projeto no seu terminal/cmd
$ cd client

# Importe os pacotes
$ go mod tidy

# Inicie a aplicação
$ go run client.go

```

## Tecnologias

As seguintes ferramentas foram usadas na construção do projeto:

-   Go
-   GORM
-   SQLite

## Licença

Este projeto esta sobe a licença [MIT](./LICENSE).

Feito por Pedro Gutierres [Entre em contato!](https://www.linkedin.com/in/pedrogabrielgutierres/)
