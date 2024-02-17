// make 4xx 5xx responses swap in as well
addEventListener('htmx:sendError', function (event) {
  document.write("A network error prevented an HTTP request from occurring, which is likely because the server is down. Check your Developer Console's Network tab to know what exactly is wrong.")
});
addEventListener('htmx:beforeOnLoad', function (event) {
  // console.log("%o", event)
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