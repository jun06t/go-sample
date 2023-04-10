import { check } from "k6";
import http from "k6/http";

export const options = {
  stages: [{ target: 100, duration: "30s" }],
};

export default function () {
  const params = {
    headers: {
      "Accept-Encoding": "gzip",
    },
  };
  const res = http.get("http://localhost:8080", params);

  check(res, {
    "status is 200": (r) => r.status === 200,
  });
}
