/* Aurora 3.1 — 渐进增强脚本：无 JS 时页面仍完整可用 */
(function () {
  'use strict';
  var root = document.documentElement;
  root.classList.add('js');

  /* ---------- 亮暗主题：auto → light → dark 三态循环 ---------- */
  function applyTheme(t) {
    if (t === 'light' || t === 'dark') root.dataset.theme = t;
    else delete root.dataset.theme;
    var metas = document.querySelectorAll('meta[name="theme-color"]');
    for (var i = 0; i < metas.length; i++) {
      if (!metas[i].dataset.def) metas[i].dataset.def = metas[i].content;
      metas[i].content = t ? (t === 'dark' ? '#0b0e13' : '#ffffff') : metas[i].dataset.def;
    }
    var btn = document.querySelector('.theme-btn');
    if (btn) {
      var label = t === 'light' ? '当前浅色模式，点击切换为深色模式' : t === 'dark' ? '当前深色模式，点击切换为跟随系统' : '当前跟随系统，点击切换为浅色模式';
      btn.setAttribute('aria-label', label);
      btn.setAttribute('title', label);
    }
  }
  window.nudTheme = function () {
    var t = null;
    try { t = localStorage.getItem('nud-theme'); } catch (e) {}
    t = !t ? 'light' : t === 'light' ? 'dark' : null;
    try {
      if (t) localStorage.setItem('nud-theme', t);
      else localStorage.removeItem('nud-theme');
    } catch (e) {}
    applyTheme(t);
  };

  document.addEventListener('DOMContentLoaded', function () {
    var savedTheme = null;
    try { savedTheme = localStorage.getItem('nud-theme'); } catch (e) {}
    applyTheme(savedTheme);

    /* ---------- 顶栏导航当前项 ---------- */
    var path = location.pathname;
    var links = document.querySelectorAll('[data-nav]');
    for (var i = 0; i < links.length; i++) {
      var href = links[i].getAttribute('href').split('?')[0];
      var match = href === '/'
        ? path === '/'
        : path === href || path.indexOf(href + '/') === 0 ||
          (href === '/apps' && (path.indexOf('/apps') === 0 || path.indexOf('/categories') === 0 || path === '/search'));
      if (links[i].getAttribute('href').indexOf('sort=hot') > -1) {
        match = location.search.indexOf('sort=hot') > -1;
      } else if (match && href === '/apps' && location.search.indexOf('sort=hot') > -1) {
        match = false;
      }
      if (match) links[i].classList.add('active');
    }

    /* ---------- 移动端抽屉 ---------- */
    var burger = document.querySelector('.nav-burger');
    if (burger) {
      function closeNav() {
        root.classList.remove('nav-open');
        burger.setAttribute('aria-expanded', 'false');
      }
      burger.addEventListener('click', function () {
        var open = root.classList.toggle('nav-open');
        burger.setAttribute('aria-expanded', open ? 'true' : 'false');
      });
      document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') closeNav();
      });
      var drawerLinks = document.querySelectorAll('.nav-drawer a');
      for (var j = 0; j < drawerLinks.length; j++) drawerLinks[j].addEventListener('click', closeNav);
    }

    /* ---------- 复制（data-copy） ---------- */
    document.addEventListener('click', function (e) {
      var btn = e.target.closest ? e.target.closest('[data-copy]') : null;
      if (!btn || !navigator.clipboard) return;
      navigator.clipboard.writeText(btn.getAttribute('data-copy')).then(function () {
        btn.classList.add('copied');
        var label = btn.querySelector('[data-copy-label]');
        if (label && !label.dataset.orig) {
          label.dataset.orig = label.textContent;
          label.textContent = '已复制';
        }
        setTimeout(function () {
          btn.classList.remove('copied');
          if (label && label.dataset.orig) {
            label.textContent = label.dataset.orig;
            delete label.dataset.orig;
          }
        }, 1600);
      });
    });

    /* ---------- 筛选表单自动提交 ---------- */
    document.addEventListener('change', function (e) {
      var el = e.target;
      if (el && el.hasAttribute && el.hasAttribute('data-autosubmit') && el.form) el.form.submit();
    });

    /* ---------- 截图灯箱 ---------- */
    var box = document.getElementById('lightbox');
    if (box) {
      document.addEventListener('click', function (e) {
        var shot = e.target.closest ? e.target.closest('.shot') : null;
        if (shot) {
          box.querySelector('img').src = shot.getAttribute('data-full');
          if (box.showModal) box.showModal();
        } else if (e.target === box || (e.target.closest && e.target.closest('[data-close]'))) {
          box.close && box.close();
        }
      });
    }

    /* ---------- 按 / 聚焦搜索 ---------- */
    document.addEventListener('keydown', function (e) {
      if (e.key !== '/' || /INPUT|TEXTAREA|SELECT/.test(document.activeElement.tagName)) return;
      var input = document.querySelector('[data-search]');
      if (input) { e.preventDefault(); input.focus(); }
    });
  });
})();
