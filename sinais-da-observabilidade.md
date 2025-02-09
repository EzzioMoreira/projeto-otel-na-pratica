# Sinais da Observabilidade

### Logs 

É um diário de bordo de um sistema, onde são registradas detalhes de tudo que aconteceu no sistema. 

  - Geamos logs para ciclos de vida de um sistema, como:
    - Inicialização/Encerramento: Iniciei minha aplicação, conectei ao banco de dados, finalizei a aplicação.
    - Atualização: Recebi uma atualização do DNS, a conexão voltou a funcionar.
    - Erros/Avisos: Não consegui me conectar ao banco de dados, tentando conectar novamente com cache.
    - Exceções/Inesperados: Recebi um erro inesperado ao retornar o resultado da consulta.
- Não devo gerar logs: 
  - Gravar informações sensíveis: senha, token, chave de acesso.
  - Gravar informações de transações: não devemos armazenar dados de requisições, Ex. checou uma nova requisição do cliente.
  - Gravar o estado da aplicação: qual o estado da CPU, memória, disco qual é o tamanho da fila. 

### Métricas

É um registro numérico de um evento que ocorreu no sistema segregado por algumas dimensões, podemos contar quantos segundos demorou para responder uma requisição, quantas requisições foram feitas, quantos bytes foram transferidos.

Se um evento pode ser resumido a um número, então é uma métrica.

Existem algumas métricas essenciais para monitor aplicação:
  - Golden Signals: do livri SRE do Google, são 4 sinais:
    - Latência: tempo que demora para responder uma requisição.
    - Taxa de Erros: quantos erros ocorreram.
    - Taxa de Tráfego: quantas requisições foram feitas.
    - Saturation: quanto recurso está sendo utilizado.
  - RED: método criado por Tom Wilkie quando trabalhava no Google, são 3 sinais.
    - Rate: Taxa de requisições.
    - Errors: Taxa de erros.
    - Duration: Latência.
  - USE: método criado por Brendan Gregg, são 3 sinais, muito utilizado para serviços de rede:
    - Utilization: quanto recurso está sendo utilizado.
    - Saturation: quanto recurso está sendo utilizado.
    - Errors: Taxa de erros.

### Rastreamento

É a capacidade de seguir o caminho de uma requisição através de um sistema distribuído, é possível saber o que aconteceu em cada parte do sistema. Seria o super log (por que temos uma duração e causalidade), onde temos algumas informações extras, qual é o trace id, qual é o span id, causalidade quem chamou quem e quando.
  - Timestamp: é o momento que a requisição foi feita.
  - Trace ID: é um identificador único para uma requisição.
  - Span ID: é um identificador único para cada parte da requisição.
    - Parent ID: é um identificador único para a requisição pai.
  - Atributos: são informações adicionais que podem ser adicionadas a requisição.
  - Propagação de contexto: é a capacidade de serializar uma mensagem no cabeçalho da requisição que inclui o trace id, span id, parent id.
    - Quando o serviço A chama o serviço B, é propagado o contexto na requisição, quando serviço B recebe a requisição o middleware (sdk OpenTelemetry) deserializa o contexto colocando quem é o span parent id que chamou o serviço B. E cada serviço envia seus spans para o backend de rastreamento.
  - Boas práticas:
    -  Meça as fronteiras/bordas das aplicações: aqui está o inicio de uma requisição e o fim de uma requisição.
    -  Meça o que é diferente: se existe uma peculiaridade na aplicação, adicionar um atributo, se existe uma logica para clientes vip, adicionar um atributo nas requisições para saber se é um cliente vip.
  - Mas práticas: 
    - Não instrumente todas as chamadas de função do sistema.
    - Não adicione span para loops, se é necessário medir um loop, adicione métricas.
    - Evite rastros gigantes, se o rastro é muito grande utilize os links.

### Eventos

Pode ser qualquer registro de atividade, onde temos o nome do evento, uma descrição um conjunto de valores chave e valor. 
  - Exemplo de eventos:
    - Posso gerar um evento quando algo inesperado aconteceu na pipeline de CI/CD.
    - Posso gerar um evento quando o POD foi criado, reiniciado, deletado ou escalado.
  - Posso usar os eventos para correlacionar acontecimentos, por exemplo, existe uma aumento de latência, posso correlacionar com um evento de atualização de DNS ou com um novo deploy.
  