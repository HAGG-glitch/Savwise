function getTheme() {
  return localStorage.getItem('savwise_theme') || 'light';
}
function setTheme(t) {
  localStorage.setItem('savwise_theme', t);
}
function applyTheme() {
  var t = getTheme();
  var html = document.documentElement;
  if (t === 'dark') {
    html.classList.add('dark');
    document.body.classList.add('bg-slate-900', 'text-slate-100');
    document.body.classList.remove('bg-slate-100', 'text-slate-900');
  } else {
    html.classList.remove('dark');
    document.body.classList.remove('bg-slate-900', 'text-slate-100');
    document.body.classList.add('bg-slate-100', 'text-slate-900');
  }
  var btn = document.getElementById('darkModeToggle');
  if (btn) btn.textContent = t === 'dark' ? '\u2600' : '\u263E';
}
function toggleTheme() {
  var t = getTheme() === 'dark' ? 'light' : 'dark';
  setTheme(t);
  applyTheme();
}
document.addEventListener('DOMContentLoaded', function() {
  applyTheme();
  var btn = document.getElementById('darkModeToggle');
  if (btn) btn.addEventListener('click', toggleTheme);
});
