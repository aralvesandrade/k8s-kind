import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 10, // número de usuários virtuais
    duration: '30s', // duração do teste
};

export default function () {
    const res = http.get('http://172.19.0.12:31523/hello');
    check(res, { 'status is 200': (r) => r.status === 200 });
    sleep(1);
}