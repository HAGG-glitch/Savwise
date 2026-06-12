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
function appendUserChat(text) {
  var log = document.getElementById('chatLog');
  if (!log) return;
  var wrapper = document.createElement('div');
  wrapper.className = 'flex justify-end';
  var div = document.createElement('div');
  div.className = 'max-w-[80%] bg-emerald-600 text-white p-3 rounded-2xl rounded-br-sm text-sm';
  div.textContent = text;
  wrapper.appendChild(div);
  log.appendChild(wrapper);
  log.scrollTop = log.scrollHeight;
}
function appendWizzChat(text, isHtml) {
  var log = document.getElementById('chatLog');
  if (!log) return;
  var wrapper = document.createElement('div');
  wrapper.className = 'flex items-start gap-2';
  var avatar = document.createElement('div');
  avatar.className = 'wizz-avatar text-xs flex-shrink-0';
  avatar.textContent = 'W';
  wrapper.appendChild(avatar);
  var bubble = document.createElement('div');
  bubble.className = 'bg-white border rounded-2xl rounded-tl-sm p-3 text-sm max-w-[85%]';
  if (isHtml) {
    bubble.innerHTML = text;
  } else {
    bubble.textContent = text;
  }
  wrapper.appendChild(bubble);
  log.appendChild(wrapper);
  log.scrollTop = log.scrollHeight;
}
function showTyping() {
  var log = document.getElementById('chatLog');
  if (!log) return;
  var wrapper = document.createElement('div');
  wrapper.id = 'typingIndicator';
  wrapper.className = 'flex items-start gap-2';
  var avatar = document.createElement('div');
  avatar.className = 'wizz-avatar text-xs flex-shrink-0';
  avatar.textContent = 'W';
  wrapper.appendChild(avatar);
  var div = document.createElement('div');
  div.className = 'bg-white border rounded-2xl rounded-tl-sm p-3 text-sm text-slate-400';
  div.innerHTML = '<span class="loading-spinner"></span> Wizz is thinking...';
  wrapper.appendChild(div);
  log.appendChild(wrapper);
  log.scrollTop = log.scrollHeight;
}
function removeTyping() {
  var el = document.getElementById('typingIndicator');
  if (el) el.remove();
}
function getTimestamp() {
  var d = new Date();
  return d.getHours().toString().padStart(2,'0') + ':' + d.getMinutes().toString().padStart(2,'0');
}
async function askCoach(message) {
  appendUserChat(message);
  showTyping();
  try {
    var res = await api('/api/coach', { method: 'POST', body: JSON.stringify({ message: message }) });
    removeTyping();
    var formatted = safeFormatWizz(res.data.response);
    var meta = '<p class="text-xs text-slate-400 mt-2">' + escapeHtml(getTimestamp()) + ' &middot; ' + escapeHtml(res.data.source) + ' &middot; Educational guidance only</p>';
    appendWizzChat(formatted + meta, true);
  } catch (err) {
    removeTyping();
    appendWizzChat('Sorry, Wizz encountered an error: ' + escapeHtml(err.message), false);
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
  appendWizzChat('<p class="font-bold">Kushe! I am Wizz.</p><p class="mt-1">I give educational budgeting and savings guidance based on your data.</p><p class="text-xs text-slate-400 mt-2">Do not share PINs, OTPs, or passwords. Educational guidance only.</p>', true);
}
