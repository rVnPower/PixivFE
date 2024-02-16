// make 4xx 5xx responses swap in as well
addEventListener('htmx:beforeOnLoad', function (event) {
  event.detail.shouldSwap = true;
  event.detail.isError = false;
});

function closeNavigationMenu() {
  document.getElementById("sidebar-toggler").checked = false
}
// browser built-in navigation
addEventListener("popstate", (event) => {
  closeNavigationMenu()
});
// htmx triggered navigation
addEventListener("htmx:pushedIntoHistory", (event) => {
  closeNavigationMenu()
});