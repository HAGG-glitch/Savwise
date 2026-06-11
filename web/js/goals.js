async function loadGoals() {
  const res = await api('/api/goals');
  renderGoals(res.data || []);
}
function renderGoals(rows) {
  const el = document.getElementById('goalsList');
  if (!el) return;
  if (!rows.length) { el.innerHTML = '<p class="text-sm text-slate-500">No savings goals yet.</p>'; return; }
  el.innerHTML = rows.map(g => `<div class="rounded-xl border p-4"><div class="flex justify-between gap-3"><h3 class="font-bold">${g.name}</h3><button class="text-red-700 font-bold text-sm" onclick="deleteGoal('${g.id}')">Delete</button></div><p class="text-sm text-slate-600 mt-2">${g.status} · ${g.estimatedCompletion}</p><div class="bg-slate-200 h-2 rounded-full mt-2"><div class="bg-emerald-600 h-2 rounded-full" style="width:${Math.min(100,g.progressPercent||0)}%"></div></div><div class="grid grid-cols-3 gap-2 text-xs mt-3"><div>Current<br><strong>${money(g.currentAmount)}</strong></div><div>Target<br><strong>${money(g.targetAmount)}</strong></div><div>Monthly<br><strong>${money(g.monthlyContribution)}</strong></div></div></div>`).join('');
}
async function deleteGoal(id) {
  if (!confirm('Delete this goal?')) return;
  await api(`/api/goals/${id}`, { method: 'DELETE' });
  showToast('Goal deleted');
  await refreshAll();
}
function bindGoals() {
  var gf = document.getElementById('goalForm');
  if (gf) {
    gf.addEventListener('submit', async e => {
      e.preventDefault();
      try {
        const body = {
          name: document.getElementById('goalName').value.trim(),
          targetAmount: Number(document.getElementById('goalTarget').value),
          currentAmount: Number(document.getElementById('goalCurrent').value),
          monthlyContribution: Number(document.getElementById('goalContribution').value),
        };
        await api('/api/goals', { method: 'POST', body: JSON.stringify(body) });
        e.target.reset();
        showToast('Goal saved');
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
}
