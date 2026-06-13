(function () {
  'use strict';

  var button = document.getElementById('mobile-menu-button');
  var menu = document.getElementById('mobile-menu');

  if (!button || !menu) return;

  if (button.dataset.mobileNavBound === 'true') return;
  button.dataset.mobileNavBound = 'true';

  var isOpen = false;

  var hamburgerSvg = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="3" y1="6" x2="21" y2="6"></line><line x1="3" y1="12" x2="21" y2="12"></line><line x1="3" y1="18" x2="21" y2="18"></line></svg>';
  var closeSvg = '<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>';

  function openMobileMenu() {
    if (isOpen) return;
    isOpen = true;
    button.setAttribute('aria-expanded', 'true');
    button.setAttribute('aria-label', 'Close navigation menu');
    menu.classList.remove('hidden');
    button.innerHTML = closeSvg;
  }

  function closeMobileMenu() {
    if (!isOpen) return;
    isOpen = false;
    button.setAttribute('aria-expanded', 'false');
    button.setAttribute('aria-label', 'Open navigation menu');
    menu.classList.add('hidden');
    button.innerHTML = hamburgerSvg;
  }

  function toggleMobileMenu() {
    if (isOpen) {
      closeMobileMenu();
    } else {
      openMobileMenu();
    }
  }

  function getActiveTabId() {
    var active = document.querySelector('.tab-btn.active');
    return active ? active.getAttribute('data-tab') : null;
  }

  function updateMobileActiveState(tabId) {
    if (!tabId) return;
    var items = menu.querySelectorAll('.mobile-nav-item[data-tab]');
    items.forEach(function (item) {
      var isActive = item.getAttribute('data-tab') === tabId;
      item.classList.toggle('active', isActive);
      if (isActive) {
        item.setAttribute('aria-current', 'page');
      } else {
        item.removeAttribute('aria-current');
      }
    });
  }

  function syncActiveFromDesktop() {
    var tabId = getActiveTabId();
    updateMobileActiveState(tabId);
  }

  button.addEventListener('click', function (e) {
    e.stopPropagation();
    toggleMobileMenu();
  });

  menu.addEventListener('click', function (e) {
    var item = e.target.closest('.mobile-nav-item');
    if (!item) return;

    var tab = item.getAttribute('data-tab');
    var action = item.getAttribute('data-action');

    if (tab) {
      closeMobileMenu();
      var tabBtn = document.querySelector('.tab-btn[data-tab="' + tab + '"]');
      if (tabBtn) tabBtn.click();
      updateMobileActiveState(tab);
    } else if (action) {
      closeMobileMenu();
      switch (action) {
        case 'profile':
          var profileBtn = document.getElementById('profileBtn');
          if (profileBtn) profileBtn.click();
          break;
        case 'dark-mode':
          var darkBtn = document.getElementById('darkModeToggle');
          if (darkBtn) darkBtn.click();
          break;
        case 'logout':
          if (confirm('Are you sure you want to logout?')) {
            clearActiveUser();
            location.reload();
          }
          break;
      }
    }
  });

  document.addEventListener('click', function (e) {
    if (!isOpen) return;
    var header = document.querySelector('.mobile-app-header');
    if (header && header.contains(e.target)) return;
    if (menu && menu.contains(e.target)) return;
    closeMobileMenu();
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' && isOpen) {
      closeMobileMenu();
      button.focus();
    }
  });

  var desktopQuery = window.matchMedia('(min-width: 768px)');
  function handleDesktopChange(event) {
    if (event.matches) {
      closeMobileMenu();
    }
  }
  if (desktopQuery.addEventListener) {
    desktopQuery.addEventListener('change', handleDesktopChange);
  } else {
    desktopQuery.addListener(handleDesktopChange);
  }

  document.addEventListener('click', function (e) {
    var tabBtn = e.target.closest('.tab-btn');
    if (tabBtn && tabBtn.getAttribute('data-tab')) {
      updateMobileActiveState(tabBtn.getAttribute('data-tab'));
    }
  });

  syncActiveFromDesktop();

})();
