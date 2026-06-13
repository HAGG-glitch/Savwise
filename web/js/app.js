window.__app = {};

async function refreshAll() {
  const uid = getActiveUserId();
  if (!uid) return;
  try {
    const profile = await api('/api/profile');
    const u = profile.data;
    window.__app.currentUser = u;
    var fn = document.getElementById('fullName');
    if (fn) fn.value = u.fullName || '';
    var mi = document.getElementById('monthlyIncome');
    if (mi) mi.value = u.monthlyIncome || 0;
    var cs = document.getElementById('currentSavings');
    if (cs) cs.value = u.currentSavings || 0;
    var et = document.getElementById('emergencyTarget');
    if (et) et.value = u.emergencyTarget || 1500;
    var pl = document.getElementById('preferredLanguage');
    if (pl) pl.value = u.preferredLanguage || 'English';
    var ca = document.getElementById('consentAccepted');
    if (ca) ca.checked = !!u.consentAccepted;
    updateConsentUI(u.consentAccepted, u.fullName);
    if (window.__profile) window.__profile.updateProfileUI();
    await Promise.all([loadDashboard(), loadTransactions(), loadGoals()]);
  } catch (err) {
    var statusEl = document.getElementById('status');
    if (statusEl) statusEl.innerHTML = '<span class="badge badge-high">Connection error</span> <span class="text-red-700 ml-2">' + err.message + '</span>';
  }
}
function updateConsentUI(accepted, userName) {
  const overlay = document.getElementById('consentOverlay');
  const appContent = document.getElementById('appContent');
  const onboarding = document.getElementById('onboardingScreen');
  const statusEl = document.getElementById('status');
  const consentStatus = document.getElementById('consentStatusDisplay');
  const userNameDisplay = document.getElementById('userNameDisplay');
  const userNameDisplay2 = document.getElementById('userNameDisplay2');

  if (consentStatus) consentStatus.textContent = accepted ? 'Yes' : 'No';
  if (userNameDisplay && userName) userNameDisplay.textContent = userName;
  if (userNameDisplay2 && userName) userNameDisplay2.textContent = userName;

  if (accepted) {
    if (overlay) overlay.classList.add('hidden');
    if (appContent) appContent.classList.remove('hidden');
    if (onboarding) onboarding.classList.add('hidden');
    var loadDemo = document.getElementById('loadDemoBtn');
    if (loadDemo) loadDemo.classList.remove('hidden');
    statusEl.innerHTML = '<span class="badge badge-positive">Consent: Yes</span> <span class="text-slate-600 ml-2">Mode: Demo / manual</span> <span class="text-slate-600 ml-2">Mobile money: No</span>';
  } else {
    if (overlay) overlay.classList.remove('hidden');
    if (appContent) appContent.classList.add('hidden');
    statusEl.innerHTML = '<span class="badge badge-medium">Consent pending</span> <span class="text-slate-600 ml-2">Accept the privacy notice to use the app.</span>';
  }
}
async function bootApp() {
  const uid = getActiveUserId();
  const email = getActiveUserEmail();
  const onboarding = document.getElementById('onboardingScreen');
  const appContent = document.getElementById('appContent');

  if (uid && email) {
    try {
      const res = await api('/api/current-user?email=' + encodeURIComponent(email));
      if (res.data && res.data.id === uid) {
        setActiveUser(uid, email);
        if (onboarding) onboarding.classList.add('hidden');
        if (appContent) appContent.classList.remove('hidden');
        bindTabs(); bindProfile(); bindConsent(); bindTransactions(); bindGoals(); bindAffordability(); bindCoach(); bindDataControls(); bindGeneral();
        await refreshAll();
        applyTheme();
        return;
      }
    } catch (e) {}
  }
  clearActiveUser();
  if (onboarding) onboarding.classList.remove('hidden');
  if (appContent) appContent.classList.add('hidden');
  bindOnboarding();
}
function bindOnboarding() {
  var form = document.getElementById('createUserForm');
  if (form) {
    form.addEventListener('submit', async function(e) {
      e.preventDefault();
      try {
        const body = {
          fullName: document.getElementById('onFullName').value.trim(),
          email: document.getElementById('onEmail').value.trim(),
          monthlyIncome: Number(document.getElementById('onMonthlyIncome').value) || 0,
          currentSavings: Number(document.getElementById('onCurrentSavings').value) || 0,
          emergencyTarget: Number(document.getElementById('onEmergencyTarget').value) || 1500,
          preferredLanguage: document.getElementById('onLanguage').value,
          consentAccepted: document.getElementById('onConsent').checked,
        };
        const res = await api('/api/users', { method: 'POST', body: JSON.stringify(body) });
        const u = res.data;
        setActiveUser(u.id, u.email);
        showToast('Welcome, ' + u.fullName + '!');
        document.getElementById('onboardingScreen').classList.add('hidden');
        document.getElementById('appContent').classList.remove('hidden');
        bindTabs(); bindProfile(); bindConsent(); bindTransactions(); bindGoals(); bindAffordability(); bindCoach(); bindDataControls(); bindGeneral();
        await refreshAll();
        applyTheme();
      } catch (err) { showToast(err.message, true); }
    });
  }

  var demoBtn = document.getElementById('loadDemoConsentBtn');
  if (demoBtn) {
    demoBtn.addEventListener('click', async function() {
      if (!confirm('Load demo data? This creates a demo user with sample records.')) return;
      try {
        var res = await api('/api/users', {
          method: 'POST',
          body: JSON.stringify({ fullName: 'Demo User', email: 'demo_' + Date.now() + '@demo', monthlyIncome: 2500, currentSavings: 800, emergencyTarget: 1500, preferredLanguage: 'English', consentAccepted: true })
        });
        var u = res.data;
        setActiveUser(u.id, u.email);
        var demoRes = await api('/api/load-demo', { method: 'POST' });
        showToast('Demo data loaded');
        document.getElementById('onboardingScreen').classList.add('hidden');
        document.getElementById('appContent').classList.remove('hidden');
        bindTabs(); bindProfile(); bindConsent(); bindTransactions(); bindGoals(); bindAffordability(); bindCoach(); bindDataControls(); bindGeneral();
        await refreshAll();
        applyTheme();
      } catch (err) { showToast(err.message, true); }
    });
  }
}
function showApp() {
  var ob = document.getElementById('onboardingScreen');
  if (ob) ob.classList.add('hidden');
  var ac = document.getElementById('appContent');
  if (ac) ac.classList.remove('hidden');
  bindTabs(); bindProfile(); bindConsent(); bindTransactions(); bindGoals(); bindAffordability(); bindCoach(); bindDataControls(); bindGeneral();
  refreshAll();
  applyTheme();
}
function bindTabs() {
  if (document.body.dataset.tabsBound === 'true') return;
  document.body.dataset.tabsBound = 'true';
  document.querySelectorAll('.tab-btn').forEach(function(btn) {
    btn.addEventListener('click', function() {
      document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
      btn.classList.add('active');
      document.querySelectorAll('.tab-section').forEach(function(s) { s.classList.add('hidden'); });
      var target = document.getElementById(btn.dataset.tab);
      if (target) target.classList.remove('hidden');
    });
  });
}
function bindProfile() {
  var pf = document.getElementById('profileForm');
  if (pf) {
    pf.addEventListener('submit', async function(e) {
      e.preventDefault();
      try {
        const body = {
          user_id: getActiveUserId(),
          fullName: document.getElementById('fullName').value.trim(),
          email: getActiveUserEmail(),
          preferredLanguage: document.getElementById('preferredLanguage').value,
          monthlyIncome: Number(document.getElementById('monthlyIncome').value),
          currentSavings: Number(document.getElementById('currentSavings').value),
          emergencyTarget: Number(document.getElementById('emergencyTarget').value),
          consentAccepted: document.getElementById('consentAccepted').checked,
        };
        await api('/api/profile', { method: 'POST', body: JSON.stringify(body) });
        showToast('Profile saved');
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
}
function bindConsent() {
  const overlayBtn = document.getElementById('overlayConsentBtn');
  if (overlayBtn) {
    overlayBtn.addEventListener('click', function() {
      var co = document.getElementById('consentOverlay');
      if (co) co.classList.add('hidden');
      var privacyTab = document.querySelector('.tab-btn[data-tab="privacy"]');
      if (privacyTab) privacyTab.click();
    });
  }
}
function bindGeneral() {
  var loadBtn = document.getElementById('loadDemoBtn');
  if (loadBtn) {
    loadBtn.addEventListener('click', async function() {
      if (!confirm('This will replace the current user records. Continue?')) return;
      try {
        await api('/api/load-demo', { method: 'POST' });
        showToast('Demo data loaded');
        await refreshAll();
      } catch (err) { showToast(err.message, true); }
    });
  }
  var refreshBtn = document.getElementById('refreshBtn');
  if (refreshBtn) refreshBtn.addEventListener('click', refreshAll);
  var switchBtn = document.getElementById('switchUserBtn');
  if (switchBtn) {
    switchBtn.addEventListener('click', function() {
      clearActiveUser();
      location.reload();
    });
  }
  var logoutBtn = document.getElementById('logoutBtn');
  if (logoutBtn) {
    logoutBtn.addEventListener('click', function() {
      clearActiveUser();
      location.reload();
    });
  }

}
document.addEventListener('DOMContentLoaded', bootApp);
