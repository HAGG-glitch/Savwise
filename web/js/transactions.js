async function loadTransactions() {
  const res = await api('/api/transactions');
  renderTransactions(res.data || []);
}
function renderTransactions(rows) {
  const el = document.getElementById('transactionsList');
  if (!el) return;
  if (!rows.length) { el.innerHTML = '<p class="text-sm text-slate-500">No transactions yet.</p>'; return; }
  const tableRows = rows.map(t => `<tr class="border-b"><td class="py-2">${t.date}</td><td class="overflow-wrap-anywhere">${t.description}</td><td>${t.type}</td><td>${t.category}</td><td class="font-bold ${t.type === 'income' ? 'text-emerald-700' : 'text-red-700'}">${money(t.amount)}</td><td><button class="text-red-700 font-bold" onclick="deleteTransaction('${t.id}')">Delete</button></td></tr>`).join('');
  const cards = rows.map(t => `<div class="tx-card"><div class="flex items-center gap-3 min-w-0"><div><div class="font-bold text-sm overflow-wrap-anywhere">${t.description}</div><div class="text-xs text-slate-500">${t.date} &middot; ${t.category}</div></div></div><div class="flex items-center gap-2 flex-shrink-0"><span class="font-bold text-sm ${t.type === 'income' ? 'text-emerald-700' : 'text-red-700'}">${money(t.amount)}</span><button class="text-red-700 text-xs font-bold" onclick="deleteTransaction('${t.id}')">&#10005;</button></div></div>`).join('');
  el.innerHTML = `<div class="tx-table-desktop"><table class="min-w-full text-sm"><thead><tr class="text-left border-b"><th class="py-2">Date</th><th>Description</th><th>Type</th><th>Category</th><th>Amount</th><th></th></tr></thead><tbody>${tableRows}</tbody></table></div><div class="tx-cards-mobile">${cards}</div>`;
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
