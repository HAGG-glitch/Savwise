function showImportPreview(data) {
  const el = document.getElementById('importPreview');
  if (!el) return;
  el.classList.remove('hidden');
  let html = '<div class="rounded-xl border p-4 bg-emerald-50"><h3 class="font-bold text-sm">Import Summary</h3>';
  html += '<div class="grid grid-cols-2 md:grid-cols-4 gap-3 mt-3 text-sm">';
  html += '<div><span class="text-slate-500">Rows found</span><br><strong>' + (data.totalRows || 0) + '</strong></div>';
  html += '<div><span class="text-slate-500">Imported</span><br><strong>' + (data.transactions || 0) + '</strong></div>';
  html += '<div><span class="text-slate-500">Income</span><br><strong>' + money(data.totalIncome || 0) + '</strong></div>';
  html += '<div><span class="text-slate-500">Expenses</span><br><strong>' + money(data.totalExpenses || 0) + '</strong></div>';
  html += '</div>';
  if (data.errors && data.errors.length > 0) {
    html += '<div class="mt-3 p-3 bg-red-50 rounded-lg"><p class="text-xs font-bold text-red-700">Errors:</p><ul class="list-disc pl-4 mt-1 text-xs text-red-600">';
    data.errors.forEach(function(e) { html += '<li>' + e + '</li>'; });
    html += '</ul></div>';
  }
  html += '</div>';
  el.innerHTML = html;
}
function bindDataControls() {
  var uid = getActiveUserId();
  var exportJson = document.getElementById('exportJsonLink');
  var exportCsv = document.getElementById('exportCsvLink');
  if (exportJson) exportJson.href = '/api/export/json?user_id=' + encodeURIComponent(uid);
  if (exportCsv) exportCsv.href = '/api/export/csv?user_id=' + encodeURIComponent(uid);

  var resetEl = document.getElementById('resetBtn');
  if (resetEl) {
    resetEl.addEventListener('click', async function() {
      if (!confirm('Delete all data for this user? This cannot be undone.')) return;
      await api('/api/reset', { method: 'DELETE' });
      showToast('Data reset');
      await refreshAll();
    });
  }
  var importJson = document.getElementById('importJson');
  if (importJson) {
    importJson.addEventListener('change', async function(e) {
      var file = e.target.files[0];
      if (!file) return;
      try {
        var text = await file.text();
        JSON.parse(text);
        await api('/api/import/json', { method: 'POST', body: text });
        showToast('JSON imported');
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
  var importCsv = document.getElementById('importCsv');
  if (importCsv) {
    importCsv.addEventListener('change', async function(e) {
      var file = e.target.files[0];
      if (!file) return;
      try {
        var formData = new FormData();
        formData.append('file', file);
        var url = '/api/import/csv?user_id=' + encodeURIComponent(getActiveUserId());
        var res = await fetch(url, { method: 'POST', body: formData });
        var data = await res.json();
        if (!res.ok) {
          showToast(data.message || 'CSV import failed', true);
          return;
        }
        showImportPreview(data.data);
        if (data.data && data.data.errors && data.data.errors.length > 0) {
          showToast('CSV imported with ' + data.data.errors.length + ' error(s)', true);
        } else {
          showToast('CSV imported');
        }
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
}
