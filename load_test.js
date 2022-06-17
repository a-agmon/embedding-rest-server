import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 500,
  duration: '5s',
};
export default function () {
  http.get('http://127.0.0.1:8080/similar?to=Eminem&topk=10');
  sleep(1);
}


// this was run using k6 run load_test.js

//ab -n 10 -c 2 http://127.0.0.1:8080/similar?to=Eminem&topk=10
//ab -p test_payload.txt -T application/json  -c 10 -n 10 http://127.0.0.1:8080/mostsimilar



