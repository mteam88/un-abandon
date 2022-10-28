function dynamicallySelectActiveTab() {
  var path = window.location.pathname;
  var activeTab = document.querySelector('a[href="' + path + '"]');
  activeTab.classList.add('active');
}
dynamicallySelectActiveTab();