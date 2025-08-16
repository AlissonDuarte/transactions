import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10,        // virtial users
  duration: '30s' // test duration
};

export default function () {
  const users = [
    { id: 1, from: 1, to: 2 }, // Alice -> Bob
    { id: 2, from: 2, to: 1 }  // Bob -> Alice
  ];

  // Select a random user
  let user = users[Math.floor(Math.random() * users.length)];

  let amount = Math.floor(Math.random() * 50) + 1;

  let payload = JSON.stringify({
    sender_id: user.from,
    sender_type: "user",
    receiver_id: user.to,
    receiver_type: "user",
    amount: amount
  });

  let params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  let res = http.post('http://localhost:8080/transactions', payload, params);

  check(res, {
    'status 200 ou 201': (r) => r.status === 200 || r.status === 201,
  });

  sleep(1);
}
