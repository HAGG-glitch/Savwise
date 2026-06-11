async function loadTransactions() {
  const res = await api('/api/transactions');
  renderTransactions(res.data || []);
}
function renderTransactions(rows) {
  const el = document.getElementById('transactionsList');
  if (!el) return;
  if (!rows.length) { el.innerHTML = '<p class="text-sm text-slate-500">No transactions yet.</p>'; return; }
  el.innerHTML = `<table class="min-w-full text-sm"><thead><tr class="text-left border-b"><th class="py-2">Date</th><th>Description</th><th>Type</th><th>Category</th><th>Amount</th><th></th></tr></thead><tbody>${rows.map(t => `<tr class="border-b"><td class="py-2">${t.date}</td><td>${t.description}</td><td>${t.type}</td><td>${t.category}</td><td class="font-bold ${t.type === 'income' ? 'text-emerald-700' : 'text-red-700'}">${money(t.amount)}</td><td><button class="text-red-700 font-bold" onclick="deleteTransaction('${t.id}')">Delete</button></td></tr>`).join('')}</tbody></table>`;
}
async function deleteTransaction(id) {
  if (!confirm('Delete this transaction?')) return;
  await api(`/api/transactions/${id}`, { method: 'DELETE' });
  showToast('Transaction deleted');
  await refreshAll();
}
function bindTransactions() {
  var txDate = document.getElementById('txDate');
  if (txDate) txDate.valueAsDate = new Date();
  var tf = document.getElementById('transactionForm');
  if (tf) {
    tf.addEventListener('submit', async e => {
      e.preventDefault();
      try {
        const body = {
          description: document.getElementById('txDescription').value.trim(),
          amount: Number(document.getElementById('txAmount').value),
          type: document.getElementById('txType').value,
          category: document.getElementById('txCategory').value,
          date: document.getElementById('txDate').value,
        };
        await api('/api/transactions', { method: 'POST', body: JSON.stringify(body) });
        e.target.reset();
        var td2 = document.getElementById('txDate');
        if (td2) td2.valueAsDate = new Date();
        showToast('Transaction saved');
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
}
