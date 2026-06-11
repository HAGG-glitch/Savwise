function getActiveUserId() {
  return localStorage.getItem('savwise_active_user_id') || '';
}
function getActiveUserEmail() {
  return localStorage.getItem('savwise_active_user_email') || '';
}
function setActiveUser(id, email) {
  localStorage.setItem('savwise_active_user_id', id);
  localStorage.setItem('savwise_active_user_email', email);
}
function clearActiveUser() {
  localStorage.removeItem('savwise_active_user_id');
  localStorage.removeItem('savwise_active_user_email');
}

async function api(path, options = {}) {
  const uid = getActiveUserId();
  const method = (options.method || 'GET').toUpperCase();

  let url = path;
  if (uid) {
    const sep = path.includes('?') ? '&' : '?';
    url = path + sep + 'user_id=' + encodeURIComponent(uid);
  }

  const fetchOpts = {
    headers: { 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options,
  };

  if (options.body && uid && (method === 'POST' || method === 'PUT')) {
    try {
      const bodyObj = JSON.parse(options.body);
      if (!bodyObj.user_id) {
        bodyObj.user_id = uid;
        fetchOpts.body = JSON.stringify(bodyObj);
      }
    } catch (e) {}
  }

  const response = await fetch(url, fetchOpts);
  const type = response.headers.get('content-type') || '';
  const data = type.includes('application/json') ? await response.json() : await response.text();
  if (!response.ok) {
    const msg = data && data.message ? data.message : 'Request failed';
    throw new Error(msg);
  }
  return data;
}

function money(value) { return 'SLE ' + Number(value || 0).toLocaleString(undefined, {maximumFractionDigits:0}); }
function pct(value) { return Number(value || 0).toFixed(1) + '%'; }
function showToast(message, error) {
  if (error === undefined) error = false;
  const el = document.getElementById('toast');
  if (!el) return;
  el.textContent = message;
  el.className = 'toast card p-4 text-sm ' + (error ? 'text-red-700 bg-red-50' : 'text-emerald-800 bg-emerald-50');
  el.classList.remove('hidden');
  setTimeout(function() { el.classList.add('hidden'); }, 3500);
}
