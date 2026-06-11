function escapeHtml(text) {
  var div = document.createElement('div');
  div.appendChild(document.createTextNode(text));
  return div.innerHTML;
}
function safeFormatWizz(text) {
  var safe = escapeHtml(text);
  safe = safe.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>');
  var lines = safe.split('\n');
  var out = [];
  var inList = false;
  for (var i = 0; i < lines.length; i++) {
    var line = lines[i].trim();
    if (line === '') {
      if (inList) { out.push('</ul>'); inList = false; }
      out.push('<p class="mb-2"></p>');
      continue;
    }
    var bulletMatch = line.match(/^[-*]\s+(.+)/);
    if (bulletMatch) {
      if (!inList) { out.push('<ul class="list-disc pl-4 mb-2 space-y-1">'); inList = true; }
      out.push('<li>' + bulletMatch[1] + '</li>');
      continue;
    }
    var numMatch = line.match(/^\d+[.)]\s+(.+)/);
    if (numMatch) {
      if (!inList) { out.push('<ol class="list-decimal pl-4 mb-2 space-y-1">'); inList = true; }
      out.push('<li>' + numMatch[1] + '</li>');
      continue;
    }
    if (inList) { out.push('</ul>'); inList = false; }
    out.push('<p class="mb-1">' + line + '</p>');
  }
  if (inList) out.push('</ul>');
  return out.join('\n');
}
function appendChat(role, text, isHtml) {
  var log = document.getElementById('chatLog');
  if (!log) return;
  var div = document.createElement('div');
  if (role === 'user') {
    div.className = 'ml-auto max-w-[80%] bg-emerald-600 text-white p-3 rounded-xl';
    div.textContent = text;
  } else {
    div.className = 'max-w-[85%] bg-white border p-3 rounded-xl';
    if (isHtml) {
      div.innerHTML = text;
    } else {
      div.textContent = text;
    }
  }
  log.appendChild(div);
  log.scrollTop = log.scrollHeight;
}
function showTyping() {
  var log = document.getElementById('chatLog');
  if (!log) return;
  var div = document.createElement('div');
  div.id = 'typingIndicator';
  div.className = 'max-w-[85%] bg-white border p-3 rounded-xl text-slate-400 text-sm';
  div.textContent = 'Wizz is thinking...';
  log.appendChild(div);
  log.scrollTop = log.scrollHeight;
}
function removeTyping() {
  var el = document.getElementById('typingIndicator');
  if (el) el.remove();
}
async function askCoach(message) {
  appendChat('user', message);
  showTyping();
  try {
    var res = await api('/api/coach', { method: 'POST', body: JSON.stringify({ message: message }) });
    removeTyping();
    var formatted = safeFormatWizz(res.data.response);
    var sourceNote = '<p class="text-xs text-slate-400 mt-2">' + escapeHtml(res.data.source) + ' · Educational guidance only</p>';
    appendChat('ai', formatted + sourceNote, true);
  } catch (err) {
    removeTyping();
    appendChat('ai', 'Sorry, Wizz encountered an error: ' + escapeHtml(err.message), false);
  }
}
function bindCoach() {
  var cf = document.getElementById('coachForm');
  if (cf) {
    cf.addEventListener('submit', async function(e) {
      e.preventDefault();
      var input = document.getElementById('coachMessage');
      if (!input) return;
      var msg = input.value.trim();
      if (!msg) return;
      input.value = '';
      await askCoach(msg);
    });
  }
  document.querySelectorAll('.coach-chip').forEach(function(btn) {
    btn.addEventListener('click', function() { askCoach(btn.dataset.q); });
  });
  appendChat('ai', 'Kushe! I am Wizz. I give educational budgeting and savings guidance based on your data. Do not share PINs, OTPs, or passwords.', false);
}
