import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  http.get('https://localhost:8080/mostsimilar');
  sleep(1);
}

//ab -n 10 -c 2 http://127.0.0.1:8080/similar?to=Eminem&topk=10

//ab -p test_payload.txt -T application/json  -c 10 -n 10 http://127.0.0.1:8080/mostsimilar



