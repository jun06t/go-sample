import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

// カウンターの定義
export let trueCount = new Counter('result_true_count');
export let falseCount = new Counter('result_false_count');

// シナリオ設定: 30秒間で1秒あたり約67リクエスト（合計約2010リクエスト）
export let options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 67,              // 1秒あたりのリクエスト数
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 50,   // 事前に確保する仮想ユーザ数
      maxVUs: 100,           // 必要に応じて増加させる上限
    },
  },
};

// ランダムなユーザIDを生成する関数（8文字の英数字）
function randomUserId() {
  return Math.random().toString(36).substring(2, 10);
}

export default function () {
  // ランダムなユーザIDを生成し、URLに設定
  let userId = randomUserId();
  let url = `http://localhost:8080/hello?userId=${userId}`;

  // GETリクエストの送信
  let res = http.get(url);

  // HTTPステータスチェック
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  // JSONレスポンスのパースとカウンタ更新
  let data;
  try {
    data = JSON.parse(res.body);
  } catch (e) {
    console.error('レスポンスがJSONではありません: ', res.body);
    return;
  }

  if (data.result === true) {
    trueCount.add(1);
  } else if (data.result === false) {
    falseCount.add(1);
  }
}
