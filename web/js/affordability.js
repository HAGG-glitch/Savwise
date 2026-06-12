function escapeHtml(text) {
  var div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}

function formatDate(dateStr) {
  if (!dateStr) return '';
  var parts = dateStr.split('-');
  if (parts.length !== 3) return dateStr;
  var months = ['January','February','March','April','May','June','July','August','September','October','November','December'];
  return parseInt(parts[2]) + ' ' + months[parseInt(parts[1])-1] + ' ' + parts[0];
}

function riskColor(level) {
  if (level === 'High') return 'danger';
  if (level === 'Medium') return 'warning';
  return 'success';
}

function riskIcon(level) {
  if (level === 'High') return '&#9888;';
  if (level === 'Medium') return '&#9888;';
  return '&#10003;';
}

function renderAffordabilityResult(r) {
  var color = riskColor(r.riskLevel);
  var icon = riskIcon(r.riskLevel);
  var targetDisplay = r.targetDate ? formatDate(r.targetDate) : 'No target date set';

  var sections = [];

  sections.push('<div class="flex items-center gap-2 mb-4"><span class="badge badge-' + color + ' text-sm px-3 py-1">' + icon + ' ' + escapeHtml(r.riskLevel) + ' Risk</span></div>');

  sections.push('<div class="mb-4"><div class="grid grid-cols-2 gap-3 text-sm"><div><span class="text-slate-500 text-xs">Item</span><p class="font-bold">' + escapeHtml(r.itemName) + '</p></div><div><span class="text-slate-500 text-xs">Price</span><p class="font-bold">' + money(r.itemPrice) + '</p></div><div><span class="text-slate-500 text-xs">Target date</span><p class="font-bold">' + targetDisplay + '</p></div><div><span class="text-slate-500 text-xs">Calculated</span><p class="font-bold text-xs">' + escapeHtml(r.calculatedAt) + '</p></div></div></div>');

  sections.push('<div class="border-t pt-3 mb-4"><h4 class="font-bold text-sm mb-2">Your numbers</h4><div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm"><div class="bg-slate-50 rounded-xl p-3"><span class="text-xs text-slate-500">Income</span><p class="font-bold text-emerald-700">' + money(r.monthlyIncome) + '</p></div><div class="bg-slate-50 rounded-xl p-3"><span class="text-xs text-slate-500">Expenses</span><p class="font-bold text-red-600">' + money(r.monthlyExpenses) + '</p></div><div class="bg-slate-50 rounded-xl p-3"><span class="text-xs text-slate-500">Surplus</span><p class="font-bold ' + (r.monthlySurplus >= 0 ? 'text-emerald-700' : 'text-red-600') + '">' + money(r.monthlySurplus) + '</p></div><div class="bg-slate-50 rounded-xl p-3"><span class="text-xs text-slate-500">Savings</span><p class="font-bold">' + money(r.currentSavings) + '</p></div></div></div>');

  if (r.targetDate && r.monthsUntilTarget > 0) {
    sections.push('<div class="border-t pt-3 mb-4"><h4 class="font-bold text-sm mb-2">Target plan</h4><div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm"><div class="rounded-xl p-3 border"><span class="text-xs text-slate-500">Funding gap</span><p class="font-bold">' + money(r.fundingGap) + '</p></div><div class="rounded-xl p-3 border"><span class="text-xs text-slate-500">Months to target</span><p class="font-bold">' + r.monthsUntilTarget + '</p></div><div class="rounded-xl p-3 border"><span class="text-xs text-slate-500">Needed monthly</span><p class="font-bold">' + money(r.requiredMonthlySaving) + '</p></div><div class="rounded-xl p-3 border"><span class="text-xs text-slate-500">Available monthly</span><p class="font-bold">' + money(r.monthlySurplus) + '</p></div></div></div>');
  }

  if (r.activeGoalCommitments > 0) {
    sections.push('<div class="border-t pt-3 mb-4"><h4 class="font-bold text-sm mb-2">Effect on goals</h4><div class="grid grid-cols-2 gap-3 text-sm"><div class="bg-amber-50 rounded-xl p-3"><span class="text-xs text-amber-700">Goal commitments</span><p class="font-bold">' + money(r.activeGoalCommitments) + ' / month</p></div><div class="' + (r.availableAfterGoals >= 0 ? 'bg-emerald-50' : 'bg-red-50') + ' rounded-xl p-3"><span class="text-xs ' + (r.availableAfterGoals >= 0 ? 'text-emerald-700' : 'text-red-700') + '">Available after goals</span><p class="font-bold">' + money(r.availableAfterGoals) + ' / month</p></div></div></div>');
  }

  sections.push('<div class="border-t pt-3 mb-4"><h4 class="font-bold text-sm mb-2">Why</h4><ul class="space-y-1 text-sm">');
  if (r.reasons && r.reasons.length) {
    r.reasons.forEach(function(reason) {
      sections.push('<li class="flex items-start gap-2"><span class="text-slate-400 mt-0.5">&bull;</span><span>' + escapeHtml(reason) + '</span></li>');
    });
  }
  sections.push('</ul></div>');

  sections.push('<div class="border-t pt-3 mb-4"><h4 class="font-bold text-sm mb-2">Next step</h4><div class="bg-' + (r.riskLevel === 'High' ? 'red' : r.riskLevel === 'Medium' ? 'amber' : 'emerald') + '-50 rounded-xl p-3 text-sm"><p>' + escapeHtml(r.recommendation) + '</p></div></div>');

  sections.push('<div class="border-t pt-3"><details class="text-sm"><summary class="cursor-pointer text-slate-500 font-bold text-xs">How this was calculated</summary><p class="mt-2 text-xs text-slate-500 leading-relaxed">' + escapeHtml(r.explanation) + '</p><p class="mt-1 text-xs text-slate-500">' + escapeHtml(r.expensePeriod) + '. Emergency target: ' + money(r.emergencyTarget) + '.</p></details></div>');

  return sections.join('\n');
}

function bindAffordability() {
  var form = document.getElementById('affordabilityForm');
  if (!form) return;
  if (form.dataset.bound === 'true') return;
  form.dataset.bound = 'true';

  var btn = form.querySelector('button[type="submit"]');
  var resultEl = document.getElementById('affordabilityResult');
  var loadingEl = document.getElementById('affordabilityLoading');

  form.addEventListener('submit', async function(e) {
    e.preventDefault();
    var inName = document.getElementById('itemName');
    var inPrice = document.getElementById('itemPrice');
    var inDate = document.getElementById('targetDate');

    var name = inName ? inName.value.trim() : '';
    var price = Number(inPrice ? inPrice.value : 0);
    var targetDate = inDate ? inDate.value || '' : '';

    if (!name) { showToast('Please enter an item name.', true); return; }
    if (isNaN(price) || price <= 0) { showToast('Please enter a valid price.', true); return; }

    if (resultEl) {
      resultEl.innerHTML = '';
      resultEl.classList.add('opacity-50');
    }
    if (loadingEl) loadingEl.classList.remove('hidden');
    if (btn) {
      btn.disabled = true;
      btn.textContent = 'Checking...';
    }

    try {
      var body = { itemName: name, itemPrice: price };
      if (targetDate) body.targetDate = targetDate;
      var res = await api('/api/affordability', { method: 'POST', body: JSON.stringify(body) });
      var r = res.data;
      if (resultEl) {
        resultEl.classList.remove('opacity-50');
        resultEl.innerHTML = renderAffordabilityResult(r);
      }
    } catch (err) {
      if (resultEl) {
        resultEl.classList.remove('opacity-50');
        resultEl.innerHTML = '<div class="bg-red-50 border border-red-200 rounded-xl p-4 text-sm text-red-700">' + escapeHtml(err.message) + '</div>';
      }
      showToast(err.message, true);
    } finally {
      if (loadingEl) loadingEl.classList.add('hidden');
      if (btn) {
        btn.disabled = false;
        btn.textContent = 'Check';
      }
    }
  });
}
