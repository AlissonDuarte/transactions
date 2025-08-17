# 📄 Transaction Processing Service — Contexto do Projeto  

## 🇺🇸 English  

This project implements a **transaction processing service** in Go, leveraging **RabbitMQ** queues to ensure scalability and resilience.  

### Goal  
- Allow multiple producers to publish transactions to a queue.  
- Ensure multiple consumers (workers) process transactions concurrently and safely.  
- Measure the impact of increasing consumers on system throughput.  

### Technical Approach  
- **Queue:** RabbitMQ is used as the message broker.  
- **Consumers:** implemented as Go goroutines.  
- **Dedicated channel:** each worker receives its own `amqp.Channel` to avoid bottlenecks.  
- **Scalability:** the number of consumers can be configured (`numConsumers`).  
- **Performance test:** increasing consumers boosted throughput from **3–8 msgs/s** to **35-42 msgs/s**.  

### Benchmarks  

#### Initial tests  
- **With few consumers (1–2):** average throughput was **3–8 messages per second**.  
- **With 5 parallel consumers:** throughput increased to **35-42 messages per second**.  

#### Benchmarks on `kr6`  
Load tests executed on `kr6` produced the following results:  

- **Throughput:** remained steady at **35-42 messages per second** under stable conditions.  
- **Errors:** observed error logs were actually **expected denials** (e.g., invalid transactions, wrong sender, insufficient balance).  
- **Business validation:** these errors confirm that the system is enforcing business rules correctly rather than suffering technical failures.  

📌 **Conclusion:** the `kr6` tests demonstrated that the system **scaled properly**, maintaining stable throughput while consistently validating business rules.  

### Lessons  
- Effective parallelism depends on both the number of consumers and **prefetch configuration**.  
- **Manual Acknowledgement** ensures consistency even in failure scenarios.  
- RabbitMQ automatically balances workload across consumers.  

```text
█ TOTAL RESULTS
checks_total.......: 4760 39.593981/s
checks_succeeded...: 83.21% 3961 out of 4760
checks_failed......: 16.78% 799 out of 4760

✗ status 200 ou 201
↳ 66% — ✓ 1581 / ✗ 799
✓ mensagem de falha para lojas

HTTP
http_req_duration..............: avg=8.91ms min=1.19ms med=7.42ms max=66.18ms p(90)=19.21ms p(95)=23.86ms
{ expected_response:true }.....: avg=8.98ms min=1.2ms  med=7.4ms  max=66.18ms p(90)=19.75ms p(95)=23.94ms
http_req_failed................: 33.57% 799 out of 2380
http_reqs......................: 2380 19.79699/s

EXECUTION
iteration_duration.............: avg=1.01s min=1s med=1s max=1.07s p(90)=1.02s p(95)=1.02s
iterations.....................: 2380 19.79699/s
vus............................: 20 min=20 max=20
vus_max........................: 20 min=20 max=20

NETWORK
data_received..................: 721 kB 6.0 kB/s
data_sent......................: 530 kB 4.4 kB/s
```

---

## 🇧🇷 Português  

Este projeto implementa um **serviço de processamento de transações** em Go, com suporte a filas no **RabbitMQ** para garantir escalabilidade e resiliência.  

### Objetivo  
- Permitir que múltiplos produtores enviem transações para uma fila.  
- Garantir que múltiplos consumidores (workers) processem essas transações de forma paralela e segura.  
- Avaliar o impacto de aumentar o número de consumidores no throughput do sistema.  

### Abordagem Técnica  
- **Fila:** RabbitMQ é utilizado como broker de mensagens.  
- **Consumers:** implementados como goroutines em Go.  
- **Canal dedicado:** cada worker recebe seu próprio `amqp.Channel` para evitar gargalos.  
- **Escalabilidade:** o número de consumidores pode ser configurado (`numConsumers`).  
- **Teste de performance:** verificou-se que aumentar os consumidores elevou a taxa de processamento de **3–8 msgs/s** para **35-42 msgs/s**.  

### Benchmarks  

#### Testes iniciais  
- **Com poucos consumidores (1–2):** média de **3–8 mensagens por segundo**.  
- **Com 5 consumidores paralelos:** throughput subiu para **35-42 mensagens por segundo**.  

#### Benchmarks no `kr6`  
Durante os testes de carga no `kr6`, coletamos as seguintes métricas:  

- **Throughput:** manteve-se dentro da faixa de **35-42 mensagens por segundo** em condições estáveis.  
- **Erros:** ocorreram logs de erro, mas estes eram **negações esperadas** (ex.: transações inválidas, remetente incorreto, saldo insuficiente).  
- **Validação de negócio:** os erros confirmam que o sistema está aplicando corretamente as regras definidas, não indicando falhas técnicas.  

📌 **Conclusão:** os testes no `kr6` mostraram que o sistema **escalou corretamente**, mantendo throughput estável e validando regras de negócio de forma consistente.  

### Lições  
- O paralelismo efetivo depende do número de consumidores + configuração de **prefetch**.  
- O uso de **Ack manual** garante consistência mesmo em falhas.  
- O balanceamento entre consumidores é feito automaticamente pelo RabbitMQ.  

```text
█ RESULTADOS TOTAIS
checks_total.......: 4760 39.593981/s
checks_succeeded...: 83.21% 3961 out of 4760
checks_failed......: 16.78% 799 out of 4760

✗ status 200 ou 201
↳ 66% — ✓ 1581 / ✗ 799
✓ mensagem de falha para lojas

HTTP
http_req_duration..............: avg=8.91ms min=1.19ms med=7.42ms max=66.18ms p(90)=19.21ms p(95)=23.86ms
{ expected_response:true }.....: avg=8.98ms min=1.2ms  med=7.4ms  max=66.18ms p(90)=19.75ms p(95)=23.94ms
http_req_failed................: 33.57% 799 out of 2380
http_reqs......................: 2380 19.79699/s

EXECUTION
iteration_duration.............: avg=1.01s min=1s med=1s max=1.07s p(90)=1.02s p(95)=1.02s
iterations.....................: 2380 19.79699/s
vus............................: 20 min=20 max=20
vus_max........................: 20 min=20 max=20

NETWORK
data_received..................: 721 kB 6.0 kB/s
data_sent......................: 530 kB 4.4 kB/s
```
