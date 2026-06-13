(function() {
  'use strict';

  const currentUser = () => window.__app?.currentUser || window.__demoUser || { username: 'Guest', email: '' };

  function getInitials(name) {
    if (!name || name === 'Guest') return '?';
    return name.split(' ').map(s => s[0]).slice(0, 2).join('').toUpperCase();
  }

  function sheetOverlay() {
    let el = document.getElementById('profileSheetOverlay');
    if (!el) {
      el = document.createElement('div');
      el.id = 'profileSheetOverlay';
      el.className = 'fixed inset-0 bg-black/50 z-40 transition-opacity';
      el.style.opacity = '0';
      el.style.pointerEvents = 'none';
      el.addEventListener('click', closeProfileMenu);
      document.body.appendChild(el);
    }
    return el;
  }

  function openProfileMenu(e) {
    e.stopPropagation();
    const isMobile = window.innerWidth < 768;
    if (isMobile) {
      renderMobileSheet();
    } else {
      renderDesktopDropdown();
    }
  }

  function closeProfileMenu() {
    const dd = document.getElementById('profileDropdown');
    if (dd) dd.remove();
    const sheet = document.getElementById('profileSheet');
    if (sheet) {
      sheet.classList.remove('open');
      setTimeout(() => sheet.remove(), 300);
    }
    const overlay = document.getElementById('profileSheetOverlay');
    if (overlay) {
      overlay.style.opacity = '0';
      overlay.style.pointerEvents = 'none';
    }
  }

  function renderDesktopDropdown() {
    closeProfileMenu();
    const user = currentUser();
    const div = document.createElement('div');
    div.id = 'profileDropdown';
    div.className = 'profile-dropdown';
    div.setAttribute('role', 'menu');
    div.innerHTML = `
      <div class="profile-dropdown-header">
        <div class="profile-btn" style="width:40px;height:40px;font-size:.75rem;cursor:default">${getInitials(user.username)}</div>
        <div>
          <h3 class="font-bold text-sm">${escapeHtml(user.username)}</h3>
          <p class="text-xs text-slate-500">${escapeHtml(user.email || 'No email')}</p>
        </div>
      </div>
      <div class="profile-dropdown-body">
        <button class="profile-dropdown-item" data-action="edit-profile" role="menuitem">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
          Edit Profile
        </button>
        <button class="profile-dropdown-item" data-action="switch-user" role="menuitem">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          Switch User
        </button>
        <button class="profile-dropdown-item danger" data-action="logout" role="menuitem">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          Logout
        </button>
      </div>
    `;
    div.addEventListener('click', e => {
      const item = e.target.closest('[data-action]');
      if (!item) return;
      e.stopPropagation();
      handleAction(item.dataset.action);
    });
    document.addEventListener('click', closeProfileMenu, { once: true });
    const btn = document.getElementById('profileBtn');
    if (btn) {
      btn.parentElement.style.position = 'relative';
      btn.parentElement.appendChild(div);
    }
  }

  function renderMobileSheet() {
    closeProfileMenu();
    const user = currentUser();
    const overlay = sheetOverlay();
    overlay.style.opacity = '1';
    overlay.style.pointerEvents = 'auto';

    const sheet = document.createElement('div');
    sheet.id = 'profileSheet';
    sheet.className = 'profile-sheet';
    sheet.innerHTML = `
      <div class="profile-sheet-handle"></div>
      <div class="profile-dropdown-header">
        <div class="profile-btn" style="width:40px;height:40px;font-size:.75rem;cursor:default">${getInitials(user.username)}</div>
        <div>
          <h3 class="font-bold text-sm">${escapeHtml(user.username)}</h3>
          <p class="text-xs text-slate-500">${escapeHtml(user.email || 'No email')}</p>
        </div>
      </div>
      <div class="profile-dropdown-body">
        <button class="profile-dropdown-item" data-action="edit-profile">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
          Edit Profile
        </button>
        <button class="profile-dropdown-item" data-action="switch-user">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          Switch User
        </button>
        <button class="profile-dropdown-item danger" data-action="logout">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/></svg>
          Logout
        </button>
      </div>
    `;
    sheet.addEventListener('click', e => {
      const item = e.target.closest('[data-action]');
      if (!item) return;
      e.stopPropagation();
      handleAction(item.dataset.action);
    });
    document.body.appendChild(sheet);
    requestAnimationFrame(() => sheet.classList.add('open'));
  }

  function handleAction(action) {
    closeProfileMenu();
    switch (action) {
      case 'edit-profile':
        showEditProfileModal();
        break;
      case 'switch-user':
        triggerSwitchUser();
        break;
      case 'logout':
        triggerLogout();
        break;
    }
  }

  function showEditProfileModal() {
    const user = currentUser();
    const overlay = document.createElement('div');
    overlay.className = 'fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4';
    overlay.style.animation = 'fadeIn .2s';
    overlay.innerHTML = `
      <div class="card p-6 w-full max-w-md" onclick="event.stopPropagation()">
        <h2 class="text-lg font-bold mb-4">Edit Profile</h2>
        <form id="editProfileForm">
          <div class="mb-4">
            <label class="block text-sm font-bold mb-1" for="editName">Name</label>
            <input type="text" id="editName" class="input" value="${escapeHtml(user.username)}" required>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-bold mb-1" for="editEmail">Email</label>
            <input type="email" id="editEmail" class="input" value="${escapeHtml(user.email || '')}">
          </div>
          <div class="flex gap-3 justify-end">
            <button type="button" class="btn btn-secondary" id="cancelEditProfile">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    `;
    document.body.appendChild(overlay);

    overlay.querySelector('#cancelEditProfile').onclick = () => overlay.remove();
    overlay.addEventListener('click', () => overlay.remove());
    overlay.querySelector('#editProfileForm').onsubmit = e => {
      e.preventDefault();
      const name = overlay.querySelector('#editName').value.trim();
      const email = overlay.querySelector('#editEmail').value.trim();
      if (!name) return;
      if (window.__app && window.__app.currentUser) {
        window.__app.currentUser.username = name;
        window.__app.currentUser.email = email;
      }
      if (window.__demoUser) {
        window.__demoUser.username = name;
        window.__demoUser.email = email;
      }
      updateProfileUI();
      overlay.remove();
      showToast('Profile updated');
    };
  }

  function triggerSwitchUser() {
    const btn = document.getElementById('switchUserBtn');
    if (btn) btn.click();
  }

  function triggerLogout() {
    const confirmed = confirm('Are you sure you want to logout?');
    if (!confirmed) return;
    const btn = document.getElementById('logoutBtn');
    if (btn) btn.click();
  }

  function updateProfileUI() {
    const user = currentUser();
    const initials = getInitials(user.username);
    const btn = document.getElementById('profileBtn');
    if (btn) btn.textContent = initials;
    const nameEl = document.getElementById('profileUserName');
    if (nameEl) nameEl.textContent = user.username;
  }

  function escapeHtml(str) {
    if (!str) return '';
    const d = document.createElement('div');
    d.textContent = str;
    return d.innerHTML;
  }

  function showToast(msg) {
    let t = document.getElementById('profileToast');
    if (t) t.remove();
    t = document.createElement('div');
    t.id = 'profileToast';
    t.className = 'toast card p-3 text-sm font-bold';
    t.textContent = msg;
    document.body.appendChild(t);
    setTimeout(() => t.remove(), 2500);
  }

  document.addEventListener('DOMContentLoaded', () => {
    const btn = document.getElementById('profileBtn');
    if (btn) btn.addEventListener('click', openProfileMenu);
    updateProfileUI();
  });

  window.__profile = { updateProfileUI, getInitials };
})();
