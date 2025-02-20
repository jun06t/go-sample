# Evaluation Context Sample

## Get Started
### Fractional Evaluation
```
$ curl http://localhost:8080/hello?userId=2
{"result":true}
$ curl http://localhost:8080/hello?userId=5
{"result":false}
```

You can confirm fraction ratio.
```
$ k6 run k6/test.js

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

     execution: local
        script: k6/test.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 1m0s max duration (incl. graceful stop):
              * constant_request_rate: 67.00 iterations/s for 30s (maxVUs: 50-100, gracefulStop: 30s)


     ✓ status is 200

     checks.........................: 100.00% 2010 out of 2010
     data_received..................: 279 kB  9.3 kB/s
     data_sent......................: 203 kB  6.8 kB/s
     http_req_blocked...............: avg=18.43µs min=1µs     med=9µs    max=1.05ms  p(90)=16.1µs p(95)=24µs
     http_req_connecting............: avg=6.35µs  min=0s      med=0s     max=828µs   p(90)=0s     p(95)=0s
     http_req_duration..............: avg=3.83ms  min=837µs   med=3.36ms max=33.92ms p(90)=5.77ms p(95)=7.27ms
       { expected_response:true }...: avg=3.83ms  min=837µs   med=3.36ms max=33.92ms p(90)=5.77ms p(95)=7.27ms
     http_req_failed................: 0.00%   0 out of 2010
     http_req_receiving.............: avg=84.09µs min=10µs    med=70µs   max=12.04ms p(90)=119µs  p(95)=143.54µs
     http_req_sending...............: avg=29.9µs  min=2µs     med=28µs   max=469µs   p(90)=44µs   p(95)=52µs
     http_req_tls_handshaking.......: avg=0s      min=0s      med=0s     max=0s      p(90)=0s     p(95)=0s
     http_req_waiting...............: avg=3.72ms  min=809µs   med=3.25ms max=33.87ms p(90)=5.62ms p(95)=7.11ms
     http_reqs......................: 2010    66.998962/s
     iteration_duration.............: avg=4.18ms  min=947.7µs med=3.7ms  max=34.39ms p(90)=6.24ms p(95)=7.84ms
     iterations.....................: 2010    66.998962/s
     result_false_count.............: 1511    50.365886/s
     result_true_count..............: 499     16.633076/s
     vus............................: 0       min=0            max=1
     vus_max........................: 50      min=50           max=50
```

### Chainable if/else

```
$ curl http://localhost:8080/color?age=5
{"color":"#b91c1c"}
$ curl http://localhost:8080/color?age=90
{"color":"#4b5563"}
```