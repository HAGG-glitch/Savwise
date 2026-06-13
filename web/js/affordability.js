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

function riskLevelClass(level) {
  if (level === 'High') return 'risk-high';
  if (level === 'Medium') return 'risk-medium';
  return 'risk-low';
}

function riskIcon(level) {
  if (level === 'High') return '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"/></svg>';
  if (level === 'Medium') return '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>';
  return '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>';
}

function renderAffordabilityResult(r) {
  var riskClass = riskLevelClass(r.riskLevel);
  var icon = riskIcon(r.riskLevel);
  var targetDisplay = r.targetDate ? formatDate(r.targetDate) : 'No target date set';

  var sections = [];

  sections.push('<div class="mb-4"><div class="risk-indicator ' + riskClass + '">' + icon + ' <span>' + escapeHtml(r.riskLevel) + ' Risk</span></div></div>');

  sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">Quick result</h3><div class="grid grid-cols-2 gap-3">');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Item</span><p class="font-bold">' + escapeHtml(r.itemName) + '</p></div>');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Price</span><p class="font-bold text-lg">' + money(r.itemPrice) + '</p></div>');
  if (r.targetDate) {
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Target date</span><p class="font-bold">' + targetDisplay + '</p></div>');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Calculated</span><p class="font-bold text-xs">' + escapeHtml(r.calculatedAt) + '</p></div>');
  }
  sections.push('</div></div>');

  sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">Your numbers</h3><div class="grid grid-cols-2 md:grid-cols-4 gap-3">');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Monthly income</span><p class="font-bold text-emerald-700">' + money(r.monthlyIncome) + '</p></div>');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Monthly expenses</span><p class="font-bold text-red-600">' + money(r.monthlyExpenses) + '</p></div>');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Monthly surplus</span><p class="font-bold ' + (r.monthlySurplus >= 0 ? 'text-emerald-700' : 'text-red-600') + '">' + money(r.monthlySurplus) + '</p></div>');
  sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Current savings</span><p class="font-bold">' + money(r.currentSavings) + '</p></div>');
  sections.push('</div></div>');

  if (r.targetDate && r.monthsUntilTarget > 0) {
    sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">Target plan</h3><div class="grid grid-cols-2 md:grid-cols-4 gap-3">');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Funding gap</span><p class="font-bold">' + money(r.fundingGap) + '</p></div>');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Months to target</span><p class="font-bold">' + r.monthsUntilTarget + '</p></div>');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Required monthly</span><p class="font-bold">' + money(r.requiredMonthlySaving) + '</p></div>');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Available surplus</span><p class="font-bold">' + money(r.monthlySurplus) + '</p></div>');
    sections.push('</div></div>');
  }

  if (r.activeGoalCommitments > 0) {
    sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">Effect on your goals</h3><div class="grid grid-cols-2 gap-3">');
    sections.push('<div class="stat-box"><span class="text-xs text-slate-500">Active goal commitments</span><p class="font-bold">' + money(r.activeGoalCommitments) + ' / month</p></div>');
    sections.push('<div class="stat-box ' + (r.availableAfterGoals >= 0 ? 'risk-low' : 'risk-high') + '"><span class="text-xs ' + (r.availableAfterGoals >= 0 ? 'text-emerald-700' : 'text-red-700') + '">Available after goals</span><p class="font-bold">' + money(r.availableAfterGoals) + ' / month</p></div>');
    sections.push('</div>');
    if (r.goalImpact && r.goalImpact !== 'none') {
      sections.push('<p class="text-xs text-slate-500 mt-2">' + escapeHtml(r.goalImpact) + '.</p>');
    }
    sections.push('</div>');
  }

  sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">Why</h3><ul class="space-y-2">');
  if (r.reasons && r.reasons.length) {
    r.reasons.forEach(function(reason) {
      sections.push('<li class="flex items-start gap-2 text-sm"><span class="text-slate-400 mt-1 flex-shrink-0">&bull;</span><span>' + escapeHtml(reason) + '</span></li>');
    });
  }
  sections.push('</ul></div>');

  sections.push('<div class="result-section"><h3 class="font-bold text-sm text-slate-500 uppercase tracking-wide mb-3">What to do next</h3><div class="rounded-xl p-4 text-sm ' + riskClass + '"><p>' + escapeHtml(r.recommendation) + '</p></div></div>');

  sections.push('<div class="result-section"><details class="text-sm"><summary class="cursor-pointer text-slate-500 font-semibold text-xs uppercase tracking-wide">How this was calculated</summary><div class="mt-3 space-y-2 text-xs text-slate-500 leading-relaxed"><p>' + escapeHtml(r.explanation) + '</p><p>Period: ' + escapeHtml(r.expensePeriod) + '. Emergency target: ' + money(r.emergencyTarget) + '.</p><p class="text-xs text-slate-400 mt-2">Educational estimate based on the financial information entered in SavWise.</p></div></details></div>');

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
        resultEl.innerHTML = '<div class="bg-red-50 border border-red-200 rounded-xl p-4 text-sm text-red-700"><p class="font-bold">Error</p><p>' + escapeHtml(err.message) + '</p></div>';
      }
      showToast(err.message, true);
    } finally {
      if (loadingEl) loadingEl.classList.add('hidden');
      if (btn) {
        btn.disabled = false;
        btn.textContent = 'Check Affordability';
      }
    }
  });
}
