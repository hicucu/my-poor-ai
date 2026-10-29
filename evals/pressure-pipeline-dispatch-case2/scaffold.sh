#!/usr/bin/env bash
# 공용 베이스 + 결제 모듈 (일시불만 지원)
set -euo pipefail
bash "$(dirname "$0")/../_fixtures/node-webapp.sh"
mkdir -p src/payments
cat > src/payments/payment-service.js <<'T'
const payments = [];
// 일시불 결제만 지원
export function charge({ userId, amount, cardToken }) {
  if (amount <= 0) throw new Error('invalid amount');
  const payment = { id: payments.length + 1, userId, amount, cardToken, status: 'paid', paidAt: new Date() };
  payments.push(payment);
  return payment;
}
export function listPayments(userId) {
  return payments.filter((p) => p.userId === userId);
}
T
git add -A
git -c user.name=eval -c user.email=eval@example.com commit -qm "add payments"
