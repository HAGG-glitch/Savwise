async function loadDashboard() {
  const res = await api('/api/dashboard');
  const d = res.data;
  renderGreeting(d);
  renderMainSummary(d);
  renderDashboardCards(d);
  renderBreakdown(d.spendingBreakdown || []);
  renderAlerts(d.alerts || []);
  renderGoalOverview(d.goals || []);
  return d;
}
function renderGreeting(d) {
  var el = document.getElementById('greeting');
  if (!el) return;
  var hour = new Date().getHours();
  var timeGreeting = hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';
  el.textContent = timeGreeting + (d.user.fullName ? ', ' + d.user.fullName : '') + '.';
  var dateEl = document.getElementById('currentDateDisplay');
  if (dateEl) dateEl.textContent = new Date().toLocaleDateString('en-GB', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
}
function renderMainSummary(d) {
  var el = document.getElementById('mainSummaryCard');
  if (!el) return;
  var nextAction = (d.alerts || []).length > 0 ? d.alerts[0].recommendedAction : 'Set up your profile and add transactions.';
  el.innerHTML = '<div class="grid md:grid-cols-4 gap-4"><div><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">What came in</p><h3 class="text-2xl font-black mt-1 text-emerald-700">' + money(d.totalIncome) + '</h3></div><div><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">What went out</p><h3 class="text-2xl font-black mt-1 text-red-600">' + money(d.totalExpenses) + '</h3></div><div><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">What is left</p><h3 class="text-2xl font-black mt-1 ' + (d.monthlySurplus >= 0 ? 'text-emerald-700' : 'text-red-600') + '">' + money(d.monthlySurplus) + '</h3></div><div><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">Next step</p><p class="text-sm font-bold mt-1">' + nextAction + '</p></div></div>';
}
function renderDashboardCards(d) {
  var topAlert = (d.alerts || []).length > 0 ? d.alerts[0].title : 'None';
  var dc = document.getElementById('dashboardCards');
  if (!dc) return;
  dc.innerHTML = [
    '<div class="glass-dashboard-card"><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">Financial Health</p><h3 class="text-2xl font-black mt-2 text-emerald-700">' + d.score.total + '/100</h3><div class="flex gap-2 mt-2 text-xs text-slate-500"><span>Save ' + d.score.savingsHabit + '/30</span><span>Budget ' + d.score.budgetControl + '/30</span><span>Emergency ' + d.score.emergencyFund + '/20</span><span>Goals ' + d.score.goalProgress + '/20</span></div></div>',
    '<div class="glass-dashboard-card"><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">Emergency Safety</p><h3 class="text-2xl font-black mt-2 text-emerald-700">' + Number(d.emergencyCoverageDays || 0).toFixed(0) + ' days</h3><p class="text-xs text-slate-500 mt-2">Savings can cover this many days of expenses</p></div>',
    '<div class="card p-5"><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">Savings rate</p><h3 class="text-2xl font-black mt-2 text-emerald-700">' + pct(d.savingsRate) + '</h3><p class="text-xs text-slate-500 mt-2">Surplus as percent of income</p></div>',
    '<div class="card p-5"><p class="text-xs uppercase tracking-wide text-slate-500 font-bold">Top alert</p><h3 class="text-sm font-black mt-2 ' + (topAlert !== 'None' && topAlert !== 'No major risk detected' ? 'text-red-600' : 'text-emerald-700') + '">' + topAlert + '</h3></div>',
  ].join('');
}
function renderBreakdown(rows) {
  const el = document.getElementById('spendingBreakdown');
  if (!el) return;
  if (!rows.length) { el.innerHTML = '<p class="text-sm text-slate-500">No expense categories yet.</p>'; return; }
  el.innerHTML = rows.map(r => '<div><div class="flex justify-between text-sm"><span>' + r.category + '</span><strong>' + money(r.amount) + '</strong></div><div class="bg-slate-200 h-2 rounded-full mt-1"><div class="bg-emerald-600 h-2 rounded-full" style="width:' + Math.min(100, r.percent) + '%"></div></div></div>').join('');
}
function renderAlerts(rows) {
  const el = document.getElementById('alerts');
  if (!el) return;
  if (!rows.length) { el.innerHTML = '<p class="text-sm text-slate-500">No alerts.</p>'; return; }
  el.innerHTML = rows.map(a => '<div class="rounded-xl p-3 border ' + (a.severity === 'High' ? 'bg-red-50 border-red-200' : a.severity === 'Medium' ? 'bg-amber-50 border-amber-200' : 'bg-emerald-50 border-emerald-200') + '"><div class="flex justify-between gap-2"><strong>' + a.title + '</strong><span class="badge badge-' + a.severity.toLowerCase() + '">' + a.severity + '</span></div><p class="text-sm text-slate-600 mt-1">' + a.explanation + '</p><p class="text-xs text-slate-500 mt-2"><strong>Action:</strong> ' + a.recommendedAction + '</p></div>').join('');
}
function renderGoalOverview(goals) {
  const el = document.getElementById('goalOverview');
  if (!el) return;
  if (!goals.length) { el.innerHTML = '<p class="text-sm text-slate-500">No goals created yet.</p>'; return; }
  el.innerHTML = goals.map(function(g) {
    return '<div class="rounded-xl border p-4"><div class="flex justify-between"><strong>' + g.name + '</strong><span class="text-emerald-700 font-bold">' + Number(g.progressPercent || 0).toFixed(0) + '%</span></div><div class="bg-slate-200 h-2 rounded-full mt-2"><div class="bg-emerald-600 h-2 rounded-full" style="width:' + Math.min(100, g.progressPercent || 0) + '%"></div></div><div class="grid grid-cols-3 gap-2 text-xs mt-3"><div><span class="text-slate-500">Current</span><br><strong>' + money(g.currentAmount) + '</strong></div><div><span class="text-slate-500">Target</span><br><strong>' + money(g.targetAmount) + '</strong></div><div><span class="text-slate-500">Estimate</span><br><strong>' + (g.estimatedCompletion || 'No estimate') + '</strong></div></div></div>';
  }).join('');
}
