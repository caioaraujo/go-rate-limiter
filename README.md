# Rate Limiter

Um rate limiter escrito em Go.

## Configuração

As configurações são definidas em variáveis em ambiente da seguinte forma:

- MAX_REQ_PERMITIDAS: número máximo de requisições permitidas por segundo;
- TEMPO_BLOQUEIO_SEC: tempo de bloqueio, em segundos, do IP ou do Token caso a quantidade de requisições tenha sido excedida;
- METODO_BLOQUEIO: método de limitação do rate limiter. Opções: IP e TOKEN;

Nota: As configurações podem ser feitas localmente no arquivo `.env` na raíz do projeto.

## Execução do web server
Na raíz do projeto, execute `go run cmd/ratelimiter/main.go`