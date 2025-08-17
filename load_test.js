import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 20,
  duration: '120s',
};

export default function () {
  const participants = [
    { id: 1, type: 'user', to: 2, to_type: 'user' },   // Alice -> Bob
    { id: 2, type: 'user', to: 1, to_type: 'user' },   // Bob -> Alice
    { id: 1, type: 'user', to: 3, to_type: 'store' },  // Alice -> Store Alpha
    { id: 2, type: 'user', to: 4, to_type: 'store' },  // Bob -> Store Beta
    { id: 3, type: 'store', to: 1, to_type: 'user' },  // Store Alpha -> Alice (should fail)
    { id: 4, type: 'store', to: 2, to_type: 'user' },  // Store Beta -> Bob (should fail)
  ];

  let sender = participants[Math.floor(Math.random() * participants.length)];
  let amount = Math.floor(Math.random() * 50) + 1;

  let payload = JSON.stringify({
    sender_id: sender.id,
    sender_type: sender.type,
    receiver_id: sender.to,
    receiver_type: sender.to_type,
    amount: amount,
  });

  let params = {
    headers: { 'Content-Type': 'application/json' },
  };

  let res = http.post('http://localhost:8080/transactions', payload, params);

  // Checks
  check(res, {
    'transação de usuário bem-sucedida': (r) => {
      return sender.type === 'user' && (r.status === 200 || r.status === 201);
    },
    'transação de loja bloqueada corretamente': (r) => {
      return sender.type === 'store' && r.status === 400 && r.json().message.includes('Invalid sender type');
    },
  });
  sleep(Math.random() * 2)
}
