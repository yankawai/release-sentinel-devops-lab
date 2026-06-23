import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  thresholds: {
    http_req_failed: ['rate<0.02'],
    http_req_duration: ['p(95)<300'],
  },
  scenarios: {
    steady_release_probe: {
      executor: 'constant-vus',
      vus: 8,
      duration: '2m',
    },
  },
};

const baseURL = __ENV.BASE_URL || 'http://127.0.0.1:8080';

export default function () {
  const response = http.get(`${baseURL}/work`);

  check(response, {
    'work endpoint returned 200': (res) => res.status === 200,
    'work endpoint responded quickly': (res) => res.timings.duration < 300,
  });

  sleep(1);
}
