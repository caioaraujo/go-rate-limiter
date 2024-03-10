# Rate Limiter

Um rate limiter escrito em Go.

## Como funciona

O rate limiter está configurado em um middleware, que verifica se um determinado IP ou Token de acesso (ver METODO_BLOQUEIO em Configuração) está autorizado a acessar a página:

- Caso o número de requisições do IP ou Token atingir o número máximo de requisições por segundo (ver MAX_REQ_PERMITIDAS em Configuração), o IP ou Token é gravado no Cache com o tempo de expiração definido (ver TEMPO_BLOQUEIO_SEC em Configuração).
- Caso o método de bloqueio esteja configurado para Token, ele pegará o valor do token da chave `API_KEY` do header da requisição.

## Configuração

As configurações são definidas em variáveis em ambiente da seguinte forma:

- MAX_REQ_PERMITIDAS: número máximo de requisições permitidas por segundo;
- TEMPO_BLOQUEIO_SEC: tempo de bloqueio, em segundos, do IP ou do Token caso a quantidade de requisições tenha sido excedida;
- METODO_BLOQUEIO: método de limitação do rate limiter. Opções: `IP` e `TOKEN`;

Nota: As configurações podem ser feitas localmente no arquivo `.env` na raíz do projeto.

## Execução do web server
Na raíz do projeto, execute `go run cmd/system/main.go`