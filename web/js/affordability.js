function bindAffordability() {
  var af = document.getElementById('affordabilityForm');
  if (af) {
    af.addEventListener('submit', async e => {
      e.preventDefault();
      try {
        var td = document.getElementById('targetDate');
        var targetDate = td ? td.value || '' : '';
        var inName = document.getElementById('itemName');
        var inPrice = document.getElementById('itemPrice');
        var body = { itemName: inName ? inName.value.trim() : '', itemPrice: Number(inPrice ? inPrice.value : 0) };
        if (targetDate) body.targetDate = targetDate;
        var res = await api('/api/affordability', { method: 'POST', body: JSON.stringify(body) });
        var r = res.data;
        var cls = r.riskLevel === 'High' ? 'badge-high' : r.riskLevel === 'Medium' ? 'badge-medium' : 'badge-low';
        var ar = document.getElementById('affordabilityResult');
        if (ar) ar.innerHTML = '<span class="badge ' + cls + '">' + r.riskLevel + ' risk</span><p class="mt-3">' + r.explanation + '</p><ul class="list-disc pl-5 mt-3 space-y-1">' + (r.reasons || []).map(function(x) { return '<li>' + x + '</li>'; }).join('') + '</ul><div class="bg-slate-50 rounded-xl p-3 mt-3"><strong>Recommendation:</strong> ' + r.recommendation + '</div>';
      } catch (err) { showToast(err.message, true); }
    });
  }
}
